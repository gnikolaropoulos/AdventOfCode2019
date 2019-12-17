package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var (
	turnLeft  = map[Vector2]Vector2{Up: Left, Left: Down, Down: Right, Right: Up}
	turnRight = map[Vector2]Vector2{Up: Right, Right: Down, Down: Left, Left: Up}
)

func main() {
	puzzleInput := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(puzzleInput, ",") {
		program = append(program, toInt64(value))
	}

	var grid []string
	var width, height int

	// Run the program and extract the camera image into grid.
	{
		input := make(chan int64)
		output := make(chan int64)
		halt := make(chan bool)

		go emulate(program, input, output, halt)

		var builder strings.Builder

	loop:
		for {
			select {
			case char := <-output:
				builder.WriteRune(rune(char))

			case <-halt:
				break loop
			}
		}

		grid = strings.Split(strings.TrimSpace(builder.String()), "\n")
		width, height = len(grid[0]), len(grid)
	}

	fmt.Println("--- Part One ---")
	sumOfAlignmentParameters := 0
	for y := 1; y+1 < height; y++ {
		for x := 1; x+1 < width; x++ {
			if grid[y][x] == '#' && grid[y-1][x] == '#' && grid[y+1][x] == '#' && grid[y][x-1] == '#' && grid[y][x+1] == '#' {
				sumOfAlignmentParameters += x * y
			}
		}
	}
	fmt.Println(sumOfAlignmentParameters)

	fmt.Println("--- Part Two ---")

	// Wake up the robot.
	program[0] = 2

	// Find the robot.
	var pos, dir Vector2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch grid[y][x] {
			case '^':
				pos = Vector2{x, y}
				dir = Up
			case 'v':
				pos = Vector2{x, y}
				dir = Down
			case '<':
				pos = Vector2{x, y}
				dir = Left
			case '>':
				pos = Vector2{x, y}
				dir = Right
			}
		}
	}

	isScaffold := func(pos Vector2) bool {
		return pos.X >= 0 && pos.Y >= 0 && pos.X < width && pos.Y < height && grid[pos.Y][pos.X] == '#'
	}

	// Gather commands to follow the path.
	// Walk straight for as long as possible, then check if we can turn left or right.
	// If we cannot do either, we have reached the end of the path.
	var path MoveList
	for {
		length := 0
		for isScaffold(pos.Add(dir)) {
			pos = pos.Add(dir)
			length++
		}
		if length != 0 {
			path = append(path, strconv.Itoa(length))
		}

		if newDir := turnLeft[dir]; isScaffold(pos.Add(newDir)) {
			dir = newDir
			path = append(path, "L")
		} else if newDir := turnRight[dir]; isScaffold(pos.Add(newDir)) {
			dir = newDir
			path = append(path, "R")
		} else {
			break
		}
	}

	result := compressPath(path, []MoveList{path}, nil)
	if len(result) == 0 {
		panic("no solution found")
	}

	input := make(chan int64, 100)
	output := make(chan int64)
	halt := make(chan bool)

	go emulate(program, input, output, halt)

	functions := result[0]
	main := strings.Join(functions[0], ",")
	a := strings.Join(functions[1], ",")
	b := strings.Join(functions[2], ",")
	c := strings.Join(functions[3], ",")

	for _, c := range fmt.Sprintf("%s\n%s\n%s\n%s\nn\n", main, a, b, c) {
		input <- int64(c)
	}

loop2:
	for {
		select {
		case char := <-output:
			if char >= 128 {
				fmt.Println(char)
			}

		case <-halt:
			break loop2
		}
	}
}

type MoveList []string

// Parameters:
//   path - full path to compress
//   fragments - parts of the path that are not yet part of a function
//   functions - functions that have already been fixed
// Return value: list of valid programs, each of which contains 4 functions (main, A, B and C)
func compressPath(path MoveList, fragments []MoveList, functions []MoveList) (result [][4]MoveList) {
	if len(functions) == 3 {
		// The main function cannot call movement commands,
		// so there must not be any commands left.
		if len(fragments) != 0 {
			return nil
		}

		// Replace path with function calls to compute main function.
		var mainFunction MoveList
		for len(path) != 0 {
			for i, function := range functions {
				if hasPrefix(path, function) {
					mainFunction = append(mainFunction, string('A'+i))
					path = path[len(function):]
				}
			}
		}

		// Check memory limit for main function.
		if len(strings.Join(mainFunction, ",")) > 20 {
			return nil
		}

		var program [4]MoveList
		program[0] = mainFunction
		program[1] = functions[0]
		program[2] = functions[1]
		program[3] = functions[2]

		result = append(result, program)
		return
	}

	if len(fragments) == 0 {
		// Add empty candidate to functions.
		newFunctions := make([]MoveList, 0, 3)
		newFunctions = append(newFunctions, functions...)
		newFunctions = append(newFunctions, MoveList{})

		subresult := compressPath(path, fragments, newFunctions)
		result = append(result, subresult...)
		return
	}

	// Checking the first fragment is enough.
	fragment := fragments[0]

	// Collect candidates.
	var candidates []MoveList
	for length := 1; length <= len(fragment); length++ {
		candidate := fragment[:length]
		text := strings.Join(candidate, ",")
		if len(text) <= 20 {
			candidates = append(candidates, candidate)
		}
	}

	// Try each candidate.
	for _, candidate := range candidates {
		// Split fragments by candidate.
		var newFragments []MoveList
		for _, fragment := range fragments {
			for {
				i := index(fragment, candidate)
				if i == -1 {
					break
				}
				if i != 0 {
					newFragments = append(newFragments, fragment[:i])
				}
				fragment = fragment[i+len(candidate):]
			}
			if len(fragment) != 0 {
				newFragments = append(newFragments, fragment)
			}
		}

		// Add candidate to functions.
		newFunctions := make([]MoveList, 0, 3)
		newFunctions = append(newFunctions, functions...)
		newFunctions = append(newFunctions, candidate)

		subresult := compressPath(path, newFragments, newFunctions)
		result = append(result, subresult...)
	}

	return
}

func hasPrefix(list, prefix MoveList) bool {
	if len(list) < len(prefix) {
		return false
	}
	for i, move := range prefix {
		if list[i] != move {
			return false
		}
	}
	return true
}

func index(list, sublist MoveList) int {
	for i := 0; i <= len(list)-len(sublist); i++ {
		if hasPrefix(list[i:], sublist) {
			return i
		}
	}
	return -1
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

func (v Vector2) Add(ov Vector2) Vector2 {
	return Vector2{
		v.X + ov.X,
		v.Y + ov.Y,
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
