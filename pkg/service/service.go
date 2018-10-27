package service

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

var (
	// ErrInvalidArgument is returned when an invalid argument was passed.
	ErrInvalidArgument = errors.New("invalid arguments")
	// ErrAuthorizationFailed is returned when authorization failue happend
	ErrAuthorizationFailed = errors.New("authorization failed")
	// ErrDomainNotFound is returned when an user ask for a domain that does not exists
	ErrDomainNotFound = errors.New("domain not found")
)

var one = big.NewInt(1)
var two = big.NewInt(2)

// Service represents the interface provided to other layers.
type Service interface {
	Register(cID *big.Int) error

	ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error)
	Challenge(ski, g, q *big.Int) (r *big.Int, err error)

	VerifyMAC(mac []byte, cID *big.Int, data ...[]byte) error

	GetMetadata(cID *big.Int) (domains []string, err error)

	Add(cID *big.Int, domain string) (err error)
	Get(cID *big.Int, domain string, bmk *big.Int, q *big.Int) (bj, qj *big.Int, err error)
}

// UserRepository represents a store for user management - need to be implemented
type UserRepository interface {
	Set(u User) error
	Get(ID *big.Int) (User, error)
}

// VaultRepository represents a store for domain management - need to be implemented
type VaultRepository interface {
	Add(d string, v Vault) error
	Get(d string) (Vault, error)
	GetDomains() ([]string, error)
}

// OnlineSphinx provides all operations needed.
type OnlineSphinx struct {
	users  UserRepository
	config Configuration
}

// New returns an Online SPHINX service - to share - pointer.
func New(users UserRepository, cfg Configuration) *OnlineSphinx {
	return &OnlineSphinx{
		users:  users,
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

	return o.users.Set(User{
		cID:    cID,
		kv:     kv,
		vaults: make(map[string]Vault),
	})
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

	u, err := o.users.Get(cID)
	if err != nil {
		return
	}
	kv = u.kv

	ski = new(big.Int)
	ski.SetBytes(crypto.HmacData(o.config.hash, kv.Bytes(), cID.Bytes(), sID.Bytes(), cNonce.Bytes(), sNonce.Bytes()))

	return
}

// Challenge decrypts the vNonce, increments it and encrypts it again.
func (o *OnlineSphinx) Challenge(ski, g, q *big.Int) (r *big.Int, err error) {
	return crypto.ExpInGroup(g, ski, q), nil
}

// GetMetadata verifies hmac and returns all domains associated with client ID
func (o *OnlineSphinx) GetMetadata(cID *big.Int) (domains []string, err error) {
	u, err := o.users.Get(cID)
	if err != nil {
		return nil, err
	}
	domains = make([]string, len(u.vaults))
	var i int
	for k := range u.vaults {
		domains[i] = k
		i++
	}

	return
}

// VerifyMAC verifies client request by calculating MAC of the request and
// comparign it with the one send by the client
func (o *OnlineSphinx) VerifyMAC(mac []byte, cID *big.Int, data ...[]byte) error {
	u, err := o.users.Get(cID)
	if err != nil {
		return err
	}

	vmac := crypto.HmacData(o.config.hash, u.kv.Bytes(), data...)

	if bytes.Compare(mac, vmac) != 0 {
		return ErrAuthorizationFailed
	}
	return nil
}

// Add by generating random keys k, qj for specific 'domain'
func (o *OnlineSphinx) Add(cID *big.Int, domain string) (err error) {

	max := new(big.Int)
	max.Exp(two, o.config.bits, nil)

	k, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	qj, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}

	u, err := o.users.Get(cID)
	if err != nil {
		return err
	}
	u.vaults[domain] = Vault{
		k:  k,
		qj: qj,
	}

	return o.users.Set(u)
}

// Get return bmk**bj and qj associated with domain
func (o *OnlineSphinx) Get(cID *big.Int, domain string, bmk, q *big.Int) (bj, qj *big.Int, err error) {

	u, err := o.users.Get(cID)
	if err != nil {
		return nil, nil, err
	}

	v, ok := u.vaults[domain]
	if !ok {
		return nil, nil, ErrDomainNotFound
	}

	return crypto.ExpInGroup(bmk, v.k, q), v.qj, nil
}
