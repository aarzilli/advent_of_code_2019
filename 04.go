package main

import (
	"fmt"
)

const input = "172930-683082"

func valid(x string) bool {
	doublefound := false
	for i := range x {
		if i+1 < len(x) {
			if x[i] > x[i+1] {
				return false
			}
			if x[i] == x[i+1] {
				doublefound = true
			}
		}
	}

	return doublefound
}

func valid2(x string) bool {
	counts := make(map[byte]int)
	for i := range x {
		if i+1 < len(x) {
			if x[i] > x[i+1] {
				return false
			}
			if x[i] == x[i+1] {
				counts[x[i]]++
			}
		}
	}

	for _, x := range counts {
		if x == 1 {
			return true
		}
	}

	return false
}

const part1test = false
const testexamples = false

func main() {
	if testexamples {
		for _, n := range []int{111111, 223450, 123789, 112233, 123444, 111122} {
			if part1test {
				fmt.Printf("%d %v\n", n, valid(fmt.Sprintf("%d", n)))
			} else {
				fmt.Printf("%d %v\n", n, valid2(fmt.Sprintf("%d", n)))
			}
		}
	}

	cnt1, cnt2 := 0, 0
	for i := 172930; i <= 683082; i++ {
		x := fmt.Sprintf("%d", i)
		if valid(x) {
			cnt1++
		}
		if valid2(x) {
			cnt2++
		}
	}
	fmt.Printf("PART 1: %d\n", cnt1)
	fmt.Printf("PART 2: %d\n", cnt2)
}
