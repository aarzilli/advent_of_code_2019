// USE 16t4.go for PART 2 this one will take days to finish

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
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

func exit(n int) {
	os.Exit(n)
}

func newpattern(n int) []int {
	n++
	pattern := make([]int, 4*n)
	for k := 0; k < n; k++ {
		pattern[k] = 0
	}
	for k := 0; k < n; k++ {
		pattern[n+k] = 1
	}
	for k := 0; k < n; k++ {
		pattern[2*n+k] = 0
	}
	for k := 0; k < n; k++ {
		pattern[3*n+k] = -1
	}
	return pattern
}

const debug = false
const debug2 = true
const debug3 = false

func step(in []int) []int {
	N := len(in)
	out := make([]int, N)
	t0 := time.Now()
	for i := 0; i < N; i++ {
		if debug2 {
			if i%1000 == 0 {
				fmt.Printf("calculating %d (expired %v)\n", i, time.Since(t0))
				t0 = time.Now()
			}
		}
		r := 0
		pattern := newpattern(i)
		for j := range in {
			if debug {
				fmt.Printf("%d*%d + ", in[j], pattern[(j+1)%len(pattern)])
			}
			r += in[j] * pattern[(j+1)%len(pattern)]
		}
		r = abs(r)
		r = r % 10

		if debug {
			fmt.Printf("= %d\n", r)
		}

		out[i] = r
	}
	return out
}

const part2 = true

func main() {
	buf, err := ioutil.ReadFile("16.txt")
	must(err)
	var in []int
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		in = make([]int, len(line))
		for i := range line {
			in[i] = int(line[i] - '0')
		}
	}

	if part2 {
		in2 := make([]int, 0, len(in)*10000)
		for i := 0; i < 10000; i++ {
			in2 = append(in2, in...)
		}
		fmt.Printf("%d %d\n", len(in), len(in2))
		in = in2
	}

	msgoff := 0

	if part2 {
		for i := 0; i < 7; i++ {
			msgoff *= 10
			msgoff += in[i]
		}
		if debug2 {
			fmt.Printf("msgoff = %d\n", msgoff)
		}
	}

	if debug {
		fmt.Printf("%v\n", in)
	}

	phases := 100
	for i := 0; i < phases; i++ {
		if debug2 {
			fmt.Printf("phase %d\n", i+1)
		}
		out := step(in)
		if debug3 {
			fmt.Printf("after phase %d: %v\n", i+1, out)
		}
		in = out
	}

	if part2 {
		fmt.Printf("PART 2: ")
	} else {
		fmt.Printf("PART 1: ")
	}
	for i := 0; i < 8; i++ {
		fmt.Printf("%d", in[msgoff+i])
	}
	fmt.Printf("\n")
}
