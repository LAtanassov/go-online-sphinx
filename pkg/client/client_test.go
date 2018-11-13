package client

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
	"github.com/pkg/errors"
)

func TestClient_Register(t *testing.T) {

	repo := NewInMemoryUserRepository()

	t.Run("should register a new user", func(t *testing.T) {
		// given
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer ts.Close()

		// when
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		err = New(http.DefaultClient, cfg, repo).Register("username")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}

		_, err = repo.Get("username")
		if err != nil {
			t.Errorf("Get() expect repo to return user but error = %v", err)
		}
	})

	t.Run("should return an error if the user exists within Online SPHINX service", func(t *testing.T) {
		// given
		ErrTest := errors.New("unit test")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			contract.MarshalError(w, ErrTest)
		}))
		defer ts.Close()

		// when
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		err = New(http.DefaultClient, cfg, repo).Register("username")

		if err == nil {
			t.Errorf("Register() error = %v wantErr = %v", err, ErrTest)
		}
	})
}

func TestClient_Login(t *testing.T) {
	// before
	user, err := newUser("username", 8)
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	t.Run("should login with password", func(t *testing.T) {
		// when
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := big.NewInt(42)

			var buf bytes.Buffer
			wr := bufio.NewWriter(&buf)
			err := contract.MarshalExpKResponse(wr, contract.ExpKResponse{
				SID:    n,
				SNonce: n,
				BD:     n,
				Q0:     n,
				KV:     n})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			wr.Flush()

			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		}))
		defer ts.Close()

		// then
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		err = New(http.DefaultClient, cfg, repo).Login("username", "password")

		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
	})

}

func TestClient_Challenge(t *testing.T) {
	// before
	user, err := newUser("username", 8)
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	sID := big.NewInt(10)
	ski := big.NewInt(10)
	mk := big.NewInt(10)

	t.Run("should return an ErrAuthenticationFailed if challenge/respond is not correct", func(t *testing.T) {
		// when
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := big.NewInt(42)

			var buf bytes.Buffer
			wr := bufio.NewWriter(&buf)
			err := contract.MarshalChallengeResponse(wr, contract.ChallengeResponse{R: n})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			wr.Flush()

			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		}))
		defer ts.Close()
		// when
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		clt := New(http.DefaultClient, cfg, repo)
		clt.session = NewSession(user, sID, ski, mk)

		err = clt.Challenge()
		if errors.Cause(err) != ErrOperationFailed {
			t.Errorf("Challenge() error = %v wantErr = %v", err, ErrOperationFailed)
		}
	})
}

func TestClient_GetMetadata(t *testing.T) {
	// before
	user, err := newUser("username", 8)
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	sID := big.NewInt(10)
	ski := big.NewInt(10)
	mk := big.NewInt(10)

	t.Run("should return domains", func(t *testing.T) {
		// given
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			wr := bufio.NewWriter(&buf)
			err := contract.MarshalMetadataResponse(wr, contract.MetadataResponse{Domains: []string{"domain"}})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			wr.Flush()

			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		}))
		defer ts.Close()
		// when
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		clt := New(http.DefaultClient, cfg, repo)
		clt.session = NewSession(user, sID, ski, mk)

		_, err = clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}
	})
}

func TestClient_Add(t *testing.T) {
	// before
	user, err := newUser("username", 8)
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	sID := big.NewInt(10)
	ski := big.NewInt(10)
	mk := big.NewInt(10)

	t.Run("should return domains", func(t *testing.T) {
		// given
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer ts.Close()
		// when
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		clt := New(http.DefaultClient, cfg, repo)
		clt.session = NewSession(user, sID, ski, mk)

		err = clt.Add("google.com")
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}
	})
}

func TestClient_Get(t *testing.T) {
	// before
	user, err := newUser("username", 8)
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	sID := big.NewInt(10)
	ski := big.NewInt(10)
	mk := big.NewInt(10)

	t.Run("should return domains", func(t *testing.T) {
		// given
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			wr := bufio.NewWriter(&buf)
			err = contract.MarshalGetResponse(wr, contract.GetResponse{Bj: big.NewInt(1), Qj: big.NewInt(1)})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			wr.Flush()

			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		}))
		defer ts.Close()
		// when
		cfg, err := NewConfiguration(ts.URL, 8, sha256.New)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		clt := New(http.DefaultClient, cfg, repo)
		clt.session = NewSession(user, sID, ski, mk)

		_, err = clt.Get("google.com")
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
	})
}
