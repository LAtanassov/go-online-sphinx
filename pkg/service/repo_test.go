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

	t.Run("should return error if an user already exists", func(t *testing.T) {
		r := NewUserRepository()
		// given
		err := r.Add(User{cID: big.NewInt(1)})
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
		// when
		err = r.Add(User{cID: big.NewInt(1)})
		if err != ErrUserAlreadyExists {
			t.Errorf("UserRepository.Add() error = %v wantErr = %v", err, ErrUserAlreadyExists)
		}
	})

	t.Run("should return an existing user", func(t *testing.T) {
		r := NewUserRepository()
		// given
		err := r.Add(User{cID: big.NewInt(1)})
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
		// when
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

func TestVaultRepository_Add(t *testing.T) {

	t.Run("should add new vault", func(t *testing.T) {
		r := NewVaultRepository()
		err := r.Add("domain", Vault{})
		if err != nil {
			t.Errorf("VaultRepository.Add() error = %v", err)
		}
	})

	t.Run("should return error if a vault already exists", func(t *testing.T) {
		r := NewVaultRepository()
		// given
		err := r.Add("domain", Vault{})
		if err != nil {
			t.Errorf("VaultRepository.Add() error = %v", err)
		}
		// when
		err = r.Add("domain", Vault{})
		if err != ErrVaultAlreadyExists {
			t.Errorf("VaultRepository.Add() error = %v wantErr = %v", err, ErrVaultAlreadyExists)
		}
	})

	t.Run("should return an existing vault", func(t *testing.T) {
		r := NewVaultRepository()
		// given
		err := r.Add("domain", Vault{})
		if err != nil {
			t.Errorf("VaultRepository.Add() error = %v", err)
		}
		// when
		_, err = r.Get("domain")
		if err != nil {
			t.Errorf("VaultRepository.Get() error = %v", err)
		}
	})

	t.Run("should return error if vault does not exist", func(t *testing.T) {
		r := NewVaultRepository()

		_, err := r.Get("id")
		if err != nil {
			t.Errorf("VaultRepository.Get() error = %v", err)
		}
	})
}
