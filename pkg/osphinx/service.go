package osphinx

import (
	"errors"
	"math/big"
)

// ErrInvalidArgument is returned when an invalid argument was passed.
var ErrInvalidArgument = errors.New("invalid arguments")

// Service represents the interface provided to other layers.
type Service interface {
	ExpK(r, q *big.Int) (*big.Int, error)
}

// OnlineSphinx provides all operations needed.
type OnlineSphinx struct {
	k  *big.Int
	Q0 *big.Int
}

// NewService returns an Online SPHINX service - to share - pointer.
func NewService(k, Q0 *big.Int) *OnlineSphinx {
	return &OnlineSphinx{
		k:  k,
		Q0: Q0,
	}
}

// ExpK returns r**k mod |2q + 1|
func (o *OnlineSphinx) ExpK(r, q *big.Int) (*big.Int, error) {
	// TODO: check preconditions, r should be not 1, q should be prime
	return ExpInGroup(r, o.k, q), nil
}
