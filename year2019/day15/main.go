package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Wall = 0
	Path = 1
)

// Direction where the robot can go
type Direction int

const (
	Nowhere Direction = iota
	North
	South
	West
	East
)

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var (
	commands   = []Direction{North, South, West, East}
	directions = []Vector2{Up, Down, Left, Right}
	reverse    = map[Direction]Direction{North: South, South: North, West: East, East: West}
	direction  = map[Direction]Vector2{North: Up, South: Down, West: Left, East: Right}
)

type QueueItem struct {
	Position Vector2
	Distance int
	Next     *QueueItem
}

func main() {
	text := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(text, ",") {
		program = append(program, toInt64(value))
	}

	input := make(chan Direction)
	output := make(chan int64)
	halt := make(chan bool)

	go emulate(program, input, output, halt)

	var pos Vector2

	grid := make(map[Vector2]int)
	grid[pos] = Path

	var oxygenPos Vector2
	var oxygenDistance int

	var queue []QueueItem
	queue = append(queue, QueueItem{Position: pos, Distance: 0})

	for len(queue) != 0 {
		item := queue[0]
		queue = queue[1:]

		pos = navigate(pos, item.Position, grid, input, output)

		for _, cmd := range commands {
			next, nextDistance := pos.Add(direction[cmd]), item.Distance+1
			if _, ok := grid[next]; !ok {
				input <- cmd
				switch <-output {
				case 0:
					grid[next] = Wall
				case 2:
					if oxygenDistance == 0 {
						oxygenDistance = nextDistance
						oxygenPos = next
					}
					fallthrough
				case 1:
					grid[next] = Path
					queue = append(queue, QueueItem{Position: next, Distance: nextDistance})
					// Command succeeded, go back to try other commands.
					input <- reverse[cmd]
					<-output
				}
			}
		}
	}

	fmt.Println("--- Part One ---")
	fmt.Println(oxygenDistance)

	var maxDistance int
	queue = append(queue, QueueItem{Position: oxygenPos, Distance: 0})

	visited := make(map[Vector2]bool)
	visited[oxygenPos] = true

	for len(queue) != 0 {
		item := queue[0]
		queue = queue[1:]

		maxDistance = max(maxDistance, item.Distance)

		for _, dir := range directions {
			next := item.Position.Add(dir)
			if !visited[next] && grid[next] == Path {
				visited[next] = true
				queue = append(queue, QueueItem{Position: next, Distance: item.Distance + 1})
			}
		}
	}

	fmt.Println("--- Part Two ---")
	fmt.Println(maxDistance)
}

func navigate(pos, target Vector2, grid map[Vector2]int, input chan Direction, output chan int64) Vector2 {
	var link *QueueItem

	// Find shortest route from target to pos (note reversed order).
	{
		var queue []QueueItem
		queue = append(queue, QueueItem{Position: target, Distance: 0})

		visited := make(map[Vector2]bool)
		visited[target] = true

		for len(queue) != 0 {
			item := queue[0]
			queue = queue[1:]

			if item.Position == pos {
				link = &item
				break
			}

			for _, dir := range directions {
				next := item.Position.Add(dir)
				if !visited[next] && grid[next] == Path {
					visited[next] = true
					queue = append(queue, QueueItem{Position: next, Distance: item.Distance + 1, Next: &item})
				}
			}
		}
	}

	// Follow path backwards from pos to target and apply correct commands.
	for link.Next != nil {
		for cmd, dir := range direction {
			if link.Position.Add(dir) == link.Next.Position {
				input <- cmd
				<-output
				break
			}
		}

		link = link.Next
	}

	return target
}

func emulate(program []int64, input <-chan Direction, output chan<- int64, halt chan<- bool) {
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
			*x = int64(<-input)
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

// Add returns the standard vector sum of v and ov (other vector).
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
