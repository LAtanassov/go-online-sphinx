package osphinx

import (
	"crypto/hmac"
	"crypto/rand"
	"hash"
	"math/big"

	"golang.org/x/crypto/openpgp/elgamal"
)

// LoginConfig contains login configuration
type LoginConfig struct {
	C      *big.Int
	pwd    string
	server []string
	h      func() hash.Hash
	q      *big.Int
	k      *elgamal.PrivateKey
}

// Login implements Online SPHINX login protocol
func Login(config LoginConfig) ([]byte, *Metadata, error) {

	CNonce, err := rand.Int(rand.Reader, config.q)

	b, kinv := blind(config.pwd, config.q, config.h)
	if err != nil {
		return nil, nil, err
	}

	c := &client{}
	S, bd, SNonce, Q0, kv, err := c.ExpK(b, config.q, CNonce)

	B0 := unblind(bd, kinv, config.q)

	mk, err := elgamal.Decrypt(config.k, B0, Q0)
	if err != nil {
		return nil, nil, err
	}

	SKi := hmacBigInt(config.h, kv, []*big.Int{config.C, S, CNonce, SNonce})

	err = c.Verify(SKi)
	if err != nil {
		return nil, nil, err
	}

	MACski := hmacBigInt(config.h, SKi, []*big.Int{config.C, S, big.NewInt(1)})
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
	g := ExpInGroup(p, two, q)

	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		return
	}

	kinv = big.NewInt(0).ModInverse(k, q)
	if kinv == nil {
		kinv = big.NewInt(0)
	}

	// blinding
	b = ExpInGroup(g, k, q)
	return
}

func unblind(bd, kinv, q *big.Int) (B0 *big.Int) {
	B0 = ExpInGroup(bd, kinv, q)
	return
}

type metadatarequest struct {
}

type client struct {
}

// ExpK runs on server
func (c *client) ExpK(b, q, CNonce *big.Int) (S, bd, SNonce, Q0, kv *big.Int, err error) {
	S = big.NewInt(0)
	d, err := rand.Int(rand.Reader, q)
	if err != nil {
		return
	}

	bd = ExpInGroup(b, d, q)
	SNonce = big.NewInt(0)
	Q0 = big.NewInt(0)
	kv = big.NewInt(0)
	return
}

func (c *client) Verify(SKi *big.Int) error {
	return nil
}

// Metadata contains information
type Metadata struct {
}

// GetMetadata request metadata
func (c *client) GetMetadata(MACski *big.Int) (*Metadata, error) {
	return &Metadata{}, nil
}
