package service

import (
	"errors"
	"math/big"
	"sync"
)

// ErrUserAlreadyExists in repostory already
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrVaultAlreadyExists in repostory already
var ErrVaultAlreadyExists = errors.New("vault already exists")

// ErrUserNotFound in repostory
var ErrUserNotFound = errors.New("user not found")

// User is an entity and contains all user related informated to implement server-side Online SPHINX.
type User struct {
	cID   *big.Int
	kv    *big.Int
	store map[string]Vault
}

// InMemoryUserRepository provides an user repository.
// This repository can also be implemented using an SQL database.UserRepository.
// It should be store for long term (replicated, shared).
// client.UserRepository is atm identical with server.UserRepository, but this might change in future
type InMemoryUserRepository struct {
	mutex sync.Mutex
	users map[string]User
}

// NewUserRepository creates and returns an inmemory user repository.
func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		mutex: sync.Mutex{},
		users: make(map[string]User),
	}
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

// InMemoryVaultRepository provides a vault repository.
type InMemoryVaultRepository struct {
	mutex  sync.Mutex
	vaults map[string]Vault
}

// NewVaultRepository creates and returns an inmemory vault repository.
func NewVaultRepository() *InMemoryVaultRepository {
	return &InMemoryVaultRepository{
		mutex:  sync.Mutex{},
		vaults: make(map[string]Vault),
	}
}

// Add new vault to vault repository if does not exists
func (r *InMemoryVaultRepository) Add(d string, v Vault) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	v, ok := r.vaults[d]
	if ok {
		return ErrVaultAlreadyExists
	}

	r.vaults[d] = v
	return nil
}

// Get an existing user
func (r *InMemoryVaultRepository) Get(d string) (Vault, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	vault, ok := r.vaults[d]
	if !ok {
		vault = Vault{}
	}

	return vault, nil
}
