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

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

var basepattern = []int{0, 1, 0, -1}

var memoize = [][]int{}

func digit(i int, phase int, in, in2 []int) int {
	if phase == 0 {
		return in[i]
	}

	patternc := i
	patterni := 0

	r := 0
	for j := 0; j < len(in); {
		if patternc <= 0 {
			patternc = i + 1
			patterni = (patterni + 1) % len(basepattern)
		}

		patternx := basepattern[patterni]

		switch patternx {
		case 0:
			j += patternc
			patternc = 0
		case 1, -1:
			first := true
			t0 := time.Now()
			for ; (patternc > 0) && (j < len(in)); patternc, j = patternc-1, j+1 {
				if memoize[phase-1][j] >= 0 {
					r += patternx * memoize[phase-1][j]
				} else {
					if first {
						if phase > 1 {
							fmt.Printf("calculating digit %d of phase %d\n", i, phase)
						}
						fmt.Printf("\twill need digits from %d to %d (n=%d) of phase %d\n", j, min(j+patternc, len(in)), min(j+patternc, len(in))-j, phase-1)
						first = false
					}
					if j%1000 == 0 {
						fmt.Printf("\t\tat %d (expired %v)\n", j, time.Since(t0))
						t0 = time.Now()
					}
					r += patternx * digit(j, phase-1, in, in2)
				}
			}
		}
	}

	r = abs(r)
	r = r % 10
	memoize[phase][i] = r
	return r
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

	memoize = make([][]int, 101)
	for i := range memoize {
		memoize[i] = make([]int, len(in))
		for j := range memoize[i] {
			memoize[i][j] = -1
		}
	}

	for i := range in {
		memoize[0][i] = in[i]
	}

	/*
		in2 := step(in)
	*/

	out := make([]int, 8)
	for i := range out {
		out[i] = digit(i+msgoff, 100, in, nil)
	}

	fmt.Printf("OUTPUT: ")
	for i := range out {
		fmt.Printf("%d", out[i])
	}
	fmt.Printf("\n")

	/*
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
		fmt.Printf("\n")*/
}
