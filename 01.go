package main

import (
	"fmt"
	"io/ioutil"
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

func mass(n int) int {
	return n/3 - 2
}

func totalmass(n int) int {
	t := 0
	x := n
	for {
		x = mass(x)
		if x < 0 {
			break
		}
		t += x
	}
	return t
}

const part2 = true

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("01.txt")
	must(err)
	var r int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !part2 {
			r += mass(atoi(line))
		} else {
			r += totalmass(atoi(line))
		}
	}
	fmt.Printf("result: %d\n", r)
}
