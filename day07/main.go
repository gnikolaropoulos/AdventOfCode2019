package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input := readFile("input.txt")

	var program []int
	for _, value := range strings.Split(input, ",") {
		program = append(program, toInt(value))
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(findBestSignal(program, []int{0, 1, 2, 3, 4}))
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(findBestSignal(program, []int{5, 6, 7, 8, 9}))
	}
}

func findBestSignal(program []int, phaseValues []int) int {
	bestSignal := 0
	for _, phaseSettings := range allPermutations(phaseValues) {
		signal := emulateAmplifiers(program, phaseSettings)
		bestSignal = max(bestSignal, signal)
	}
	return bestSignal
}

func emulateAmplifiers(program []int, phaseSettings []int) int {
	// Set up the channels connecting the amplifiers.
	ea := make(chan int, 1) // must be buffered to receive final result
	ab := make(chan int)
	bc := make(chan int)
	cd := make(chan int)
	de := make(chan int)

	// This channel will receive a value each time an amplifier halts.
	halt := make(chan bool)

	// Start amplifiers in parallel.
	go emulate(program, ea, ab, halt)
	go emulate(program, ab, bc, halt)
	go emulate(program, bc, cd, halt)
	go emulate(program, cd, de, halt)
	go emulate(program, de, ea, halt)

	// Provide phase settings.
	ea <- phaseSettings[0]
	ab <- phaseSettings[1]
	bc <- phaseSettings[2]
	cd <- phaseSettings[3]
	de <- phaseSettings[4]

	// Send initial input signal.
	ea <- 0

	// Wait for all amplifiers to halt.
	for i := 0; i < 5; i++ {
		<-halt
	}

	// Read the final result.
	return <-ea
}

func emulate(program []int, input <-chan int, output chan<- int, halt chan<- bool) {
	// Copy the program into memory, so that we do not modify the original.
	memory := make([]int, len(program))
	copy(memory, program)

	ip := 0
	for {
		instruction := memory[ip]
		a := instruction / 10000
		b := (instruction - a*10000) / 1000
		c := (instruction - a*10000 - b*1000) / 100
		opcode := instruction % 100

		switch opcode {
		case 1:
			if a != 0 {
				panic("Instruction writes to an immediate mode")
			}

			x := fetchValue(c, memory, ip+1)
			y := fetchValue(b, memory, ip+2)
			z := memory[ip+3]
			memory[z] = x + y
			ip += 4

		case 2:
			if a != 0 {
				panic("Instruction writes to an immediate mode")
			}

			x := fetchValue(c, memory, ip+1)
			y := fetchValue(b, memory, ip+2)
			z := memory[ip+3]
			memory[z] = x * y
			ip += 4

		case 3:
			x := memory[ip+1]
			memory[x] = <-input
			ip += 2

		case 4:
			x := fetchValue(c, memory, ip+1)
			output <- x
			ip += 2

		case 5: // JUMP IF TRUE
			x := fetchValue(c, memory, ip+1)
			y := fetchValue(b, memory, ip+2)
			if x != 0 {
				ip = y
			} else {
				ip += 3
			}

		case 6: // JUMP IF FALSE
			x := fetchValue(c, memory, ip+1)
			y := fetchValue(b, memory, ip+2)
			if x == 0 {
				ip = y
			} else {
				ip += 3
			}

		case 7: // LESS THAN
			x := fetchValue(c, memory, ip+1)
			y := fetchValue(b, memory, ip+2)
			z := memory[ip+3]
			if x < y {
				memory[z] = 1
			} else {
				memory[z] = 0
			}

			ip += 4

		case 8: // EQUAL
			x := fetchValue(c, memory, ip+1)
			y := fetchValue(b, memory, ip+2)
			z := memory[ip+3]
			if x == y {
				memory[z] = 1
			} else {
				memory[z] = 0
			}

			ip += 4

		case 99: // HALT
			halt <- true
			return

		default:
			panic(fmt.Sprintf("error: invalid opcode: ip=%d, instruction=%d, opcode=%d", ip, instruction, opcode))
		}
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSpace(string(bytes))
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchValue(mode int, memory []int, position int) int {
	if mode == 0 {
		return memory[memory[position]]
	}

	return memory[position]
}

func allPermutations(values []int) (result [][]int) {
	if len(values) == 1 {
		result = append(result, values)
		return
	}
	for i, current := range values {
		others := make([]int, 0, len(values)-1)
		others = append(others, values[:i]...)
		others = append(others, values[i+1:]...)
		for _, route := range allPermutations(others) {
			result = append(result, append(route, current))
		}
	}
	return
}

func max(x, y int) int {
	if y > x {
		return y
	}
	return x
}
