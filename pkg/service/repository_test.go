package service

import (
	"math/big"
	"testing"
)

func TestUserRepository_Add(t *testing.T) {

	t.Run("should add new user", func(t *testing.T) {
		r := NewUserRepository()
		err := r.Add(User{cID: big.NewInt(1)})
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
	})

	t.Run("should return error if an existing user is added again", func(t *testing.T) {
		r := NewUserRepository()
		r.Add(User{cID: big.NewInt(1)})
		err := r.Add(User{cID: big.NewInt(1)})
		if err != ErrUserAlreadyExists {
			t.Errorf("UserRepository.Add() error = %v wantErr = %v", err, ErrUserAlreadyExists)
		}
	})

	t.Run("should return an existing user", func(t *testing.T) {
		r := NewUserRepository()
		err := r.Add(User{cID: big.NewInt(1)})
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}

		u, err := r.Get(big.NewInt(1).Text(16))
		if err != nil {
			t.Errorf("UserRepository.Get() error = %v", err)
		}

		if u.cID.Text(16) != big.NewInt(1).Text(16) {
			t.Errorf("UserRepository.Get() error = %v", err)
		}
	})

	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := NewUserRepository()

		_, err := r.Get(big.NewInt(1).Text(16))
		if err != ErrUserNotFound {
			t.Errorf("UserRepository.Get() error = %v wantError = %v", err, ErrUserNotFound)
		}

	})
}

func TestDomainRepository_Add(t *testing.T) {

	t.Run("should add new user", func(t *testing.T) {
		r := NewDomainRepository()
		err := r.Add("id", Domain{})
		if err != nil {
			t.Errorf("DomainRepository.Add() error = %v", err)
		}
	})

	t.Run("should return error if an existing user is added again", func(t *testing.T) {
		r := NewDomainRepository()
		err := r.Add("id", Domain{})
		if err != nil {
			t.Errorf("DomainRepository.Add() error = %v", err)
		}

		err = r.Add("id", Domain{})
		if err != ErrDomainAlreadyExists {
			t.Errorf("DomainRepository.Add() error = %v wantErr = %v", err, ErrDomainAlreadyExists)
		}
	})

	t.Run("should return an existing domain", func(t *testing.T) {
		r := NewDomainRepository()
		err := r.Add("id", Domain{})
		if err != nil {
			t.Errorf("DomainRepository.Add() error = %v", err)
		}

		domains, err := r.Get("id")
		if err != nil {
			t.Errorf("DomainRepository.Get() error = %v", err)
		}

		if domains == nil || len(domains) != 1 {
			t.Errorf("DomainRepository.Get() error = %v", err)
		}
	})

	t.Run("should return error if domain does not exist", func(t *testing.T) {
		r := NewDomainRepository()

		domains, err := r.Get("id")
		if err != nil {
			t.Errorf("DomainRepository.Get() error = %v", err)
		}

		if domains == nil && len(domains) != 0 {
			t.Errorf("DomainRepository.Get() error = %v", err)
		}
	})
}
