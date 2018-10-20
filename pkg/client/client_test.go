package client

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Register(t *testing.T) {

	repo := NewInMemoryUserRepository()

	cfg := NewConfiguration()

	t.Run("should register a new user", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer ts.Close()

		cfg.registerPath = ts.URL
		err := New(http.DefaultClient, cfg, repo).Register("username")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}

		_, err = repo.Get("username")
		if err != nil {
			t.Errorf("Register() expect repo to return user but error = %v", err)
		}
	})

	t.Run("should return an error if the user exists within Online SPHINX service", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
		}))
		defer ts.Close()

		cfg.registerPath = ts.URL
		err := New(http.DefaultClient, cfg, repo).Register("username")
		if err != ErrUserAlreadyExists {
			t.Errorf("Register() error = %v wantErr = %v", err, ErrUserAlreadyExists)
		}
	})
}

func TestClient_Login(t *testing.T) {
	user, err := NewUser("username")
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	cfg := NewConfiguration()

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

		cfg.expkPath = ts.URL
		err := New(http.DefaultClient, cfg, repo).Login("username", "password")

		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

}

func TestClient_Verify(t *testing.T) {
	user, err := NewUser("username")
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	cfg := NewConfiguration()

	sID := big.NewInt(10)
	ski := big.NewInt(10)
	mk := big.NewInt(10)

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

		cfg.verifyPath = ts.URL

		clt := New(http.DefaultClient, cfg, repo)
		clt.session = NewSession(user, sID, ski, mk)

		err = clt.Verify()
		if err != ErrAuthenticationFailed {
			t.Errorf("Verify() error = %v wantErr = %v", err, ErrAuthenticationFailed)
		}
	})
}

func TestClient_GetMetadata(t *testing.T) {
	user, err := NewUser("username")
	if err != nil {
		t.Errorf("before test started - error = %v", err)
	}
	repo := NewInMemoryUserRepository()
	repo.Add(user)

	cfg := NewConfiguration()
	sID := big.NewInt(10)
	ski := big.NewInt(10)
	mk := big.NewInt(10)

	t.Run("should return domains", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, err := marshalMetadataResponse([]Domain{Domain{}}, nil)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(buf)
		}))
		defer ts.Close()

		cfg.metadataPath = ts.URL
		clt := New(http.DefaultClient, cfg, repo)
		clt.session = NewSession(user, sID, ski, mk)

		_, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}
	})
}
