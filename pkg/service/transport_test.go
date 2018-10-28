package service

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
)

func TestMakeRegisterHandler(t *testing.T) {
	s := New(NewUserRepository(), Configuration{
		sID:  big.NewInt(1),
		k:    big.NewInt(1),
		q0:   big.NewInt(1),
		bits: big.NewInt(1),
		hash: sha256.New,
	})
	ct := "application/json"

	t.Run("should register a user", func(t *testing.T) {

		ts := httptest.NewServer(MakeRegisterHandler(s))
		defer ts.Close()

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		contract.MarshalRegisterRequest(w, contract.RegisterRequest{CID: big.NewInt(1)})
		w.Flush()
		r := bufio.NewReader(&buf)

		_, err := http.Post(ts.URL+"/v1/register", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})

}

func TestMakeExpKHandler(t *testing.T) {
	s := New(NewUserRepository(), Configuration{
		sID:  big.NewInt(1),
		k:    big.NewInt(1),
		q0:   big.NewInt(1),
		bits: big.NewInt(1),
		hash: sha256.New,
	})
	ct := "application/json"

	t.Run("should exponent to k given blinded secret", func(t *testing.T) {

		ts := httptest.NewServer(MakeExpKHandler(s))
		defer ts.Close()

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		contract.MarshalExpKRequest(w, contract.ExpKRequest{
			CID:    big.NewInt(1),
			CNonce: big.NewInt(2),
			B:      big.NewInt(3),
			Q:      big.NewInt(4),
		})
		w.Flush()
		r := bufio.NewReader(&buf)

		_, err := http.Post(ts.URL+"/v1/login/expk", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeChallengeHandler(t *testing.T) {
	s := New(NewUserRepository(), Configuration{
		sID:  big.NewInt(1),
		k:    big.NewInt(1),
		q0:   big.NewInt(1),
		bits: big.NewInt(1),
		hash: sha256.New,
	})
	ct := "application/json"

	t.Run("should respond to a given challenge", func(t *testing.T) {

		ts := httptest.NewServer(MakeChallengeHandler(s))
		defer ts.Close()

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		contract.MarshalChallengeRequest(w, contract.ChallengeRequest{
			G: big.NewInt(3),
			Q: big.NewInt(4),
		})
		w.Flush()
		r := bufio.NewReader(&buf)

		_, err := http.Post(ts.URL+"/v1/login/challenge", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeMetadataHandler(t *testing.T) {
	s := New(NewUserRepository(), Configuration{
		sID:  big.NewInt(1),
		k:    big.NewInt(1),
		q0:   big.NewInt(1),
		bits: big.NewInt(1),
		hash: sha256.New,
	})
	ct := "application/json"

	t.Run("should return metadata", func(t *testing.T) {

		ts := httptest.NewServer(MakeMetadataHandler(s))
		defer ts.Close()

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		contract.MarshalMetadataRequest(w, contract.MetadataRequest{
			MAC: []byte("mac"),
		})
		w.Flush()
		r := bufio.NewReader(&buf)

		_, err := http.Post(ts.URL+"/v1/metadata", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeAddHandler(t *testing.T) {
	s := New(NewUserRepository(), Configuration{
		sID:  big.NewInt(1),
		k:    big.NewInt(1),
		q0:   big.NewInt(1),
		bits: big.NewInt(1),
		hash: sha256.New,
	})
	ct := "application/json"

	t.Run("should add vault", func(t *testing.T) {

		ts := httptest.NewServer(MakeAddHandler(s))
		defer ts.Close()

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		contract.MarshalAddRequest(w, contract.AddRequest{
			MAC:    []byte("mac"),
			Domain: "domain",
		})
		w.Flush()
		r := bufio.NewReader(&buf)

		_, err := http.Post(ts.URL+"/v1/add", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}

func TestMakeGetHandler(t *testing.T) {
	s := New(NewUserRepository(), Configuration{
		sID:  big.NewInt(1),
		k:    big.NewInt(1),
		q0:   big.NewInt(1),
		bits: big.NewInt(1),
		hash: sha256.New,
	})
	ct := "application/json"

	t.Run("should add vault", func(t *testing.T) {

		ts := httptest.NewServer(MakeGetHandler(s))
		defer ts.Close()

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		contract.MarshalGetRequest(w, contract.GetRequest{
			MAC:    []byte("mac"),
			Domain: "domain",
			BMK:    big.NewInt(1),
			Q:      big.NewInt(2),
		})
		w.Flush()
		r := bufio.NewReader(&buf)

		_, err := http.Post(ts.URL+"/v1/get", ct, r)
		if err != nil {
			t.Errorf("http.Post() error = %v", err)
		}
	})
}
