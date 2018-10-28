package client

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"net/http"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

var (
	// ErrRegistrationFailed ...
	ErrRegistrationFailed = errors.New("registration failed")
	// ErrAuthenticationFailed ...
	ErrAuthenticationFailed = errors.New("authentication failed")
)

// New creates a new Online SPHINX Client.
func New(pst Poster, cfg Configuration, repo Repository) *Client {
	return &Client{
		poster: pst,
		config: cfg,
		repo:   repo,
	}
}

// Client represents an Online SPHINX Client
type Client struct {
	poster  Poster
	config  Configuration
	repo    Repository
	session *Session
}

// Poster provides a Post operation used e.g. http.DefaultClient
type Poster interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// Repository provides a basic user configuration repository interface
type Repository interface {
	Add(u User) error
	Get(username string) (User, error)
}

// Register will register a new user.
func (clt *Client) Register(username string) error {

	user, err := NewUser(username)
	if err != nil {
		return err
	}

	err = clt.repo.Add(user)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = contract.MarshalRegisterRequest(w, contract.RegisterRequest{CID: user.cID})
	if err != nil {
		return err
	}
	w.Flush()
	rd := bufio.NewReader(&buf)

	r, err := clt.poster.Post(clt.config.registerPath, clt.config.contentType, rd)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusCreated {
		return ErrRegistrationFailed
	}
	return nil
}

// Login runs the Online SPHINX login protocol
func (clt *Client) Login(username, pwd string) error {

	user, err := clt.repo.Get(username)
	if err != nil {
		return err
	}

	g := crypto.HashInGroup(pwd, clt.config.hash, user.q)
	cNonce, b, kinv, err := crypto.Blind(g, user.q, clt.config.bits)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = contract.MarshalExpKRequest(w, contract.ExpKRequest{CID: user.cID, CNonce: cNonce, B: b, Q: user.q})
	if err != nil {
		return err
	}
	w.Flush()
	rd := bufio.NewReader(&buf)

	r, err := clt.poster.Post(clt.config.expkPath, clt.config.contentType, rd)
	if err != nil {
		return err
	}

	expKResp, err := contract.UnmarshalExpKResponse(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	B0 := crypto.Unblind(expKResp.BD, kinv, user.q)
	SKi := new(big.Int)
	SKi.SetBytes(crypto.HmacData(clt.config.hash, expKResp.KV.Bytes(), user.cID.Bytes(), expKResp.SID.Bytes(), cNonce.Bytes(), expKResp.SNonce.Bytes()))
	mk := new(big.Int)
	mk.Mul(crypto.ExpInGroup(B0, user.k, user.q), expKResp.Q0)

	clt.session = NewSession(user, expKResp.SID, SKi, mk)

	return nil
}

// Challenge session key SKi
func (clt *Client) Challenge() error {

	if clt.session == nil {
		return ErrAuthenticationFailed
	}

	g, err := rand.Int(rand.Reader, clt.session.user.q)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = contract.MarshalChallengeRequest(w, contract.ChallengeRequest{G: g, Q: clt.session.user.q})
	if err != nil {
		return err
	}
	w.Flush()

	rd := bufio.NewReader(&buf)
	r, err := clt.poster.Post(clt.config.verifyPath, clt.config.contentType, rd)
	if err != nil {
		return err
	}

	response, err := contract.UnmarshalChallengeResponse(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	verifier := crypto.ExpInGroup(g, clt.session.ski, clt.session.user.q)
	if response.R.Cmp(verifier) != 0 {
		return ErrAuthenticationFailed
	}

	return nil
}

// GetMetadata ...
func (clt *Client) GetMetadata() ([]string, error) {

	mac := crypto.HmacData(clt.config.hash, clt.session.ski.Bytes(), clt.session.user.cID.Bytes(), clt.session.sID.Bytes())

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := contract.MarshalMetadataRequest(w, contract.MetadataRequest{MAC: mac})
	if err != nil {
		return nil, err
	}
	w.Flush()
	rd := bufio.NewReader(&buf)

	r, err := clt.poster.Post(clt.config.metadataPath, clt.config.contentType, rd)
	if err != nil {
		return nil, err
	}

	metaResp, err := contract.UnmarshalMetadataResponse(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return metaResp.Domains, nil
}

// Add ...
func (clt *Client) Add(domain string) error {
	mac := crypto.HmacData(clt.config.hash, clt.session.ski.Bytes(), clt.session.user.cID.Bytes(), []byte(domain))

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := contract.MarshalAddRequest(w, contract.AddRequest{
		Domain: domain,
		MAC:    mac,
	})
	w.Flush()

	rd := bufio.NewReader(&buf)
	_, err = clt.poster.Post(clt.config.addPath, clt.config.contentType, rd)
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (clt *Client) Get(domain string) error {
	return nil
}
