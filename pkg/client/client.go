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
)

// ErrUserNotCreated in repostory
var ErrUserNotCreated = errors.New("user not created")

// ErrAuthenticationFailed covers several authentication issues
var ErrAuthenticationFailed = errors.New("authentication failed")

// often used big.Int
var one = big.NewInt(1)
var two = big.NewInt(2)

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
	k    *big.Int

	bits         *big.Int
	contentType  string
	baseURL      string
	registerPath string
	expkPath     string
	verifyPath   string
}

// New returns a Online SPHINX client
func New(p Poster, c Configuration) *Client {
	return &Client{http: p, config: c}
}

// Login an user
// $> curl -d '{"username":"hans"}' -H "Content-Type: application/json" -X POST http://localhost:8080/v1/register
// $> curl -d '{"username":"hans", "cNonce": "43", "b": "17b", "q": "d3"}' -H "Content-Type: application/json" -X POST http://localhost:8080/v1/login/expk
func (c *Client) Login(username, password string) error {

	max := new(big.Int)
	max.Exp(two, c.config.bits, nil)

	cNonce, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	k, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	p := new(big.Int)
	p.SetBytes(c.config.hash().Sum([]byte(password)))

	g := crypto.ExpInGroup(p, two, c.config.q)

	kinv := new(big.Int)
	kinv.ModInverse(k, c.config.q)

	if kinv == nil {
		kinv = big.NewInt(0)
	}

	// blinding
	b := crypto.ExpInGroup(g, k, c.config.q)

	jsonReq, err := json.Marshal(&expKRequest{
		Username: username,
		CNonce:   cNonce.Text(16),
		B:        b.Text(16),
		Q:        c.config.q.Text(16),
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
	bd.SetString(jsonResp.BD, 16)

	// unblinding
	B0 := crypto.ExpInGroup(bd, kinv, c.config.q)

	Q0 := new(big.Int)
	Q0.SetString(jsonResp.Q0, 16)

	mk := new(big.Int)
	mk.Mul(crypto.ExpInGroup(B0, c.config.k, c.config.q), Q0)

	kv := new(big.Int)
	kv.SetString(jsonResp.KV, 16)

	cID := new(big.Int)
	cID.SetString(c.config.cID, 16)

	sID := new(big.Int)
	sID.SetString(jsonResp.SID, 16)

	sNonce := new(big.Int)
	sNonce.SetString(jsonResp.SNonce, 16)

	mac := hmac.New(c.config.hash, kv.Bytes())
	for _, d := range []*big.Int{cID, sID, cNonce, sNonce} {
		mac.Write(d.Bytes())
	}

	SKi := new(big.Int)
	SKi.SetBytes(mac.Sum(nil))

	return c.verify(SKi)
}

// Register a new user
// $> curl -d '{"username":"username"}' -H "Content-Type: application/json" -X POST http://localhost:8080/v1/register
func (c *Client) Register(username, password string) error {

	b, err := json.Marshal(&registerRequest{Username: username})
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

// Logout an logged in user
func (c *Client) Logout() error {
	return nil
}

// Verify session key SKi
func (c *Client) verify(SKi *big.Int) error {

	max := new(big.Int)
	max.Exp(two, c.config.bits, nil)

	vNonce, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	jsonReq, err := json.Marshal(&verifyRequest{
		VNonce: vNonce.Text(16),
	})

	if err != nil {
		return err
	}

	u, err := url.Parse(c.config.baseURL)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, c.config.verifyPath)

	resp, err := c.http.Post(u.String(), c.config.contentType, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonResp := verifyResponse{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return err
	}

	w := new(big.Int)
	w.SetString(jsonResp.WNonce, 16)

	wc := new(big.Int)
	if wc.Add(one, vNonce).Cmp(w) != 0 {
		return ErrAuthenticationFailed
	}

	return nil
}

// Metadata contains information
type Metadata struct {
}

// GetMetadata request metadata
func (c *Client) GetMetadata(MACski *big.Int) (*Metadata, error) {
	return &Metadata{}, nil
}
