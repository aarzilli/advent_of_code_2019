package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
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

const debug = false
const debug3 = false

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

var basepattern = []int{0, 1, 0, -1}

var memoize = [][]int{}

func digit(i int, phase int, in []int) int {
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
			compute(phase, i, patternc, j, in)
			end := min(len(in), j+patternc)

			r += patternx * sum(phase-1, j, end)

			j = end

			patternc = 0
		}
	}

	r = abs(r)
	r = r % 10
	memoize[phase][i] = r
	return r
}

func sum(phase int, start, end int) int {
	r := 0
	for i := start; i < end; {
		done := false
		for _, factor := range []int{100000, 10000, 1000, 100} {
			if i%factor != 0 {
				continue
			}
			if i+factor >= end {
				continue
			}
			r += memosum(phase, factor, i)
			i += factor
			done = true
			break
		}
		if !done {
			r += memoize[phase][i]
			i++
		}
	}
	return r
}

type memosumIn struct {
	phase, factor, start int
}

var memosumMap = map[memosumIn]int{}

func memosum(phase, factor, start int) int {
	k := memosumIn{phase, factor, start}
	if r, ok := memosumMap[k]; ok {
		return r
	}

	r := sum(phase, start, start+factor)

	memosumMap[k] = r
	return r
}

func compute(phase, i, patternc, j int, in []int) {
	if phase == 1 {
		return
	}
	first := true
	t0 := time.Now()
	for ; (patternc > 0) && (j < len(in)); patternc, j = patternc-1, j+1 {
		if memoize[phase-1][j] < 0 {
			if first {
				fmt.Printf("calculating digit %d of phase %d\n", i, phase)
				fmt.Printf("\twill need digits from %d to %d (n=%d) of phase %d\n", j, min(j+patternc, len(in)), min(j+patternc, len(in))-j, phase-1)
				first = false
			}
			if j%1000 == 0 {
				fmt.Printf("\t\tdigit %d of phase %d (remaining %d) at %d (expired %v)\n", i, phase, min(patternc, len(in)-j), j, time.Since(t0))
				t0 = time.Now()
			}
			digit(j, phase-1, in)
		}
	}
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

	out := make([]int, 8)
	for i := range out {
		out[i] = digit(i+msgoff, 100, in)
	}

	fmt.Printf("OUTPUT: ")
	for i := range out {
		fmt.Printf("%d", out[i])
	}
	fmt.Printf("\n")
}
