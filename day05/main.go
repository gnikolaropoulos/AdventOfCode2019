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

	fmt.Println("--- Part One ---")
	output := emulate(program, []int{1})
	for i := 0; i < len(output)-1; i++ {
		if output[i] != 0 {
			panic(fmt.Sprintf("test failure: %v", output))
		}
	}

	fmt.Println(output[len(output)-1])

	fmt.Println("--- Part Two ---")
	output = emulate(program, []int{5})
	if len(output) != 1 {
		panic(fmt.Sprintf("unexpected output: %v", output))
	}

	fmt.Println(output[0])
}

func emulate(program []int, input []int) (output []int) {
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
			memory[x] = input[0]
			input = input[1:]
			ip += 2

		case 4:
			x := fetchValue(c, memory, ip+1)
			output = append(output, x)
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
			return output
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
