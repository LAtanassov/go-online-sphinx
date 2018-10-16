package service

import (
	"crypto/sha256"
	"math/big"
	"os"
	"testing"
)

func TestOnlineSphinx_ExpK(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := New(NewInMemoryRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		_, _, _, _, _, err := r.ExpK(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1))
		if err != ErrUserNotFound {
			t.Errorf("Service.ExpK() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})

	t.Run("should return no error if user exists", func(t *testing.T) {
		r := New(NewInMemoryRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})
		r.Register(big.NewInt(1))

		_, _, _, _, _, err := r.ExpK(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1))
		if err != nil {
			t.Errorf("Service.ExpK() error = %v", err)
		}
	})
}

func TestOnlineSphinx_Verify(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		s := New(NewInMemoryRepository(), Configuration{
			sID:  big.NewInt(1),
			k:    big.NewInt(1),
			q0:   big.NewInt(1),
			bits: big.NewInt(1),
			hash: sha256.New,
		})

		os.Setenv("SKi", big.NewInt(13).Text(16))

		_, err := s.Verify(big.NewInt(24), big.NewInt(52))
		if err != nil {
			t.Errorf("Service.ExpK() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})
}
