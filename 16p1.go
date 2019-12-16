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

func newpattern(n int) []int {
	n++
	pattern := make([]int, 4*n)
	for k := 0; k < n; k++ {
		pattern[k] = 0
	}
	for k := 0; k < n; k++ {
		pattern[n+k] = 1
	}
	for k := 0; k < n; k++ {
		pattern[2*n+k] = 0
	}
	for k := 0; k < n; k++ {
		pattern[3*n+k] = -1
	}
	return pattern
}

const debug = false

func step(in []int) []int {
	N := len(in)
	out := make([]int, N)
	for i := 0; i < N; i++ {
		r := 0
		pattern := newpattern(i)
		for j := range in {
			if debug {
				fmt.Printf("%d*%d + ", in[j], pattern[(j+1)%len(pattern)])
			}
			r += in[j] * pattern[(j+1)%len(pattern)]
		}
		r = abs(r)
		r = r % 10

		if debug {
			fmt.Printf("= %d\n", r)
		}

		out[i] = r
	}
	return out
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("16.txt")
	must(err)
	var in []int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		in = make([]int, len(line))
		for i := range line {
			in[i] = int(line[i] - '0')
		}
	}

	fmt.Printf("%v\n", in)

	for i := 0; i < 100; i++ {
		out := step(in)
		in = out
	}

	fmt.Printf("PART 1: ")
	for i := 0; i < 8; i++ {
		fmt.Printf("%d", in[i])
	}
	fmt.Printf("\n")
}
