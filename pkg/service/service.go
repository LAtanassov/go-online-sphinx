package service

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

var (
	// ErrMacMismatch is returned when the MAC of the request and the MAC within the request does not match
	ErrMacMismatch = errors.New("MAC mismatch")
	// ErrDomainNotFound is returned an existing user does not
	ErrDomainNotFound = errors.New("domain not found")
)
var (
	one = big.NewInt(1)
	two = big.NewInt(2)
)

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

// Middleware is a chainable behavior modifier for Service.
type Middleware func(Service) Service

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

// Register an user with its cID.
// Returns error if user with same cID already exists,
// or if could not set user to repository.
func (o *OnlineSphinx) Register(cID *big.Int) error {

	kv, err := rand.Int(rand.Reader, o.config.max)
	if err != nil {
		return errors.Wrap(err, "Register: failed to generate random int kv")
	}

	_, err = o.users.Get(cID)
	switch err {
	case ErrUserNotFound:
		// continue
	default:
		return errors.Wrapf(err, "Register: failed to register user with cID=%v", cID)
	}

	return errors.Wrapf(
		o.users.Set(User{
			cID:    cID,
			kv:     kv,
			vaults: make(map[string]Vault),
		}), "Register: failed to users.set() with ID %v", cID)
}

// ExpK returns r**k mod |2q + 1|
func (o *OnlineSphinx) ExpK(cID, cNonce, b, q *big.Int) (ski, sID, sNonce, bd, q0, kv *big.Int, err error) {
	sID = o.config.sID
	q0 = o.config.q0

	bd = crypto.ExpInGroup(b, o.config.k, q)

	sNonce, err = rand.Int(rand.Reader, o.config.max)
	if err != nil {
		err = errors.Wrap(err, "ExpK: failed to generate random int sNonce")
		return
	}

	u, err := o.users.Get(cID)
	if err != nil {
		err = errors.Wrapf(err, "ExpK: failed to users.get() user with cID=%v", cID)
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
		return nil, errors.Wrapf(err, "GetMetadata: failed to users.get() user with cID=%v", cID)
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
func (o *OnlineSphinx) VerifyMAC(mac []byte, ski *big.Int, data ...[]byte) error {

	vmac := crypto.HmacData(o.config.hash, ski.Bytes(), data...)

	if bytes.Compare(mac, vmac) != 0 {
		return errors.Wrapf(ErrMacMismatch, "VerifyMAC: given mac=%v and calculated mac=%v are different", mac, vmac)
	}
	return nil
}

// Add by generating random keys k, qj for specific 'domain'
func (o *OnlineSphinx) Add(cID *big.Int, domain string) error {

	k, err := rand.Int(rand.Reader, o.config.max)
	if err != nil {
		return errors.Wrap(err, "Add: failed to generate random int k")
	}

	qj, err := rand.Int(rand.Reader, o.config.max)
	if err != nil {
		return errors.Wrap(err, "Add: failed to generate random int qj")
	}

	u, err := o.users.Get(cID)
	if err != nil {
		return errors.Wrapf(err, "Add: failed to users.get() user with cID=%v", cID)
	}
	u.vaults[domain] = Vault{
		k:  k,
		qj: qj,
	}

	return errors.Wrapf(o.users.Set(u), "Add: failed to users.add() user with cID=%v and domain=%v", cID, domain)
}

// Get return bmk**bj and qj associated with domain
func (o *OnlineSphinx) Get(cID *big.Int, domain string, bmk, q *big.Int) (bj, qj *big.Int, err error) {

	u, err := o.users.Get(cID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Get: failed to users.get() user with cID=%v", cID)
	}

	v, ok := u.vaults[domain]
	if !ok {
		return nil, nil, errors.Wrapf(ErrDomainNotFound, "Get: failed to get user with cID=%v and domain=%v", cID, domain)
	}

	return crypto.ExpInGroup(bmk, v.k, q), v.qj, nil
}
