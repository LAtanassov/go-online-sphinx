package crypto

import (
	"crypto/sha256"
	"hash"
	"math/big"
	"reflect"
	"testing"
)

func TestCrypto_ExpInGroup(t *testing.T) {
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

func TestCrypto_HashInGroup(t *testing.T) {
	type args struct {
		password string
		newHash  func() hash.Hash
		q        *big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			"should hash in group 'password'",
			args{
				password: "password",
				newHash:  sha256.New,
				q:        big.NewInt(42),
			},
			big.NewInt(59),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashInGroup(tt.args.password, tt.args.newHash, tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HashInGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHmacData(t *testing.T) {
	type args struct {
		h    func() hash.Hash
		key  []byte
		data [][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []byte
	}{
		{"should compare HMAC's", args{sha256.New, []byte("key"), [][]byte{}}, true, []byte("another-hash")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HmacData(tt.args.h, tt.args.key, tt.args.data...)

			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("HmacData() = %v, want %v", got, tt.want)
			}
		})
	}
}
