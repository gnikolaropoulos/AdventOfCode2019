package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var firstSpringscript = "NOT A T\nOR T J\nNOT B T\nOR T J\nNOT C T\nOR T J\nAND D J\nWALK\n"

var secondSpringscript = `NOT A J
NOT B T
OR T J
NOT C T
OR T J
AND D J
NOT E T
NOT T T
OR H T
AND T J
RUN
`

func main() {
	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	fmt.Println("--- Part One ---")
	result, output := execute(program, firstSpringscript)
	if result != 0 {
		fmt.Println(result)
	} else {
		fmt.Println(output)
	}

	fmt.Println("--- Part Two ---")
	result, output = execute(program, secondSpringscript)
	if result != 0 {
		fmt.Println(result)
	} else {
		fmt.Println(output)
	}
}

func execute(program []int64, script string) (int64, string) {
	input := make(chan int64, len(script))
	output := make(chan int64)
	halt := make(chan bool, 1)

	go emulate(program, input, output, halt)

	for _, c := range script {
		input <- int64(c)
	}

	var result int64
	var message strings.Builder

	for {
		select {
		case char := <-output:
			if char >= 128 {
				result = char
			} else {
				message.WriteRune(rune(char))
			}

		case <-halt:
			return result, message.String()
		}
	}
}

func emulate(program []int64, input <-chan int64, output chan<- int64, halt chan<- bool) {
	memory := make([]int64, 5000)
	copy(memory, program)

	var ip, relativeBase int64

	for {
		instruction := memory[ip]
		a := instruction / 10000
		b := (instruction - a*10000) / 1000
		c := (instruction - a*10000 - b*1000) / 100
		opcode := instruction % 100

		switch opcode {
		case 1: // ADD
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			y := fetchPointerToMemory(b, &memory, ip+2, relativeBase)
			z := fetchPointerToMemory(a, &memory, ip+3, relativeBase)
			*z = *x + *y
			ip += 4

		case 2: // MULTIPLY
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			y := fetchPointerToMemory(b, &memory, ip+2, relativeBase)
			z := fetchPointerToMemory(a, &memory, ip+3, relativeBase)
			*z = *x * *y
			ip += 4

		case 3: // INPUT
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			*x = <-input
			ip += 2

		case 4: // OUTPUT
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			output <- *x
			ip += 2

		case 5: // JUMP IF TRUE
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			y := fetchPointerToMemory(b, &memory, ip+2, relativeBase)
			if *x != 0 {
				ip = *y
			} else {
				ip += 3
			}

		case 6: // JUMP IF FALSE
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			y := fetchPointerToMemory(b, &memory, ip+2, relativeBase)
			if *x == 0 {
				ip = *y
			} else {
				ip += 3
			}

		case 7: // LESS THAN
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			y := fetchPointerToMemory(b, &memory, ip+2, relativeBase)
			z := fetchPointerToMemory(a, &memory, ip+3, relativeBase)

			if *x < *y {
				*z = 1
			} else {
				*z = 0
			}

			ip += 4

		case 8: // EQUAL
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			y := fetchPointerToMemory(b, &memory, ip+2, relativeBase)
			z := fetchPointerToMemory(a, &memory, ip+3, relativeBase)
			if *x == *y {
				*z = 1
			} else {
				*z = 0
			}

			ip += 4

		case 9: // ADJUST RELATIVE BASE
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			relativeBase += *x
			ip += 2

		case 99: // HALT
			halt <- true
			return
		default:
			panic(fmt.Sprintf("error: invalid opcode: ip=%d, instruction=%d, opcode=%d", ip, instruction, opcode))
		}
	}
}

func fetchPointerToMemory(mode int64, memory *[]int64, position int64, relativeBase int64) *int64 {
	switch mode {
	case 0:
		index := (*memory)[position]
		return &(*memory)[index]
	case 1:
		return &(*memory)[position]
	case 2:
		index := (*memory)[position] + relativeBase
		return &(*memory)[index]
	default:
		panic("error in mode")
	}
}

func toInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	check(err)
	return result
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(bytes))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
