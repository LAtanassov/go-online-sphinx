package client_test

import (
	"crypto/sha256"
	"hash"
	"net/http"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"
)

func BenchmarkClient_Register_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Register(b, 2048, sha256.New)
}
func BenchmarkClient_Register_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Register(b, 2048, sha256.New)
}

func BenchmarkClient_Register_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Register(b, 3072, sha256.New)
}

func benchmarkClient_Register(b *testing.B, bits int, hash func() hash.Hash) {

	baseURL := "http://localhost:8080"

	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err := clt.Register("registered-user")
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkClient_Login_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Login(b, 1024, sha256.New)
}

func BenchmarkClient_Login_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Login(b, 2048, sha256.New)
}

func BenchmarkClient_Login_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Login(b, 3072, sha256.New)
}

func benchmarkClient_Login(b *testing.B, bits int, hash func() hash.Hash) {

	baseURL := "http://localhost:8080"

	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	err := clt.Register("login-user")
	if err != nil {
		b.Errorf("Register() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		clt.Login("login-user", "password")
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkClient_Add_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Add(b, 1024, sha256.New)
}

func BenchmarkClient_Add_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Add(b, 2048, sha256.New)
}

func BenchmarkClient_Add_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Add(b, 3072, sha256.New)
}

func benchmarkClient_Add(b *testing.B, bits int, hash func() hash.Hash) {

	baseURL := "http://localhost:8080"

	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	err := clt.Register("add-domain-user")
	if err != nil {
		b.Errorf("Register() error = %v", err)
	}

	clt.Login("add-domain-user", "password")
	if err != nil {
		b.Errorf("Login() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = clt.Add(string(n))
		if err != nil {
			b.Errorf("Add() error = %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkClient_Get_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Get(b, 1024, sha256.New)
}

func BenchmarkClient_Get_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Get(b, 2048, sha256.New)
}

func BenchmarkClient_Get_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Get(b, 3072, sha256.New)
}

var pwd string

func benchmarkClient_Get(b *testing.B, bits int, hash func() hash.Hash) {

	baseURL := "http://localhost:8080"

	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	err := clt.Register("get-domain-user")
	if err != nil {
		b.Errorf("Register() error = %v", err)
	}

	err = clt.Login("get-domain-user", "password")
	if err != nil {
		b.Errorf("Login() error = %v", err)
	}

	err = clt.Add("new-domain")
	if err != nil {
		b.Errorf("Add() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		gpwd, err := clt.Get("new-domain")
		if err != nil {
			b.Errorf("Add() error = %v", err)
		}
		pwd = gpwd
	}

	b.StopTimer()
}
