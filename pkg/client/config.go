package client

import (
	"hash"
)

// Configuration ...
type Configuration struct {
	hash          func() hash.Hash
	Bits          int
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
func NewConfiguration(baseURL string, bits int, hashFn func() hash.Hash) Configuration {
	return Configuration{
		hash:          hashFn,
		Bits:          8,
		contentType:   "application/json",
		baseURL:       baseURL,
		registerPath:  "/v1/register",
		expkPath:      "/v1/login/expk",
		challengePath: "/v1/login/challenge",
		metadataPath:  "/v1/metadata",
		addPath:       "/v1/add",
		getPath:       "/v1/get",
	}
}
