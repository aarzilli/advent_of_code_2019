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

func getints(in string, hasneg bool) []int {
	v := getnums(in, hasneg, false)
	return vatoi(v)
}

func getnums(in string, hasneg, hasdot bool) []string {
	r := []string{}
	start := -1

	flush := func(end int) {
		if start < 0 {
			return
		}
		hasdigit := false
		for i := start; i < end; i++ {
			if in[i] >= '0' && in[i] <= '9' {
				hasdigit = true
				break
			}
		}
		if hasdigit {
			r = append(r, in[start:end])
		}
		start = -1
	}

	for i, ch := range in {
		isnumch := false

		switch {
		case hasneg && (ch == '-'):
			isnumch = true
		case hasdot && (ch == '.'):
			isnumch = true
		case ch >= '0' && ch <= '9':
			isnumch = true
		}

		if start >= 0 {
			if !isnumch {
				flush(i)
			}
		} else {
			if isnumch {
				start = i
			}
		}
	}
	flush(len(in))
	return r
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Point struct {
	x, y, z int
}

type Moon struct {
	pos Point
	v   Point
}

var points []Moon

func gravity1(a, b int) int {
	if a > b {
		return -1
	} else if a < b {
		return +1
	}
	return 0
}

func gravity() {
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			gx := gravity1(points[i].pos.x, points[j].pos.x)
			points[i].v.x += gx
			points[j].v.x -= gx

			gy := gravity1(points[i].pos.y, points[j].pos.y)
			points[i].v.y += gy
			points[j].v.y -= gy

			gz := gravity1(points[i].pos.z, points[j].pos.z)
			points[i].v.z += gz
			points[j].v.z -= gz
		}
	}
}

func velocity() {
	for i := range points {
		points[i].pos.x += points[i].v.x
		points[i].pos.y += points[i].v.y
		points[i].pos.z += points[i].v.z
	}
}

func step() {
	gravity()
	velocity()
}

func energy() int {
	e := 0
	for i := range points {
		pot := abs(points[i].pos.x) + abs(points[i].pos.y) + abs(points[i].pos.z)
		kin := abs(points[i].v.x) + abs(points[i].v.y) + abs(points[i].v.z)
		e += pot * kin
	}
	return e
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("12.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := getints(line, true)
		points = append(points, Moon{pos: Point{x: v[0], y: v[1], z: v[2]}})
	}
	fmt.Printf("%v\n", points)

	N := 1000

	for i := 0; i < N; i++ {
		step()
	}

	fmt.Printf("%v %d\n", points, energy())
}
