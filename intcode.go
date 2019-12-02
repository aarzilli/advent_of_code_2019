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

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
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

func instr4(p []int, pc *int) (int, int, int) {
	a1 := p[*pc+1]
	a2 := p[*pc+2]
	d := p[*pc+3]
	*pc += 4
	return a1, a2, d
}

var InstrLen = map[int]int{
	1:  4,
	2:  4,
	99: 1,
}

func cpu(p []int) {
	pc := 0

evalLoop:
	for pc < len(p) {
		opcode := p[pc]
		n := InstrLen[opcode]
		a := p[pc : pc+n][1:]
		switch opcode {
		case 1:
			p[a[2]] = p[a[0]] + p[a[1]]
		case 2:
			p[a[2]] = p[a[0]] * p[a[1]]
		case 99:
			break evalLoop
		}
		pc += n
	}
}

func pretty(p []int, start int) {
	pc := start
	for pc < len(p) {
		opcode := p[pc]
		n := InstrLen[opcode]
		pc += n
		if pc == 99 {
			break
		}
	}

	ramstart := pc

	dirty := make([]bool, len(p))

	dirty[1] = true
	dirty[2] = true

	constOrInd := func(a int) string {
		if a < ramstart && !dirty[a] {
			return fmt.Sprintf("%d", p[a])
		}
		return fmt.Sprintf("[%d]", a)
	}

	pc = start

instrLoop:
	for pc < len(p) {
		opcode := p[pc]
		n := InstrLen[opcode]
		a := p[pc : pc+n][1:]
		fmt.Printf("%04d\t", pc)
		switch opcode {
		case 1:
			a0 := constOrInd(a[0])
			a1 := constOrInd(a[1])
			fmt.Printf("%s + %s -> [%d]\n", a0, a1, a[2])
			dirty[a[2]] = true
		case 2:
			a0 := constOrInd(a[0])
			a1 := constOrInd(a[1])
			fmt.Printf("%s * %s -> [%d]\n", a0, a1, a[2])
			dirty[a[2]] = true
		case 99:
			fmt.Printf("END\n")
			pc++
			break instrLoop
		}
		pc += n
	}

	for pc < len(p) {
		fmt.Printf("%04d\t%d\n", pc, p[pc])
		pc++
	}
}

func copyprog(p []int) []int {
	q := make([]int, len(p))
	copy(q, p)
	return q
}

const part1 = false
const debug = false

func main() {
	buf, err := ioutil.ReadFile("02.txt")
	must(err)
	var p []int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p = append(p, vatoi(splitandclean(line, ",", -1))...)
	}

	if part1 {
		if debug {
			pretty(p, 0)
		}

		p[1] = 12
		p[2] = 2
		cpu(p)
		fmt.Printf("%d\n", p[0])
	} else {
		for in1 := 0; in1 < 100; in1++ {
			for in2 := 0; in2 < 100; in2++ {
				q := copyprog(p)
				q[1] = in1
				q[2] = in2
				cpu(q)
				if q[0] == 19690720 {
					fmt.Printf("%d%d\n", in1, in2)
					break
				}
			}
		}
	}
}
