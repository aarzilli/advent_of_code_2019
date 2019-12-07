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

func cpu(p []int, phase int, inch <-chan int, outch chan<- int) int {
	p = copyprog(p)
	pc := 0

	phaseRead := false
	var out int

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
			if !phaseRead {
				p[a[0]] = phase
				phaseRead = true
			} else {
				p[a[0]] = <-inch
			}
		case 4: // output
			out = arg(0)
			outch <- out
			//fmt.Printf("OUT: %d\n", arg(0))

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

	return out
}

func copyprog(p []int) []int {
	q := make([]int, len(p))
	copy(q, p)
	return q
}

func runsequence(p []int, seq []int) int {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)
	ch5 := make(chan int)

	var out int

	done := make(chan bool)
	firstdone := make(chan bool)

	go func() {
		cpu(p, seq[0], ch1, ch2)
		close(firstdone)
	}()
	go cpu(p, seq[1], ch2, ch3)
	go cpu(p, seq[2], ch3, ch4)
	go cpu(p, seq[3], ch4, ch5)
	go func() {
		cpu(p, seq[4], ch5, ch1)
		close(done)
	}()

	ch1 <- 0
	<-firstdone
	out = <-ch1
	<-done

	return out
}

var program []int
var maxout int

const part1 = false

func enum(set []int, seen []bool, seq []int) {
	if len(seq) == 5 {
		if part1 {
			out := runsequence(program, seq)
			if out > maxout {
				maxout = out
			}
		} else {
			fmt.Printf("%v\n", seq)
			out := runsequence(program, seq)
			if out > maxout {
				maxout = out
			}
		}
		return
	}

	for i := range set {
		if seen[i] {
			continue
		}
		seen[i] = true
		enum(set, seen, append(seq, set[i]))
		seen[i] = false
	}
}

func main() {
	fmt.Printf("hello\n")
	program = readprog("07.txt")

	//fmt.Printf("%d\n", runsequence(program, []int{ 9,8,7,6,5 }))

	set := []int{5, 6, 7, 8, 9}
	seen := make([]bool, len(set))
	enum(set, seen, make([]int, 0, 6))
	fmt.Printf("PART 2: %d\n", maxout)
}
