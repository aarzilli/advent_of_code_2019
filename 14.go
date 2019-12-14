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

const debug = false

func do(depth string, tgt string, quantity int, residual map[string]int) int {
	if residual[tgt] > 0 {
		if residual[tgt] > quantity {
			residual[tgt] -= quantity
			return 0
		} else {
			quantity -= residual[tgt]
			residual[tgt] = 0
		}
	}
	if tgt == "ORE" {
		return quantity
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
	
	r := 0
	for _, comp := range comps.c {
		if debug {
			fmt.Printf("%s\t%d %s\n", depth, k*comp.n, comp.kind)
		}
		r += do(depth+"\t", comp.kind, k*comp.n, residual)
	}
	
	
	if outq > quantity {
		residual[tgt] = outq - quantity
	}
	return r
}

const part1 = false
const debug2 = true

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
	
	if part1 {
		fmt.Printf("%v\n", components)
		
		residual := map[string]int{}
		orecount := do("", "FUEL", 1, residual)
		fmt.Printf("PART 1: %d\n", orecount)
		exit(0)
	}
	
	/*
	ore := 1000000000000
	residual := map[string]int{}
	
	factor := 100
	
	cnt := 0
	for ore > 0 {
		if ore < 10000000000 {
			factor = 1
		}
		if debug2 {
			fmt.Printf("%d\n", ore)
			//fmt.Printf("%v\n", residual)
		}
		orecount := do("", "FUEL", 1, residual)
		if ore - (orecount*factor) < 0 {
			if ore - orecount < 0 {
				break
			} else {
				ore -= orecount
				cnt++
			}
		} else {
			for k := range residual {
				residual[k] *= factor
			}
			ore -= orecount * factor
			cnt = cnt + factor
		}
	}
	
	fmt.Printf("PART 2: %d\n", cnt)
	
	fmt.Printf("residual %v\n", residual)
	*/
	
	start := 6226160
	factor := 1
	
	for x := start; x > 0; x -= factor {
		residual := map[string]int{}
		orecount := do("", "FUEL", x, residual)
		fmt.Printf("FUEL %d ORE %d\n", x, orecount)
		if orecount < 1000000000000 {
			break
		}
	}
}

// 6226142