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

var M [][]byte

func showmatrix() {
	for i := range M {
		for j := range M[i] {
			fmt.Printf("%c", M[i][j])
		}
		fmt.Printf("\n")
	}
}

type state struct {
	i, j int
	keys [25]bool
	doors [25]bool
}

// finds the closest node in the fringe, lastmin is an optimization, if we find a node that is at that distance we return it immediately (there can be nothing that's closer)
func minimum(fringe map[state]int, lastmin int) state {
	var mink state
	first := true
	for k, d := range fringe {
		if first {
			mink = k
			first = false
		}
		if d == lastmin {
			return k
		}
		if d < fringe[mink] {
			mink = k
		}
	}
	return mink
}

func search(starti, startj int) {
	fringe := map[state]int{state{i: starti, j: startj}: 0}   // nodes discovered but not visited (start at node 0 with distance 0)
	seen := map[state]bool{state{i: starti, j: startj}: true} // nodes already visited (we know the minimum distance of those)

	lastmin := 0

	cnt := 0

	for len(fringe) > 0 {
		cur := minimum(fringe, lastmin)

		if cnt%1000 == 0 {
			fmt.Printf("fringe %d (min dist %d)\n", len(fringe), fringe[cur])
		}
		cnt++

		if finished(cur) {
			fmt.Printf("%s %d\n", showstate(cur), fringe[cur])
			return
		}

		distcur := fringe[cur]
		lastmin = distcur
		delete(fringe, cur)
		seen[cur] = true

		maybeadd := func(nb state) {
			// check if we can add the node
			if nb.i < 0 || nb.i >= len(M) {
				return
			}
			if nb.j < 0 || nb.j >= len(M[nb.i]) {
				return
			}
			
			if M[nb.i][nb.j] == '#' {
				return
			} else if n, ok := isdoor(M[nb.i][nb.j]); ok {
				if !nb.keys[n] {
					return
				}
				nb.doors[n] = true
			} else if n, ok := iskey(M[nb.i][nb.j]); ok {
				nb.keys[n] = true
			}
			
			
			// if we can add the node add it to the fringe
			// but first check that it's either a new node or we improved its distance
			if d, ok := fringe[nb]; !ok || distcur+1 < d {
				fringe[nb] = distcur + 1
			}
		}

		// try to add all possible neighbors
		maybeadd(state{cur.i - 1, cur.j,  cur.keys, cur.doors})
		maybeadd(state{cur.i + 1, cur.j, cur.keys, cur.doors})
		maybeadd(state{cur.i, cur.j - 1, cur.keys, cur.doors})
		maybeadd(state{cur.i, cur.j + 1, cur.keys, cur.doors})
	}
}

func isdoor(b byte) (n int, ok bool) {
	if b < 'A' || b > 'Z' {
		return 0, false
	}
	
	return int(b - 'A'), true
}

func iskey(b byte) (n int, ok bool) {
	if b < 'a' || b > 'z' {
		return 0, false
	}
	return int(b - 'a'), true
}

var tocollect [25]bool

func finished(cur state) bool {
	for i := range cur.keys {
		if tocollect[i] && !cur.keys[i] {
			return false
		}
	}
	return true
}

func showstate(cur state) string {
	s := ""
	s += fmt.Sprintf("{%d, %d ", cur.i, cur.j)
	for i := range cur.keys {
		if tocollect[i] {
			if cur.keys[i] {
				s += fmt.Sprintf("%c", i+'a')
			} else {
				s += "."
			}
		}
	}
	return s+ "}"
}

func main() {
	buf, err := ioutil.ReadFile("18.example2")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}
	showmatrix()
	
	var starti, startj int
	
	for i := range M {
		for j := range M[i] {
			if M[i][j] == '@' {
				starti = i
				startj = j
			}
			if n, ok := iskey(M[i][j]); ok {
				tocollect[n] = true
			}
		}
	}
	
	search(starti, startj)
}
