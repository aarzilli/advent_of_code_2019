package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// convert string to integer
func atoi(in string) int {
	n, err := strconv.Atoi(in)
	must(err)
	return n
}

const dealwithincrement = "deal with increment"
const cut = "cut"

const part2 = true
const debug = false

func simplify(instr []string, Len int) (*big.Int, *big.Int) {
	a := big.NewInt(1)
	b := big.NewInt(0)
	for _, line := range instr {
		switch {
		case line == "deal into new stack":
			a.Mul(a, big.NewInt(-1))
			b.Mul(b, big.NewInt(-1))
			b.Add(b, big.NewInt(int64(Len-1)))

		case strings.HasPrefix(line, dealwithincrement):
			n := atoi(strings.TrimSpace(line[len(dealwithincrement):]))
			a.Mul(a, big.NewInt(int64(n)))
			b.Mul(b, big.NewInt(int64(n)))

		case strings.HasPrefix(line, cut):
			n := atoi(strings.TrimSpace(line[len(cut):]))
			b.Add(b, big.NewInt(int64(Len-n)))

		default:
			fmt.Printf("unknown instruction %q\n", line)
			panic("blah")
		}
	}
	return a, b
}

func pow2(a, b *big.Int, M int, Len int) (*big.Int, *big.Int) {
	bigLen := big.NewInt(int64(Len))
	r := big.NewInt(1)
	r.Sub(r, a) // r = 1 - a

	if r.Exp(r, big.NewInt(-1), bigLen) == nil { // r = 1/(1-a)
		panic("wtf")
	}

	a.Exp(a, big.NewInt(int64(M)), bigLen) // a = a**M

	q := big.NewInt(1)
	q.Sub(q, a) // q = 1 - a**M

	q.Mul(q, r) // q = (1 - a**M) / (1 - a)
	q.Mod(q, bigLen)

	b.Mul(b, q) // b = b * (1 - a**M) / (1 - a)
	b.Mod(b, bigLen)
	a.Mod(a, bigLen)

	return a, b
}

func solve(a, b *big.Int, Len int) *big.Int {
	// find x such that (a*x + b) mod Len == 2020
	bigLen := big.NewInt(int64(Len))
	y := big.NewInt(2020)
	y.Sub(y, b) // a*x mod Len == (2020 - b) % mod Len
	y.Mod(y, bigLen)

	var x big.Int
	x.ModInverse(a, bigLen) // finds x such sthat a*x mod Len == 1
	x.Mul(&x, y)            // multiplying previous x by ((2020 - b) % mod Len) gives us the solution
	x.Mod(&x, bigLen)
	return &x
}

func main() {
	var deck []int
	var deck2 []int

	in := "22.txt"

	N := 10006

	if !part2 {
		if strings.Contains(in, "example") {
			N = 9
		}

		for i := 0; i <= N; i++ {
			deck = append(deck, i)
		}

		deck2 = make([]int, len(deck))
	} else {
		N = 119315717514046
	}

	instr := []string{}

	buf, err := ioutil.ReadFile(in)
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		instr = append(instr, line)

	}

	if !part2 {
		const M = 5

		deck3 := make([]int, len(deck))

		var a, b *big.Int
		{
			a, b = simplify(instr, N+1)
			a, b = pow2(a, b, M, N+1)

			for i := 0; i < len(deck); i++ {
				var r big.Int
				r.Mul(a, big.NewInt(int64(i)))
				r.Add(&r, b)
				r.Mod(&r, big.NewInt(int64(len(deck))))
				j := int(r.Int64())
				if j < 0 {
					j = len(deck) + j
				}
				deck3[j] = deck[i]
			}
		}

		for cnt := 0; cnt < M; cnt++ {
			if cnt%100 == 0 {
				fmt.Printf("at %d\n", cnt)
			}
			for _, line := range instr {
				switch {
				case line == "deal into new stack":
					for i, j := 0, len(deck)-1; i < j; i, j = i+1, j-1 {
						deck[i], deck[j] = deck[j], deck[i]
					}
				case strings.HasPrefix(line, dealwithincrement):
					n := atoi(strings.TrimSpace(line[len(dealwithincrement):]))

					for i, j := 0, 0; i < len(deck); i, j = i+1, (j+n)%len(deck) {
						deck2[j] = deck[i]
					}
					deck, deck2 = deck2, deck

				case strings.HasPrefix(line, cut):
					n := atoi(strings.TrimSpace(line[len(cut):]))
					if n > 0 {
						for i := n; i < len(deck); i++ {
							deck2[i-n] = deck[i]
						}
						for i := 0; i < n; i++ {
							deck2[(len(deck)-n)+i] = deck[i]
						}
					} else {
						for i, j := len(deck)+n, 0; i < len(deck); i, j = i+1, j+1 {
							deck2[j] = deck[i]
						}
						for i := 0; i < len(deck)+n; i++ {
							deck2[i-n] = deck[i]
						}
					}
					deck, deck2 = deck2, deck

				default:
					fmt.Printf("unknown instruction %q\n", line)
					panic("blah")
				}
			}
		}

		if N < 20 {
			fmt.Printf("Result %v\n", deck)
		} else {
			if debug {
				fmt.Printf("normal: %v\n", deck[:10])
				fmt.Printf("calculated: %v\n", deck3[:10])
			}

			for i := range deck {
				if deck[i] == 2019 {
					fmt.Printf("PART 1: %d\n", i)
					break
				}
			}

			if debug {
				x := solve(a, b, N+1)
				fmt.Printf("%v %d\n", x, deck[2020])
			}
		}
	} else {
		a, b := simplify(instr, N+1)
		a, b = pow2(a, b, 101741582076661, N+1)

		x := solve(a, b, N+1)
		fmt.Printf("PART 2: %d\n", x)
	}
}

// 101741582076661 (rep)
// 119315717514047 (deck)
