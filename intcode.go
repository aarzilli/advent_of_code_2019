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

type Opcode struct {
	Len  int
	Name string
}

var Opcodes = map[int]Opcode{
	1:  {4, "ADD"},
	2:  {4, "MUL"},
	3:  {2, "INPUT"},
	4:  {2, "OUTPUT"},
	5:  {3, "JNZ"},
	6:  {3, "JZ"},
	7:  {4, "LT"},
	8:  {4, "EQ"},
	99: {1, "END"},
}

const TRACECPU = false

func prettyInstr(p []int, pc, mode, opcode int, a []int) {
	modev := make([]int, 3)
	modev[0] = mode % 10
	modev[1] = (mode / 10) % 10
	modev[2] = (mode / 100) % 10
	symop := fmt.Sprintf("? %d", opcode+mode*100)
	if oc, ok := Opcodes[opcode]; ok {
		symop = oc.Name
	}
	fmt.Printf("%04d\tmode=%03d %s\t", pc, mode, symop)
	for i := range a {
		if modev[i] == 0 {
			fmt.Printf(" [%d]=%d", a[i], p[a[i]])
		} else {
			fmt.Printf(" %d", a[i])
		}
	}
	fmt.Printf("\n")
}

func cpu(p []int, input int) []int {
	p = copyprog(p)
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

		n := Opcodes[opcode].Len
		a := p[pc : pc+n][1:]

		arg := func(n int) int {
			if modev[n] == 0 {
				return p[a[n]]
			}
			return a[n]
		}

		jumped := false

		if TRACECPU {
			prettyInstr(p, pc, mode, opcode, a)
		}

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
			p[a[0]] = input
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

	return p
}

func pretty(p []int, start int) {
	pc := start
	for pc < len(p) {
		opcode := p[pc]
		oc, ok := Opcodes[opcode]
		n := oc.Len
		if !ok {
			n = 1
		}
		pc += n
		if opcode == 99 {
			break
		}
	}

	dirty := make([]bool, len(p))

	dirty[1] = true
	dirty[2] = true

	pc = start

	for pc < len(p) {
		opcode := p[pc]
		mode := opcode / 100
		opcode = opcode % 100
		oc, ok := Opcodes[opcode]
		n := oc.Len
		if !ok {
			n = 1
		}
		a := p[pc : pc+n][1:]
		fmt.Printf("%04d\t", pc)
		prettyInstr(p, pc, mode, opcode, a)
		switch opcode {
		case 1, 2, 7, 8: // add, mul, lt, eq
			dirty[a[2]] = true
		case 3: // input
			dirty[a[0]] = true
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

func readprog(path string) []int {
	buf, err := ioutil.ReadFile(path)
	must(err)
	var p []int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p = append(p, vatoi(splitandclean(line, ",", -1))...)
	}
	return p
}

const disassemble = false

func main() {
	////////////////////////////////////
	// DAY 2

	p2 := readprog("02.txt")
	p2[1] = 12
	p2[2] = 2
	p2out := cpu(p2, 0)
	fmt.Printf("DAY 2 PART 1: %d\n", p2out[0])

	for in1 := 0; in1 < 100; in1++ {
		for in2 := 0; in2 < 100; in2++ {
			p2[1] = in1
			p2[2] = in2
			q := cpu(p2, 0)
			if q[0] == 19690720 {
				fmt.Printf("DAY 2 PART 2: %d%d\n", in1, in2)
				break
			}
		}
	}

	fmt.Printf("\n")

	////////////////////////////////////
	// DAY 5

	p5 := readprog("05.txt")

	if disassemble {
		pretty(p5, 0)
	}

	fmt.Printf("DAY 5 PART 1:\n")
	cpu(p5, 1)
	fmt.Printf("DAY 5 PART 2:\n")
	cpu(p5, 5)
}
