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

// convert string to integer
func atoi(in string) int {
	n, err := strconv.Atoi(in)
	must(err)
	return n
}

func mass(n int) int {
	return n/3 - 2
}

func totalmass(n int) int {
	t := 0
	x := n
	for {
		x = mass(x)
		if x < 0 {
			break
		}
		t += x
	}
	return t
}

const part2 = true

func main() {
	buf, err := ioutil.ReadFile("01.txt")
	must(err)
	var part1, part2 int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		part1 += mass(atoi(line))
		part2 += totalmass(atoi(line))
	}
	fmt.Printf("PART 1: %d\n", part1)
	fmt.Printf("PART 2: %d\n", part2)
}

// Part 1: 3229279
// Part 2: 4841054
