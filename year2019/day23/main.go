package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	in, out := make([]chan int64, 50), make([]chan int64, 50)
	var i int64
	for i = 0; i < 50; i++ {
		in[i], out[i] = make(chan int64), make(chan int64)
		go emulate(program, in[i], out[i])
		in[i] <- i
		in[i] <- -1
	}

	idle := 0
	var old, nat [2]int64

	for i := 0; ; i = (i + 1) % 50 {
		select {
		case addr := <-out[i]:
			if addr == 255 {
				new := [2]int64{<-out[i], <-out[i]}
				if nat == [2]int64{} {
					fmt.Println("--- Part One ---")
					fmt.Println(new[1])
				}

				nat = new
			} else {
				in[addr] <- <-out[i]
				in[addr] <- <-out[i]
			}

			idle = 0
		case in[i] <- -1:
			idle++
		}

		if idle >= 50 {
			if nat[1] == old[1] {
				fmt.Println("--- Part Two ---")
				fmt.Println(nat[1])
				return
			}

			in[0] <- nat[0]
			in[0] <- nat[1]
			old = nat
			idle = 0
		}
	}
}

func emulate(program []int64, input <-chan int64, output chan<- int64) {
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
