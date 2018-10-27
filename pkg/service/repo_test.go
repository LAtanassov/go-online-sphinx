package service

import (
	"math/big"
	"reflect"
	"testing"
)

func TestUserRepository_Add(t *testing.T) {

	cID := big.NewInt(1)

	t.Run("should add new user", func(t *testing.T) {
		r := NewUserRepository()
		err := r.Set(User{cID: cID})
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}

		u, err := r.Get(cID)
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
		if u.cID.Cmp(cID) != 0 {
			t.Errorf("UserRepository.Add() want = %v but found %v", cID, u.cID)
		}
	})

	t.Run("should override existing user", func(t *testing.T) {
		oldU := User{cID: cID, kv: big.NewInt(1)}
		newU := User{cID: cID, kv: big.NewInt(1)}
		r := NewUserRepository()
		// given
		err := r.Set(oldU)
		if err != nil {
			t.Errorf("UserRepository.Set() error = %v", err)
		}
		// when
		err = r.Set(newU)
		if err != nil {
			t.Errorf("UserRepository.Set() error = %v", err)
		}

		got, err := r.Get(cID)
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
		if !reflect.DeepEqual(got, newU) {
			t.Errorf("GetResponse = %v, want %v", got, newU)
		}
	})

	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := NewUserRepository()

		_, err := r.Get(cID)
		if err != ErrUserNotFound {
			t.Errorf("UserRepository.Get() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})
}
