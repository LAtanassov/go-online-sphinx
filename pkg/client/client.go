package client

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"path"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
	"github.com/pkg/errors"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

var two = big.NewInt(2)

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
func (clt *Client) Register(user User) error {

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := contract.MarshalRegisterRequest(w, contract.RegisterRequest{CID: user.cID})
	if err != nil {
		return errors.Wrap(err, "failed to marshal RegisterRequest")
	}
	w.Flush()
	rd := bufio.NewReader(&buf)

	u, err := url.Parse(clt.config.baseURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse baseURL")
	}
	u.Path = path.Join(u.Path, clt.config.registerPath)
	r, err := clt.poster.Post(u.String(), clt.config.contentType, rd)
	if err != nil {
		return errors.Wrap(err, "failed to post RegisterRequest")
	}

	err = contract.UnmarshalIfError(r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal error from response")
	}

	err = clt.repo.Add(user)
	if err != nil {
		return errors.Wrap(err, "failed to add new user to repo")
	}

	return nil
}

// Login runs the Online SPHINX login protocol
func (clt *Client) Login(username, pwd string) error {

	user, err := clt.repo.Get(username)
	if err != nil {
		return errors.Wrap(err, "failed to get user from repo")
	}

	g := crypto.HashInGroup(pwd, clt.config.hash, user.q)

	cNonce, err := rand.Int(rand.Reader, user.q)
	if err != nil {
		return errors.Wrap(err, "failed to generate random cNonce")
	}

	k, err := rand.Int(rand.Reader, user.q)
	if err != nil {
		return errors.Wrap(err, "failed to generate random k")
	}

	kinv := new(big.Int)
	kinv.ModInverse(k, user.q)

	if kinv == nil {
		kinv = big.NewInt(0)
	}

	b := crypto.ExpInGroup(g, k, user.q)

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = contract.MarshalExpKRequest(w, contract.ExpKRequest{CID: user.cID, CNonce: cNonce, B: b, Q: user.q})
	if err != nil {
		return errors.Wrap(err, "failed to marshal ExpKRequest")
	}
	w.Flush()
	rd := bufio.NewReader(&buf)

	u, err := url.Parse(clt.config.baseURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse baseURL")
	}
	u.Path = path.Join(u.Path, clt.config.expkPath)
	r, err := clt.poster.Post(u.String(), clt.config.contentType, rd)
	if err != nil {
		return errors.Wrap(err, "failed to post ExpKRequest")
	}

	err = contract.UnmarshalIfError(r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal error")
	}

	expKResp, err := contract.UnmarshalExpKResponse(r.Body)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ExpKResponse")
	}
	defer r.Body.Close()

	B0 := crypto.ExpInGroup(expKResp.BD, kinv, user.q)

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
		return errors.Wrap(contract.ErrAuthenticationFailed, "client session missing")
	}

	g, err := rand.Int(rand.Reader, clt.session.user.q)
	if err != nil {
		return errors.Wrap(err, "failed to generate random g")
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = contract.MarshalChallengeRequest(w, contract.ChallengeRequest{G: g, Q: clt.session.user.q})
	if err != nil {
		return errors.Wrap(err, "failed to marshal ChallengeRequest")
	}
	w.Flush()

	rd := bufio.NewReader(&buf)

	u, err := url.Parse(clt.config.baseURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse baseURL")
	}
	u.Path = path.Join(u.Path, clt.config.challengePath)
	r, err := clt.poster.Post(u.String(), clt.config.contentType, rd)
	if err != nil {
		return errors.Wrap(err, "failed to post ChallengeRequest")
	}

	err = contract.UnmarshalIfError(r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal error")
	}

	response, err := contract.UnmarshalChallengeResponse(r.Body)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ChallengeResponse")
	}
	defer r.Body.Close()

	verifier := crypto.ExpInGroup(g, clt.session.ski, clt.session.user.q)
	if response.R.Cmp(verifier) != 0 {
		return errors.Wrap(contract.ErrAuthenticationFailed, "challenge-response with session key ski failed")
	}

	return nil
}

