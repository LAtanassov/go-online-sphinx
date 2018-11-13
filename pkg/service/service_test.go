package service

import (
	"crypto/sha256"
	"math/big"
	"reflect"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/crypto"
)

func TestOnlineSphinx_ExpK(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		// given
		r := New(
			NewUserRepository(),
			NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
		)

		// when
		_, _, _, _, _, _, err := r.ExpK(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1))
		// then
		if err == nil {
			t.Errorf("Service.ExpK() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})

	t.Run("should b ** k mod q ", func(t *testing.T) {
		// given
		config := NewConfiguration(big.NewInt(1), big.NewInt(13), big.NewInt(1), big.NewInt(1), sha256.New)
		r := New(NewUserRepository(), config)
		r.Register(one)
		cID := one
		cNonce := one
		b := big.NewInt(23)
		q := big.NewInt(31)
		want := crypto.ExpInGroup(b, config.k, q)

		// when
		_, _, _, bd, _, _, err := r.ExpK(cID, cNonce, b, q)
		if err != nil {
			t.Errorf("Service.ExpK() error = %v", err)
		}

		if !reflect.DeepEqual(bd, want) {
			t.Errorf("Service.ExpK() want = %v but got %v", want, bd)
		}
	})
}

func TestOnlineSphinx_Challenge(t *testing.T) {
	t.Run("should return g ** k mod q", func(t *testing.T) {
		s := New(
			NewUserRepository(),
			NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
		)
		g := big.NewInt(42)
		ski := big.NewInt(24)
		q := big.NewInt(52)

		want := crypto.ExpInGroup(g, ski, q)

		got, err := s.Challenge(ski, g, q)
		if err != nil {
			t.Errorf("Service.Challenge() error = %v wantError = %v", err, ErrUserNotFound)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Service.Challenge() want = %v but got %v", want, got)
		}
	})
}

func TestOnlineSphinx_GetMetadata(t *testing.T) {
	t.Run("should return all domains", func(t *testing.T) {
		// given
		s := New(
			NewUserRepository(),
			NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
		)
		cID := big.NewInt(1)

		s.Register(cID)
		// when
		domains, err := s.GetMetadata(cID)
		if err != nil {
			t.Errorf("Service.GetMetadata() error = %v", err)
		}
		if len(domains) != 0 {
			t.Errorf("Service.GetMetadata() expect no domains but got %v domains", len(domains))
		}
	})
}

func TestOnlineSphinx_AddVault(t *testing.T) {
	t.Run("should add vault", func(t *testing.T) {
		s := New(
			NewUserRepository(),
			NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
		)
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
		// given
		s := New(
			NewUserRepository(),
			NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
		)

		cID := big.NewInt(1)
		s.Register(cID)
		err := s.Add(cID, "domain")
		// when
		_, _, err = s.Get(cID, "domain", big.NewInt(1), big.NewInt(2))
		if err != nil {
			t.Errorf("Service.AddVault() error = %v", err)
		}
	})
}

func TestOnlineSphinx_VerifyMAC(t *testing.T) {
	t.Run("should verify MAC", func(t *testing.T) {
		s := New(
			NewUserRepository(),
			NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
		)

		ski := big.NewInt(31)
		want := crypto.HmacData(sha256.New, ski.Bytes(), []byte("data"))
		err := s.VerifyMAC(want, ski, []byte("data"))

		if err != nil {
			t.Errorf("Service.AddVault() error = %v", err)
		}
	})
}
