package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

// ErrInvalidArgument is returned when an invalid argument was passed.
var ErrInvalidArgument = errors.New("invalid arguments")

var one = big.NewInt(1)
var two = big.NewInt(2)

// Service represents the interface provided to other layers.
type Service interface {
	Register(cID *big.Int) error
	ExpK(cID, cNonce, b, q *big.Int) (sID, sNonce, bd, q0, kv *big.Int, err error)
	Verify(g, q *big.Int) (r *big.Int, err error)

	GetMetadata(cID *big.Int, mac []byte) (domains []Domain, err error)

	//AddVault(u string) (err error)
	//GetVault(u string, bmk *big.Int) (bj, qj *big.Int, err error)
}

// Repository represents a store for user management - need to be implemented
type Repository interface {
	Add(u User) error
	Get(ID string) (User, error)
}

// User is an entity and contains all user related informated
type User struct {
	cID   *big.Int
	kv    *big.Int
	store map[string]Vault
}

// Vault ...
type Vault struct {
	k *big.Int
	q *big.Int
}

// OnlineSphinx provides all operations needed.
type OnlineSphinx struct {
	repo   Repository
	config Configuration
}

// New returns an Online SPHINX service - to share - pointer.
func New(repo Repository, cfg Configuration) *OnlineSphinx {
	return &OnlineSphinx{
		repo:   repo,
		config: cfg,
	}
}

// Register an user with its id
func (o *OnlineSphinx) Register(cID *big.Int) error {
	max := new(big.Int)
	max.Exp(two, o.config.bits, nil)

	kv, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	return o.repo.Add(User{cID: cID, kv: kv})
}

// ExpK returns r**k mod |2q + 1|
func (o *OnlineSphinx) ExpK(cID, cNonce, b, q *big.Int) (sID, sNonce, bd, q0, kv *big.Int, err error) {
	sID = o.config.sID
	q0 = o.config.q0

	bd = crypto.ExpInGroup(b, o.config.k, q)

	max := new(big.Int)
	max.Exp(two, o.config.bits, nil)

	sNonce, err = rand.Int(rand.Reader, max)
	if err != nil {
		return
	}

	u, err := o.repo.Get(cID.Text(16))
	if err != nil {
		return
	}
	kv = u.kv

	SKi := new(big.Int)
	SKi.SetBytes(crypto.HmacData(o.config.hash, kv.Bytes(), cID.Bytes(), sID.Bytes(), cNonce.Bytes(), sNonce.Bytes()))

	os.Setenv("SKi", SKi.Text(16))

	return
}

// Verify decrypts the vNonce, increments it and encrypts it again.
func (o *OnlineSphinx) Verify(g, q *big.Int) (r *big.Int, err error) {

	SKi := new(big.Int)
	SKi.SetString(os.Getenv("SKi"), 16)

	return crypto.ExpInGroup(g, SKi, q), nil
}

func (o *OnlineSphinx) GetMetadata(cID *big.Int, mac []byte) (domains []Domain, err error) {
	return nil, nil
}
