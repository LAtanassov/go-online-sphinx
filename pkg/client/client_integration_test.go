// +build integration

package client_test

import (
	"crypto/sha256"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"
)

var baseURL = "http://localhost:8080"
var bits = 8
var hashFn = sha256.New

func TestITClient_Register(t *testing.T) {

	t.Run("should register a new user ID", func(t *testing.T) {

		cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		clt := client.New(
			&http.Client{},
			cfg,
			client.NewInMemoryUserRepository())

		err = clt.Register("registered-user")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

	t.Run("should not be able to register with an existing user ID", func(t *testing.T) {
		cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
		if err != nil {
			t.Errorf("NewConfiguration() error = %v", err)
		}
		clt := client.New(
			&http.Client{},
			cfg,
			client.NewInMemoryUserRepository())

		err = clt.Register("double-registered-user")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
		// when
		err = clt.Register("double-registered-user")
		if err == nil {
			t.Errorf("Register() no error but got err = %v", err)
		}
	})

}

func TestITClient_Login(t *testing.T) {

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		t.Errorf("cookiejar.New() error = %v", err)
	}
	httpClient := &http.Client{
		Jar: cookieJar,
	}

	cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
	if err != nil {
		t.Errorf("NewConfiguration() error = %v", err)
	}
	clt := client.New(
		httpClient,
		cfg,
		client.NewInMemoryUserRepository())

	err = clt.Register("login-username")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should login with a valid password", func(t *testing.T) {
		err := clt.Login("login-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Challenge()
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})

	t.Run("should recv. common error if wrong password", func(t *testing.T) {
		err := clt.Login("login-username", "wrong-password")
		if err == nil {
			t.Errorf("Login() want error but got wantErr = %v", err)
		}
		err = clt.Challenge()
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})
}

func TestITClient_GetMetadata(t *testing.T) {

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		t.Errorf("cookiejar.New() error = %v", err)
	}

	httpClient := &http.Client{
		Jar: cookieJar,
	}
	cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
	if err != nil {
		t.Errorf("NewConfiguration() error = %v", err)
	}
	clt := client.New(
		httpClient,
		cfg,
		client.NewInMemoryUserRepository())

	err = clt.Register("get-metadata-username")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should have no domains", func(t *testing.T) {
		err := clt.Login("get-metadata-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		err = clt.Challenge()
		if err != nil {
			t.Errorf("Challenge() error = %v", err)
		}

		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if len(domains) != 0 {
			t.Errorf("domains = %v wantDomains = %v", domains, []string{})
		}
	})

	t.Run("should have google.com domain", func(t *testing.T) {
		// given
		wantDomains := []string{"google.com"}
		err := clt.Login("get-metadata-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		err = clt.Challenge()
		if err != nil {
			t.Errorf("Challenge() error = %v", err)
		}

		err = clt.Add(wantDomains[0])
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if !reflect.DeepEqual(domains, wantDomains) {
			t.Errorf("domains = %v wantDomains = %v", domains, wantDomains)
		}
	})
}

func TestITClient_Add(t *testing.T) {

	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
	}
	cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
	if err != nil {
		t.Errorf("NewConfiguration() error = %v", err)
	}
	clt := client.New(
		httpClient,
		cfg,
		client.NewInMemoryUserRepository())

	err = clt.Register("add-domain-username")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should add a new domain with name google.com", func(t *testing.T) {
		// given
		wantDomains := []string{"google.com"}
		err := clt.Login("add-domain-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		err = clt.Challenge()
		if err != nil {
			t.Errorf("Challenge() error = %v", err)
		}

		err = clt.Add(wantDomains[0])
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if !reflect.DeepEqual(domains, wantDomains) {
			t.Errorf("domains = %v wantDomains = %v", domains, wantDomains)
		}
	})
}

func TestITClient_Get(t *testing.T) {

	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: cookieJar,
	}
	cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
	if err != nil {
		t.Errorf("NewConfiguration() error = %v", err)
	}
	clt := client.New(
		httpClient,
		cfg,
		client.NewInMemoryUserRepository())

	err = clt.Register("get-domain-username")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should get the same password within one session", func(t *testing.T) {
		// given
		domain := "google.com"
		err := clt.Login("get-domain-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		err = clt.Challenge()
		if err != nil {
			t.Errorf("Challenge() error = %v", err)
		}

		err = clt.Add(domain)
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		pwda, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		pwdb, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		if !reflect.DeepEqual(pwda, pwdb) {
			t.Errorf("pwda = %v pwdb = %v", pwda, pwdb)
		}

		clt.Logout()
	})

	t.Run("should get the same password from two different sessions", func(t *testing.T) {
		// given
		domain := "google.com"
		err := clt.Login("get-domain-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		err = clt.Challenge()
		if err != nil {
			t.Errorf("Challenge() error = %v", err)
		}

		err = clt.Add(domain)
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		pwda, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		clt.Logout()

		err = clt.Login("get-domain-username", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		err = clt.Challenge()
		if err != nil {
			t.Errorf("Challenge() error = %v", err)
		}

		pwdb, err := clt.Get(domain)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}

		if !reflect.DeepEqual(pwda, pwdb) {
			t.Errorf("pwda = %v pwdb = %v", pwda, pwdb)
		}
	})
}
