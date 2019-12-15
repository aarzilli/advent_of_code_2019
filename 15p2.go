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

var M = [][]byte{}

type Point struct {
	i, j int
}

func showmap() {
	for i := range M {
		fmt.Printf("%s\n", string(M[i]))
	}

}

func done() bool {
	for i := range M {
		for j := range M[i] {
			if M[i][j] == '.' {
				return false
			}
		}
	}
	return true
}

const debug = false

func main() {
	buf, err := ioutil.ReadFile("15.addenda")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}

		M = append(M, []byte(line))
	}

	for i := range M {
		for j := range M[i] {
			if M[i][j] == '2' {
				M[i][j] = 'O'
			}
		}
	}

	if debug {
		showmap()
	}

	fringe := make(map[Point]bool)

	for i := range M {
		for j := range M[i] {
			if M[i][j] == 'O' {
				fringe[Point{i, j}] = true
			}
		}
	}

	cnt := 0
	for {
		newfringe := make(map[Point]bool)

		expand := func(i, j int) {
			if i < 0 || i >= len(M) {
				return
			}
			if j < 0 || j >= len(M[i]) {
				return
			}
			if M[i][j] == '.' {
				newfringe[Point{i, j}] = true
				M[i][j] = 'O'
			}
		}

		for p := range fringe {
			expand(p.i-1, p.j)
			expand(p.i, p.j-1)
			expand(p.i, p.j+1)
			expand(p.i+1, p.j)
		}
		cnt++

		fringe = newfringe
		if debug {
			showmap()
			fmt.Printf("\n")
		}

		if done() {
			fmt.Printf("PART 2: %d\n", cnt)
			break
		}
	}
}
