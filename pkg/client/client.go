package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"path"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"

	"golang.org/x/crypto/openpgp/elgamal"
)

// ErrUserNotCreated in repostory
var ErrUserNotCreated = errors.New("user not created")

// ErrAuthenticationFailed covers several authentication issues
var ErrAuthenticationFailed = errors.New("authentication failed")

// Client represent Online Sphinx Client
type Client struct {
	http   Poster
	config Configuration
}

// Poster represents an interface to do POST requests
type Poster interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// Configuration of an Online SPHINX client.
type Configuration struct {
	cID  string
	hash func() hash.Hash
	q    *big.Int
	k    *elgamal.PrivateKey

	contentType  string
	baseURL      string
	registerPath string
	expkPath     string
}

// New returns a Online SPHINX client
func New(p Poster, c Configuration) *Client {
	return &Client{http: p, config: c}
}

// Login an user
func (c *Client) Login(username, password string) error {

	cNonce, err := rand.Int(rand.Reader, c.config.q)

	b, kinv := blind(password, c.config.q, c.config.hash)
	if err != nil {
		return err
	}

	jsonReq, err := json.Marshal(&expKRequest{
		username: username,
		cNonce:   cNonce.Text(16),
		b:        b.Text(16),
		q:        c.config.q.Text(16),
	})
	if err != nil {
		return err
	}

	u, err := url.Parse(c.config.baseURL)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, c.config.expkPath)

	resp, err := c.http.Post(u.String(), c.config.contentType, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonResp := expKResponse{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return err
	}

	bd := new(big.Int)
	bd.SetString(jsonResp.bd, 16)

	B0 := unblind(bd, kinv, c.config.q)

	Q0 := new(big.Int)
	Q0.SetString(jsonResp.bd, 16)

	_, err = elgamal.Decrypt(c.config.k, B0, Q0)
	if err != nil {
		return err
	}

	kv := new(big.Int)
	kv.SetString(jsonResp.kv, 16)

	cID := new(big.Int)
	cID.SetString(c.config.cID, 16)

	sID := new(big.Int)
	sID.SetString(jsonResp.sID, 16)

	sNonce := new(big.Int)
	sNonce.SetString(jsonResp.sNonce, 16)

	SKi := hmacBigInt(c.config.hash, kv, []*big.Int{cID, sID, cNonce, sNonce})

	err = c.verify(SKi)
	if err != nil {
		return err
	}

	MACski := hmacBigInt(c.config.hash, SKi, []*big.Int{cID, sID, big.NewInt(1)})
	_, err = c.GetMetadata(MACski)
	if err != nil {
		return err
	}

	return nil
}

// Register a new user
func (c *Client) Register(username, password string) error {

	b, err := json.Marshal(&registerRequest{username: username})
	if err != nil {
		return err
	}

	u, err := url.Parse(c.config.baseURL)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, c.config.registerPath)

	_, err = c.http.Post(u.String(), c.config.contentType, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	return nil
}

type registerRequest struct {
	username string
}

// Logout an logged in user
func (c *Client) Logout() error {
	return nil
}

// Verify session key SKi
func (c *Client) verify(SKi *big.Int) error {
	return nil
}

// Metadata contains information
type Metadata struct {
}

// GetMetadata request metadata
func (c *Client) GetMetadata(MACski *big.Int) (*Metadata, error) {
	return &Metadata{}, nil
}

type expKRequest struct {
	username string
	cNonce   string
	b        string
	q        string
}

type expKResponse struct {
	sID    string
	sNonce string
	bd     string
	q0     string
	kv     string
}

type metadataRequest struct {
}

var one = big.NewInt(1)
var two = big.NewInt(2)

func hmacBigInt(h func() hash.Hash, key *big.Int, data []*big.Int) (m *big.Int) {
	mac := hmac.New(h, key.Bytes())
	for _, d := range data {
		mac.Write(d.Bytes())
	}
	m = big.NewInt(0)
	m.SetBytes(mac.Sum(nil))
	return
}

func blind(pwd string, q *big.Int, h func() hash.Hash) (b, kinv *big.Int) {
	p := big.NewInt(0).SetBytes(h().Sum([]byte(pwd)))
	g := crypto.ExpInGroup(p, two, q)

	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		return
	}

	kinv = big.NewInt(0).ModInverse(k, q)
	if kinv == nil {
		kinv = big.NewInt(0)
	}

	// blinding
	b = crypto.ExpInGroup(g, k, q)
	return
}

func unblind(bd, kinv, q *big.Int) (B0 *big.Int) {
	B0 = crypto.ExpInGroup(bd, kinv, q)
	return
}
