package main

import (
	"io/ioutil"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
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
	9:  {2, "REL"},
	99: {1, "END"},
}
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

func main() {
	p := readprog(os.Argv[1])
}
