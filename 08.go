package main

import (
	"fmt"
	"io/ioutil"
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

const WIDTH = 25
const HEIGHT = 6

func count(layer []byte, v byte) int {
	cnt := 0
	for i := range layer {
		if layer[i] == v {
			cnt++
		}
	}
	return cnt
}

func main() {
	buf, err := ioutil.ReadFile("08.txt")
	must(err)
	var in []byte
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		in = []byte(line)
	}
	for i := range in {
		in[i] = in[i] - '0'
	}
	layers := [][]byte{}
	var sz int
	if len(in) < 30 {
		sz = 2 * 2
	} else {
		sz = WIDTH * HEIGHT
	}
	rem := in
	for len(rem) > 0 {
		layers = append(layers, rem[:sz])
		rem = rem[sz:]
	}

	minnum0 := 10000
	minnum0i := 0
	for i, layer := range layers {
		num0 := count(layer, 0)
		if num0 < minnum0 {
			minnum0 = num0
			minnum0i = i
		}
	}

	fmt.Printf("PART 1: %d\n", count(layers[minnum0i], 1)*count(layers[minnum0i], 2))

	clayer := make([]byte, sz)

	copy(clayer, layers[0])

	for i := 1; i < len(layers); i++ {
		for j := range clayer {
			if clayer[j] == 2 {
				clayer[j] = layers[i][j]
			}
		}
	}

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			if clayer[i*WIDTH+j] != 0 {
				fmt.Printf("%d", clayer[i*WIDTH+j])
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}
