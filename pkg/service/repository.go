package service

import (
	"errors"
	"math/big"
	"sync"
)

// ErrUserAlreadyExists in repostory already
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrDomainAlreadyExists in repostory already
var ErrDomainAlreadyExists = errors.New("domain already exists")

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

// InMemoryDomainRepository provides an domain repository.
type InMemoryDomainRepository struct {
	mutex   sync.Mutex
	domains map[string][]Domain
}

// NewDomainRepository creates and returns an inmemory user repository.
func NewDomainRepository() *InMemoryDomainRepository {
	return &InMemoryDomainRepository{
		mutex:   sync.Mutex{},
		domains: make(map[string][]Domain),
	}
}

// Add new user to user repository if does not exists
func (r *InMemoryDomainRepository) Add(id string, d Domain) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	domains, ok := r.domains[id]
	if ok {
		return ErrDomainAlreadyExists
	}

	r.domains[id] = append(domains, d)
	return nil
}

// Get an existing user
func (r *InMemoryDomainRepository) Get(id string) ([]Domain, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	domains, ok := r.domains[id]
	if !ok {
		domains = []Domain{}
	}

	return domains, nil
}
