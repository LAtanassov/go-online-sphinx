package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"hash"
	"math/big"
)

var one = big.NewInt(1)
var two = big.NewInt(2)

// ExpInGroup - exponential in cyclic group returns g**k mod |2*q + 1|
func ExpInGroup(g, k, q *big.Int) *big.Int {
	var p = big.NewInt(0)
	var r = big.NewInt(0)

	p.Add(p.Mul(two, q), one)

	return r.Exp(g, k, p)
}

// HashInGroup takes an arbitrary string transform into a group element
func HashInGroup(password string, newHash func() hash.Hash, q *big.Int) *big.Int {
	p := new(big.Int)
	p.SetBytes(newHash().Sum([]byte(password)))

	return ExpInGroup(p, two, q)
}

// Blind ...
func Blind(g, q, bits *big.Int) (cNonce, b, kinv *big.Int, err error) {

	max := new(big.Int)
	max.Exp(two, bits, nil)

	cNonce, err = rand.Int(rand.Reader, max)
	if err != nil {
		return
	}

	k, err := rand.Int(rand.Reader, max)
	if err != nil {
		return
	}

	kinv = new(big.Int)
	kinv.ModInverse(k, q)

	if kinv == nil {
		kinv = big.NewInt(0)
	}

	b = ExpInGroup(g, k, q)
	return
}

// Unblind wraps ExpInGroup for readibility
func Unblind(bd, kinv, q *big.Int) *big.Int {
	return ExpInGroup(bd, kinv, q)
}

// HmacData ...
func HmacData(h func() hash.Hash, key []byte, data ...[]byte) []byte {
	mac := hmac.New(h, key)
	for _, d := range data {
		mac.Write(d)
	}

	return mac.Sum(nil)
}
