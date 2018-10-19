package client

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"net/http"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

// ErrRegistrationFailed ...
var ErrRegistrationFailed = errors.New("registration failed")

// ErrAuthenticationFailed ...
var ErrAuthenticationFailed = errors.New("authentication failed")

// Client represents an Online SPHINX Client
type Client struct {
	poster  Poster
	config  Configuration
	user    User
	session Session
}

// Poster provides a Post operation used e.g. http.DefaultClient
type Poster interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// New creates a new Online SPHINX Client.
func New(pst Poster, cfg Configuration, usr User) *Client {
	return &Client{
		poster:  pst,
		config:  cfg,
		user:    usr,
		session: Session{},
	}
}

// Register will register a new user.
func (clt *Client) Register(username string) error {
	buf, err := marshalRegisterRequest(username)
	if err != nil {
		return err
	}

	r, err := clt.poster.Post(clt.config.registerPath, clt.config.contentType, bytes.NewBuffer(buf))
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

	// TODO load user from local config
	g := crypto.HashInGroup(pwd, clt.config.hash, clt.user.q)
	cNonce, b, kinv, err := crypto.Blind(g, clt.user.q, clt.config.bits)
	if err != nil {
		return err
	}

	buf, err := marshalExpKRequest(clt.user.cID, cNonce, b, clt.user.q)
	if err != nil {
		return err
	}

	r, err := clt.poster.Post(clt.config.expkPath, clt.config.contentType, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	sID, sNonce, bd, kv, q0, err := unmarsalExpKResponse(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	B0 := crypto.Unblind(bd, kinv, clt.user.q)

	clt.session.ski = new(big.Int)
	clt.session.ski.SetBytes(crypto.HmacData(clt.config.hash, kv.Bytes(), clt.user.cID.Bytes(), sID.Bytes(), cNonce.Bytes(), sNonce.Bytes()))

	clt.session.mk = new(big.Int)
	clt.session.mk.Mul(crypto.ExpInGroup(B0, clt.user.k, clt.user.q), q0)
	clt.session.sID = sID
	return nil
}

// Verify session key SKi
func (clt *Client) Verify() error {

	g, err := rand.Int(rand.Reader, clt.user.q)
	if err != nil {
		return err
	}

	challenge, err := marshalVerifyRequest(g, clt.user.q)
	if err != nil {
		return err
	}

	r, err := clt.poster.Post(clt.config.verifyPath, clt.config.contentType, bytes.NewBuffer(challenge))
	if err != nil {
		return err
	}

	response, err := unmarsalVerifyResponse(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	verifier := crypto.ExpInGroup(g, clt.session.ski, clt.user.q)
	if response.Cmp(verifier) != 0 {
		return ErrAuthenticationFailed
	}

	return nil
}

// GetMetadata ...
func (clt *Client) GetMetadata() ([]Domain, error) {

	mac := crypto.HmacData(clt.config.hash, clt.session.ski.Bytes(), clt.user.cID.Bytes(), clt.session.sID.Bytes())
	req, err := marshalMetadataRequest(clt.user.cID.Text(16), hex.EncodeToString(mac))

	r, err := clt.poster.Post(clt.config.metadataPath, clt.config.contentType, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	domains, err := unmarsalMetadataResponse(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return domains, nil
}

// Add ...
func (clt *Client) Add() error {
	return nil
}

// Get ...
func (clt *Client) Get() error {
	return nil
}
