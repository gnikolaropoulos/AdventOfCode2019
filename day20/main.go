package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cell struct {
	Pos         Vector2
	Connections []Connection
}

type Connection struct {
	Neighbor    *Cell
	LevelOffset int
}

type Label struct {
	Name  string
	Cells []*Cell
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
	width, height := len(lines[0]), len(lines)

	cells := make(map[Vector2]*Cell)
	labels := make(map[Vector2]*Label)
	labelsByName := make(map[string]*Label)

	// Create cells and labels.
	for y, line := range lines {
		for x, char := range line {
			if char != '.' {
				continue
			}

			pos := Vector2{x, y}
			cell := &Cell{Pos: pos}
			cells[pos] = cell

			for _, dir := range directions {
				next := pos.Add(dir)
				nextchar := lines[next.Y][next.X]

				if !(nextchar >= 'A' && nextchar <= 'Z') {
					continue
				}

				beyondnext := next.Add(dir)
				beyondnextchar := lines[beyondnext.Y][beyondnext.X]

				name := string(beyondnextchar) + string(nextchar)
				if dir.X >= 0 && dir.Y >= 0 {
					name = string(nextchar) + string(beyondnextchar)
				}

				label := labelsByName[name]
				if label == nil {
					label = &Label{Name: name}
					labelsByName[name] = label
				}

				labels[next] = label
				label.Cells = append(label.Cells, cell)
			}
		}
	}

	// Connect cells.
	for _, cell := range cells {
		for _, dir := range directions {
			next := cell.Pos.Add(dir)
			if nextCell := cells[next]; nextCell != nil {
				cell.Connections = append(cell.Connections, Connection{nextCell, 0})
			}
			if nextLabel := labels[next]; nextLabel != nil {
				offset := 1
				if next.X <= 2 || next.X >= width-2 || next.Y <= 2 || next.Y >= height-2 {
					offset = -1
				}
				for _, nextCell := range nextLabel.Cells {
					if nextCell != cell {
						cell.Connections = append(cell.Connections, Connection{nextCell, offset})
						break
					}
				}
			}
		}
	}

	start := labelsByName["AA"].Cells[0]
	target := labelsByName["ZZ"].Cells[0]

	{
		fmt.Println("--- Part One ---")
		fmt.Println(findDistance(start, target, false))
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(findDistance(start, target, true))
	}
}

func findDistance(start, target *Cell, recursive bool) int {
	type Position struct {
		Cell  *Cell
		Level int
	}

	type Item struct {
		Position Position
		Distance int
	}

	startItem := Item{Position{start, 0}, 0}

	var open []Item
	open = append(open, startItem)

	visited := make(map[Position]bool)
	visited[startItem.Position] = true

	for {
		item := open[0]
		open = open[1:]

		if item.Position.Cell == target && item.Position.Level == 0 {
			return item.Distance
		}

		for _, conn := range item.Position.Cell.Connections {
			nextPosition := Position{conn.Neighbor, item.Position.Level}
			if recursive {
				nextPosition.Level += conn.LevelOffset
			}
			if nextPosition.Level >= 0 && !visited[nextPosition] {
				visited[nextPosition] = true
				open = append(open, Item{nextPosition, item.Distance + 1})
			}
		}
	}
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
