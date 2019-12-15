package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"os"
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

type Point struct {
	i, j int
}

func showmap() {
	for i := range M {
		fmt.Printf("%s\n", string(M[i]))
	}
	
}

func done() bool {
	for i := range M {
		for j := range M[i] {
			if M[i][j] == '.' {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("15.addenda")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		//line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		M = append(M, []byte(line))
	}
	
	showmap()
	
	fringe := make(map[Point]bool)
	
	for i := range M {
		for j := range M[i] {
			if M[i][j] == 'O' {
				fringe[Point{i,j}] = true
			}
		}
	}
	
	cnt := 0
	for {
		newfringe := make(map[Point]bool)
		
		expand := func(i, j int) {
			if i < 0 || i >= len(M) {
				return
			}
			if j < 0 || j >= len(M[i]) {
				return
			}
			if M[i][j] == '.' {
				newfringe[Point{i,j}] = true
				M[i][j] = 'O'
			}	
		}
		
		for p := range fringe {
			expand(p.i-1, p.j)
			expand(p.i, p.j-1)
			expand(p.i, p.j+1)
			expand(p.i+1, p.j)
		}
		cnt++
		
		fringe = newfringe
		showmap()
		fmt.Printf("\n")
		
		if done() {
			fmt.Printf("steps %d\n", cnt)
			break
		}
	}
}
