package main

import (
	"fmt"
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

const input = "172930-683082"

func valid(x string) bool {
	doublefound := false
	for i := range x {
		if i+1 < len(x) {
			if x[i] > x[i+1] {
				return false
			}
			if x[i] == x[i+1] {
				doublefound = true
			}
		}
	}

	return doublefound
}

func main() {
	for _, n := range []int{111111, 223450, 123789} {
		fmt.Printf("%d %v\n", n, valid(fmt.Sprintf("%d", n)))
	}

	cnt := 0
	for i := 172930; i <= 683082; i++ {
		x := fmt.Sprintf("%d", i)
		if valid(x) {
			//fmt.Printf("PART 1: %d\n", i)
			cnt++
		}
	}
	fmt.Printf("PART 1: %d\n", cnt)
}
