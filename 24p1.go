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

var M = [][]byte{}
var M2 = [][]byte{}

func showmatrix() {
	for i := range M {
		for j := range M[i] {
			fmt.Printf("%c", M[i][j])
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
}

func at(i, j int) int {
	if i < 0 || i >= len(M) {
		return 0
	}
	if j < 0 || j >= len(M[i]) {
		return 0
	}
	if M[i][j] == '#' {
		return 1
	}
	return 0
}

func step() {
	for i := range M {
		for j := range M[i] {
			n := at(i-1, j) + at(i, j-1) + at(i, j+1) + at(i+1, j)
			if M[i][j] == '#' {
				if n != 1 {
					M2[i][j] = '.'
				} else {
					M2[i][j] = '#'
				}
			} else {
				switch n {
				case 1, 2:
					M2[i][j] = '#'
				default:
					M2[i][j] = '.'
				}
			}
		}
	}

	M, M2 = M2, M
}

var buf []byte

func stringify() string {
	buf = buf[:0]
	for i := range M {
		buf = append(buf, M[i]...)
	}
	return string(buf)
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("24.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	M2 = make([][]byte, len(M))
	for i := range M2 {
		M2[i] = make([]byte, len(M[i]))
	}

	m := map[string]bool{}

	for cnt := 0; cnt < 100; cnt++ {
		if cnt%1000 == 0 {
			fmt.Printf("at %d\n", cnt)
		}
		s := stringify()
		if m[s] {
			break
		}
		m[s] = true
		step()
	}

	showmatrix()

	biod := 0
	p := 1

	for i := range M {
		for j := range M[i] {
			if M[i][j] == '#' {
				biod += p
			}
			p *= 2
		}
	}

	fmt.Printf("PART 1: %d\n", biod)
}
