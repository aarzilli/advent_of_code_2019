package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// returns x without the last character
func nolast(x string) string {
	return x[:len(x)-1]
}

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
}

// convert string to integer
func atoi(in string) int {
	n, err := strconv.Atoi(in)
	must(err)
	return n
}

// convert vector of strings to integer
func vatoi(in []string) []int {
	r := make([]int, len(in))
	for i := range in {
		var err error
		r[i], err = strconv.Atoi(in[i])
		must(err)
	}
	return r
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func exit(n int) {
	os.Exit(n)
}

var deck []int
var deck2 []int

const dealwithincrement = "deal with increment"
const cut = "cut"

const part2 = true

var memo = map[int]int{}

func findidx(endidx int, N int, instr []string) int {
	if r, ok := memo[endidx]; ok {
		fmt.Printf("memoized %d\n", endidx)
		return r
	}

	idx := endidx

	for i := len(instr) - 1; i >= 0; i-- {
		line := instr[i]

		switch {
		case line == "deal into new stack":
			oldidx := idx
			idx = (N + 1) - idx - 1
			if idx >= (N + 1) {
				fmt.Printf("%d\n", oldidx)
				panic("cockup")
			}

		case strings.HasPrefix(line, dealwithincrement):
			n := atoi(strings.TrimSpace(line[len(dealwithincrement):]))
			oldidx := idx
			// faster way?
			for k := 0; k < n; k++ {
				if (idx+k*(N+1))%n == 0 {
					idx = (idx + k*(N+1)) / n
					break
				}
			}
			if (idx*n)%(N+1) != oldidx {
				fmt.Printf("newold=%d old=%d\n", (idx*n)%(N+1), oldidx%(N+1))
				panic("wtf")
			}

		case strings.HasPrefix(line, cut):
			n := atoi(strings.TrimSpace(line[len(cut):]))
			idx = idx - (N + 1 - n)
			if idx < 0 {
				for idx < 0 {
					idx = (N + 1) + idx
				}
			} else if idx >= (N + 1) {
				idx = idx % (N + 1)
			}
			if idx >= (N+1) || idx < 0 {
				fmt.Printf("blah %d\n", idx)
				panic("cockup")
			}
		}
	}
	memo[endidx] = idx
	return idx
}

func main() {
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
		const M = 1
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
			for i := range deck {
				if deck[i] == 2019 {
					fmt.Printf("PART 1: %d\n", i)
					break
				}
			}
		}

		endidx := 2020
		if endidx > len(deck) {
			endidx = 0
		}

		idx := endidx
		for cnt := 0; cnt < M; cnt++ {
			idx = findidx(idx, N, instr)
		}

		fmt.Printf("%d %d\n", idx, deck[endidx])
	} else {
		idx := 2020
		for cnt := 0; cnt < 101741582076661; cnt++ {
			if cnt%10000 == 0 {
				fmt.Printf("at %d %g%%\n", cnt, (float64(cnt)/101741582076661)*100.0)
			}
			oldidx := idx
			const debug = true
			if debug {
				fmt.Printf("%f -> ", float64(idx)/1000000000)
			}
			idx = findidx(idx, N, instr)
			if debug {
				fmt.Printf("%f (%d)\n", float64(idx)/1000000000, idx-oldidx)
			}
		}
		fmt.Printf("%d\n", idx)
	}

}

// 101741582076661 (rep)
// 119315717514047 (deck)
