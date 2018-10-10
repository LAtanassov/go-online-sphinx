package client

import (
	"crypto/hmac"
	"crypto/rand"
	"errors"
	"hash"
	"io"
	"math/big"
	"net/http"
	"net/url"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"

	"golang.org/x/crypto/openpgp/elgamal"
)

// ErrUserNotCreated in repostory
var ErrUserNotCreated = errors.New("user not created")

var one = big.NewInt(1)
var two = big.NewInt(2)

// LoginConfig contains login configuration
type LoginConfig struct {
	cID    *big.Int
	pwd    string
	server []string
	h      func() hash.Hash
	q      *big.Int
	k      *elgamal.PrivateKey
}

// Login implements Online SPHINX login protocol
func Login(config LoginConfig) ([]byte, *Metadata, error) {

	cNonce, err := rand.Int(rand.Reader, config.q)

	b, kinv := blind(config.pwd, config.q, config.h)
	if err != nil {
		return nil, nil, err
	}

	c := &Client{}
	sID, sNonce, bd, Q0, kv, err := c.ExpK(config.cID, cNonce, b, config.q)

	B0 := unblind(bd, kinv, config.q)

	mk, err := elgamal.Decrypt(config.k, B0, Q0)
	if err != nil {
		return nil, nil, err
	}

	SKi := hmacBigInt(config.h, kv, []*big.Int{config.cID, sID, cNonce, sNonce})

	err = c.Verify(SKi)
	if err != nil {
		return nil, nil, err
	}

	MACski := hmacBigInt(config.h, SKi, []*big.Int{config.cID, sID, big.NewInt(1)})
	meta, err := c.GetMetadata(MACski)
	if err != nil {
		return nil, nil, err
	}

	return mk, meta, nil
}

func hmacBigInt(h func() hash.Hash, key *big.Int, data []*big.Int) (m *big.Int) {
	mac := hmac.New(h, key.Bytes())
	for _, d := range data {
		mac.Write(d.Bytes())
	}
	m = big.NewInt(0)
	m.SetBytes(mac.Sum(nil))
	return
}

// runs on client
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

type metadatarequest struct {
}

// Client represent Online Sphinx Client
type Client struct {
	p Poster
}

// New returns a Online SPHINX client
func New(p Poster, c Configuration) *Client {
	return &Client{p: p}
}

// Poster represents an interface to do POST requests
type Poster interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// Configuration of an Online SPHINX client.
type Configuration struct {
	cID    string
	server *url.URL
	hash   func() hash.Hash
	q      *big.Int
	k      *big.Int
}

// Login an user
func (c *Client) Login(username, password string) error {
	return nil
}

// Register a new user
func (c *Client) Register(username string) error {
	return nil
}

// ExpK runs on server
func (c *Client) ExpK(cID, cNonce, b, q *big.Int) (sID, sNonce, bd, Q0, kv *big.Int, err error) {
	sID = big.NewInt(0)
	d, err := rand.Int(rand.Reader, q)
	if err != nil {
		return
	}

	bd = crypto.ExpInGroup(b, d, q)
	sNonce = big.NewInt(0)
	Q0 = big.NewInt(0)
	kv = big.NewInt(0)
	return
}

// Verify session key SKi
func (c *Client) Verify(SKi *big.Int) error {
	return nil
}

// Metadata contains information
type Metadata struct {
}

// GetMetadata request metadata
func (c *Client) GetMetadata(MACski *big.Int) (*Metadata, error) {
	return &Metadata{}, nil
}
