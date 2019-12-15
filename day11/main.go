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
	fmt.Println(len(emulateEmergencyHullPaintingRobot(program, 0)))

	fmt.Println("--- Part Two ---")

	grid := emulateEmergencyHullPaintingRobot(program, 1)

	var min, max Vector2
	for pos := range grid {
		min = min.Min(pos)
		max = max.Max(pos)
	}

	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			if grid[Vector2{x, y}] == 1 {
				fmt.Print("")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func emulateEmergencyHullPaintingRobot(program []int64, initialPanel int64) map[Vector2]int64 {
	up := Vector2{0, -1}
	right := Vector2{1, 0}
	down := Vector2{0, 1}
	left := Vector2{-1, 0}

	input := make(chan int64, 1)
	output := make(chan int64)
	halt := make(chan bool)

	go emulate(program, input, output, halt)

	grid := make(map[Vector2]int64)
	pos, dir := Vector2{0, 0}, up

	grid[pos] = initialPanel

	for {
		input <- grid[pos]

		select {
		case value := <-output:
			grid[pos] = value

			if turn := <-output; turn == 1 {
				// turn right
				switch dir {
				case up:
					dir = right
				case right:
					dir = down
				case down:
					dir = left
				case left:
					dir = up
				}
			} else {
				// turn left
				switch dir {
				case up:
					dir = left
				case left:
					dir = down
				case down:
					dir = right
				case right:
					dir = up
				}
			}

			pos = pos.Add(dir)

		case <-halt:
			return grid
		}
	}
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

type Vector2 struct {
	X, Y int
}

// Sub returns the standard vector difference of v and ov (other vector).
func (v Vector2) Sub(ov Vector2) Vector2 {
	return Vector2{
		v.X - ov.X,
		v.Y - ov.Y,
	}
}

// Sub returns the standard vector sum of v and ov (other vector).
func (v Vector2) Add(ov Vector2) Vector2 {
	return Vector2{
		v.X + ov.X,
		v.Y + ov.Y,
	}
}

// Min returns a Vector2 with the minimum coordiates of v and ov (other vector).
func (v Vector2) Min(ov Vector2) Vector2 {
	return Vector2{
		min(v.X, ov.X),
		min(v.Y, ov.Y),
	}
}

// Max returns a Vector2 with the maximum coordiates of v and ov (other vector).
func (v Vector2) Max(ov Vector2) Vector2 {
	return Vector2{
		max(v.X, ov.X),
		max(v.Y, ov.Y),
	}
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

func min(x, y int) int {
	if y < x {
		return y
	}

	return x
}

func max(x, y int) int {
	if y > x {
		return y
	}

	return x
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
