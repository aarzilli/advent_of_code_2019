package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
	"encoding/binary"
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

func main() {
	p := readprog(os.Args[1])
	fmt.Printf("%d\n", p)
	fh, err := os.Create(os.Args[2])
	must(err)
	for i := range p {
		must(binary.Write(fh, binary.LittleEndian, int64(p[i])))
	}
	must(fh.Close())
}
