package crypto

import (
	"math/big"
)

var one = big.NewInt(1)
var two = big.NewInt(2)

// ExpInGroup - exponential in cyclic group returns g**k mod |2*q + 1|
func ExpInGroup(g, k, q *big.Int) *big.Int {
	var p = big.NewInt(0)
	var r = big.NewInt(0)

	p.Add(p.Mul(two, q), one)

	return r.Exp(g, k, p)
}
