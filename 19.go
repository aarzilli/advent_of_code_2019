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

const TRACECPU = false

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

type Cpustate struct {
	p            []int
	mem          map[int]int
	pc           int
	relativeBase int
	input        func() int // if nil will suspend instead
	output       func(int)  // if nil will suspend instead
	// update copycpu if you change this
}

func newCpustate(p []int) *Cpustate {
	return &Cpustate{
		p:            copyprog(p),
		mem:          make(map[int]int),
		pc:           0,
		relativeBase: 0,
	}
}

func cpu(s *Cpustate, input int, inputValid bool) (int, int) {
	modev := make([]int, 3)

	if TRACECPU {
		fmt.Printf("PC\tBASE\tOPCODE\tARGS\n")
	}

	if s.input != nil && inputValid {
		panic("two input systems provided")
	}

	for s.pc < len(s.p) {
		opcode := s.p[s.pc]
		mode := opcode / 100
		opcode = opcode % 100

		modev[0] = mode % 10
		modev[1] = (mode / 10) % 10
		modev[2] = (mode / 100) % 10

		n := Opcodes[opcode].Len
		a := s.p[s.pc : s.pc+n][1:]

		arg := func(n int) int {
			var addr int
			switch modev[n] {
			case 0:
				addr = a[n]
			case 1:
				return a[n]
			case 2:
				addr = a[n] + s.relativeBase
			default:
				panic("wtf")
			}
			if addr < len(s.p) {
				return s.p[addr]
			} else {
				return s.mem[addr]
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
				addr = a[n] + s.relativeBase
			default:
				panic("wtf")
			}
			if addr < len(s.p) {
				s.p[addr] = out
			} else {
				s.mem[addr] = out
			}
		}

		jumped := false

		if TRACECPU {
			prettyInstr(s.p, s.mem, s.pc, mode, opcode, a, s.relativeBase)
		}

		switch opcode {
		case 1: // ADD
			save(2, arg(0)+arg(1))
		case 2: // MUL
			save(2, arg(0)*arg(1))

		case 3: // input
			if s.input == nil {
				if inputValid {
					save(0, input)
					if TRACECPU {
						fmt.Printf("\tinput was %d\n", input)
					}
					inputValid = false
				} else {
					// suspend CPU waiting for input
					return 3, 0
				}
			} else {
				in := s.input()
				save(0, input)
				if TRACECPU {
					fmt.Printf("\tinput was %d\n", in)
				}
			}
		case 4: // output
			if s.output == nil {
				// return output and suspend CPU
				s.pc += n
				return 4, arg(0)
			} else {
				s.output(arg(0))
			}

		case 5: // JNZ
			if arg(0) != 0 {
				s.pc = arg(1)
				jumped = true
			}

		case 6: // JZ
			if arg(0) == 0 {
				s.pc = arg(1)
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
			s.relativeBase += arg(0)

		case 99: // END
			return 99, 0
		}
		if !jumped {
			s.pc += n
		} else {
			if TRACECPU {
				fmt.Printf("\tjumped\n")
			}
		}
	}

	//TODO: do I need to run from memory?

	panic("ran out of instructions")
}

func copyprog(p []int) []int {
	q := make([]int, len(p))
	copy(q, p)
	return q
}

func copycpu(s *Cpustate) *Cpustate {
	r := *s
	r.p = copyprog(r.p)
	r.mem = make(map[int]int)

	for addr, v := range s.mem {
		r.mem[addr] = v
	}

	return &r
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

func exit(n int) {
	os.Exit(n)
}

const part2 = true

type Point struct {
	x, y int
}

var M = map[Point]int{}

const debug = false

func check(radius int) {
	for y := 0; y < radius; y++ {
		for x := 0; x < radius; x++ {
			if M[Point{x, y}] != 1 {
				continue
			}
			if M[Point{x + 99, y}] != 1 {
				continue
			}
			if M[Point{x, y + 99}] != 1 {
				continue
			}
			if M[Point{x + 99, y + 99}] != 1 {
				continue
			}
			fmt.Printf("Found at x=%d y=%d (PART 2: %d)\n", x, y, x*10000+y)
			exit(0)
		}
	}
}

func main() {
	p := readprog("19.txt")

	if !part2 {
		cnt := 0
		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				s := newCpustate(p)
				instr, _ := cpu(s, x, true)
				if instr != 3 {
					panic("wtf")
				}
				instr, n := cpu(s, y, true)
				if instr != 4 {
					panic("wtf2")
				}
				if n == 1 {
					cnt++
				}
			}
		}

		fmt.Printf("PART 1: %d\n", cnt)
		exit(0)
	}

	firstseenx := 0
	for y := 0; y < 10000; y++ {
		if y%100 == 0 {
			fmt.Printf("%d\n", y)
		}
		if y%1000 == 0 {
			fmt.Printf("check\n")
			check(y)
			fmt.Printf("check done\n")
		}
		first := true
		if debug {
			for x := 0; x < firstseenx; x++ {
				fmt.Printf(" ")
			}
		}
		for x := firstseenx; x < 10000; x++ {
			s := newCpustate(p)
			instr, _ := cpu(s, x, true)
			if instr != 3 {
				panic("wtf")
			}
			instr, n := cpu(s, y, true)
			if instr != 4 {
				panic("wtf2")
			}
			M[Point{x, y}] = n
			if debug {
				if n == 1 {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			}
			if n == 1 && first {
				firstseenx = x
				first = false
			}
			if x > 0 && n == 0 && M[Point{x - 1, y}] == 1 {
				break
			}
		}
		if debug {
			fmt.Printf("\n")
		}
	}
}
