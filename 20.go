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

var M [][]byte

func showmatrix() {
	for i := range M {
		for j := range M[i] {
			fmt.Printf("%c", M[i][j])
		}
		fmt.Printf("\n")
	}
}

func isletter(i, j int) bool {
	if i < 0 || i >= len(M) {
		return false
	}
	if j < 0 || j >= len(M[i]) {
		return false
	}
	ch := M[i][j]
	return ch >= 'A' && ch <= 'Z'
}

type Point struct {
	i, j int
	z    int
}

type FixedPoint struct {
	i, j int
}

type Portal struct {
	i, j   int
	deltaz int
}

var portalpoint = map[FixedPoint]string{}
var portal = map[string][]Portal{}

func isportal(p Point) string {
	return portalpoint[FixedPoint{p.i, p.j}]
}

func search(part2 bool) {
	aa := portal["AA"][0]
	fringe := []Point{Point{aa.i, aa.j, 0}}
	S := make(map[Point]int)

	cnt := 0
	for len(fringe) > 0 {
		cur := fringe[0]
		fringe = fringe[1:]

		if cnt%1000 == 0 && debug {
			fmt.Printf("%d fringe %d\n", cnt, len(fringe))
		}
		cnt++

		if isportal(cur) == "ZZ" && cur.z == 0 {
			part := 1
			if part2 {
				part = 2
			}
			fmt.Printf("PART %d: %d\n", part, S[cur])
			break
		}

		curdist := S[cur]

		add := func(i, j, deltaz int) {
			if cur.z+deltaz < 0 {
				return
			}
			if i < 0 || i >= len(M) {
				return
			}
			if j < 0 || j >= len(M[i]) {
				return
			}
			if M[i][j] != '.' {
				return
			}
			p := Point{i, j, cur.z + deltaz}
			if _, ok := S[p]; ok {
				return
			}
			fringe = append(fringe, p)
			S[p] = curdist + 1
		}

		add(cur.i-1, cur.j, 0)
		add(cur.i, cur.j-1, 0)
		add(cur.i, cur.j+1, 0)
		add(cur.i+1, cur.j, 0)

		if s := isportal(cur); s != "" {
			portals := portal[s]
			if len(portals) != 1 {
				if len(portals) != 2 {
					panic("wtf")
				}
				var curportal, nextportal Portal
				f1, f2 := false, false
				for _, nb := range portal[s] {
					if nb.i == cur.i && nb.j == cur.j {
						if f1 {
							panic("wtf")
						}
						curportal = nb
						f1 = true
					} else {
						if f2 {
							panic("wtf")
						}
						nextportal = nb
						f2 = true
					}
				}
				if !f1 || !f2 {
					panic("wtf")
				}
				deltaz := curportal.deltaz
				if !part2 {
					deltaz = 0
				}
				add(nextportal.i, nextportal.j, deltaz)
			}
		}
	}
}

const debug = false

func main() {
	buf, err := ioutil.ReadFile("20.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	if debug {
		showmatrix()
	}

	for i := range M {
		for j := range M[i] {
			if isletter(i, j) && isletter(i+1, j) {
				var p FixedPoint
				if i+2 < len(M) && M[i+2][j] == '.' {
					p.i = i + 2
					p.j = j
				}
				if i-1 > 0 && M[i-1][j] == '.' {
					p.i = i - 1
					p.j = j
				}
				s := string(M[i][j]) + string(M[i+1][j])
				portalpoint[p] = s
			}
			if isletter(i, j) && isletter(i, j+1) {
				var p FixedPoint
				if j+2 < len(M[i]) && M[i][j+2] == '.' {
					p.i = i
					p.j = j + 2
				}
				if j-1 >= 0 && M[i][j-1] == '.' {
					p.i = i
					p.j = j - 1
				}
				s := string(M[i][j]) + string(M[i][j+1])
				portalpoint[p] = s
			}
		}
	}

	for p, s := range portalpoint {
		var q Portal
		q.i = p.i
		q.j = p.j

		if p.i == 2 || p.i == len(M)-3 {
			q.deltaz = -1
		} else if p.j == 2 || p.j == len(M[p.i])-3 {
			q.deltaz = -1
		} else {
			q.deltaz = +1
		}
		portal[s] = append(portal[s], q)
	}

	if len(portal["AA"]) != 1 {
		panic("wtf")
	}

	search(false)
	search(true)
}
