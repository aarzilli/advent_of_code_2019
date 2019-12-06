package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
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

func find(path []string, dest string, cur *Node) []string {
	if cur.name == dest {
		return path
	}
	for _, n2 := range cur.child {
		r := find(append(path, n2.name), dest, n2)
		if r != nil {
			return r
		}
	}
	return nil
}

const debug = false

func main() {
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

	path2you := find(make([]string, 0, 100), "YOU", lookup("COM"))
	path2santa := find(make([]string, 0, 100), "SAN", lookup("COM"))
	if debug {
		fmt.Printf("%#v\n", path2you)
		fmt.Printf("%#v\n", path2santa)
	}

	commonanc := ""
	dist2you := 0
	dist2santa := 0

	for i := range path2you {
		if path2you[i] != path2santa[i] {
			commonanc = path2you[i-1]
			dist2you = len(path2you) - i
			dist2santa = len(path2santa) - i
			break
		}
	}

	if debug {
		fmt.Printf("%s %d %d\n", commonanc, dist2you, dist2santa)
	}
	fmt.Printf("PART 2: %d\n", dist2you+dist2santa-2)
}
