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

const part2 = true

type periodFinder struct {
	name   string
	m      map[string]int
	period []string
	start  int
	found  bool
}

func newPeriodFinder(name string) *periodFinder {
	return &periodFinder{name: name, m: make(map[string]int)}
}

func (pf *periodFinder) findPeriod(i int, x string) {
	if pf.m[x] > 0 {
		//fmt.Printf("%s repetition found %d %s (period %d)\n", pf.name, i, x, i-pf.m[x])
		if pf.period == nil {
			pf.period = make([]string, i-pf.m[x])
			pf.start = pf.m[x]
			for x, j := range pf.m {
				if j < pf.start {
					continue
				}
				pf.period[j-pf.start] = x
			}
		}
		//fmt.Printf("\t%s\n", pf.period[(i-pf.start) % len(pf.period)])
		if pf.period[(i-pf.start)%len(pf.period)] != x {
			fmt.Printf("FAIL\n")
			os.Exit(1)
		}
		pf.found = true
	} else {
		pf.m[x] = i
	}
}

func (pf *periodFinder) get(i int) string {
	return pf.period[(i-pf.start)%len(pf.period)]
}

func (pf *periodFinder) mod(n int) int {
	return (n - pf.start) % len(pf.period)
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

	if !part2 {
		N := 1000

		for i := 0; i < N; i++ {
			step()
		}

		fmt.Printf("%v %d\n", points, energy())
	} else {
		pfx := newPeriodFinder("x")
		pfy := newPeriodFinder("y")
		pfz := newPeriodFinder("z")
		//my := map[string]int{}
		//mz := map[string]int{}

		var i = 0
		for ; true; i++ {
			step()

			x := fmt.Sprintf("%d %d %d %d %d %d %d %d",
				points[0].pos.x, points[0].v.x,
				points[1].pos.x, points[1].v.x,
				points[2].pos.x, points[2].v.x,
				points[3].pos.x, points[3].v.x)

			pfx.findPeriod(i, x)

			y := fmt.Sprintf("%d %d %d %d %d %d %d %d",
				points[0].pos.y, points[0].v.y,
				points[1].pos.y, points[1].v.y,
				points[2].pos.y, points[2].v.y,
				points[3].pos.y, points[3].v.y)

			pfy.findPeriod(i, y)

			z := fmt.Sprintf("%d %d %d %d %d %d %d %d",
				points[0].pos.z, points[0].v.z,
				points[1].pos.z, points[1].v.z,
				points[2].pos.z, points[2].v.z,
				points[3].pos.z, points[3].v.z)

			pfz.findPeriod(i, z)

			if pfx.found && pfy.found && pfz.found {
				i++
				break
			}
		}

		fmt.Printf("x start=%d len=%d y start=%d len=%d z start=%d len=%d\n", pfx.start, len(pfx.period), pfy.start, len(pfy.period), pfz.start, len(pfz.period))

		fmt.Printf("%d %d %d\n", pfx.mod(1), pfy.mod(1), pfz.mod(1))

		//const magic = (2 * 2 * 46507 * 115807 * 2 * 2 * 25589) + 1
		const magic = 551272644867044 + 1
		fmt.Printf("%d %d %d\n", pfx.mod(magic), pfy.mod(magic), pfz.mod(magic))
		fmt.Printf("magic: %d\n", magic)

		/*
			for ; true; i++ {
				step()

				x := fmt.Sprintf("%d %d %d %d %d %d %d %d",
					points[0].pos.x, points[0].v.x,
					points[1].pos.x, points[1].v.x,
					points[2].pos.x, points[2].v.x,
					points[3].pos.x, points[3].v.x)

				x2 := pfx.get(i)
				if x != x2 {
					fmt.Printf("ERROR %s | %s\n", x, x2)
					panic("blah")
				}

			}*/
	}
}

// x repetition found 186029 -4 -6 -7 3 -3 -2 -5 5 (period 186028)
// 2 x 2 x 46507
// y repetition found 231615 2 -1 -2 6 -2 -6 4 1 (period 231614)
// 2 x 115807
// z repetition found 102357 1 -2 6 6 3 2 1 -6 (period 102356)
// 2 x 2 x 25589

// 2205090579468177 too high
// 551272644867045 too high (off by one, I'm retarded) (it's the least common multiplier of the 3 periods)
