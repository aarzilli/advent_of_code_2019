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

const WIDTH = 25
const HEIGHT = 6

func count(layer []byte, v byte) int {
	cnt := 0
	for i := range layer {
		if layer[i] == v {
			cnt++
		}
	}
	return cnt
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("08.txt")
	must(err)
	var in []byte
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		in = []byte(line)
	}
	for i := range in {
		in[i] = in[i] - '0'
	}
	layers := [][]byte{}
	sz := WIDTH * HEIGHT
	rem := in
	for len(rem) > 0 {
		layers = append(layers, rem[:sz])
		rem = rem[sz:]
	}

	minnum0 := 10000
	minnum0i := 0
	for i, layer := range layers {
		num0 := count(layer, 0)
		if num0 < minnum0 {
			minnum0 = num0
			minnum0i = i
		}
	}

	fmt.Printf("%d\n", minnum0i)
	fmt.Printf("PART 1: %d\n", count(layers[minnum0i], 1)*count(layers[minnum0i], 2))
}
