package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var day string
	if len(os.Args) == 2 {
		n, err := strconv.Atoi(os.Args[1])
		must(err)
		day = fmt.Sprintf("%02d", n)
	} else {
		t := time.Now()
		if t.Month() != time.December {
			panic("not december")
		}
		day = fmt.Sprintf("%02d", t.Day())
	}
	p := fmt.Sprintf("%s.example", day)
	cnt := 2
	for {
		_, err := os.Stat(p)
		if err != nil {
			break
		}
		p = fmt.Sprintf("%s.example%d", day, cnt)
		cnt++
	}
	fmt.Printf("WRITING TO: %s\n", p)

	clip, err := exec.Command("xclip", "-o").CombinedOutput()
	must(err)
	fmt.Printf("%s\n", string(clip))

	fh, err := os.Create(p)
	must(err)
	n, err := fh.Write(clip)
	must(err)
	if n != len(clip) {
		panic("short write")
	}
	must(fh.Close())
}
