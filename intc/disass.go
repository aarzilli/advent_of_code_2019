package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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
	3:  {2, "IN"},
	4:  {2, "OUT"},
	5:  {3, "JNZ"},
	6:  {3, "JZ"},
	7:  {4, "LT"},
	8:  {4, "EQ"},
	9:  {2, "ADDBASE"},
	99: {1, "END"},
}

func prettyInstr(w io.Writer, p []int, pc, mode, opcode int, a []int, dirty map[int]bool, label map[int]int) {
	if label != nil && label[pc] > 0 {
		fmt.Fprintf(w, "L%d:\t", label[pc])
	} else {
		fmt.Fprintf(w, "\t")
	}
	fmt.Fprintf(w, "%04d\t%d", pc, opcode+(mode*100))

	for i := range a {
		fmt.Fprintf(w, ",")
		fmt.Fprintf(w, "%d", a[i])
	}

	fmt.Fprintf(w, "\t")

	oc, ok := Opcodes[opcode]
	if !ok {
		fmt.Fprintf(w, "?\n")
		return
	}
	symop := oc.Name

	modev := make([]int, 3)
	modev[0] = mode % 10
	modev[1] = (mode / 10) % 10
	modev[2] = (mode / 100) % 10

	fmt.Fprintf(w, "%03d\t%s\t", mode, symop)
	for i := range a {
		switch modev[i] {
		case 0:
			if a[i] < len(p) {
				fmt.Fprintf(w, "[%d]=%d", a[i], p[a[i]])
				if dirty != nil && dirty[a[i]] {
					fmt.Fprintf(w, "!")
				}
			} else {
				fmt.Fprintf(w, "[%d]", a[i])
			}
		case 1:
			if symop[0] == 'J' && label != nil && label[a[i]] > 0 {
				fmt.Fprintf(w, "L%d", label[a[i]])
			} else {
				fmt.Fprintf(w, "%d", a[i])
			}
		case 2:
			fmt.Fprintf(w, "[BASE%+d]", a[i])
		}
		if i != len(a)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, "\n")
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

type prettyInstrArgs struct {
	pc, mode, opcode int
	a                []int
}

func pretty(w io.Writer, p []int, start int) {
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

	dirty := make(map[int]bool)
	label := make(map[int]int)
	labelcnt := 1

	instrs := []prettyInstrArgs{}

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
		instrs = append(instrs, prettyInstrArgs{
			pc:     pc,
			mode:   mode,
			opcode: opcode,
			a:      a,
		})

		switch opcode {
		case 1, 2, 7, 8: // add, mul, lt, eq
			dirty[a[2]] = true
		case 3: // input
			dirty[a[0]] = true
		case 5, 6:
			dstmode := (mode / 10) % 10
			if dstmode == 1 {
				label[a[1]] = labelcnt
				labelcnt++
			}
		}
		pc += n
	}

	fmt.Fprintf(w, "LBL\tPC\tINTS\tMODE\tOPCODE\tARGS\n")

	for _, pia := range instrs {
		prettyInstr(w, p, pia.pc, pia.mode, pia.opcode, pia.a, dirty, label)
	}

}

func main() {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	defer w.Flush()
	p := readprog(os.Args[1])
	pretty(w, p, 0)
}
