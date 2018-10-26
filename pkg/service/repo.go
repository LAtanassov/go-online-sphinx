package service

import (
	"errors"
	"math/big"
	"sync"
)

var (
	// ErrUserAlreadyExists in repostory already
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrVaultAlreadyExists in repostory already
	ErrVaultAlreadyExists = errors.New("vault already exists")
	// ErrUserNotFound in repostory
	ErrUserNotFound = errors.New("user not found")
)

// NewUserRepository creates and returns an inmemory user repository.
func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		mutex: sync.Mutex{},
		users: make(map[string]User),
	}
}

// InMemoryUserRepository provides an user repository.
// This repository can also be implemented using an SQL database.UserRepository.
// It should be store for long term (replicated, shared).
// client.UserRepository is atm identical with server.UserRepository, but this might change in future
type InMemoryUserRepository struct {
	mutex sync.Mutex
	users map[string]User
}

// User is an entity and contains all user related informated to implement server-side Online SPHINX.
type User struct {
	cID *big.Int
	kv  *big.Int
}

// Add new user to user repository if does not exists
func (r *InMemoryUserRepository) Add(u User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.users[u.cID.Text(16)]
	if ok {
		return ErrUserAlreadyExists
	}

	r.users[u.cID.Text(16)] = u
	return nil
}

// Get an existing user
func (r *InMemoryUserRepository) Get(cID string) (User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	u, ok := r.users[cID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

// NewVaultRepository creates and returns an inmemory vault repository.
func NewVaultRepository() *InMemoryVaultRepository {
	return &InMemoryVaultRepository{
		mutex:  sync.Mutex{},
		vaults: make(map[string]Vault),
	}
}

// InMemoryVaultRepository provides a vault repository.
type InMemoryVaultRepository struct {
	mutex  sync.Mutex
	vaults map[string]Vault
}

// Vault ...
type Vault struct {
	k  *big.Int
	qj *big.Int
}

// Add new vault to vault repository if does not exists
func (r *InMemoryVaultRepository) Add(domain string, vault Vault) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.vaults[domain]
	if ok {
		return ErrVaultAlreadyExists
	}

	r.vaults[domain] = vault
	return nil
}

// Get an existing user
func (r *InMemoryVaultRepository) Get(domain string) (Vault, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	vault, ok := r.vaults[domain]
	if !ok {
		vault = Vault{}
	}

	return vault, nil
}

// GetDomains an existing user
func (r *InMemoryVaultRepository) GetDomains() ([]string, error) {
	return []string{}, nil
}
