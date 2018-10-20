package client

import (
	"crypto/sha256"
	"errors"
	"hash"
	"math/big"
	"sync"
)

// ErrUserAlreadyExists in repostory already
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserNotFound in repostory
var ErrUserNotFound = errors.New("user not found")

// User specific configuration contains
// a client ID and important login-specific variables like prime q and secret k.
type User struct {
	username string
	cID      *big.Int
	q        *big.Int
	k        *big.Int
}

// NewUser generates new user with username.
func NewUser(username string) (User, error) {
	// TODO: generate keys here
	return User{
		username: username,
		cID:      big.NewInt(42),
		k:        big.NewInt(42),
		q:        big.NewInt(42),
	}, nil
}

// Configuration ...
type Configuration struct {
	hash         func() hash.Hash
	bits         *big.Int
	contentType  string
	baseURL      string
	registerPath string
	expkPath     string
	verifyPath   string
	metadataPath string
}

// NewConfiguration return default configuration.
func NewConfiguration() Configuration {
	return Configuration{
		hash: sha256.New,
		bits: big.NewInt(8),
	}
}

// UserRepository contains user specific configuration.
// Adds new user to the repository when registered.
// Load existing user from the repository before the login process.
// UserRepository SHOULD be able to store this configuration as file ,
// so that users can easily copy and transfer those files.
// client.UserRepository is atm identical with server.UserRepository, but this might change in future
type UserRepository struct {
	mutex sync.Mutex
	users map[string]User
}

// NewInMemoryUserRepository return an in memory UserRepository.
// using pointer semantic allocated in heap once for sharing
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

	_, ok := r.users[u.username]
	if ok {
		return ErrUserAlreadyExists
	}

	r.users[u.username] = u
	return nil
}

// Get an existing user
func (r *UserRepository) Get(username string) (User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	u, ok := r.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}