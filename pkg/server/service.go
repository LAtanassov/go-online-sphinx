package server

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

// ErrInvalidArgument is returned when an invalid argument was passed.
var ErrInvalidArgument = errors.New("invalid arguments")

// Service represents the interface provided to other layers.
type Service interface {
	Register(id string) error
	ExpK(uID string, r, q *big.Int) (sID string, sNonce, bd, q0, kv *big.Int, err error)
}

// Repository represents a store for user management - need to be implemented
type Repository interface {
	Add(u User) error
	Get(ID string) (User, error)
}

// User is an entity and contains all user related informated
type User struct {
	id string
	kv *big.Int
}

// OnlineSphinx provides all operations needed.
type OnlineSphinx struct {
	sID string

	k  *big.Int
	q0 *big.Int

	repo Repository
}

// NewService returns an Online SPHINX service - to share - pointer.
func NewService(sID string, k, q0 *big.Int, repo Repository) *OnlineSphinx {
	return &OnlineSphinx{
		sID: sID,

		k:  k,
		q0: q0,

		repo: repo,
	}
}

// Register an user with its id
func (o *OnlineSphinx) Register(id string) error {
	return o.repo.Add(User{id: id, kv: big.NewInt(0)})
}

// ExpK returns r**k mod |2q + 1|
func (o *OnlineSphinx) ExpK(uID string, r, q *big.Int) (sID string, sNonce, bd, q0, kv *big.Int, err error) {
	sID = o.sID
	q0 = o.q0
	bd = crypto.ExpInGroup(r, o.k, q)

	sNonce, err = rand.Int(rand.Reader, q)
	if err != nil {
		return
	}

	u, err := o.repo.Get(uID)
	if err != nil {
		return
	}
	kv = u.kv

	return
}
