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

func exit(n int) {
	os.Exit(n)
}

type Comp struct {
	n    int
	kind string
}

type Comps struct {
	outn int
	c    []Comp
}

var components = map[string]Comps{}

func parsecomp(in string) Comp {
	v := splitandclean(in, " ", 2)
	return Comp{
		n:    atoi(v[0]),
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
	if k*comps.outn < quantity {
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

const part1 = true
const debug2 = false

func main() {
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
			c:    comps,
		}
	}

	if part1 {
		if debug {
			fmt.Printf("%v\n", components)
		}

		residual := map[string]int{}
		orecount := do("", "FUEL", 1, residual)
		fmt.Printf("PART 1: %d\n", orecount)
	}

	const ORETGT = 1000000000000

	fuel := ORETGT
	searchSpeed := ORETGT / 10
	lastGoodFuel := fuel

	for {
		residual := map[string]int{}
		orecount := do("", "FUEL", fuel, residual)
		if debug2 {
			fmt.Printf("FUEL %d ORE %d\n", fuel, orecount)
		}
		if orecount < ORETGT {
			if searchSpeed == 1 {
				fmt.Printf("PART 2: %d\n", fuel+1)
				exit(0)
			}
			searchSpeed = searchSpeed / 10
			fuel = lastGoodFuel
		} else {
			lastGoodFuel = fuel
			fuel -= searchSpeed
		}
	}
}
