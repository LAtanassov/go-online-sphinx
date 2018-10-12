package server

import (
	"testing"
)

func TestInMemoryRepository_Add(t *testing.T) {

	t.Run("should add new user", func(t *testing.T) {
		r := NewInMemoryRepository()
		err := r.Add(User{username: "hans"})
		if err != nil {
			t.Errorf("InMemoryRepository.Add() error = %v", err)
		}
	})

	t.Run("should return error if an existing user is added again", func(t *testing.T) {
		r := NewInMemoryRepository()
		r.Add(User{username: "hans"})
		err := r.Add(User{username: "hans"})
		if err != ErrUserAlreadyExists {
			t.Errorf("InMemoryRepository.Add() error = %v wantErr = %v", err, ErrUserAlreadyExists)
		}
	})

	t.Run("should return an existing user", func(t *testing.T) {
		r := NewInMemoryRepository()
		err := r.Add(User{username: "hans"})
		if err != nil {
			t.Errorf("InMemoryRepository.Add() error = %v", err)
		}

		u, err := r.Get("hans")
		if err != nil {
			t.Errorf("InMemoryRepository.Get() error = %v", err)
		}

		if u.username != "hans" {
			t.Errorf("InMemoryRepository.Get() error = %v", err)
		}
	})

	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := NewInMemoryRepository()

		_, err := r.Get("hans")
		if err != ErrUserNotFound {
			t.Errorf("InMemoryRepository.Get() error = %v wantError = %v", err, ErrUserNotFound)
		}

	})
}
