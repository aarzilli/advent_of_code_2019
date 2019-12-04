package main

import (
	"fmt"
	"runtime"
	"sync"
)

func numconv(x []byte, n int) []byte {
	for n > 0 {
		x = append(x, byte(n%10)+'0')
		n = n / 10
	}
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
	return x
}

func tonum(x []byte) int {
	r := 0
	for i := range x {
		r *= 10
		r += int(x[i] - '0')
	}
	return r
}

const selftest = false

func problem(start, end int) (cnt1, cnt2 int) {
	cnt1, cnt2 = 0, 0
	x := numconv(make([]byte, 0, 10), start)
	n := start
	for n <= end {
		valid1, valid2 := false, false
		{ // validity check
			groupdigit, grouplen := byte(10), 0
			incrementing := true
			for i := range x {
				if i+1 < len(x) {
					if x[i] > x[i+1] {
						incrementing = false
						break
					}
				}
				if x[i] == groupdigit {
					grouplen++
				} else {
					if grouplen == 2 {
						valid2 = true
					}
					if grouplen >= 2 {
						valid1 = true
					}
					groupdigit = x[i]
					grouplen = 1
				}
			}

			if incrementing {
				if grouplen == 2 {
					valid2 = true
				}
				if grouplen >= 2 {
					valid1 = true
				}
			}
		}

		if valid1 {
			cnt1++
		}
		if valid2 {
			cnt2++
		}

		{ // increment x
			ok := false
			for i := len(x) - 1; i >= 0; i-- {
				x[i]++
				if x[i] <= '9' {
					ok = true
					for j := i + 1; j < len(x); j++ {
						x[j] = x[i]
					}
					if i == len(x)-1 {
						n++
					} else {
						n = tonum(x)
					}
					break
				}
				x[i] = '0'
			}
			if !ok {
				x = append(x, 0)
				x[0] = '1'
				for i := 1; i < len(x); i++ {
					x[i] = '0'
				}
				n = tonum(x)
			}
		}

	}
	return cnt1, cnt2
}

func problemParallel(start, end int) (cnt1, cnt2 int) {
	cpu := runtime.NumCPU()

	var mu sync.Mutex
	var wall sync.WaitGroup
	m := (end - start) / cpu

	s0 := start

	for i := 0; i < cpu; i++ {
		e0 := s0 + m
		if i == cpu-1 {
			e0 = end
		}
		wall.Add(1)
		go func(s0, e0 int) {
			c1, c2 := problem(s0, e0)
			mu.Lock()
			cnt1 += c1
			cnt2 += c2
			mu.Unlock()
			wall.Done()
		}(s0, e0)
		s0 += m + 1
	}

	wall.Wait()

	return cnt1, cnt2
}

func main() {
	if selftest {
		cnt1, cnt2 := problemParallel(172930, 683082)
		if cnt1 != 1675 || cnt2 != 1142 {
			fmt.Printf("self test failed %d %d\n", cnt1, cnt2)
			return
		}
	}
	cnt1, cnt2 := problemParallel(100000000, 1000000000)

	fmt.Printf("PART 1: %d\n", cnt1)
	fmt.Printf("PART 2: %d\n", cnt2)
}
