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

type Node struct {
	name  string
	child []*Node
}

var M = map[string]*Node{}

func lookup(name string) *Node {
	if n, ok := M[name]; ok {
		return n
	}
	n := &Node{name: name}
	M[name] = n
	return n
}

var part1cnt int

func part1(depth int, n *Node) {
	part1cnt += depth
	for _, n2 := range n.child {
		part1(depth+1, n2)
	}
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("06.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := splitandclean(line, ")", -1)
		n1 := lookup(v[0])
		n1.child = append(n1.child, lookup(v[1]))
	}

	part1(0, lookup("COM"))
	fmt.Printf("PART 1: %d\n", part1cnt)
}
