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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func instr4(p []int, pc *int) (int, int, int) {
	a1 := p[*pc+1]
	a2 := p[*pc+2]
	d := p[*pc+3]
	*pc += 4
	return a1, a2, d
}

var InstrLen = map[int]int{
	1: 4,
	2: 4,
	3: 2,
	4: 2,

	5: 3,
	6: 3,
	7: 4,
	8: 4,

	99: 1,
}

func cpu(p []int) {
	pc := 0

	modev := make([]int, 3)

evalLoop:
	for pc < len(p) {
		opcode := p[pc]
		mode := opcode / 100
		opcode = opcode % 100

		modev[0] = mode % 10
		modev[1] = (mode / 10) % 10
		modev[2] = (mode / 100) % 10

		n := InstrLen[opcode]
		a := p[pc : pc+n][1:]

		arg := func(n int) int {
			if modev[n] == 0 {
				return p[a[n]]
			}
			return a[n]
		}

		jumped := false

		fmt.Printf("%04d\t%03d %d %d %v\n", pc, mode, opcode, n, a)

		switch opcode {
		case 1:
			if modev[2] != 0 {
				panic("wtf")
			}
			p[a[2]] = arg(0) + arg(1)
		case 2:
			if modev[2] != 0 {
				panic("wtf")
			}
			p[a[2]] = arg(0) * arg(1)

		case 3: // input
			if modev[0] != 0 {
				panic("wtf")
			}
			p[a[0]] = 5 // PROBLEM DEPENDENT!!!
		case 4: // output
			fmt.Printf("OUT: %d\n", arg(0))

		case 5:
			if arg(0) != 0 {
				pc = arg(1)
				jumped = true
			}

		case 6:
			if arg(0) == 0 {
				pc = arg(1)
				jumped = true
			}

		case 7:
			if modev[2] != 0 {
				panic("wtf")
			}
			if arg(0) < arg(1) {
				p[a[2]] = 1
			} else {
				p[a[2]] = 0
			}
		case 8:
			if modev[2] != 0 {
				panic("wtf")
			}
			if arg(0) == arg(1) {
				p[a[2]] = 1
			} else {
				p[a[2]] = 0
			}

		case 99:
			break evalLoop
		}
		if !jumped {
			pc += n
		}
	}
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("05.txt")
	must(err)
	var p []int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p = append(p, vatoi(splitandclean(line, ",", -1))...)
	}

	cpu(p)
}
