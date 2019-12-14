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

type Comp struct {
	n int
	kind string
}

type Comps struct {
	outn int
	c []Comp
}

var components = map[string]Comps{}

func parsecomp(in string) Comp {
	v := splitandclean(in, " ", 2)
	return Comp{
		n: atoi(v[0]),
		kind: v[1],
	}
}

const debug = true

var orecount int

func do(depth string, tgt string, quantity int, residual map[string]int) {
	if residual[tgt] > 0 {
		if residual[tgt] > quantity {
			residual[tgt] -= quantity
			return
		} else {
			quantity -= residual[tgt]
			residual[tgt] = 0
		}
	}
	if tgt == "ORE" {
		orecount += quantity
		return
	}
	comps := components[tgt]
	
	k := quantity / comps.outn
	if k * comps.outn < quantity {
		k++
	}
	
	outq := k * comps.outn
	
	if debug {
		fmt.Printf("%sConsume to produce %d %s\n", depth, outq, tgt)
	}
	for _, comp := range comps.c {
		if debug {
			fmt.Printf("%s\t%d %s\n", depth, k*comp.n, comp.kind)
		}
		do(depth+"\t", comp.kind, k*comp.n, residual)
	}
	
	
	if outq > quantity {
		residual[tgt] = outq - quantity
	}
	return
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("14.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := splitandclean(line, "=>", 2)
		out := parsecomp(v[1])
		v2 := splitandclean(v[0], ",", -1)
		comps := make([]Comp, len(v2))
		for i := range v2 {
			comps[i] = parsecomp(v2[i])
		}
		if _, ok := components[out.kind]; ok {
			panic("blah")
		}
		components[out.kind] = Comps{
			outn: out.n,
			c: comps,
		}
	}
	
	fmt.Printf("%v\n", components)
	
	do("", "FUEL", 1, map[string]int{})
	fmt.Printf("PART 1: %d\n", orecount)
}
