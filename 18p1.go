package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const N = 26

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

type state struct {
	ch   byte
	keys [N]bool
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

var tocollect [N]bool

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
	s += fmt.Sprintf("{%c ", cur.ch)
	for i := range cur.keys {
		if tocollect[i] {
			if cur.keys[i] {
				s += fmt.Sprintf("%c", i+'a')
			} else {
				s += "."
			}
		}
	}
	return s + "}"
}

type Point struct {
	i, j int
}

func position(ch byte) Point {
	for i := range M {
		for j := range M[i] {
			if M[i][j] == ch {
				return Point{i, j}
			}
		}
	}
	panic("wtf")
}

type distqel struct {
	p    Point
	keys string
}

func dist(start, end Point) (int, string) {
	S := make(map[Point]int)

	queue := []distqel{{p: start, keys: ""}}
	S[start] = 0

	for {
		if len(queue) == 0 {
			panic("not found")
		}
		p := queue[0]
		queue = queue[1:]

		if p.p == end {
			return S[end], p.keys
		}

		add := func(p2 Point) {
			if p2.i < 0 || p2.i >= len(M) {
				return
			}
			if p2.j < 0 || p2.j >= len(M[p2.i]) {
				return
			}
			keys := p.keys
			if M[p2.i][p2.j] == '#' {
				return
			} else if _, ok := isdoor(M[p2.i][p2.j]); ok {
				keys = keys + strings.ToLower(string(M[p2.i][p2.j]))
			}
			if _, ok := S[p2]; ok {
				if S[p2] <= S[p.p]+1 {
					return
				}
			}
			S[p2] = S[p.p] + 1
			queue = append(queue, distqel{p: p2, keys: keys})
		}

		add(Point{p.p.i - 1, p.p.j})
		add(Point{p.p.i, p.p.j - 1})
		add(Point{p.p.i, p.p.j + 1})
		add(Point{p.p.i + 1, p.p.j})
	}

	panic("not found 2")
}

var G = map[byte][]E{}

type E struct {
	dest     byte
	dist     int
	needkeys string
}

type dstate struct {
	path string
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

func search() {
	fringe := map[state]int{state{ch: '@'}: 0}   // nodes discovered but not visited (start at node 0 with distance 0)
	seen := map[state]bool{state{ch: '@'}: true} // nodes already visited (we know the minimum distance of those)

	lastmin := 0

	cnt := 0

	for len(fringe) > 0 {
		cur := minimum(fringe, lastmin)

		if cnt%1000 == 0 {
			fmt.Printf("fringe %d (min dist %d)\n", len(fringe), fringe[cur])
		}
		cnt++

		if finished(cur) {
			fmt.Printf("%v %d\n", showstate(cur), fringe[cur])
			return
		}

		distcur := fringe[cur]
		lastmin = distcur
		delete(fringe, cur)
		seen[cur] = true

		maybeadd := func(e E) {
			// check if we can add the node
			if cur.keys[e.dest-'a'] {
				return
			}
			for _, k := range e.needkeys {
				if !cur.keys[k-'a'] {
					return
				}
			}

			keys := cur.keys
			keys[e.dest-'a'] = true
			nb := state{ch: e.dest, keys: keys}

			// if we can add the node add it to the fringe
			// but first check that it's either a new node or we improved its distance
			if d, ok := fringe[nb]; !ok || distcur+e.dist < d {
				fringe[nb] = distcur + e.dist
			}
		}

		// try to add all possible neighbors
		for _, e := range G[cur.ch] {
			maybeadd(e)
		}
	}
}

func main() {
	buf, err := ioutil.ReadFile("18.txt")
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

	for i := 0; i < N; i++ {
		if !tocollect[i] {
			continue
		}
		posi := position(byte(i + 'a'))

		{
			d, keys := dist(Point{starti, startj}, posi)
			G['@'] = append(G['@'], E{dest: byte(i + 'a'), dist: d, needkeys: keys})
		}

		for j := 0; j < N; j++ {
			if i == j || !tocollect[j] {
				continue
			}
			posj := position(byte(j + 'a'))
			d, keys := dist(posi, posj)
			G[byte(i+'a')] = append(G[byte(i+'a')], E{dest: byte(j + 'a'), dist: d, needkeys: keys})
		}
	}

	search()
}
