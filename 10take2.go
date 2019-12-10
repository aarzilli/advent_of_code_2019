package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var M = [][]byte{}

type Point struct {
	x, y int
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}

const debug = false

func listvisible(x0, y0 int) []Point {
	candidates := map[Point]bool{}
	for y := 0; y < len(M); y++ {
		for x := 0; x < len(M[y]); x++ {
			if M[y][x] != '.' {
				candidates[Point{x: x, y: y}] = true
			}
		}
	}

	if debug {
		fmt.Printf("point %d %d\n", x0, y0)
	}
	delete(candidates, Point{x: x0, y: y0})

	for p := range candidates {
		dx := p.x - x0
		dy := p.y - y0

		if debug {
			fmt.Printf("occluded by %v\n", p)
		}

		if dx == 0 {
			for q := range candidates {
				if q.x != x0 {
					continue
				}
				dy2 := q.y - y0
				if sign(dy2) == sign(dy) && abs(dy2) > abs(dy) {
					if debug {
						fmt.Printf("\tdeleting candidate %v (a)\n", q)
					}
					delete(candidates, q)
				}
			}
		} else if dy == 0 {
			for q := range candidates {
				if q.y != y0 {
					continue
				}
				dx2 := q.x - x0
				if sign(dx2) == sign(dx) && abs(dx2) > abs(dx) {
					if debug {
						fmt.Printf("\tdeleting candidate %v (b)\n", q)
					}
					delete(candidates, q)
				}
			}
		} else {
			m := float64(dy) / float64(dx)

			for q := range candidates {
				if p == q {
					continue
				}
				dx2 := q.x - x0
				dy2 := q.y - y0
				if dx2 == 0 {
					continue
				}
				m2 := float64(dy2) / float64(dx2)

				if debug {
					fmt.Printf("\t%v %g %g\n", q, m, m2)
				}
				if math.Abs(m-m2) < 0.0001 {
					if sign(dx2) == sign(dx) && abs(dx2) > abs(dx) {
						if debug {
							fmt.Printf("\tdeleting candidate %v (c)\n", q)
						}
						delete(candidates, q)
					}
				}
			}
		}
	}

	r := make([]Point, 0, len(candidates))

	for p := range candidates {
		r = append(r, p)
	}

	return r
}

func calcm(v []Point, z Point) []float64 {
	ms := make([]float64, 0, len(v))
	for _, p := range v {
		dx := p.x - z.x
		dy := p.y - z.y
		m := float64(dy) / float64(dx)
		ms = append(ms, m)
	}
	return ms
}

type Pointm struct {
	p Point
	m float64
}

const debug2 = false

func main() {
	buf, err := ioutil.ReadFile("10.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	for y := 0; y < len(M); y++ {
		for x := 0; x < len(M[y]); x++ {
			fmt.Printf("%c", M[y][x])
		}
		fmt.Printf("\n")
	}

	//fmt.Printf("output: %d\n", countvisible(2, 2))

	maxpoint := Point{}
	maxcnt := -1
	for y := 0; y < len(M); y++ {
		for x := 0; x < len(M[y]); x++ {
			if M[y][x] == '#' {
				cnt := len(listvisible(x, y))
				if debug {
					fmt.Printf("y=%d,x=%d %d\n", y, x, cnt)
				}
				if cnt > maxcnt {
					maxcnt = cnt
					maxpoint = Point{x: x, y: y}
				}
			}
		}
	}

	fmt.Printf("maximum %d,%d with %d\n", maxpoint.x, maxpoint.y, maxcnt)

	n := 1

	visible := listvisible(maxpoint.x, maxpoint.y)

	vaporize := func(p Point) {
		if debug2 {
			fmt.Printf("%d is %v\n", n, p)
			M[p.y][p.x] = fmt.Sprintf("%d", n%10)[0]
		} else {
			M[p.y][p.x] = '.'
			if n == 200 {
				fmt.Printf("%d,%d\n", p.x, p.y)
			}
		}
		n++
		if debug2 {
			if n%10 == 0 && n > 0 {
				for y := 0; y < len(M); y++ {
					for x := 0; x < len(M[y]); x++ {
						fmt.Printf("%c", M[y][x])
					}
					fmt.Printf("\n")
				}
				for y := 0; y < len(M); y++ {
					for x := 0; x < len(M[y]); x++ {
						if M[y][x] != '#' {
							M[y][x] = '.'
						}
					}
				}
			}
		}
	}

	// up
	for _, p := range visible {
		dx := p.x - maxpoint.x
		if dx == 0 && p.y < maxpoint.y {
			vaporize(p)
		}
	}

	topright := []Pointm{}
	bottomright := []Pointm{}
	bottomleft := []Pointm{}
	topleft := []Pointm{}
	for _, p := range visible {
		dx := p.x - maxpoint.x
		dy := p.y - maxpoint.y
		if dx == 0 || dy == 0 {
			continue
		}
		m := float64(dy) / float64(dx)
		if dx > 0 && dy < 0 {
			topright = append(topright, Pointm{p, m})
		} else if dx > 0 && dy > 0 {
			bottomright = append(bottomright, Pointm{p, m})
		} else if dx < 0 && dy > 0 {
			bottomleft = append(bottomleft, Pointm{p, m})
		} else if dx < 0 && dy < 0 {
			topleft = append(topleft, Pointm{p, m})
		}
	}

	sort.Slice(topright, func(i, j int) bool {
		return topright[i].m < topright[j].m
	})

	sort.Slice(bottomright, func(i, j int) bool {
		return bottomright[i].m < bottomright[j].m
	})

	sort.Slice(bottomleft, func(i, j int) bool {
		return bottomleft[i].m < bottomleft[j].m
	})

	sort.Slice(topleft, func(i, j int) bool {
		return topleft[i].m < topleft[j].m
	})

	// top-right
	for _, p := range topright {
		if debug2 {
			fmt.Printf("%v %g\n", p.p, p.m)
		}
		vaporize(p.p)
	}

	// right
	for _, p := range visible {
		dy := p.y - maxpoint.y
		if dy == 0 && p.x > maxpoint.x {
			vaporize(p)
		}
	}

	// bottom-right
	for _, p := range bottomright {
		if debug2 {
			fmt.Printf("%v %g\n", p.p, p.m)
		}
		vaporize(p.p)
	}

	// down
	for _, p := range visible {
		dx := p.x - maxpoint.x
		if dx == 0 && p.y > maxpoint.y {
			vaporize(p)
		}
	}

	// bottom-left
	for _, p := range bottomleft {
		if debug2 {
			fmt.Printf("%v %g\n", p.p, p.m)
		}
		vaporize(p.p)
	}

	// left
	for _, p := range visible {
		dy := p.y - maxpoint.y
		if dy == 0 && p.x < maxpoint.x {
			vaporize(p)
		}
	}

	// top-left
	for _, p := range topleft {
		if debug2 {
			fmt.Printf("%v %g\n", p.p, p.m)
		}
		vaporize(p.p)
	}

	if debug2 {
		for y := 0; y < len(M); y++ {
			for x := 0; x < len(M[y]); x++ {
				fmt.Printf("%c", M[y][x])
			}
			fmt.Printf("\n")
		}
	}
}
