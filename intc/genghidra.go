package main

import (
	"fmt"
)

type Opcode struct {
	Len  int
	Name string
	Num int
}

var Opcodes = []Opcode{
	{4, "ADD", 1},
	{4, "MUL", 2},
	{2, "IN", 3},
	{2, "OUT", 4},
	{3, "JNZ", 5},
	{3, "JZ", 6},
	{4, "LT", 7},
	{4, "EQ", 8},
	{2, "ADDBASE", 9},
	{1, "END", 99},
}

func incmode(modev []int) bool {
	for i := 0; i < len(modev); i++ {
		modev[i]++
		if modev[i] > 2 {
			modev[i] = 0
		} else {
			return true
		}
	}
	return false
}

func main() {
	fmt.Printf(`define endian=little;
define alignment=8;
define space ram type=ram_space size=4 default wordsize=8;
define space register type=register_space size=8;
define register offset=0 size=8 [BAS];

define token instr(64)
	opcode = (0, 62)
;

define token arg(64)
	arg1 = (0, 63) signed dec
	arg2 = (0, 63) signed dec
	arg3 = (0, 63) signed dec
;

dest: abs is arg2 [ abs = arg2 + 0; ] {
	export *[ram]:8 abs;
}

mem1: abs is arg1 [ abs = arg1 + 0; ] {
	export *[ram]:8 abs;
}

mem2 : abs is arg2 [ abs = arg2 + 0; ] {
	export *[ram]:8 abs;
}

mem3 : abs is arg3 [ abs = arg3 + 0; ] {
	export *[ram]:8 abs;
}

`)
	for _, opcode := range Opcodes {
		modev := make([]int, opcode.Len-1)
		for {
			mode := 0
			for i := len(modev)-1; i >= 0; i-- {
				mode = mode * 10
				mode += modev[i]
			}
			opn := mode*100 + opcode.Num
			//fmt.Printf("%s %v %d\n", opcode.Name, modev, opn)
			
			fmt.Printf(":%s ", opcode.Name)
			
			for i := range modev {
				if i != 0 {
					fmt.Printf(", ")
				}
				switch modev[i] {
				case 0:
					//fmt.Printf(`"["^arg%d^"]"`, i+1)
					fmt.Printf("mem%d", i+1)
				case 1:
					if opcode.Name[0] == 'J' && i == 1 {
						fmt.Printf("dest")
					} else {
						fmt.Printf("arg%d", i+1)
					}
				case 2:
					fmt.Printf(`"[BAS+"^arg%d^"]"`, i+1)
				}
			}
			
			fmt.Printf(" is opcode=%d", opn)
			
			for i := range modev {
				if opcode.Name[0] == 'J' && i == 1 && modev[i] == 1 {
					fmt.Printf(" ; dest ")
				} else if modev[i] == 0 {
					fmt.Printf(" ; mem%d", i+1)
				} else {
					fmt.Printf(" ; arg%d", i+1) 
				}
			}
			
			fmt.Printf(" {\n")
			
			fmt.Printf("}\n\n")
			
			if !incmode(modev) {
				break
			}
		}
	}
}
