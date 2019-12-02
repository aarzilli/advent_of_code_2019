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

var program []int

const part1 = false

func main() {
	fmt.Printf("hello\n")
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
	if part1 {
		program[1] = 12
		program[2] = 2
	}
	for pc < len(program) {
		opcode := program[pc]
		if opcode == 99 {
			fmt.Printf("%d\t%d\n", pc, opcode)
			for pc < len(program) {
				fmt.Printf("%d\t%d\n", pc, program[pc])
				pc++
			}
			break
		}
		arg1 := program[pc+1]
		arg2 := program[pc+2]
		dest := program[pc+3]
		pc += 4
		fmt.Printf("%d\t%d,%d,%d,%d ", pc-4, opcode, arg1, arg2, dest)
		switch opcode {
		case 1:
			fmt.Printf("[%d] <- [%d] + [%d]\n", dest, arg1, arg2)
		case 2:
			fmt.Printf("[%d] <- [%d] * [%d]\n", dest, arg1, arg2)
		default:
		}
	}
}
