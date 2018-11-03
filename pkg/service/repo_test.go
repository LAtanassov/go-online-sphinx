package service

import (
	"math/big"
	"reflect"
	"testing"
)

func TestUserRepository_Add(t *testing.T) {

	cID := big.NewInt(1)

	t.Run("should add new user", func(t *testing.T) {
		// given
		r := NewUserRepository()
		wantUser := User{cID: cID}
		// when
		err := r.Set(wantUser)
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
		// then expect
		gotUser, err := r.Get(cID)
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}

		if !reflect.DeepEqual(wantUser, gotUser) {
			t.Errorf("UserRepository.Add() wantUser = %v but gotUser = %v", wantUser, gotUser)
		}
	})

	t.Run("should override existing user", func(t *testing.T) {
		oldUser := User{cID: cID, kv: big.NewInt(1)}
		newUser := User{cID: cID, kv: big.NewInt(1)}
		r := NewUserRepository()
		// given
		err := r.Set(oldUser)
		if err != nil {
			t.Errorf("UserRepository.Set() error = %v", err)
		}
		// when
		err = r.Set(newUser)
		if err != nil {
			t.Errorf("UserRepository.Set() error = %v", err)
		}

		gotUser, err := r.Get(cID)
		if err != nil {
			t.Errorf("UserRepository.Add() error = %v", err)
		}
		if !reflect.DeepEqual(newUser, gotUser) {
			t.Errorf("Get() wantUser = %v, gotUser %v", newUser, gotUser)
		}
	})

	t.Run("should return ErrUserNotFound if user does not exist", func(t *testing.T) {
		r := NewUserRepository()

		_, err := r.Get(cID)
		if err != ErrUserNotFound {
			t.Errorf("UserRepository.Get() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})
}
