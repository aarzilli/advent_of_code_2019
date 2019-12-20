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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func exit(n int) {
	os.Exit(n)
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

var portalpoint = map[Point]string{}
var portal = map[string][]Point{}

func isportal(p Point) string {
	return portalpoint[Point{p.i, p.j, 0}]
}

func search() map[Point]int {
	aa := portal["AA"][0]
	fringe := []Point{Point{aa.i, aa.j, 0}}
	S := make(map[Point]int)

	cnt := 0
	for len(fringe) > 0 {
		cur := fringe[0]
		fringe = fringe[1:]

		if cnt%1000 == 0 {
			fmt.Printf("%d fringe %d\n", cnt, len(fringe))
		}
		cnt++

		if isportal(cur) == "ZZ" && cur.z == 0 {
			fmt.Printf("PART 1: %d\n", S[cur])
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
			/*if isportal(cur) != "" && deltaz != 0 {
				fmt.Printf("adding from %v to %v (%s)\n", cur, p, isportal(cur))
			}*/
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
				var curportal, nextportal Point
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
				add(nextportal.i, nextportal.j, curportal.z)
			}
		}
	}

	return S
}

func main() {
	buf, err := ioutil.ReadFile("20.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		//line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	showmatrix()

	for i := range M {
		for j := range M[i] {
			if isletter(i, j) && isletter(i+1, j) {
				var p Point
				if i+2 < len(M) && M[i+2][j] == '.' {
					p.i = i + 2
					p.j = j
				}
				if i-1 > 0 && M[i-1][j] == '.' {
					p.i = i - 1
					p.j = j
				}
				s := string(M[i][j]) + string(M[i+1][j])
				//fmt.Printf("point %d,%d is portal %s\n", p.i, p.j, s)
				portalpoint[p] = s
			}
			if isletter(i, j) && isletter(i, j+1) {
				var p Point
				if j+2 < len(M[i]) && M[i][j+2] == '.' {
					p.i = i
					p.j = j + 2
				}
				if j-1 >= 0 && M[i][j-1] == '.' {
					p.i = i
					p.j = j - 1
				}
				s := string(M[i][j]) + string(M[i][j+1])
				//fmt.Printf("point %d,%d is portal %s\n", p.i, p.j, s)
				portalpoint[p] = s
			}
		}
	}

	for p, s := range portalpoint {
		q := p
		if p.i == 2 || p.i == len(M)-3 {
			q.z = -1
		} else if p.j == 2 || p.j == len(M[p.i])-3 {
			q.z = -1
		} else {
			q.z = +1
		}
		portal[s] = append(portal[s], q)
		//fmt.Printf("point %v portal %s delta z %d\n", p, s, q.z)
	}

	if len(portal["AA"]) != 1 {
		panic("wtf")
	}

	S := search()
	_ = S
}
