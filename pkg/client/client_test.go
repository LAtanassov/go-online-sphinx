// +build integration

package client

import (
	"crypto/rand"
	"crypto/sha256"
	"net/http"
	"testing"
)

func TestClient_Register(t *testing.T) {

	t.Run("should register a new user ID", func(t *testing.T) {
		err := New(&http.Client{}, Configuration{baseURL: "http://localhost:8080", registerPath: "/v1/register"}).Register("new-user", "password")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

	t.Run("should be able to register an existing user ID", func(t *testing.T) {
		err := New(&http.Client{}, Configuration{baseURL: "http://localhost:8080", registerPath: "/v1/register"}).Register("another-new-user", "password")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}

		err = New(&http.Client{}, Configuration{baseURL: "http://localhost:8080", registerPath: "/v1/register"}).Register("another-new-user", "password")
		if err == ErrUserNotCreated {
			t.Errorf("Register() error = %v wantErr = %v", err, ErrUserNotCreated)
		}
	})

}

func TestLogin(t *testing.T) {

	q, err := rand.Prime(rand.Reader, 8)
	if err != nil {
		t.Fatal(err)
	}

	c := New(&http.Client{}, Configuration{
		baseURL:      "http://localhost:8080",
		registerPath: "/v1/register",
		q:            q,
		hash:         sha256.New,
	})

	err = c.Register("user", "password")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should login with an valid password", func(t *testing.T) {
		err := c.Login("user", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		c.Logout()
	})

	t.Run("should recv. common error if login failed because of wrong password", func(t *testing.T) {
		err := c.Login("user", "wrong-password")
		if err == ErrAuthenticationFailed {
			t.Errorf("Login() error = %v wantErr = %v", err, ErrAuthenticationFailed)
		}
		c.Logout()
	})

	t.Run("should recv. error if configuration is invalid", func(t *testing.T) {
		err := c.Login("wrong-user", "password")
		if err == ErrAuthenticationFailed {
			t.Errorf("Login() error = %v wantErr = %v", err, ErrAuthenticationFailed)
		}
		c.Logout()
	})
}
