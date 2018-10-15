package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"log"
	"math/big"
	"testing"
)

func BenchmarkLogin_SHA256_8Bits(b *testing.B) {
	benchmarkLogin(b, sha256.New(), 8)
}

func BenchmarkLogin_SHA256_16Bits(b *testing.B) {
	benchmarkLogin(b, sha256.New(), 16)
}

func BenchmarkLogin_SHA256_32Bits(b *testing.B) {
	benchmarkLogin(b, sha256.New(), 32)
}

func BenchmarkLogin_SHA256_64Bits(b *testing.B) {
	benchmarkLogin(b, sha256.New(), 64)
}

func BenchmarkLogin_SHA256_128Bits(b *testing.B) {
	benchmarkLogin(b, sha256.New(), 128)
}

func BenchmarkLogin_SHA256_256Bits(b *testing.B) {
	benchmarkLogin(b, sha256.New(), 256)
}

func BenchmarkLogin_SHA512_512Bits(b *testing.B) {
	benchmarkLogin(b, sha512.New(), 512)
}

func BenchmarkLogin_SHA512_1024Bits(b *testing.B) {
	benchmarkLogin(b, sha512.New(), 1024)
}

func login(h hash.Hash, bits int) *big.Int {
	pwd := big.NewInt(0).SetBytes(h.Sum([]byte("Ford Kaliski Password to Random")))

	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	q, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}

	k, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}

	d, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}

	kinv := big.NewInt(0).ModInverse(k, q)
	if kinv == nil {
		kinv = big.NewInt(0)
	}

	g := ExpInGroup(pwd, two, q)

	// blinding
	b := ExpInGroup(g, k, q)

	// exp with secret
	bd := ExpInGroup(b, d, q)

	// unblinding
	r := ExpInGroup(bd, kinv, q)
	return r
}

var result *big.Int

func benchmarkLogin(b *testing.B, h hash.Hash, bits int) {
	var r *big.Int
	for n := 0; n < b.N; n++ {
		r = login(h, bits)
	}
	result = r
}
