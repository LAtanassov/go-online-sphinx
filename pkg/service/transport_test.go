package service

import (
	"crypto/sha256"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
	"github.com/go-kit/kit/log"
)

func TestMakeRegisterHandler(t *testing.T) {
	s := New(
		NewUserRepository(),
		NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
	)
	ct := "application/json"

	t.Run("should register a user", func(t *testing.T) {

		ts := httptest.NewServer(NewHTTPTransport(s, log.NewNopLogger()).MakeRegisterHandler())
		defer ts.Close()

		r, err := contract.MarshalRegisterRequest(contract.RegisterRequest{CID: big.NewInt(1)})
		if err != nil {
			t.Errorf("MarshalRegisterRequest() error = %v", err)
		}

		_, err = http.Post(ts.URL+"/v1/register", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})

}

func TestMakeExpKHandler(t *testing.T) {
	s := New(
		NewUserRepository(),
		NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
	)
	ct := "application/json"

	t.Run("should exponent to k given blinded secret", func(t *testing.T) {

		ts := httptest.NewServer(NewHTTPTransport(s, log.NewNopLogger()).MakeExpKHandler())
		defer ts.Close()

		r, err := contract.MarshalExpKRequest(contract.ExpKRequest{
			CID:    big.NewInt(1),
			CNonce: big.NewInt(2),
			B:      big.NewInt(3),
			Q:      big.NewInt(4),
		})
		if err != nil {
			t.Errorf("contract.MarshalExpKRequest() error = %v", err)
		}

		_, err = http.Post(ts.URL+"/v1/login/expk", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeChallengeHandler(t *testing.T) {
	s := New(
		NewUserRepository(),
		NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
	)
	ct := "application/json"

	t.Run("should respond to a given challenge", func(t *testing.T) {

		ts := httptest.NewServer(NewHTTPTransport(s, log.NewNopLogger()).MakeChallengeHandler())
		defer ts.Close()

		r, err := contract.MarshalChallengeRequest(contract.ChallengeRequest{
			G: big.NewInt(3),
			Q: big.NewInt(4),
		})
		if err != nil {
			t.Errorf("contract.MarshalChallengeRequest() error = %v", err)
		}

		_, err = http.Post(ts.URL+"/v1/login/challenge", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeMetadataHandler(t *testing.T) {
	s := New(
		NewUserRepository(),
		NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
	)
	ct := "application/json"

	t.Run("should return metadata", func(t *testing.T) {

		ts := httptest.NewServer(NewHTTPTransport(s, log.NewNopLogger()).MakeMetadataHandler())
		defer ts.Close()

		r, err := contract.MarshalMetadataRequest(contract.MetadataRequest{
			MAC: []byte("mac"),
		})
		if err != nil {
			t.Errorf("contract.MarshalMetadataRequest() error = %v", err)
		}

		_, err = http.Post(ts.URL+"/v1/metadata", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeAddHandler(t *testing.T) {
	s := New(
		NewUserRepository(),
		NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
	)
	ct := "application/json"

	t.Run("should add vault", func(t *testing.T) {

		ts := httptest.NewServer(NewHTTPTransport(s, log.NewNopLogger()).MakeAddHandler())
		defer ts.Close()

		r, err := contract.MarshalAddRequest(contract.AddRequest{
			MAC:    []byte("mac"),
			Domain: "domain",
		})
		if err != nil {
			t.Errorf("contract.MarshalAddRequest() error = %v", err)
		}

		_, err = http.Post(ts.URL+"/v1/add", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeGetHandler(t *testing.T) {
	s := New(
		NewUserRepository(),
		NewConfiguration(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), sha256.New),
	)
	ct := "application/json"

	t.Run("should add vault", func(t *testing.T) {

		ts := httptest.NewServer(NewHTTPTransport(s, log.NewNopLogger()).MakeGetHandler())
		defer ts.Close()

		r, err := contract.MarshalGetRequest(contract.GetRequest{
			MAC:    []byte("mac"),
			Domain: "domain",
			BMK:    big.NewInt(1),
			Q:      big.NewInt(2),
		})
		if err != nil {
			t.Errorf("contract.MarshalGetRequest() error = %v", err)
		}

		_, err = http.Post(ts.URL+"/v1/get", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}
