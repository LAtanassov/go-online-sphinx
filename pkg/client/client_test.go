package client

import (
	"crypto/sha256"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClient_Login(t *testing.T) {
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
		}, User{
			username: "username",
			cID:      big.NewInt(42),
			k:        big.NewInt(42),
			q:        big.NewInt(42),
		}).Login("username", "password")

		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})
}

func TestClient_Register(t *testing.T) {

	t.Run("should register a new user", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer ts.Close()

		err := New(http.DefaultClient, Configuration{
			registerPath: ts.URL,
		}, User{}).Register("username")
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
		}, User{}).Register("username")
		if err != ErrRegistrationFailed {
			t.Errorf("Register() error = %v wantErr = %v", err, ErrRegistrationFailed)
		}
	})
}

func TestClient_Verify(t *testing.T) {

	t.Run("should return an ErrAuthenticationFailed if challenge/respond is not correct", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := big.NewInt(42)
			buf, err := marshalVerifyResponse(n, nil)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(buf)
		}))
		defer ts.Close()
		os.Setenv("SKi", "A")

		clt := New(http.DefaultClient, Configuration{
			verifyPath: ts.URL,
		}, User{
			q: big.NewInt(1),
		})
		clt.session.ski = big.NewInt(10)

		err := clt.Verify()
		if err != ErrAuthenticationFailed {
			t.Errorf("Verify() error = %v wantErr = %v", err, ErrAuthenticationFailed)
		}
	})
}

func TestClient_GetMetadata(t *testing.T) {

	t.Run("should return domains", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, err := marshalMetadataResponse([]Domain{NewDomain()}, nil)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(buf)
		}))
		defer ts.Close()

		clt := New(http.DefaultClient, Configuration{
			metadataPath: ts.URL,
			hash:         sha256.New,
		}, User{
			cID: big.NewInt(1),
			q:   big.NewInt(1),
		})
		clt.session.ski = big.NewInt(10)
		clt.session.sID = big.NewInt(10)

		_, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}
	})
}
