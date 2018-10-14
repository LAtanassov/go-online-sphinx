package service

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

// ErrInvalidArgument is returned when an invalid argument was passed.
var ErrInvalidArgument = errors.New("invalid arguments")

var one = big.NewInt(1)
var two = big.NewInt(2)

// Service represents the interface provided to other layers.
type Service interface {
	Register(id string) error
	ExpK(uID string, b, q *big.Int) (sID string, sNonce, bd, q0, kv *big.Int, err error)
	Verify(v *big.Int) (w *big.Int, err error)
}

// Repository represents a store for user management - need to be implemented
type Repository interface {
	Add(u User) error
	Get(ID string) (User, error)
}

// User is an entity and contains all user related informated
type User struct {
	username string
	kv       *big.Int
}

// OnlineSphinx provides all operations needed.
type OnlineSphinx struct {
	sID string

	k  *big.Int
	q0 *big.Int

	bits *big.Int
	repo Repository
}

// New returns an Online SPHINX service - to share - pointer.
func New(sID string, k, q0, bits *big.Int, repo Repository) *OnlineSphinx {
	return &OnlineSphinx{
		sID: sID,

		k:  k,
		q0: q0,

		bits: bits,
		repo: repo,
	}
}

// Register an user with its id
func (o *OnlineSphinx) Register(username string) error {
	max := new(big.Int)
	max.Exp(two, o.bits, nil)

	kv, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	return o.repo.Add(User{username: username, kv: kv})
}

// ExpK returns r**k mod |2q + 1|
func (o *OnlineSphinx) ExpK(uID string, b, q *big.Int) (sID string, sNonce, bd, q0, kv *big.Int, err error) {
	sID = o.sID
	q0 = o.q0
	bd = crypto.ExpInGroup(b, o.k, q)

	max := new(big.Int)
	max.Exp(two, o.bits, nil)

	sNonce, err = rand.Int(rand.Reader, max)
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

// Verify decrypts the vNonce, increments it and encrypts it again.
func (o *OnlineSphinx) Verify(v *big.Int) (w *big.Int, err error) {
	w = new(big.Int)
	return w.Add(v, one), nil
}
