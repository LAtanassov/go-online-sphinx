package service

import (
	"math/big"
	"testing"
)

func TestInMemoryUserRepository_Add(t *testing.T) {

	t.Run("should add new user", func(t *testing.T) {
		r := NewInMemoryUserRepository()
		err := r.Add(User{cID: big.NewInt(1)})
		if err != nil {
			t.Errorf("InMemoryRepository.Add() error = %v", err)
		}
	})

	t.Run("should return error if an existing user is added again", func(t *testing.T) {
		r := NewInMemoryUserRepository()
		r.Add(User{cID: big.NewInt(1)})
		err := r.Add(User{cID: big.NewInt(1)})
		if err != ErrUserAlreadyExists {
			t.Errorf("InMemoryRepository.Add() error = %v wantErr = %v", err, ErrUserAlreadyExists)
		}
	})

	t.Run("should return an existing user", func(t *testing.T) {
		r := NewInMemoryUserRepository()
		err := r.Add(User{cID: big.NewInt(1)})
		if err != nil {
			t.Errorf("InMemoryRepository.Add() error = %v", err)
		}

		u, err := r.Get(big.NewInt(1).Text(16))
		if err != nil {
			t.Errorf("InMemoryRepository.Get() error = %v", err)
		}

		if u.cID.Text(16) != big.NewInt(1).Text(16) {
			t.Errorf("InMemoryRepository.Get() error = %v", err)
		}
	})

	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := NewInMemoryUserRepository()

		_, err := r.Get(big.NewInt(1).Text(16))
		if err != ErrUserNotFound {
			t.Errorf("InMemoryRepository.Get() error = %v wantError = %v", err, ErrUserNotFound)
		}

	})
}
