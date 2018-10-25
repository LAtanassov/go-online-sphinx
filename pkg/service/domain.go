package service

import (
	"hash"
	"math/big"
)

// Configuration ...
type Configuration struct {
	sID *big.Int

	k  *big.Int
	q0 *big.Int

	hash func() hash.Hash
	bits *big.Int
}

// NewConfiguration ...
func NewConfiguration(sID, k, q0, bits *big.Int, hash func() hash.Hash) Configuration {
	return Configuration{
		sID:  sID,
		k:    k,
		q0:   q0,
		bits: bits,
		hash: hash,
	}
}
