package client

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/pkg/errors"
)

var (
	// ErrUserAlreadyExists in repostory already
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound in repostory
	ErrUserNotFound = errors.New("user not found")
)

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

// User specific configuration contains
// a client ID and important login-specific variables like prime q and secret k.
type User struct {
	username string
	cID      *big.Int
	q        *big.Int
	k        *big.Int
}

// NewUser generates new user with username.
func NewUser(username string, bits int) (User, error) {
	q, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		return User{}, errors.Wrap(err, "NewUser: failed to generate radnom prime")
	}

	cID, err := rand.Int(rand.Reader, q)
	if err != nil {
		return User{}, errors.Wrap(err, "NewUser: failed to generate random int")
	}

	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		return User{}, errors.Wrap(err, "NewUser: failed to generate random int")
	}
	return User{
		username: username,
		cID:      cID,
		k:        k,
		q:        q,
	}, nil
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
		return errors.Wrap(ErrUserAlreadyExists, "Add")
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
		return User{}, errors.Wrap(ErrUserNotFound, "Get")
	}
	return u, nil
}
