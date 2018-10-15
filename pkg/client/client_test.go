package client

import (
	"crypto/sha256"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Online SPHINX Client should

// + login an existing user
// + logout an logged in user

func TestClientLogin(t *testing.T) {
	t.Run("should login with password", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := big.NewInt(42)
			buf, err := marshalExpKResponse(n, n, n, n, n)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(buf)
		}))
		defer ts.Close()

		err := New(http.DefaultClient, Configuration{
			expkPath: ts.URL,
			hash:     sha256.New,
			bits:     big.NewInt(8),
		}).Login(User{
			username: "username",
			k:        big.NewInt(42),
			q:        big.NewInt(42),
		}, "password")

		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})
}

func TestClientRegister(t *testing.T) {

	t.Run("should register a new user", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer ts.Close()

		err := New(http.DefaultClient, Configuration{
			registerPath: ts.URL,
		}).Register(User{
			username: "username",
		})
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

	t.Run("should return an error if the user exists within Online SPHINX service", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
		}))
		defer ts.Close()

		err := New(http.DefaultClient, Configuration{
			registerPath: ts.URL,
		}).Register(User{
			username: "username",
		})
		if err != ErrUserNotCreated {
			t.Errorf("Register() error = %v wantErr = %v", err, ErrUserNotCreated)
		}
	})
}
