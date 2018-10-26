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
}

// NewConfiguration return Configuration
func NewConfiguration(sID, k, q0, bits *big.Int, hash func() hash.Hash) Configuration {
	return Configuration{
		sID:  sID,
		k:    k,
		q0:   q0,
		bits: bits,
		hash: hash,
	}
}
