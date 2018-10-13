package service

import (
	"math/big"
	"testing"
)

func TestOnlineSphinx_ExpK(t *testing.T) {
	t.Run("should return error if user does not exist", func(t *testing.T) {
		r := New("sID", big.NewInt(1), big.NewInt(1), big.NewInt(1), NewInMemoryRepository())

		_, _, _, _, _, err := r.ExpK("uID", big.NewInt(1), big.NewInt(1))
		if err != ErrUserNotFound {
			t.Errorf("Service.ExpK() error = %v wantError = %v", err, ErrUserNotFound)
		}
	})

	t.Run("should return no error if user exists", func(t *testing.T) {
		r := New("sID", big.NewInt(1), big.NewInt(1), big.NewInt(1), NewInMemoryRepository())
		r.Register("username")

		_, _, _, _, _, err := r.ExpK("username", big.NewInt(1), big.NewInt(1))
		if err != nil {
			t.Errorf("Service.ExpK() error = %v", err)
		}
	})
}