// GetMetadata ...
func (clt *Client) GetMetadata() ([]string, error) {

	if clt.session == nil {
		return nil, errors.Wrap(contract.ErrAuthenticationFailed, "client session missing")
	}

	mac := crypto.HmacData(clt.config.hash, clt.session.ski.Bytes(), []byte("metadata"))

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := contract.MarshalMetadataRequest(w, contract.MetadataRequest{MAC: mac})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal MetadataRequest")
	}
	w.Flush()
	rd := bufio.NewReader(&buf)

	u, err := url.Parse(clt.config.baseURL)
	u.Path = path.Join(u.Path, clt.config.metadataPath)
	r, err := clt.poster.Post(u.String(), clt.config.contentType, rd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to post MetadataRequest")
	}

	err = contract.UnmarshalIfError(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal error")
	}

	metaResp, err := contract.UnmarshalMetadataResponse(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal MetadataResponse")
	}
	defer r.Body.Close()
	return metaResp.Domains, nil

}

// Add ...
func (clt *Client) Add(domain string) error {

	if clt.session == nil {
		return errors.Wrap(contract.ErrAuthenticationFailed, "client session missing")
	}

	mac := crypto.HmacData(clt.config.hash, clt.session.ski.Bytes(), []byte(domain))

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err := contract.MarshalAddRequest(w, contract.AddRequest{
		Domain: domain,
		MAC:    mac,
	})
	w.Flush()
	if err != nil {
		return errors.Wrap(err, "failed to marshal AddRequest")
	}

	rd := bufio.NewReader(&buf)

	u, err := url.Parse(clt.config.baseURL)
	u.Path = path.Join(u.Path, clt.config.addPath)
	r, err := clt.poster.Post(u.String(), clt.config.contentType, rd)
	if err != nil {
		return errors.Wrap(err, "failed to post AddRequest")
	}

	err = contract.UnmarshalIfError(r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal error")
	}

	if r.StatusCode != http.StatusCreated {
		return errors.Wrap(contract.ErrAddVaultFailed, "failed to add vault")
	}

	return nil
}

// Get ...
func (clt *Client) Get(domain string) (string, error) {

	if clt.session == nil {
		return "", errors.Wrap(contract.ErrAuthenticationFailed, "client session missing")
	}

	k, err := rand.Int(rand.Reader, clt.session.user.q)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random k")
	}
	kinv := new(big.Int)
	kinv.ModInverse(k, clt.session.user.q)

	if kinv == nil {
		kinv = big.NewInt(0)
	}

	bmk := crypto.ExpInGroup(clt.session.mk, k, clt.session.user.q)

	mac := crypto.HmacData(clt.config.hash, clt.session.ski.Bytes(), bmk.Bytes())

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = contract.MarshalGetRequest(w, contract.GetRequest{
		Domain: domain,
		MAC:    mac,
		BMK:    bmk,
		Q:      clt.session.user.q,
	})
	w.Flush()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal GetRequest")
	}

	rd := bufio.NewReader(&buf)

	u, err := url.Parse(clt.config.baseURL)
	u.Path = path.Join(u.Path, clt.config.getPath)
	r, err := clt.poster.Post(u.String(), clt.config.contentType, rd)
	if err != nil {
		return "", errors.Wrap(err, "failed to post GetRequest")
	}

	err = contract.UnmarshalIfError(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal error")
	}

	getResp, err := contract.UnmarshalGetResponse(r.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal GetResponse")
	}
	defer r.Body.Close()

	B0 := crypto.ExpInGroup(getResp.Bj, kinv, clt.session.user.q)

	rwd := new(big.Int)
	rwd.Mul(crypto.ExpInGroup(B0, clt.session.user.k, clt.session.user.q), getResp.Qj)

	return rwd.Text(16), nil
}

// Logout ...
func (clt *Client) Logout() {
	clt.session = nil
}
