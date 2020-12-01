package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	Empty  = 0
	Wall   = 1
	Block  = 2
	Paddle = 3
	Ball   = 4
)

// Warning: For my input, this outputs about 150k lines.
var printFlag = flag.Bool("print", false, "print game state before each input is provided")

func main() {
	flag.Parse()

	input := readFile("input.txt")

	var program []int64
	for _, value := range strings.Split(input, ",") {
		program = append(program, toInt64(value))
	}

	fmt.Println("--- Part One ---")
	fmt.Println(countBlocks(program))

	fmt.Println("--- Part Two ---")
	fmt.Println(emulateArcadeCabinet(program))
}

func countBlocks(program []int64) (count int) {
	input := make(chan int64)
	messages := make(chan Message)

	go emulate(program, input, messages)

	grid := make(map[Vector2]int64)

	for {
		message := <-messages
		switch message.Kind {
		case MessageOutput:
			var pos Vector2
			pos.X = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			pos.Y = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			grid[pos] = message.Value

		case MessageHalt:
			for _, tile := range grid {
				if tile == Block {
					count++
				}
			}
			return

		default:
			panic("unexpected message")
		}
	}
}

func emulateArcadeCabinet(program []int64) int64 {
	// Insert quarters.
	program[0] = 2

	input := make(chan int64)
	messages := make(chan Message)

	go emulate(program, input, messages)

	grid := make(map[Vector2]int64)
	var score int64

	for {
		message := <-messages
		switch message.Kind {
		case MessageWaitingForInput:
			if *printFlag {
				var min, max Vector2
				for pos := range grid {
					min = min.Min(pos)
					max = max.Max(pos)
				}

				for y := min.Y; y <= max.Y; y++ {
					for x := min.X; x <= max.X; x++ {
						switch grid[Vector2{x, y}] {
						case Empty:
							fmt.Print(" ")
						case Wall:
							fmt.Print("█")
						case Block:
							fmt.Print("X")
						case Paddle:
							fmt.Print("-")
						case Ball:
							fmt.Print("O")
						}
					}
					fmt.Println()
				}
				fmt.Println("Score: ", score)
			}

			// Find the ball and the paddle, then move the paddle closer to the ball.
			// Once they have the same x position, this will track the ball perfectly,
			// since they both move at the same speed (1 tile / frame).
			var ball, paddle Vector2
			for pos, tile := range grid {
				switch tile {
				case Ball:
					ball = pos
				case Paddle:
					paddle = pos
				}
			}
			input <- int64(sign(ball.X - paddle.X))

		case MessageOutput:
			var pos Vector2
			pos.X = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			pos.Y = int(message.Value)

			message = <-messages
			if message.Kind != MessageOutput {
				panic("unexpected message")
			}
			if pos.X == -1 && pos.Y == 0 {
				score = message.Value
			} else {
				grid[pos] = message.Value
			}

		case MessageHalt:
			return score

		default:
			panic("unexpected message")
		}
	}
}

const (
	MessageWaitingForInput = iota
	MessageOutput
	MessageHalt
)

type Message struct {
	Kind  int
	Value int64
}

func emulate(program []int64, input <-chan int64, messages chan<- Message) {
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
			messages <- Message{Kind: MessageWaitingForInput}
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			*x = <-input
			ip += 2

		case 4: // OUTPUT
			x := fetchPointerToMemory(c, &memory, ip+1, relativeBase)
			messages <- Message{Kind: MessageOutput, Value: *x}
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
			messages <- Message{Kind: MessageHalt}
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

func sign(x int) int {
	if x > 0 {
		return 1
	}

	if x < 0 {
		return -1
	}

	return 0
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
