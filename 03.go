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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func dist(p0, p1 Point) int {
	return abs(p0.i-p1.i) + abs(p0.j-p1.j)
}

type Point struct {
	i, j int
}

var wirepaths = [][]string{}

var M = map[Point]int{}
var D = []map[Point]int{}
var coll = []Point{}

func exec(instr string, cur *Point, steps *int, typ int) {
	var otyp int
	if typ == 1 {
		otyp = 2
	} else {
		otyp = 1
	}
	n := atoi(instr[1:])
	for k := 0; k < n; k++ {
		switch instr[0] {
		case 'U':
			cur.i--
		case 'D':
			cur.i++
		case 'L':
			cur.j--
		case 'R':
			cur.j++
		}
		if M[*cur] == otyp {
			coll = append(coll, *cur)
		}
		*steps++
		if _, ok := D[typ-1][*cur]; !ok {
			D[typ-1][*cur] = *steps
		}
		M[*cur] = typ
	}
}

func main() {
	buf, err := ioutil.ReadFile("03.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		wirepath := splitandclean(line, ",", -1)
		wirepaths = append(wirepaths, wirepath)
	}

	M[Point{0, 0}] = 2
	D = make([]map[Point]int, 2)
	D[0] = make(map[Point]int)
	D[1] = make(map[Point]int)

	for i, wirepath := range wirepaths {
		if i >= 2 {
			panic("what")
		}
		cur := Point{0, 0}
		steps := 0
		for _, instr := range wirepath {
			exec(instr, &cur, &steps, i+1)
		}
	}

	// Part 1
	mindist := -1
	for _, p := range coll {
		d := dist(p, Point{0, 0})
		if mindist < 0 || d < mindist {
			mindist = d
		}
	}
	fmt.Printf("PART 1: %d\n", mindist)

	// Part 2
	minsteps := -1
	for _, p := range coll {
		steps := D[0][p] + D[1][p]
		if minsteps < 0 || steps < minsteps {
			minsteps = steps
		}
	}
	fmt.Printf("PART 2: %d\n", minsteps)
}
