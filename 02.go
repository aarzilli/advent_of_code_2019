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

var program []int

const debug = false

func main() {
	buf, err := ioutil.ReadFile("02.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		program = append(program, vatoi(splitandclean(line, ",", -1))...)
	}
	pc := 0
	program[1] = 12
	program[2] = 2
evalLoop:
	for pc < len(program) {
		opcode := program[pc]
		arg1 := program[pc+1]
		arg2 := program[pc+2]
		dest := program[pc+3]
		pc += 4
		switch opcode {
		case 1:
			program[dest] = program[arg1] + program[arg2]
			if debug {
				fmt.Printf("%d %d is now %d\n", pc-4, dest, program[dest])
			}
		case 2:
			program[dest] = program[arg1] * program[arg2]
			if debug {
				fmt.Printf("%d %d is now %d\n", pc-4, dest, program[dest])
			}
		default:
			break evalLoop
		}
		/*if debug {
			fmt.Printf("%d -> %v\n", pc-4, program)
		}*/
	}
	fmt.Printf("Part 1: %d\n", program[0])
}
