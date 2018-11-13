package service

import (
	"hash"
	"math/big"
)

// Configuration contains cryptographical key material needed for an Online SPHINX service.
type Configuration struct {
	sID  *big.Int         // service ID
	k    *big.Int         // service key k
	q0   *big.Int         // common ElGammal component Q_0
	hash func() hash.Hash // hash function
	bits *big.Int         // bits used in crypthographic directives
	max  *big.Int
}

// NewConfiguration initialize and returns a Configuration
func NewConfiguration(sID, k, q0, bits *big.Int, hash func() hash.Hash) Configuration {
	max := new(big.Int)
	max.Exp(big.NewInt(2), bits, nil)
	return Configuration{
		sID:  sID,
		k:    k,
		q0:   q0,
		bits: bits,
		hash: hash,
		max:  max,
	}
}
