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

type Matrix struct {
	M    [][]byte
	M2   [][]byte
	Down *Matrix
	Up   *Matrix
}

func newmatrix(M [][]byte) *Matrix {
	m := &Matrix{}
	if M != nil {
		if len(M) != 5 {
			panic("wtf")
		}
		for i := range M {
			if len(M[i]) != 5 {
				panic("wtf")
			}
		}
		if part2 {
			if M[2][2] != '?' {
				panic("wtf")
			}
		}
		m.M = M
	} else {
		m.M = make([][]byte, 5)
		for i := range m.M {
			m.M[i] = make([]byte, 5)
			for j := range m.M[i] {
				m.M[i][j] = '.'
			}
		}
		m.M[2][2] = '?'
	}

	m.M2 = make([][]byte, len(m.M))
	for i := range m.M {
		m.M2[i] = make([]byte, len(m.M[i]))
	}
	return m
}

func showmatrix(m *Matrix) {
	for i := range m.M {
		for j := range m.M[i] {
			fmt.Printf("%c", m.M[i][j])
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
}

func squaretoint(m *Matrix, i, j int) int {
	if m.M[i][j] == '?' {
		panic("wtf")
	}
	if m.M[i][j] == '#' {
		return 1
	}
	return 0
}

func at(m *Matrix, i, j, oi, oj int) int {
	if i < 0 {
		if !part2 {
			return 0
		}
		if m.Down == nil {
			return 0
		}
		return squaretoint(m.Down, 1, 2)
	}
	if i >= len(m.M) {
		if !part2 {
			return 0
		}
		if m.Down == nil {
			return 0
		}
		return squaretoint(m.Down, 3, 2)
	}
	if j < 0 {
		if !part2 {
			return 0
		}
		if m.Down == nil {
			return 0
		}
		return squaretoint(m.Down, 2, 1)

	}
	if j >= len(m.M[i]) {
		if !part2 {
			return 0
		}
		if m.Down == nil {
			return 0
		}
		return squaretoint(m.Down, 2, 3)

	}
	if m.M[i][j] == '?' {
		if !part2 {
			panic("wtf")
		}
		if m.Up == nil {
			return 0
		}
		r := 0
		switch {
		case oi == 1 && oj == 2:
			// top row
			for j := range m.Up.M[0] {
				r += squaretoint(m.Up, 0, j)
			}
		case oi == 2 && oj == 1:
			// left column
			for i := range m.Up.M {
				r += squaretoint(m.Up, i, 0)
			}
		case oi == 2 && oj == 3:
			// right column
			for i := range m.Up.M {
				r += squaretoint(m.Up, i, 4)
			}
		case oi == 3 && oj == 2:
			// bottom row
			for j := range m.Up.M {
				r += squaretoint(m.Up, 4, j)
			}
		default:
			panic("wtf")
		}
		return r

	}
	return squaretoint(m, i, j)
}

func neighborsum(m *Matrix, i, j int) int {
	r := at(m, i-1, j, i, j) + at(m, i, j-1, i, j) + at(m, i, j+1, i, j) + at(m, i+1, j, i, j)
	return r
}

func stepsingle(m *Matrix) {
	for i := range m.M {
		for j := range m.M[i] {
			if m.M[i][j] == '?' {
				m.M2[i][j] = '?'
				continue
			}
			n := neighborsum(m, i, j)
			if m.M[i][j] == '#' {
				if n != 1 {
					m.M2[i][j] = '.'
				} else {
					m.M2[i][j] = '#'
				}
			} else {
				switch n {
				case 1, 2:
					m.M2[i][j] = '#'
				default:
					m.M2[i][j] = '.'
				}
			}
		}
	}
}

func step(m *Matrix) {
	stepsingle(m)
	var lastd, lastu *Matrix
	for d := m.Up; d != nil; d = d.Up {
		stepsingle(d)
		lastd = d
	}
	for u := m.Down; u != nil; u = u.Down {
		stepsingle(u)
		lastu = u
	}

	if lastd == nil {
		lastd = m
	}

	// switch to new state
	for c := lastd; c != nil; c = c.Down {
		c.M, c.M2 = c.M2, c.M
	}

	if part2 {
		if lastd == nil || lastu == nil {
			panic("wtf")
		}
		if infected(lastd) {
			lastd.Up = newmatrix(nil)
			lastd.Up.Down = lastd
		}
		if infected(lastu) {
			lastu.Down = newmatrix(nil)
			lastu.Down.Up = lastu
		}
	}
}

func infected(m *Matrix) bool {
	for i := range m.M {
		for j := range m.M[i] {
			if m.M[i][j] == '#' {
				return true
			}
		}
	}
	return false
}

var buf []byte

func stringify(m *Matrix) string {
	buf = buf[:0]
	for i := range m.M {
		buf = append(buf, m.M[i]...)
	}
	return string(buf)
}

func showmatrix2(m *Matrix) {
	var lastd *Matrix
	var cnt int
	for d := m.Up; d != nil; d = d.Up {
		lastd = d
		cnt++
	}

	d := lastd.Down
	cnt--
	for d != nil {
		fmt.Printf("Depth %d\n", -cnt)
		showmatrix(d)
		cnt--
		d = d.Down
	}

}

func countbugssingle(M [][]byte) int {
	cnt := 0
	for i := range M {
		for j := range M[i] {
			if M[i][j] == '#' {
				cnt++
			}
		}
	}
	return cnt
}

func countbugs(m *Matrix) int {
	var lastd *Matrix
	for d := m.Up; d != nil; d = d.Up {
		lastd = d
	}

	r := 0
	for d := lastd; d != nil; d = d.Down {
		r += countbugssingle(d.M)
	}
	return r
}

const part2 = true
const debug = false

func main() {
	in := "24.txt"
	buf, err := ioutil.ReadFile(in)
	must(err)

	var mat *Matrix
	{
		var M [][]byte
		for _, line := range strings.Split(string(buf), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			M = append(M, []byte(line))
		}

		if part2 {
			M[2][2] = '?'
		}

		mat = newmatrix(M)
		if part2 {
			mat.Down = newmatrix(nil)
			mat.Down.Up = mat
			mat.Up = newmatrix(nil)
			mat.Up.Down = mat
		}
	}

	if !part2 {
		m := map[string]bool{}

		for cnt := 0; cnt < 100; cnt++ {
			if cnt%1000 == 0 {
				fmt.Printf("at %d\n", cnt)
			}
			if debug {
				showmatrix(mat)
			}
			s := stringify(mat)
			if m[s] {
				break
			}
			m[s] = true
			step(mat)
		}

		showmatrix(mat)

		biod := 0
		p := 1

		for i := range mat.M {
			for j := range mat.M[i] {
				if mat.M[i][j] == '#' {
					biod += p
				}
				p *= 2
			}
		}

		fmt.Printf("PART 1: %d\n", biod)
	} else {
		N := 200
		if strings.Contains(in, "example") {
			N = 10
		}
		if debug {
			fmt.Printf("Initial state:\n")
			showmatrix2(mat)
		}

		for cnt := 0; cnt < N; cnt++ {
			step(mat)
			if debug {
				fmt.Printf("After minute %d\n", cnt+1)
				showmatrix2(mat)
			}
		}

		fmt.Printf("PART 2: %d\n", countbugs(mat))
	}
}
