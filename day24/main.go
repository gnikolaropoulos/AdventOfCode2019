package main

import (
	"bufio"
	"fmt"
	"os"
)

type Vector2 struct {
	X, Y int
}

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var directions = []Vector2{Up, Down, Left, Right}

func main() {
	lines := readLines("input.txt")
	var grid [5][5]bool
	for i, line := range lines {
		for j := 0; j < 5; j++ {
			grid[i][j] = line[j] == '#'
		}
	}

	seen := map[[5][5]bool]bool{}
	for {
		if seen[grid] {
			break
		}

		seen[grid] = true

		var nextGrid [5][5]bool
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				count := 0
				currentPosition := Vector2{i, j}
				for _, direction := range directions {
					neighbor := currentPosition.Add(direction)
					if neighbor.X < 0 || neighbor.X >= 5 || neighbor.Y < 0 || neighbor.Y >= 5 {
						continue
					}

					if grid[neighbor.X][neighbor.Y] {
						count++
					}
				}

				if grid[currentPosition.X][currentPosition.Y] {
					nextGrid[currentPosition.X][currentPosition.Y] = count == 1
				} else {
					nextGrid[currentPosition.X][currentPosition.Y] = count == 1 || count == 2
				}
			}
		}

		grid = nextGrid
	}

	rating := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			points := 1 << int((i*5)+j)
			if grid[i][j] {
				rating += points
			}
		}
	}

	fmt.Println("--- Part One ---")
	fmt.Println(rating)
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Add returns the standard vector sum of v and ov (other vector).
func (v Vector2) Add(ov Vector2) Vector2 {
	return Vector2{
		v.X + ov.X,
		v.Y + ov.Y,
	}
}
