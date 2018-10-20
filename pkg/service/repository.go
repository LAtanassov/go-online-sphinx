package service

import (
	"errors"
	"math/big"
	"sync"
)

// ErrUserAlreadyExists in repostory already
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserNotFound in repostory
var ErrUserNotFound = errors.New("user not found")

// User is an entity and contains all user related informated to implement server-side Online SPHINX.
type User struct {
	cID   *big.Int
	kv    *big.Int
	store map[string]Vault
}

// UserRepository provides an user repository.
// This repository can also be implemented using an SQL database.UserRepository.
// It should be store for long term (replicated, shared).
// client.UserRepository is atm identical with server.UserRepository, but this might change in future
type UserRepository struct {
	mutex sync.Mutex
	users map[string]User
}

// NewInMemoryUserRepository creates and returns an inmemory user repository.
func NewInMemoryUserRepository() *UserRepository {
	return &UserRepository{
		mutex: sync.Mutex{},
		users: make(map[string]User),
	}
}

// Add new user to user repository if does not exists
func (r *UserRepository) Add(u User) error {
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
func (r *UserRepository) Get(cID string) (User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	u, ok := r.users[cID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
