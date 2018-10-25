package service

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

func TestOnlineSphinx_ExpK(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		_, _, _, _, _, _, err := r.ExpK(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1))
		if err != ErrUserNotFound {
			t.Errorf("Service.ExpK() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})

	t.Run("should return no error if user exists", func(t *testing.T) {
		r := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})
		r.Register(big.NewInt(1))

		_, _, _, _, _, _, err := r.ExpK(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1))
		if err != nil {
			t.Errorf("Service.ExpK() error = %v", err)
		}
	})
}

func TestOnlineSphinx_Verify(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		s := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		_, err := s.Verify(big.NewInt(42), big.NewInt(24), big.NewInt(52))
		if err != nil {
			t.Errorf("Service.Verify() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})
}

func TestOnlineSphinx_GetMetadata(t *testing.T) {
	// TODO: this might expose information for attacker
	t.Run("should return ErrUserNotFound if user does not exists", func(t *testing.T) {
		s := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		_, err := s.GetMetadata(big.NewInt(1), []byte("test"))
		if err != ErrUserNotFound {
			t.Errorf("Service.GetMetadata() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})

	t.Run("should return ErrAuthorizationFailed if user does not exists", func(t *testing.T) {
		s := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		s.Register(big.NewInt(1))
		_, err := s.GetMetadata(big.NewInt(1), []byte("test"))
		if err != ErrAuthorizationFailed {
			t.Errorf("Service.GetMetadata() error = %v wantError = %v", err, ErrAuthorizationFailed)
		}
	})
}

func TestOnlineSphinx_AddVault(t *testing.T) {
	// TODO: this might expose information for attacker
	t.Run("should return ErrUserNotFound if user does not exists", func(t *testing.T) {
		s := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		err := s.AddVault(big.NewInt(1), []byte("test"), "domain")
		if err != ErrUserNotFound {
			t.Errorf("Service.AddVault() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})

	t.Run("should return ErrAuthorizationFailed if user does not exists", func(t *testing.T) {
		s := New(NewUserRepository(), NewVaultRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		s.Register(big.NewInt(1))
		err := s.AddVault(big.NewInt(1), []byte("test"), "domain")
		if err != ErrAuthorizationFailed {
			t.Errorf("Service.AddVault() error = %v wantError = %v", err, ErrAuthorizationFailed)
		}
	})

	t.Run("should add vault", func(t *testing.T) {
		users := NewUserRepository()
		one := big.NewInt(1)
		domain := "test"
		s := New(users, NewVaultRepository(), Configuration{
			sID:  one,
			k:    one,
			q0:   one,
			bits: one,
			hash: sha256.New,
		})

		s.Register(one)

		u, err := users.Get(one.Text(16))
		if err != nil {
			t.Errorf("UserRepository.Get() error = %v", err)
		}
		mac := crypto.HmacData(sha256.New, u.kv.Bytes(), one.Bytes(), one.Bytes(), []byte(domain))

		err = s.AddVault(one, mac, domain)
		if err != nil {
			t.Errorf("Service.AddVault() error = %v", err)
		}
	})
}
