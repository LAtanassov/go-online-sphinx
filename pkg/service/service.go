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
	Challenge(ski, g, q *big.Int) (r *big.Int, err error)

	VerifyMAC(mac []byte, key []byte, data ...[]byte) error

	GetMetadata() (domains []string, err error)

	Add(domain string) (err error)
	Get(domain string, bmk *big.Int) (bj, qj *big.Int, err error)
}

// UserRepository represents a store for user management - need to be implemented
type UserRepository interface {
	Add(u User) error
	Get(ID string) (User, error)
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
	vaults VaultRepository
	config Configuration
}

// New returns an Online SPHINX service - to share - pointer.
func New(users UserRepository, vaults VaultRepository, cfg Configuration) *OnlineSphinx {
	return &OnlineSphinx{
		users:  users,
		vaults: vaults,
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

// Challenge decrypts the vNonce, increments it and encrypts it again.
func (o *OnlineSphinx) Challenge(ski, g, q *big.Int) (r *big.Int, err error) {
	return crypto.ExpInGroup(g, ski, q), nil
}

// GetMetadata verifies hmac and returns all domains associated with client ID
func (o *OnlineSphinx) GetMetadata() (domains []string, err error) {
	domains, err = o.vaults.GetDomains()
	return
}

// VerifyMAC verifies client request by calculating MAC of the request and
// comparign it with the one send by the client
func (o *OnlineSphinx) VerifyMAC(mac []byte, key []byte, data ...[]byte) error {
	vmac := crypto.HmacData(o.config.hash, key, data...)

	if bytes.Compare(mac, vmac) != 0 {
		return ErrAuthorizationFailed
	}
	return nil
}

// Add by generating random keys k, qj for specific 'domain'
func (o *OnlineSphinx) Add(domain string) (err error) {

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

	err = o.vaults.Add(domain, Vault{
		k:  k,
		qj: qj,
	})

	if err != nil {
		return err
	}

	return nil
}

// Get return bmk**bj and qj associated with domain
func (o *OnlineSphinx) Get(domain string, bmk *big.Int) (bj, qj *big.Int, err error) {

	v, err := o.vaults.Get(domain)
	if err != nil {
		return nil, nil, err
	}

	return v.k, v.qj, nil
}
