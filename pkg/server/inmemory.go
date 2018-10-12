package server

import (
	"errors"
	"sync"
)

// ErrUserAlreadyExists in repostory already
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserNotFound in repostory
var ErrUserNotFound = errors.New("user not found")

// InMemoryRepository provides an inmemory user repository
type InMemoryRepository struct {
	mutex sync.Mutex
	users map[string]User
}

// NewInMemoryRepository creates and returns an inmemory user repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		mutex: sync.Mutex{},
		users: make(map[string]User),
	}
}

// Add new user to user repository if does not exists
func (r *InMemoryRepository) Add(u User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.users[u.username]
	if ok {
		return ErrUserAlreadyExists
	}

	r.users[u.username] = u
	return nil
}

// Get an existing user
func (r *InMemoryRepository) Get(username string) (User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	u, ok := r.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
