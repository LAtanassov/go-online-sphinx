package service

import (
	"errors"
	"math/big"
	"sync"
)

// ErrUserNotFound is returned when an user with a given cID does not exists
var ErrUserNotFound = errors.New("user repo: user not found")

// NewUserRepository creates and returns an inmemory user repository.
func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		mutex: sync.Mutex{},
		users: make(map[string]User),
	}
}

// InMemoryUserRepository provides an user repository.
// This repository can also be implemented using an SQL database.UserRepository.
// It should be stored for long term (replicated, shared).
// client.UserRepository is atm identical with server.UserRepository, but this might change in future
type InMemoryUserRepository struct {
	mutex sync.Mutex
	users map[string]User
}

// User is an entity and contains all user related informated to implement server-side Online SPHINX.
type User struct {
	cID    *big.Int
	kv     *big.Int
	vaults map[string]Vault
}

// Vault is an entity and contains cryptographic material necessary to retrieve the vault password
type Vault struct {
	k  *big.Int
	qj *big.Int
}

// Set new or overrides existing user to user repository
func (r *InMemoryUserRepository) Set(u User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[u.cID.Text(16)] = u
	return nil
}

// Get an existing user
func (r *InMemoryUserRepository) Get(cID *big.Int) (User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	u, ok := r.users[cID.Text(16)]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
