package service

import (
	"crypto/sha256"
	"math/big"
	"testing"
)

func TestOnlineSphinx_ExpK(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := New(NewUserRepository(), Configuration{
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
		r := New(NewUserRepository(), Configuration{
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

func TestOnlineSphinx_Challenge(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		s := New(NewUserRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		_, err := s.Challenge(big.NewInt(42), big.NewInt(24), big.NewInt(52))
		if err != nil {
			t.Errorf("Service.Challenge() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})
}

func TestOnlineSphinx_GetMetadata(t *testing.T) {
	t.Run("should return all domains", func(t *testing.T) {
		s := New(NewUserRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})
		cID := big.NewInt(1)

		s.Register(cID)
		_, err := s.GetMetadata(cID)
		if err != nil {
			t.Errorf("Service.GetMetadata() error = %v", err)
		}
	})
}

func TestOnlineSphinx_AddVault(t *testing.T) {
	t.Run("should add vault", func(t *testing.T) {
		s := New(NewUserRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})
		cID := big.NewInt(1)

		s.Register(cID)
		err := s.Add(cID, "domain")
		if err != nil {
			t.Errorf("Service.AddVault() error = %v", err)
		}
	})
}

func TestOnlineSphinx_GetVault(t *testing.T) {
	t.Run("should get vault", func(t *testing.T) {
		s := New(NewUserRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		cID := big.NewInt(1)
		s.Register(cID)
		err := s.Add(cID, "domain")
		_, _, err = s.Get(cID, "domain", big.NewInt(1), big.NewInt(2))
		if err != nil {
			t.Errorf("Service.AddVault() error = %v", err)
		}
	})
}

func TestOnlineSphinx_VerifyMAC(t *testing.T) {
	t.Run("should get vault", func(t *testing.T) {
		s := New(NewUserRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})
		cID := big.NewInt(1)
		s.Register(cID)

		err := s.VerifyMAC([]byte("mac"), cID, []byte("data"))

		if err != ErrAuthorizationFailed {
			t.Errorf("Service.AddVault() error = %v wantErr = %v", err, ErrAuthorizationFailed)
		}
	})
}
