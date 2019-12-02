package main

import "fmt"

func attempt(in2 int) (out bool) {
	defer func() {
		recover()
	}()
	x := 19690720

	deadd := func(y int) {
		fmt.Printf("%d - %d -> %d\n", x, y, x-y)
		x -= y
	}

	demul := func(y int) {
		if x%y != 0 {
			fmt.Printf("no multiplo %d\n", y)
			panic("")
		}
		fmt.Printf("%d / %d -> %d\n", x, y, x/y)

		x = x / y
	}

	x -= 3
	fmt.Printf("Trying in2 == %d\n", in2)
	deadd(in2)
	demul(5)
	deadd(3)
	demul(5)
	demul(5)
	deadd(9)
	demul(2)
	deadd(3)
	demul(5)
	deadd(1)
	demul(3)
	deadd(6)
	demul(4)
	deadd(1)
	demul(5)
	deadd(2)
	demul(4)
	fmt.Printf("successo %d\n", x)
	out = true
	return
}

func main() {
	in2 := 0
	for {
		if attempt(in2) {
			break
		}
		in2++
	}
}
