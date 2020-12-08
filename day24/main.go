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

		printGrid(nextGrid)
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

	{
		fmt.Println("--- Part Two ---")

		var grid [5][5]bool
		for y, line := range lines {
			for x, char := range line {
				if char == '#' {
					grid[y][x] = true
				}
			}
		}

		state := make(map[int][5][5]bool)
		state[0] = grid
		min, max := 0, 0

		for minute := 0; minute < 200; minute++ {
			next := make(map[int][5][5]bool)

			for index := min - 1; index <= max+1; index++ {
				var nextGrid [5][5]bool

				for y := 0; y < 5; y++ {
					for x := 0; x < 5; x++ {
						if x == 2 && y == 2 {
							continue
						}

						neighbors := 0

						if x > 0 && state[index][y][x-1] {
							neighbors++
						}
						if x < 4 && state[index][y][x+1] {
							neighbors++
						}
						if y > 0 && state[index][y-1][x] {
							neighbors++
						}
						if y < 4 && state[index][y+1][x] {
							neighbors++
						}

						if x == 0 && state[index-1][2][1] {
							neighbors++
						}
						if x == 4 && state[index-1][2][3] {
							neighbors++
						}
						if y == 0 && state[index-1][1][2] {
							neighbors++
						}
						if y == 4 && state[index-1][3][2] {
							neighbors++
						}

						if x == 1 && y == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][i][0] {
									neighbors++
								}
							}
						}

						if x == 3 && y == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][i][4] {
									neighbors++
								}
							}
						}

						if y == 1 && x == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][0][i] {
									neighbors++
								}
							}
						}

						if y == 3 && x == 2 {
							for i := 0; i < 5; i++ {
								if state[index+1][4][i] {
									neighbors++
								}
							}
						}

						nextGrid[y][x] = (state[index][y][x] && neighbors == 1) || (!state[index][y][x] && neighbors >= 1 && neighbors <= 2)
					}
				}

				next[index] = nextGrid
			}

			state = next
			min, max = min-1, max+1
		}

		bugs := 0
		for _, layer := range state {
			bugs += count(layer)
		}

		fmt.Println(bugs)
	}
}

func count(grid [5][5]bool) (bugs int) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if grid[y][x] {
				bugs++
			}
		}
	}

	return
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

func printGrid(grid [5][5]bool) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if grid[i][j] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println()
	}

	fmt.Println()
}
