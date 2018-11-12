package client_test

import (
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"math"
	"math/big"
	"net/http"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"
)

func BenchmarkClient_Register_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Register(b, 32, sha256.New)
}

func BenchmarkClient_Register_SHA256_512Bits(b *testing.B) {
	benchmarkClient_Register(b, 512, sha256.New)
}

func BenchmarkClient_Register_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Register(b, 1024, sha256.New)
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

	users := generator(b, b.N)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err := clt.Register(users[n])
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkClient_Login_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Login(b, 32, sha256.New)
}

func BenchmarkClient_Login_SHA256_512Bits(b *testing.B) {
	benchmarkClient_Login(b, 512, sha256.New)
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

	users := generator(b, b.N)

	for i := 0; i < b.N; i++ {
		err := clt.Register(users[i])
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err := clt.Login(users[n], "password")
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
		err = clt.Challenge()
		if err != nil {
			b.Errorf("Challenge() error = %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkClient_Add_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Add(b, 32, sha256.New)
}

func BenchmarkClient_Add_SHA256_512Bits(b *testing.B) {
	benchmarkClient_Add(b, 512, sha256.New)
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
		b.Errorf("Register() error = %+v", err)
	}

	err = clt.Login("add-domain-user", "password")
	if err != nil {
		b.Errorf("Login() error = %+v", err)
	}

	err = clt.Challenge()
	if err != nil {
		b.Errorf("Challenge() error = %+v", err)
	}

	domains := generator(b, b.N)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = clt.Add(domains[n])
		if err != nil {
			b.Errorf("Add() error = %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkClient_Get_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Get(b, 32, sha256.New)
}

func BenchmarkClient_Get_SHA256_512Bits(b *testing.B) {
	benchmarkClient_Get(b, 512, sha256.New)
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

	err = clt.Challenge()
	if err != nil {
		b.Errorf("Challenge() error = %v", err)
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

func generator(b *testing.B, x int) []string {
	sample := make([]string, x)
	for i := 0; i < x; i++ {
		j, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			b.Errorf("generator() err=%v", err)
		}
		sample[i] = j.Text(16)
	}
	return sample
}
