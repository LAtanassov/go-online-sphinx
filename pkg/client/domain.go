package client

import (
	"hash"
	"math/big"
)

// User ...
type User struct {
	username string
	cID      *big.Int
	q        *big.Int
	k        *big.Int
}

// NewUser ...
func NewUser(username string) User {
	return User{
		username: username,
	}
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
	metadataPath string
}

// Domain ...
type Domain struct {
}

type Session struct {
	ski *big.Int
	mk  *big.Int
	sID *big.Int
}

// NewDomain ...
func NewDomain() Domain {
	return Domain{}
}
