package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
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
				save(0, in)
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

var mu sync.Mutex
var queue = make([][]int, 50)
var first = make([]bool, 50)
var cpus = make([]*Cpustate, 50)
var idle = make([]bool, 50)
var nat = make([]int, 2)

const part1 = false
const debug = true

func makecpu(p []int, idx int) {
	idlecnt := 0
	cpus[idx] = newCpustate(p)
	cpus[idx].input = func() int {
		mu.Lock()
		defer mu.Unlock()
		if first[idx] {
			first[idx] = false
			return idx
		}
		if len(queue[idx]) == 0 {
			idlecnt++
			if idlecnt == 10 {
				idle[idx] = true
				if debug {
					fmt.Printf("CPU %d is idle\n", idx)
				}
			}
			if idx == 0 {
				allidle := true
				for i := range idle {
					if !idle[i] || len(queue[i]) > 0 {
						allidle = false
						break
					}
				}
				if allidle {
					fmt.Printf("NAT sends %d %d\n", nat[0], nat[1])
					queue[idx] = append(queue[idx], nat[0])
					queue[idx] = append(queue[idx], nat[1])
				}
			}
			if len(queue[idx]) == 0 {
				return -1
			}
		}
		if debug && idlecnt > 2 {
			fmt.Printf("CPU %d unhidled\n", idx)
		}
		idlecnt = 0
		idle[idx] = false
		r := queue[idx][0]
		queue[idx] = queue[idx][1:]
		return r
	}
	outch := make(chan int)
	cpus[idx].output = func(n int) {
		outch <- n
	}
	go func() {
		for {
			dest, ok := <-outch
			if !ok {
				return
			}
			p1 := <-outch
			p2 := <-outch
			if dest == 255 {
				if part1 {
					fmt.Printf("PART 1: %d\n", p2)
					exit(0)
				} else {
					fmt.Printf("Sent to NAT %d %d\n", p1, p2)
					mu.Lock()
					fmt.Printf("sent to NAT %d %d\n", nat[0], nat[1])
					nat[0] = p1
					nat[1] = p2
					mu.Unlock()
				}
			} else {
				mu.Lock()
				queue[dest] = append(queue[dest], p1)
				queue[dest] = append(queue[dest], p2)
				mu.Unlock()
			}
		}
	}()
	go func() {
		cpu(cpus[idx], 0, false)
	}()
}

func main() {
	p := readprog("23.txt")

	for i := range first {
		first[i] = true
	}
	for i := range cpus {
		makecpu(p, i)
	}

	select {}
}

// 99289
