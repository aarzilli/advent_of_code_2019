package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"os"
)

const N = 26

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
	chs [4]byte
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
	s += fmt.Sprintf("{%c %c %c %c ", cur.chs[0], cur.chs[1], cur.chs[2], cur.chs[3])
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

type Point struct {
	i, j int
}

type reachout struct {
	key byte
	dist int
}

func reachablekeys(starti, startj int, keys [N]bool) []reachout {
	fringe := make(map[Point]bool)
	fringe[Point{starti, startj}] = true
	seen := make(map[Point]bool)
	r := []reachout{}
	
	cnt := 0
	for len(fringe) > 0 {
		newfringe := make(map[Point]bool)
		
		expand := func(i, j int) {
			if seen[Point{i, j}] {
				return
			}
			if i < 0 || i >= len(M) {
				return
			}
			if j < 0 || j >= len(M[i]) {
				return
			}
			if M[i][j] == '#' {
				return
			} else if n, ok := iskey(M[i][j]); ok {
				if !keys[n] {
					r = append(r, reachout{ M[i][j], cnt+1})
				}
			} else if n, ok := isdoor(M[i][j]); ok {
				if !keys[n] {
					return
				}
			}
			newfringe[Point{i, j}] = true
		}
		
		for p := range fringe {
			seen[p] = true
			expand(p.i-1, p.j)
			expand(p.i, p.j-1)
			expand(p.i, p.j+1)
			expand(p.i+1, p.j)
		}
		
		fringe = newfringe
		cnt++
	}
	
	return r
}

/*
var minpathlen int

func explore(depth int, pathlen int, s state) {
	if finished(s) {
		if minpathlen == 0 || pathlen < minpathlen {
			fmt.Printf("min path found %d\n", pathlen)
			minpathlen = pathlen
		}
		return
	}
	if minpathlen > 0 && pathlen > minpathlen {
		return
	}
	if depth == 0 {
		fmt.Printf("at %s\n", showstate(s))
	}
	reachkeys := reachablekeys(s.i, s.j, s.keys)
	sort.Slice(reachkeys, func(i, j int) bool {
		return reachkeys[i].dist < reachkeys[j].dist
	})
	for _, ro := range reachkeys {
		if s.keys[ro.key-'a'] {
			continue
		}
		//fmt.Printf("\tgoing to %c\n", key)
		keypos := position(ro.key)
		//newpathlen := pathlen + dist(Point{s.i, s.j}, keypos, s.keys)
		newpathlen := pathlen + ro.dist
		newkeys := s.keys
		newkeys[ro.key - 'a'] = true
		explore(depth+1, newpathlen, state{ keypos.i, keypos.j, newkeys })
	}
}
*/

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
	p Point
	keys string
}

func dist(start, end Point) (bool, int, string) {
	S := make(map[Point]int)

	queue := []distqel{{ p: start, keys: "" }}
	S[start] = 0

	for {
		if len(queue) == 0 {
			return false, -1, ""
		}
		p := queue[0]
		queue = queue[1:]

		if p.p == end {
			return true, S[end], p.keys
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
			queue = append(queue, distqel{ p: p2, keys: keys })
		}

		add(Point{p.p.i - 1, p.p.j})
		add(Point{p.p.i, p.p.j - 1})
		add(Point{p.p.i, p.p.j + 1})
		add(Point{p.p.i + 1, p.p.j})
	}

	panic("not found 2")
}

func rotomap(rk []reachout) map[byte]int {
	m := make(map[byte]int)
	for _, ro := range rk {
		m[ro.key] = ro.dist
	}
	return m
}

var G = map[byte][]E{}

type E struct {
	dest byte
	dist int
	needkeys string
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
	fringe := map[state]int{state{chs:[4]byte{  '@', '$', '%', '~'  } }: 0}   // nodes discovered but not visited (start at node 0 with distance 0)
	seen := map[state]bool{state{chs:[4]byte{  '@', '$', '%', '~'  }}: true} // nodes already visited (we know the minimum distance of those)

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
		
		//fmt.Printf("current %s %d\n", showstate(cur), distcur)

		maybeadd := func(robotidx int, e E) {
			// check if we can add the node
			for _, k := range e.needkeys {
				if !cur.keys[k-'a'] {
					return
				}
			}
			if cur.keys[e.dest-'a'] {
				return
			}
			
			nb := cur
			nb.keys[e.dest-'a'] = true
			nb.chs[robotidx] = e.dest
			
			// if we can add the node add it to the fringe
			// but first check that it's either a new node or we improved its distance
			if d, ok := fringe[nb]; !ok || distcur+e.dist < d {
				fringe[nb] = distcur + e.dist
			}
		}
		
		// try to add all possible neighbors
		for robotidx := 0; robotidx < 4; robotidx++ {
			for _, e := range G[cur.chs[robotidx]] {
				maybeadd(robotidx, e)
			}
		}
	}
}

func main() {
	buf, err := ioutil.ReadFile("18p2.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}
	
	robots := []byte{ '@', '$', '%', '~' }
	robotsidx := 0
		
	for i := range M {
		for j := range M[i] {
			if M[i][j] == '@' {
				M[i][j] = robots[robotsidx]
				robotsidx++
			}
			if n, ok := iskey(M[i][j]); ok {
				tocollect[n] = true
			}
		}
	}
	
	showmatrix()
	
	if robotsidx != 4 {
		panic("use 18t2p1.go this is for part 2")
	}
	
	for i := 0; i < N; i++ {
		if !tocollect[i] {
			continue
		}
		posi := position(byte(i+'a'))
		
		for _, robot := range robots {
			ok , d, keys := dist(position(robot), posi)
			if ok {
				G[robot] = append(G[robot], E{ dest: byte(i+'a'), dist: d, needkeys: keys })
				fmt.Printf("from %c to %c dist %d needs %q\n", robot, i+'a', d, keys)
			}
		}
		
		for j := 0; j < N; j++ {
			if i == j || !tocollect[j] {
				continue
			}
			posj := position(byte(j+'a'))
			ok, d, keys := dist(posi, posj)
			if ok {
				G[byte(i+'a')] = append(G[byte(i+'a')], E{ dest: byte(j+'a'), dist: d, needkeys: keys })
				fmt.Printf("from %c to %c dist %d needs %q\n", i+'a', j+'a', d, keys)
			}
		}
	}
	
	search()
	
	//fmt.Printf("%v\n", reachablekeys(starti, startj, [N]bool{}))
	/*explore(0, 0, state{ i: starti, j: startj })
	fmt.Printf("final minpathlen: %d\n", minpathlen)*/
}


// Not 3242