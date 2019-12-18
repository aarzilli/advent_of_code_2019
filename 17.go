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

func cpu(p []int, input func() int, output func(int)) []int {
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
			x := input()
			save(0, x)
			if TRACECPU {
				fmt.Printf("\tinput was %d\n", x)
			}
		case 4: // output
			output(arg(0))

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

func exit(n int) {
	os.Exit(n)
}

const part2 = true

var M [][]byte

type Pathelem struct {
	Turn byte
	Step int
}

func main() {
	p := readprog("17.txt")
	if part2 {
		p[0] = 2
	}
	first := true
	phase2 := false
	phase3 := false
	phase3done := make(chan bool)
	inchan := make(chan int)
	done := make(chan bool)
	var buf []byte
	go func() {
		cpu(p, func() int {
			if first {
				close(done)
				first = false
			}
			return <- inchan
		}, func(n int) {
			if phase3 {
			} else if phase2 {
				if n > 0xff {
					fmt.Printf("OUT %d\n", n)
				} else {
					fmt.Printf("OUT %c\n", n)
				}
			} else {
				buf = append(buf, byte(n))
			}
		})
		close(phase3done)
	}()
	<-done

	lines := strings.Split(string(buf), "\n")
	M = make([][]byte, len(lines))
	for i := range M {
		M[i] = []byte(lines[i])
	}
	
	filled := func(i, j int) bool {
		if i < 0 || i >= len(M) {
			return false
		}
		if j < 0 || j >= len(M[i]) {
			return false
		}
		return M[i][j] == '#' || M[i][j] == 'O' || M[i][j] == 'X' || M[i][j] == '^'
	}
	
	var pos [2]int
	var dir byte
	
	r := 0
	for i := range M {
		for j := range M[i] {
			if filled(i, j) && filled(i-1, j) && filled(i, j-1) && filled(i, j+1) && filled(i+1, j) {
				r += i*j
			}
			if M[i][j] == '^' || M[i][j] == '<' || M[i][j] == '>' || M[i][j] == 'v' {
				pos[0] = i
				pos[1] = j
				dir = M[i][j]
			}
		}
	}
	
	_ = dir
	
	var endpos [2]int
	
	for i := range M {
		for j := range M[i] {
			if !filled(i, j) || M[i][j] == '^' {
				continue
			}
			n := 0
			if filled(i-1, j) {
				n++
			}
			if filled(i, j-1) {
				n++
			}
			if filled(i, j+1) {
				n++
			}
			if filled(i+1, j) {
				n++
			}
			if n == 1 {
				endpos = [2]int{ i, j }
			}
		}
	}
	
	for i := range M {
		fmt.Printf("%v\n", string(M[i]))
	}
	
	fmt.Printf("PART 1: %d\n", r)
	
	
	path := []Pathelem{ }
	
	
	for {		
		if pos == endpos {
			break
		}
		var nextpos [2]int
		switch dir {
		case '^':
			nextpos = [2]int{ pos[0] - 1, pos[1] }
		case '<':
			nextpos = [2]int{ pos[0], pos[1] - 1}
		case '>':
			nextpos = [2]int{ pos[0], pos[1] + 1 }
		case 'v':
			nextpos = [2]int{ pos[0] + 1, pos[1] }
		}
		
		if filled(nextpos[0], nextpos[1]) {
			path[len(path)-1].Step++
			pos = nextpos
		} else {			
			var nextdir byte
			
			switch dir {
			case '^':
				switch {
				case filled(pos[0], pos[1]-1):
					path = append(path, Pathelem{ 'L', 0 })
					nextdir = '<'
				case filled(pos[0], pos[1]+1):
					path = append(path, Pathelem{ 'R', 0 })
					nextdir = '>'
				}
			case '<':
				switch {
				case filled(pos[0]-1, pos[1]):
					path = append(path, Pathelem{ 'R', 0 })
					nextdir = '^'
				case filled(pos[0]+1, pos[1]):
					path = append(path, Pathelem{ 'L', 0 })
					nextdir = 'v'
				}
				
			case '>':
				switch {
				case filled(pos[0]-1, pos[1]):
					path = append(path, Pathelem{ 'L', 0 })
					nextdir = '^'
				case filled(pos[0]+1, pos[1]):
					path = append(path, Pathelem{ 'R', 0 })
					nextdir = 'v'
				}
				
			case 'v':
				switch {
				case filled(pos[0], pos[1]-1):
					path = append(path, Pathelem{ 'R', 0 })
					nextdir = '<'
				case filled(pos[0], pos[1]+1):
					path = append(path, Pathelem{ 'L', 0 })
					nextdir = '>'
				}
			}
			if nextdir == 0 {
				panic("wtf")
			}
			dir = nextdir
		}
	}
	
	var spath string
	
	for i := range path {
		spath += fmt.Sprintf("%c,%d", path[i].Turn, path[i].Step)
		if i < len(path)-1 {
			spath += ","
		}
	}
	
	const programA = "L,6,L,4,L,12"
	const programB = "R,10,L,8,L,4,R,10"
	const programC = "L,12,L,8,R,10,R,10"
	
	spath = strings.Replace(spath, programA, "A", -1)
	spath = strings.Replace(spath, programB, "B", -1)
	spath = strings.Replace(spath, programC, "C", -1)
	
	fmt.Printf("A: %s %d\n", programA, len(programA))
	fmt.Printf("B: %s %d\n", programB, len(programB))
	fmt.Printf("C: %s %d\n", programC, len(programC))
	fmt.Printf("%s %d\n", spath, len(spath))
	
	phase2 = true
	
	sendline := func(s string) {
		for i := range s {
			inchan <- int(s[i])
		}
		inchan <- '\n'
	}
	
	sendline(spath)
	sendline(programA)
	sendline(programB)
	sendline(programC)
	
	sendline("n")
	
	<- phase3done
	//TODO: send movement function
	//TODO: send movement sub-functions
	
}
