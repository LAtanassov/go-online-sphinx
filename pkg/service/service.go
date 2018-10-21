package service

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

// ErrInvalidArgument is returned when an invalid argument was passed.
var ErrInvalidArgument = errors.New("invalid arguments")

// ErrAuthorizationFailed is returned when authorization failue happend
var ErrAuthorizationFailed = errors.New("authorization failed")

var one = big.NewInt(1)
var two = big.NewInt(2)

// Service represents the interface provided to other layers.
type Service interface {
	Register(cID *big.Int) error
	ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error)
	Verify(ski, g, q *big.Int) (r *big.Int, err error)

	GetMetadata(cID *big.Int, mac []byte) (domains []Domain, err error)

	//AddVault(u string) (err error)
	//GetVault(u string, bmk *big.Int) (bj, qj *big.Int, err error)
}

// UserRepository represents a store for user management - need to be implemented
type UserRepository interface {
	Add(u User) error
	Get(ID string) (User, error)
}

// DomainRepository represents a store for domain management - need to be implemented
type DomainRepository interface {
	Add(id string, d Domain) error
	Get(id string) ([]Domain, error)
}

// Vault ...
type Vault struct {
	k *big.Int
	q *big.Int
}

// OnlineSphinx provides all operations needed.
type OnlineSphinx struct {
	users   UserRepository
	domains DomainRepository
	config  Configuration
}

// New returns an Online SPHINX service - to share - pointer.
func New(users UserRepository, domains DomainRepository, cfg Configuration) *OnlineSphinx {
	return &OnlineSphinx{
		users:   users,
		domains: domains,
		config:  cfg,
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

	return o.users.Add(User{cID: cID, kv: kv})
}

// ExpK returns r**k mod |2q + 1|
func (o *OnlineSphinx) ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error) {
	sID = o.config.sID
	q0 = o.config.q0

	bd = crypto.ExpInGroup(b, o.config.k, q)

	max := new(big.Int)
	max.Exp(two, o.config.bits, nil)

	sNonce, err = rand.Int(rand.Reader, max)
	if err != nil {
		return
	}

	u, err := o.users.Get(cID.Text(16))
	if err != nil {
		return
	}
	kv = u.kv

	ski = new(big.Int)
	ski.SetBytes(crypto.HmacData(o.config.hash, kv.Bytes(), cID.Bytes(), sID.Bytes(), cNonce.Bytes(), sNonce.Bytes()))

	return
}

// Verify decrypts the vNonce, increments it and encrypts it again.
func (o *OnlineSphinx) Verify(ski, g, q *big.Int) (r *big.Int, err error) {
	return crypto.ExpInGroup(g, ski, q), nil
}

// GetMetadata verifies hmac and returns all domains associated with client ID
func (o *OnlineSphinx) GetMetadata(cID *big.Int, mac []byte) (domains []Domain, err error) {

	u, err := o.users.Get(cID.Text(16))
	if err != nil {
		return
	}

	vmac := crypto.HmacData(o.config.hash, u.kv.Bytes(), cID.Bytes(), o.config.sID.Bytes())

	if bytes.Compare(mac, vmac) != 0 {
		err = ErrAuthorizationFailed
		return
	}

	domains, err = o.domains.Get(cID.Text(16))

	return
}
