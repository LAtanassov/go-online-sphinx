package client

import (
	"hash"
	"math/big"
)

func NewUser(username string) User {
	return User{
		username: username,
	}
}

// User ...
type User struct {
	username string
	q        *big.Int
	k        *big.Int
}

// Configuration ...
type Configuration struct {
	hash         func() hash.Hash
	bits         *big.Int
	contentType  string
	baseURL      string
	registerPath string
	expkPath     string
	verifyPath   string
}

type metadata struct {
}
