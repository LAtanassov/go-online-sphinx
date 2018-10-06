package osphinx

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"log"
	"math/big"
	"reflect"
	"testing"
)

func TestExpInGroup(t *testing.T) {
	type args struct {
		g *big.Int
		k *big.Int
		q *big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "should calc 3**8 mod 27 = 26",
			args: args{
				g: big.NewInt(8),
				k: big.NewInt(9),
				q: big.NewInt(13),
			},
			want: big.NewInt(26),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpInGroup(tt.args.g, tt.args.k, tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpInGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func fordKaliskiPasswordToRandom(h hash.Hash, bits int) *big.Int {
	pwd := big.NewInt(0).SetBytes(h.Sum([]byte("Ford Kaliski Password to Random")))

	q, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}

	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		log.Fatal(err)
	}

	d, err := rand.Int(rand.Reader, q)
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

func benchmarkFordKaliskiPasswordToRandom(b *testing.B, h hash.Hash, bits int) {
	var r *big.Int
	for n := 0; n < b.N; n++ {
		r = fordKaliskiPasswordToRandom(h, bits)
	}
	result = r
}

func BenchmarkFordKaliskiPasswordToRandom_SHA256_8Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha256.New(), 8)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA256_16Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha256.New(), 16)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA256_32Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha256.New(), 32)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA256_64Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha256.New(), 64)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA256_128Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha256.New(), 128)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA256_256Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha256.New(), 256)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA512_512Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha512.New(), 512)
}

func BenchmarkFordKaliskiPasswordToRandom_SHA512_1024Bits(b *testing.B) {
	benchmarkFordKaliskiPasswordToRandom(b, sha512.New(), 1024)
}
