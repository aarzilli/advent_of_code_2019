package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
			in := input()
			save(0, in)
			if TRACECPU {
				fmt.Printf("\tinput was %d\n", in)
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

const (
	North = 1
	South = 2
	West  = 3
	East  = 4
)

type Point struct {
	i, j int
}

var M = map[Point]byte{}
var Pos Point

func bounds() (mini, maxi, minj, maxj int) {
	for p := range M {
		if p.i < mini {
			mini = p.i
		}
		if p.i > maxi {
			maxi = p.i
		}
		if p.j < minj {
			minj = p.j
		}
		if p.j > maxj {
			maxj = p.j
		}
	}
	return
}

const debug = true
const debugclear = true
const debugsleep = true
const debug2 = true

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func dist(p1, p2 Point) int {
	return abs(p1.i-p2.i) + abs(p1.j-p2.j)
}

func explore() map[Point]int {
	ok, dest := findexploretgt(Pos)
	if !ok {
		return nil
	}

	if debug {
		fmt.Printf("new destination %v\n", dest)
	}

	return calcsteps(dest, Pos)
}

func calcsteps(start, end Point) map[Point]int {
	S := make(map[Point]int)

	queue := []Point{start}
	S[start] = 0

	for {
		if len(queue) == 0 {
			return nil
		}
		p := queue[0]
		queue = queue[1:]

		if p == end {
			break
		}

		add := func(p2 Point) {
			if M[p2] == '#' {
				return
			}
			if _, ok := S[p2]; ok {
				if S[p2] <= S[p]+1 {
					return
				}
			}
			S[p2] = S[p] + 1
			queue = append(queue, p2)
		}

		add(Point{p.i - 1, p.j})
		add(Point{p.i, p.j - 1})
		add(Point{p.i, p.j + 1})
		add(Point{p.i + 1, p.j})
	}

	return S
}

func findexploretgt(start Point) (bool, Point) {
	seen := make(map[Point]bool)
	queue := []Point{start}

	for {
		if len(queue) == 0 {
			return false, Point{}
		}
		p := queue[0]
		queue = queue[1:]

		if _, ok := M[p]; !ok {
			return true, p
		}

		add := func(p2 Point) {
			if M[p2] == '#' {
				return
			}
			if seen[p2] {
				return
			}
			seen[p2] = true
			queue = append(queue, p2)
		}

		add(Point{p.i - 1, p.j})
		add(Point{p.i, p.j - 1})
		add(Point{p.i, p.j + 1})
		add(Point{p.i + 1, p.j})
	}
}

var lastw, lasth int

func showmap(w io.Writer, showcursor bool) {
	mini, maxi, minj, maxj := bounds()

	wd := maxj - minj
	ht := maxi - mini

	if debug && debugclear && (wd > lastw || ht > lasth) {
		lastw = wd
		lasth = ht
		fmt.Printf("\x1b[2J\x1b[H")
	}

	for i := mini; i <= maxi; i++ {
		for j := minj; j <= maxj; j++ {
			p := Point{i, j}
			if p == Pos && showcursor {
				fmt.Fprintf(w, "X")
			} else {
				c := M[p]
				if c == 0 {
					c = ' '
				}
				fmt.Fprintf(w, "%c", c)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func explorer(input chan<- int, output <-chan int, cont chan<- bool) Point {
	var oxygenpos Point

	M[Pos] = '.'
	S := explore()
	for {
		if debug {
			if debugclear {
				//fmt.Printf("\x1b[2J\x1b[H")
				fmt.Printf("\x1b[H")
			}
			showmap(os.Stdout, true)
			if debugclear {
				fmt.Printf("\x1b[0J")
			}
			if debugsleep {
				time.Sleep(20 * time.Millisecond)
			}
		}

		// Select a direction to go into
		var dir int
		var nextpos Point
		for i, np := range []Point{{Pos.i - 1, Pos.j}, {Pos.i + 1, Pos.j}, {Pos.i, Pos.j - 1}, {Pos.i, Pos.j + 1}} {
			nps, ok := S[np]
			if !ok {
				continue
			}
			if nps < S[Pos] {
				nextpos = np
				dir = i + 1
				break
			}
		}

		// Go to that direction
		if debug {
			fmt.Printf("at %v going %d to %v\n", Pos, dir, nextpos)
		}
		cont <- true
		input <- dir

		// See wtf happened
		status := <-output

		switch status {
		case 0:
			if M[nextpos] == '.' {
				fmt.Printf("desync")
				exit(1)
			}
			M[nextpos] = '#'
			if debug {
				fmt.Printf("hit wall\n")
			}
			S = explore()
			if S == nil {
				return oxygenpos
			}

		case 1:
			M[nextpos] = '.'
			Pos = nextpos

		case 2:
			if debug2 {
				fmt.Printf("goal found at %v\n", nextpos)
			}
			M[nextpos] = '2'
			Pos = nextpos
			oxygenpos = nextpos
		}

		if S[Pos] == 0 {
			S = explore()
		}
	}
}

func main() {
	p := readprog("15.txt")
	input := make(chan int)
	output := make(chan int)
	cont := make(chan bool)
	go func() {
		<-cont
		cpu(p, func() int {
			return <-input
		}, func(n int) {
			output <- n
			<-cont
		})
	}()

	if debug && debugclear {
		fmt.Printf("\x1b[2J")
	}

	dest := explorer(input, output, cont)

	S := calcsteps(dest, Point{0, 0})
	if S == nil {
		panic("fuck?")
	}
	fmt.Printf("PART 1: %d\n", S[Point{0, 0}])

	// Part 2
	fh, err := os.Create("15.addenda")
	must(err)
	showmap(fh, false)
	fh.Close()

	cmd := exec.Command("go", "run", "15p2.go")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
