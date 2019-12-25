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

var TRACECPU = false

var seen = map[int]bool{}
var newinstr = map[int]bool{}

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
		if !seen[s.pc] {
			newinstr[s.pc] = true
		}
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

type item struct {
	name string
	path string
}

var items = []item{
	//{ "infinite loop", "south" },
	{"fixed point", "west"},
	{"sand", "west north"},
	{"asterisk", "east"},
	//{ "escape pod", "east east south"},
	//{ "photons","east east south east" },
	//{ "giant electromagnet", "east north" },
	{"hypercube", "east north north"},
	{"coin", "east north north north"},
	{"easter egg", "east north north north north"},
	{"spool of cat6", "east north west north"},
	//{ "molten lava", "east north west south" },
	{"shell", "east north west north north"},
}

func invstep(step string) string {
	switch step {
	case "north":
		return "south"
	case "west":
		return "east"
	case "east":
		return "west"
	case "south":
		return "north"
	}
	panic("wtf")
}

func main() {
	p := readprog("25.txt")
	s := newCpustate(p)
	inchan := make(chan int)
	s.input = func() int {
		return <-inchan
	}
	outbuf := []byte{}
	s.output = func(n int) {
		TRACECPU = false
		if len(newinstr) > 0 {
			/*fmt.Printf("new instructions: ")
			for pc := range newinstr {
				fmt.Printf("%d ", pc)
			}
			fmt.Printf("\n")*/
			for pc := range newinstr {
				seen[pc] = true
				delete(newinstr, pc)
			}
		}
		outbuf = append(outbuf, byte(n))
		fmt.Printf("%c", n)
	}
	go func() {
		cpu(s, 0, false)
		fmt.Printf("FINISHED\n")
		exit(0)
	}()

	sendline := func(txt string) {
		outbuf = outbuf[:0]
		for i := range txt {
			inchan <- int(txt[i])
		}
		inchan <- int('\n')
	}

	getitem := func(it item) {
		pathv := strings.Split(it.path, " ")
		for _, step := range pathv {
			sendline(step)
		}
		sendline(fmt.Sprintf("take %s", it.name))
		for i := len(pathv) - 1; i >= 0; i-- {
			step := invstep(pathv[i])
			sendline(step)
		}
		sendline(fmt.Sprintf("drop %s", it.name))
	}

	trycnt := 0
	trypart2 := false
	didgetall := false
	var ittoget []item

	buf := make([]byte, 100)
intloop:
	for {
		n, err := os.Stdin.Read(buf)
		must(err)
		in := strings.TrimSpace(string(buf[:n]))
		v := strings.SplitN(in, " ", 2)
		switch v[0] {
		case "trace":
			TRACECPU = true
			continue intloop
		case "get":
			var theitem item
			found := false
			for _, item := range items {
				if item.name == v[1] {
					fmt.Printf("want %s\n", item.name)
					theitem = item
					found = true
					break
				}
			}
			if found {
				fmt.Printf("path %s\n", theitem.path)
				getitem(theitem)
			}
			continue intloop
		case "getall":
			didgetall = true
			for _, it := range items {
				getitem(it)
			}
			continue intloop
		case "":
			fmt.Printf("auto\n")
			if trypart2 {
				pathv := strings.Split("east north west north north west north", " ")
				for i := len(pathv) - 1; i >= 0; i-- {
					step := pathv[i]
					sendline(invstep(step))
				}
				for i := range ittoget {
					sendline(fmt.Sprintf("drop %s", ittoget[i].name))
				}
				trypart2 = false

				continue intloop
			}
			if didgetall {
				x := trycnt
				ittoget = []item{}
				for i := range items {
					if x&1 != 0 {
						ittoget = append(ittoget, items[i])
					}
					x = x >> 1
				}
				fmt.Printf("getting: ")
				for i := range ittoget {
					fmt.Printf(" %s", ittoget[i].name)
				}
				fmt.Printf("\n")
				trycnt++

				for i := range ittoget {
					sendline(fmt.Sprintf("take %s", ittoget[i].name))
				}

				for _, step := range strings.Split("east north west north north west north north", " ") {
					sendline(step)
				}

				trypart2 = true
			}
			continue intloop
		}
		for i := range buf[:n] {
			inchan <- int(buf[i])
		}
	}
}
