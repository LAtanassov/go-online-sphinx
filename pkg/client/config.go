package client

import (
	"crypto/sha256"
	"hash"
)

// Configuration ...
type Configuration struct {
	hash          func() hash.Hash
	bits          int
	contentType   string
	baseURL       string
	registerPath  string
	expkPath      string
	challengePath string
	metadataPath  string
	addPath       string
	getPath       string
}

// NewConfiguration return default configuration.
func NewConfiguration() Configuration {
	return Configuration{
		hash: sha256.New,
		bits: 8,
	}
}
