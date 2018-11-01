package client

import (
	"testing"

	"github.com/pkg/errors"
)

func TestUserRepository_Add(t *testing.T) {
	t.Run("should add new user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		user, err := NewUser("username", 8)
		if err != nil {
			t.Errorf("NewUser() failed error = %v", err)
		}
		err = repo.Add(user)
		if err != nil {
			t.Errorf("InMemoryRepository.Add() error = %v", err)
		}
	})

	t.Run("should return error if an existing user is added again", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		user, err := NewUser("username", 8)
		if err != nil {
			t.Errorf("NewUser() failed error = %v", err)
		}

		err = repo.Add(user)
		if err != nil {
			t.Errorf("Add() failed error = %v", err)
		}

		err = repo.Add(user)
		if errors.Cause(err) != ErrUserAlreadyExists {
			t.Errorf("InMemoryRepository.Add() error = %v wantErr = %v", err, ErrUserAlreadyExists)
		}
	})

	t.Run("should return an existing user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		expUser, err := NewUser("username", 8)
		if err != nil {
			t.Errorf("NewUser() failed error = %v", err)
		}

		err = repo.Add(expUser)
		if err != nil {
			t.Errorf("InMemoryRepository.Add() error = %v", err)
		}

		actUser, err := repo.Get(expUser.username)
		if err != nil {
			t.Errorf("InMemoryRepository.Get() error = %v", err)
		}

		if actUser.username != expUser.username {
			t.Errorf("InMemoryRepository.Get() error = %v", err)
		}
	})
}

func TestUserRepository_Get(t *testing.T) {
	r := NewInMemoryUserRepository()

	_, err := r.Get("username")
	if errors.Cause(err) != ErrUserNotFound {
		t.Errorf("InMemoryRepository.Get() error = %v wantError = %v", err, ErrUserNotFound)
	}
}
