package client_test

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"hash"
	"math"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"
)

func BenchmarkClient_Register_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Register(b, 32, sha256.New)
}

func BenchmarkClient_Register_SHA256_128Bits(b *testing.B) {
	benchmarkClient_Register(b, 128, sha256.New)
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

func benchmarkClient_Register(b *testing.B, bits int, hashfn func() hash.Hash) {

	clt, err := newBenchOscli(bits, hashfn)
	if err != nil {
		b.Errorf("creating oscli() error = %v", err)
	}

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
func BenchmarkClient_Login_SHA256_128Bits(b *testing.B) {
	benchmarkClient_Login(b, 128, sha256.New)
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

func benchmarkClient_Login(b *testing.B, bits int, hashfn func() hash.Hash) {

	clt, err := newBenchOscli(bits, hashfn)
	if err != nil {
		b.Errorf("creating oscli() error = %v", err)
	}

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
		clt.Logout()
	}

	b.StopTimer()
}

func BenchmarkClient_Add_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Add(b, 32, sha256.New)
}
func BenchmarkClient_Add_SHA256_128Bits(b *testing.B) {
	benchmarkClient_Add(b, 128, sha256.New)
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

func benchmarkClient_Add(b *testing.B, bits int, hashfn func() hash.Hash) {

	clt, err := newBenchOscli(bits, hashfn)
	if err != nil {
		b.Errorf("creating oscli() error = %v", err)
	}

	err = clt.Register("add-domain-user")
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

func BenchmarkClient_GetMetadata_SHA256_32Bits(b *testing.B) {
	benchmarkClient_GetMetadata(b, 32, sha256.New)
}
func BenchmarkClient_GetMetadata_SHA256_128Bits(b *testing.B) {
	benchmarkClient_GetMetadata(b, 128, sha256.New)
}
func BenchmarkClient_GetMetadata_SHA256_512Bits(b *testing.B) {
	benchmarkClient_GetMetadata(b, 512, sha256.New)
}

func BenchmarkClient_GetMetadata_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_GetMetadata(b, 1024, sha256.New)
}

func BenchmarkClient_GetMetadata_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_GetMetadata(b, 2048, sha256.New)
}

func BenchmarkClient_GetMetadata_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_GetMetadata(b, 3072, sha256.New)
}

var domains []string

func benchmarkClient_GetMetadata(b *testing.B, bits int, hashfn func() hash.Hash) {

	clt, err := newBenchOscli(bits, hashfn)
	if err != nil {
		b.Errorf("creating oscli() error = %v", err)
	}

	err = clt.Register("add-domain-user")
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

	err = clt.Add("google.com")
	if err != nil {
		b.Errorf("Add() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ds, err := clt.GetMetadata()
		if err != nil {
			b.Errorf("GetMetadata() error = %v", err)
		}
		domains = ds
	}

	b.StopTimer()
}

func BenchmarkClient_Get_SHA256_32Bits(b *testing.B) {
	benchmarkClient_Get(b, 32, sha256.New)
}
func BenchmarkClient_Get_SHA256_128Bits(b *testing.B) {
	benchmarkClient_Get(b, 128, sha256.New)
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

func benchmarkClient_Get(b *testing.B, bits int, hashfn func() hash.Hash) {

	clt, err := newBenchOscli(bits, hashfn)
	if err != nil {
		b.Errorf("creating oscli() error = %v", err)
	}

	err = clt.Register("get-domain-user")
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

func newBenchOscli(bits int, hashFn func() hash.Hash) (*client.Client, error) {
	var baseURL = "https://localhost"

	cfg, err := client.NewConfiguration(baseURL, bits, hashFn)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cli := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	repo := client.NewInMemoryUserRepository()
	return client.New(
		cli,
		cfg,
		repo,
	), nil
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
