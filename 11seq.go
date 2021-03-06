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

type cpustate struct {
	p []int
	mem map[int]int
	pc int
	relativeBase int
}

func newCpustate(p []int) *cpustate {
	return &cpustate{
		p: copyprog(p),
		mem: make(map[int]int),
		pc: 0,
		relativeBase: 0,
	}
}

func cpu(s *cpustate, input int, inputValid bool) (int, int) {
	modev := make([]int, 3)

	if TRACECPU {
		fmt.Printf("PC\tBASE\tOPCODE\tARGS\n")
	}

evalLoop:
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
			if inputValid {
				save(0, input)
				if TRACECPU {
					fmt.Printf("\tinput was %d\n", input)
				}
				inputValid = false
			} else {
				return 3, 0
			}
		case 4: // output
			//fmt.Printf("OUT: %d\n", arg(0))
			s.pc += n
			return 4, arg(0)

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
			break evalLoop
		}
		if !jumped {
			s.pc += n
		} else {
			if TRACECPU {
				fmt.Printf("\tjumped\n")
			}
		}
	}

	return 99, 0
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

type Point struct {
	i, j int
}

type Direction uint8

const (
	Up = iota
	Right
	Down
	Left
)

var M = map[Point]int{}

const part2 = true

func main() {
	p := readprog(os.Args[1])

	pos := Point{0, 0}
	dir := Up

	if part2 {
		M[pos] = 1
	}
	
	s := newCpustate(p)
	
	input := 0
	inputValid := false

roboloop:
	for {
		op, x := cpu(s, input, inputValid)
		
		switch op {
		case 3: // input
			input = M[pos]
			inputValid = true
			
		case 4: // output
			M[pos] = x
			op, x := cpu(s, 0, false)
			if op != 4 {
				panic("wtf")
			}
			if x == 0 {
				dir = (dir - 1) % 4
			} else {
				dir = (dir + 1) % 4
			}
			if dir < 0 {
				dir = Left
			}
			switch dir {
			case Up:
				pos.i--
			case Down:
				pos.i++
			case Left:
				pos.j--
			case Right:
				pos.j++
			default:
				panic(fmt.Errorf("wtf %d", dir))
			}
			
		case 99:
			fmt.Printf("finished\n")
			break roboloop
		}
	}

	fmt.Printf("PART 1: %d\n", len(M))

	if part2 {
		var mini, minj, maxi, maxj = 0, 0, 0, 0
		for pos := range M {
			if pos.i < mini {
				mini = pos.i
			}
			if pos.j < minj {
				minj = pos.j
			}
			if pos.i > maxi {
				maxi = pos.i
			}
			if pos.j > maxj {
				maxj = pos.j
			}
		}

		for i := mini; i <= maxi; i++ {
			for j := minj; j <= maxj; j++ {
				if M[Point{i, j}] == 0 {
					fmt.Printf(".")
				} else {
					fmt.Printf("#")
				}
			}
			fmt.Printf("\n")
		}
	}
}
