package client

import (
	"crypto/sha256"
	"hash"
	"math/big"
)

// Configuration ...
type Configuration struct {
	hash         func() hash.Hash
	bits         *big.Int
	contentType  string
	baseURL      string
	registerPath string
	expkPath     string
	verifyPath   string
	metadataPath string
	addPath      string
	getPath      string
}

// NewConfiguration return default configuration.
func NewConfiguration() Configuration {
	return Configuration{
		hash: sha256.New,
		bits: big.NewInt(8),
	}
}
