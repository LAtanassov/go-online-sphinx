package client

import (
	"math/big"
)

// Domain ...
type Domain struct {
}

// Session contains cryptographical key material used to associate
// several HTTP request with an authenticated user.
type Session struct {
	ski  *big.Int
	mk   *big.Int
	sID  *big.Int
	user User
}

// NewSession ...
func NewSession(user User, sID, ski, mk *big.Int) *Session {
	return &Session{
		ski:  ski,
		mk:   mk,
		sID:  sID,
		user: user,
	}
}
