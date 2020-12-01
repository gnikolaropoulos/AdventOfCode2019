package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var program []int64

func main() {
	text := readFile("input.txt")

	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}
	fmt.Println("--- Part One ---")

	count := 0

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if probe(x, y) {
				count++
			}
		}
	}

	fmt.Println(count)

	fmt.Println("--- Part Two ---")
	startX, startY := 0, 0
	for {
		if !probe(startX, startY) {
			startX++
		}

		x, y := startX, startY

		for {
			if probe(x, y) && probe(x+99, y) && probe(x, y+99) {
				fmt.Println(x*10000 + y)
				return
			}

			x++

			if !probe(x+99, y) {
				startY++
				break
			}
		}
	}
}

func probe(x, y int) bool {
	input := make(chan int64)
	output := make(chan int64)
	halt := make(chan bool, 1)

	go emulate(program, input, output, halt)

	input <- int64(x)
	input <- int64(y)

	return <-output == 1
}

func emulate(program []int64, input <-chan int64, output chan<- int64, halt chan<- bool) {
	// We will get pointers as for write operators we might get relative mode
	// so z=memory(ip+3) is not correct anymore
	memory := make([]int64, 3000)
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

func fetchPointerToMemory(mode int64, memory *[]int64, position int64, relativeBase int64) *int64 {
	switch mode {
	case 0:
		index := (*memory)[position]
		for int64(len(*memory)) <= index {
			*memory = append(*memory, 0)
		}
		return &(*memory)[index]
	case 1:
		return &(*memory)[position]
	case 2:
		index := (*memory)[position] + relativeBase
		for int64(len(*memory)) <= index {
			*memory = append(*memory, 0)
		}
		return &(*memory)[index]
	default:
		panic("error in mode")
	}
}
