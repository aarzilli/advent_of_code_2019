package main

import (
	"fmt"
	"os"
)

type Entry struct {
	counter int
	next int
	mistery int
}

func main() {
	v := make([]Entry, 1000)
	v[0] = Entry{ counter: 27, next: 915, mistery: 0 }
	
	i := 0
	for {
		i++
		fmt.Printf("%v\n", v[:i])
		if v[i-1].counter < 3 {
			weirdpart:
			i--
			switch v[i].next {
			case 942:
				v[i-1].mistery = v[i].counter
				v[i].counter = v[i-1].counter-3
				v[i].next = 957
			case 957:
				v[i-1].counter = v[i].counter + v[i-1].mistery
				goto weirdpart
			case 915:
				fmt.Printf("DONE: %d\n", v[i].counter + 24405)
				os.Exit(0)
			default:
				fmt.Printf("not implemented %d\n", v[i].next)
				os.Exit(1)
			}
		} else {
			v[i].counter = v[i-1].counter-1
			v[i].next = 942
		}
	}
}
