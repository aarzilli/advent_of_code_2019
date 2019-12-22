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

func main() {
	in := "22.txt"

	N := 10006

	if strings.Contains(in, "example") {
		N = 9
	}

	for i := 0; i <= N; i++ {
		deck = append(deck, i)
	}

	deck2 = make([]int, len(deck))

	buf, err := ioutil.ReadFile(in)
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

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
}
