package main

import (
	"fmt"
	"io"
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

func splitany(in, sep string, n int) []string {
	out := []string{}
	for len(in) > 0 {
		idx := strings.IndexAny(in, sep)
		if idx < 0 {
			out = append(out, in)
			break
		}
		if in[:idx] != "" {
			out = append(out, in[:idx])
		}
		in = in[idx+1:]
		if n > 0 && len(out) == n-1 {
			out = append(out, in)
			break
		}
	}
	return out
}

type Patch struct {
	idx    int
	sym    string
	neg    bool
	lineno int
}

const debug = false

var Program = []int{}
var SymTab = map[string]int{}
var PatchTab = []Patch{}

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

func symbolicNameToOpcode(name string) (int, Opcode) {
	name = strings.ToUpper(name)
	if name == "ADDBAS" {
		name = "ADDBASE"
	}
	for i := range Opcodes {
		if Opcodes[i].Name == name {
			return i, Opcodes[i]
		}
	}
	return -1, Opcode{}
}

func printprogram(out io.Writer) {
	for i := range Program {
		fmt.Fprintf(out, "%d", Program[i])
		if i != len(Program)-1 {
			fmt.Fprintf(out, ",")
		}
	}
	fmt.Fprintf(out, "\n")
}

func ass1(op int, opcode Opcode, args []string, fail func(reason string), lineno int) {
}

func main() {
	infile := os.Args[1]
	b, err := ioutil.ReadFile(infile)
	must(err)

	out := os.Stdout
	if len(os.Args) > 2 {
		out, err = os.Create(os.Args[2])
		must(err)
		defer out.Close()
	}

	usestack := false
	localnum, argnum := 0, 0

	for lineno, line := range strings.Split(string(b), "\n") {
		lineClean := line
		if cmt := strings.Index(lineClean, "//"); cmt >= 0 {
			lineClean = lineClean[:cmt]
		}

		lineno++ // 1-based lines in output
		fields := splitany(lineClean, " \t", 2)
		if len(fields) == 0 {
			continue
		}

		fail := func(reason string) {
			fmt.Fprintf(os.Stderr, "%s:%d: %s %q\n", infile, lineno, reason, line)
			os.Exit(1)
		}

		symdef := func(name string, val int) {
			if _, ok := SymTab[name]; ok {
				fail("symbol already defined")
			}

			SymTab[name] = val
		}

		ass1 := func(op int, opcode Opcode, args []string) {
			if opcode.Len != len(args)+1 {
				fail("wrong number of arguments")
			}

			const (
				basePrefix = "[base"
				fbPrefix   = "[fb"
			)

			modev := make([]int, 3)
			for i, arg := range args {
				if arg[0] == '[' {
					if strings.HasPrefix(arg, basePrefix) {
						modev[i] = 2
					} else if strings.HasPrefix(arg, fbPrefix) {
						if !usestack {
							fail("fb syntax used without usestack")
						}
						modev[i] = 2
					} else {
						modev[i] = 0
					}
				} else {
					modev[i] = 1
				}
			}

			mode := modev[0] + (10 * modev[1]) + (100 * modev[2])
			Program = append(Program, op+(mode*100))

			for _, arg := range args {
				nomodearg := ""
				isfb := false
				if arg[0] == '[' {
					if strings.HasPrefix(arg, basePrefix) {
						nomodearg = arg[len(basePrefix):]
						nomodearg = nomodearg[:len(nomodearg)-1]
					} else if strings.HasPrefix(arg, fbPrefix) {
						nomodearg = arg[len(fbPrefix):]
						nomodearg = nomodearg[:len(nomodearg)-1]
						isfb = true
					} else {
						nomodearg = arg[1 : len(arg)-1]
					}
				} else {
					nomodearg = arg
				}
				n, err := strconv.Atoi(nomodearg)
				if err == nil {
					if isfb {
						if n < 0 {
							if -n > argnum {
								fail("fb syntax invalid, accessing undeclared argument")
							}
						} else {
							if n > localnum {
								fail("fb syntax invalid, accessing undeclared local variable")
							}
						}
						n -= localnum
					}
					Program = append(Program, n)
					continue
				}
				if isfb {
					fail("fb syntax without a number")
				}
				neg := false
				switch nomodearg[0] {
				case '+':
					nomodearg = nomodearg[1:]
				case '-':
					neg = true
					nomodearg = nomodearg[1:]
				}
				PatchTab = append(PatchTab, Patch{
					idx:    len(Program),
					neg:    neg,
					sym:    nomodearg,
					lineno: lineno,
				})
				Program = append(Program, 0)
			}
		}

		assSym := func(name string, args []string) {
			op, opcode := symbolicNameToOpcode(name)
			ass1(op, opcode, args)
		}

		if fields[0][len(fields[0])-1] == ':' {
			name := fields[0][:len(fields[0])-1]
			if debug {
				fmt.Printf("%s:%d: defining label %s at %d\n", infile, lineno, name, len(Program))
			}
			symdef(name, len(Program))
			if len(fields) == 1 {
				continue
			}
			fields = splitany(fields[1], " \t", 2)
		}
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case ".const":
			args := splitany(fields[1], " \t", -1)
			if len(args) != 2 {
				fail("wrong number of arguments")
			}
			name := args[0]
			val, err := strconv.Atoi(args[1])
			if err != nil {
				fail("non-numeric second argument")
			}
			if debug {
				fmt.Printf("%s:%d: defining symbolic constant %s value %d\n", infile, lineno, name, val)
			}
			symdef(name, val)

		case ".ord":
			n, err := strconv.Atoi(fields[1])
			if err != nil {
				fail("non-numeric argument")
			}
			if n < len(Program) {
				fail(fmt.Sprintf("program is already %d ints, .ord can't go backwards", len(Program)))
			}
			for len(Program) < n {
				Program = append(Program, 0)
			}

		case ".var":
			vars := splitany(fields[1], ",", -1)
			for i := range vars {
				vars[i] = strings.TrimSpace(vars[i])
			}
			for i := range vars {
				if debug {
					fmt.Printf("%s:%d: defining symbolic variable %s at %d\n", infile, lineno, vars[i], len(Program))
				}
				symdef(vars[i], len(Program))
				Program = append(Program, 0)
			}

		case ".arr":
			args := splitany(fields[1], " \t", -1)
			if len(args) != 2 {
				fail("wrong number of arguments")
			}
			name := args[0]
			n, err := strconv.Atoi(args[1])
			if err != nil {
				var ok bool
				n, ok = SymTab[args[1]]
				if !ok {
					fail("non-numeric argument")
				}
			}
			if n < 0 {
				fail("negative length")
			}
			if debug {
				fmt.Printf("%s:%d: defining symbolic variable %s at %d (length %d)\n", infile, lineno, name, len(Program), n)
			}
			symdef(name, len(Program))
			for i := 0; i < n; i++ {
				Program = append(Program, 0)
			}

		case ".usestack":
			usestack = true
			assSym("ADDBASE", []string{"stack"})

		case "push":
			assSym("ADDBASE", []string{"1"})
			localnum++
			assSym("ADD", []string{fields[1], "0", "[base+0]"})

		case "pop":
			assSym("ADD", []string{"[base+0]", "0", fields[1]})
			assSym("ADDBASE", []string{"-1"})
			localnum--

		case "call":
			if len(fields) == 1 {
				fail("wrong number of arguments")
			}
			args := splitany(fields[1], ",", -1)
			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}
			if len(args) != 2 {
				fail("wrong number of arguments")
			}
			an, err := strconv.Atoi(args[1])
			if err != nil {
				fail("non-numeric argument")
			}
			assSym("ADDBASE", []string{"1"})
			localnum++
			assSym("ADD", []string{fmt.Sprintf("%d", len(Program)+7), "0", "[base+0]"})
			assSym("JZ", []string{"0", args[0]})
			assSym("ADDBASE", []string{fmt.Sprintf("-%d", an+1)})
			localnum -= an + 1

		case "frame":
			if len(fields) == 1 {
				fail("wrong number of arguments")
			}
			args := splitany(fields[1], ",", -1)
			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}
			if len(args) != 2 {
				fail("wrong number of arguments")
			}
			argnum, err = strconv.Atoi(args[0])
			if err != nil {
				fail("non-numeric argument")
			}
			localnum, err = strconv.Atoi(args[1])
			if err != nil {
				fail("non-numeric argument")
			}
			assSym("ADDBASE", []string{fmt.Sprintf("%d", localnum)})

		case "ret":
			assSym("ADDBASE", []string{fmt.Sprintf("%d", -localnum)})
			assSym("JZ", []string{"0", "[base+0]"})

		default:
			op, opcode := symbolicNameToOpcode(fields[0])
			if opcode.Name == "" {
				fail("unknown opcode")
			}

			var args []string
			if len(fields) > 1 {
				args = splitany(fields[1], ",", -1)
			}

			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}

			if opcode.Name == "ADDBASE" && usestack {
				fail("addbase not allowed when usestack is active")
			}

			ass1(op, opcode, args)
		}
	}

	if debug {
		fmt.Printf("before patches:\n")
		printprogram(os.Stdout)
	}

	if usestack {
		if _, ok := SymTab["result"]; ok {
			fmt.Fprintf(os.Stderr, "result redefined in usestack mode\n")
		}
		SymTab["result"] = len(Program)
		Program = append(Program, 0)
		if _, ok := SymTab["stack"]; ok {
			fmt.Fprintf(os.Stderr, "stack redefined in usestack mode\n")
		}
		SymTab["stack"] = len(Program)
	}

	for _, patch := range PatchTab {
		val, ok := SymTab[patch.sym]
		if !ok {
			fmt.Fprintf(os.Stderr, "%s:%d: symbol %s not defined\n", infile, patch.lineno, patch.sym)
			os.Exit(1)
		}
		if patch.neg {
			val = -val
		}
		Program[patch.idx] = val
	}

	if debug {
		fmt.Printf("\n\n")
	}
	printprogram(out)
}
