package client

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

// ErrRegistrationFailed ...
var ErrRegistrationFailed = errors.New("registration failed")

// ErrAuthenticationFailed ...
var ErrAuthenticationFailed = errors.New("authentication failed")

// Client represents an Online SPHINX Client
type Client struct {
	poster Poster
	config Configuration
}

// Poster provides a Post operation used e.g. http.DefaultClient
type Poster interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// New creates a new Online SPHINX Client.
func New(pst Poster, cfg Configuration) *Client {
	return &Client{
		poster: pst,
		config: cfg,
	}
}

// Register will register a new user.
func (clt *Client) Register(usr User) error {
	buf, err := marshalRegisterRequest(usr.username)
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
func (clt *Client) Login(usr User, pwd string) error {

	g := crypto.HashInGroup(pwd, clt.config.hash, usr.q)
	cNonce, b, kinv, err := crypto.Blind(g, usr.q, clt.config.bits)
	if err != nil {
		return err
	}

	buf, err := marshalExpKRequest(usr.cID, cNonce, b, usr.q)
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

	B0 := crypto.Unblind(bd, kinv, usr.q)

	mk := new(big.Int)
	mk.Mul(crypto.ExpInGroup(B0, usr.k, usr.q), q0)

	SKi := new(big.Int)
	SKi.SetBytes(crypto.HmacData(clt.config.hash, kv.Bytes(), usr.cID.Bytes(), sID.Bytes(), cNonce.Bytes(), sNonce.Bytes()))

	os.Setenv("SKi", SKi.Text(16))
	os.Setenv("mk", mk.Text(16))
	os.Setenv("sID", sID.Text(16))
	os.Setenv("cID", usr.cID.Text(16))

	return nil
}

// Verify session key SKi
func (clt *Client) Verify(usr User) error {
	SKi := new(big.Int)
	SKi.SetString(os.Getenv("SKi"), 16)

	g, err := rand.Int(rand.Reader, usr.q)
	if err != nil {
		return err
	}

	challenge, err := marshalVerifyRequest(g, usr.q)
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

	verifier := crypto.ExpInGroup(g, SKi, usr.q)
	if response.Cmp(verifier) != 0 {
		return ErrAuthenticationFailed
	}

	return nil
}

// GetMetadata ...
func (clt *Client) GetMetadata(usr User) ([]Domain, error) {

	SKi := new(big.Int)
	SKi.SetString(os.Getenv("SKi"), 16)

	sID := new(big.Int)
	sID.SetString(os.Getenv("sID"), 16)

	mac := crypto.HmacData(clt.config.hash, SKi.Bytes(), usr.cID.Bytes(), sID.Bytes())
	req, err := marshalMetadataRequest(usr.cID.Text(16), hex.EncodeToString(mac))

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
