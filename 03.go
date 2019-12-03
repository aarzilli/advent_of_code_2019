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
	set := func() {
		if M[*cur] == otyp {
			coll = append(coll, *cur)
		}
		*steps++
		if _, ok := D[typ-1][*cur]; !ok {
			D[typ-1][*cur] = *steps
		}
		M[*cur] = typ
	}
	n := atoi(instr[1:])
	switch instr[0] {
	case 'U':
		for k := 0; k < n; k++ {
			cur.i--
			set()
		}
	case 'D':
		for k := 0; k < n; k++ {
			cur.i++
			set()
		}
	case 'L':
		for k := 0; k < n; k++ {
			cur.j--
			set()
		}
	case 'R':
		for k := 0; k < n; k++ {
			cur.j++
			set()
		}
	}
}

const debug = false
const part1 = false

func main() {
	fmt.Printf("hello\n")
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

	if debug {
		for i := -7; i < 2; i++ {
			for j := -1; j < 10; j++ {
				if M[Point{i, j}] != 0 {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
	}

	if part1 {
		mindist := -1
		for _, p := range coll {
			d := dist(p, Point{0, 0})
			if mindist < 0 || d < mindist {
				mindist = d
			}
		}
		fmt.Printf("PART 1: %d\n", mindist)
	} else {
		minsteps := -1
		for _, p := range coll {
			steps := D[0][p] + D[1][p]
			if minsteps < 0 || steps < minsteps {
				minsteps = steps
			}
		}
		fmt.Printf("PART2: %d\n", minsteps)
	}
}
