package client

import (
	"hash"
	"net/url"
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
	logoutPath    string
}

// NewConfiguration return default configuration.
func NewConfiguration(baseURL string, bits int, hashFn func() hash.Hash) (Configuration, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return Configuration{}, err
	}

	c := Configuration{
		hash:        hashFn,
		bits:        bits,
		contentType: "application/json",
		baseURL:     baseURL,
	}
	u.Path = "/v1/register"
	c.registerPath = u.String()

	u.Path = "/v1/login/expk"
	c.expkPath = u.String()

	u.Path = "/v1/login/challenge"
	c.challengePath = u.String()

	u.Path = "/v1/metadata"
	c.metadataPath = u.String()

	u.Path = "/v1/add"
	c.addPath = u.String()

	u.Path = "/v1/get"
	c.getPath = u.String()

	u.Path = "/v1/logout"
	c.logoutPath = u.String()

	return c, nil
}
