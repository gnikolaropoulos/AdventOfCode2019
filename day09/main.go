package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(input, ",") {
		program = append(program, toInt64(value))
	}

	fmt.Println("--- Part One ---")
	output := emulate(program, []int64{1})
	for i := 0; i < len(output)-1; i++ {
		if output[i] != 0 {
			panic(fmt.Sprintf("test failure: %v", output))
		}
	}

	fmt.Println(output[len(output)-1])

	fmt.Println("--- Part Two ---")
	output = emulate(program, []int64{2})
	if len(output) != 1 {
		panic(fmt.Sprintf("unexpected output: %v", output))
	}

	fmt.Println(output[0])
}

func emulate(program []int64, input []int64) (output []int64) {
	// Copy the program into memory, so that we do not modify the original.
	memory := make([]int64, len(program))
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
			x := fetchValue(c, &memory, ip+1, relativeBase)
			y := fetchValue(b, &memory, ip+2, relativeBase)
			z := fetchValue(a, &memory, ip+3, relativeBase)
			*z = *x + *y
			ip += 4

		case 2: // MULTIPLY
			x := fetchValue(c, &memory, ip+1, relativeBase)
			y := fetchValue(b, &memory, ip+2, relativeBase)
			z := fetchValue(a, &memory, ip+3, relativeBase)
			*z = *x * *y
			ip += 4

		case 3: // INPUT
			x := fetchValue(c, &memory, ip+1, relativeBase)
			*x = input[0]
			input = input[1:]
			ip += 2

		case 4: // OUTPUT
			x := fetchValue(c, &memory, ip+1, relativeBase)
			output = append(output, *x)
			ip += 2

		case 5: // JUMP IF TRUE
			x := fetchValue(c, &memory, ip+1, relativeBase)
			y := fetchValue(b, &memory, ip+2, relativeBase)
			if *x != 0 {
				ip = *y
			} else {
				ip += 3
			}

		case 6: // JUMP IF FALSE
			x := fetchValue(c, &memory, ip+1, relativeBase)
			y := fetchValue(b, &memory, ip+2, relativeBase)
			if *x == 0 {
				ip = *y
			} else {
				ip += 3
			}

		case 7: // LESS THAN
			x := fetchValue(c, &memory, ip+1, relativeBase)
			y := fetchValue(b, &memory, ip+2, relativeBase)
			z := fetchValue(a, &memory, ip+3, relativeBase)

			if *x < *y {
				*z = 1
			} else {
				*z = 0
			}

			ip += 4

		case 8: // EQUAL
			x := fetchValue(c, &memory, ip+1, relativeBase)
			y := fetchValue(b, &memory, ip+2, relativeBase)
			z := fetchValue(a, &memory, ip+3, relativeBase)
			if *x == *y {
				*z = 1
			} else {
				*z = 0
			}

			ip += 4

		case 9: // ADJUST RELATIVE BASE
			x := fetchValue(c, &memory, ip+1, relativeBase)
			relativeBase += *x
			ip += 2

		case 99: // HALT
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

func toInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchValue(mode int64, memory *[]int64, position int64, relativeBase int64) *int64 {

	switch mode {
	case 0:
		index := (*memory)[position]
		return fetchDataFromMemory(memory, index)
	case 1:
		return fetchDataFromMemory(memory, position)
	case 2:
		index := (*memory)[position] + relativeBase
		return fetchDataFromMemory(memory, index)
	default:
		panic("error in mode")
	}
}

func fetchDataFromMemory(memory *[]int64, index int64) *int64 {
	for int64(len(*memory)) <= index {
		*memory = append(*memory, 0)
	}
	return &(*memory)[index]
}
