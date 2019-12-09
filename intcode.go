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
	9:  {2, "ADDBAS"},
	99: {1, "END"},
}

const TRACECPU = true

func prettyInstr(p []int, mem map[int]int, pc, mode, opcode int, a []int, base int) {
	modev := make([]int, 3)
	modev[0] = mode % 10
	modev[1] = (mode / 10) % 10
	modev[2] = (mode / 100) % 10
	symop := fmt.Sprintf("? %d", opcode+mode*100)
	if oc, ok := Opcodes[opcode]; ok {
		symop = oc.Name
	}
	fmt.Printf("%04d\t%04d\tmode=%03d %s\t", pc, base, mode, symop)
	for i := range a {
		switch modev[i] {
		case 0:
			if a[i] < len(p) {
				fmt.Printf(" [%d]=%d", a[i], p[a[i]])
			} else {
				fmt.Printf(" [%d]=%d", a[i], mem[a[i]])
			}
		case 1:
			fmt.Printf(" %d", a[i])
		case 2:
			addr := a[i] + base
			var n int
			if addr >= 0 && addr < len(p) {
				n = p[addr]
			} else {
				n = mem[addr]
			}
			fmt.Printf(" [BASE%+d]=%d", a[i], n)
		}
	}
	fmt.Printf("\n")
}

func cpu(p []int, input int) []int {
	p = copyprog(p)
	mem := make(map[int]int)
	pc := 0
	relativeBase := 0

	modev := make([]int, 3)

	if TRACECPU {
		fmt.Printf("PC\tBASE\tOPCODE\tARGS\n")
	}

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
			var addr int
			switch modev[n] {
			case 0:
				addr = a[n]
			case 1:
				return a[n]
			case 2:
				addr = a[n] + relativeBase
			default:
				panic("wtf")
			}
			if addr < len(p) {
				return p[addr]
			} else {
				return mem[addr]
			}

		}

		save := func(n, out int) {
			addr := 0
			switch modev[n] {
			case 0:
				addr = a[n]
			case 1:
				panic("wtf")
			case 2:
				addr = a[n] + relativeBase
			default:
				panic("wtf")
			}
			if addr < len(p) {
				p[addr] = out
			} else {
				mem[addr] = out
			}
		}

		jumped := false

		if TRACECPU {
			prettyInstr(p, mem, pc, mode, opcode, a, relativeBase)
		}

		switch opcode {
		case 1: // ADD
			save(2, arg(0)+arg(1))
		case 2: // MUL
			save(2, arg(0)*arg(1))

		case 3: // input
			save(0, input)
			if TRACECPU {
				fmt.Printf("\tinput was %d\n", input)
			}
		case 4: // output
			fmt.Printf("OUT: %d\n", arg(0))

		case 5: // JNZ
			if arg(0) != 0 {
				pc = arg(1)
				jumped = true
			}

		case 6: // JZ
			if arg(0) == 0 {
				pc = arg(1)
				jumped = true
			}

		case 7: // LT
			if arg(0) < arg(1) {
				save(2, 1)
			} else {
				save(2, 0)
			}
		case 8: // EQ
			if arg(0) == arg(1) {
				save(2, 1)
			} else {
				save(2, 0)
			}
		case 9: // ADDBAS
			relativeBase += arg(0)

		case 99: // END
			break evalLoop
		}
		if !jumped {
			pc += n
		} else {
			if TRACECPU {
				fmt.Printf("\tjumped\n")
			}
		}
	}

	return p
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

func main() {
	p := readprog("09.txt")
	fmt.Printf("PART 1:\n")
	cpu(p, 1)
	fmt.Printf("PART 2:\n")
	cpu(p, 2)
}
