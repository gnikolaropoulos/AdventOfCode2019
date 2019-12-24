package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var (
	Up    = Vector2{0, -1}
	Down  = Vector2{0, 1}
	Left  = Vector2{-1, 0}
	Right = Vector2{1, 0}
)

var directions = []Vector2{Up, Down, Left, Right}

func main() {
	fmt.Println("--- Part One ---")
	findMinimumSteps()

	fmt.Println("--- Part Two ---")
	findMinimumStepsForMultipleVaults()
}

func findMinimumSteps() {
	lines := readLines("input.txt")

	var startingPosition Vector2
	var keyPositions [26]Vector2

	for i, line := range lines {
		for j := 0; j < len(line); j++ {
			if line[j] == '@' {
				startingPosition = Vector2{i, j}
			} else if line[j] >= 'a' && line[j] <= 'z' {
				keyPositions[line[j]-'a'] = Vector2{i, j}
			}
		}
	}

	pathsMemory := map[KeyInfo]int{}

	pathFind := func(keyInfo KeyInfo) int {
		if distance, ok := pathsMemory[keyInfo]; ok {
			return distance
		}

		workItems := []WorkItem{WorkItem{keyInfo.Position, 0}}
		visited := map[Vector2]bool{}
		target := keyPositions[keyInfo.Key]

		for len(workItems) > 0 {
			topItem := workItems[0]
			workItems = workItems[1:]
			currentPosition := topItem.Position
			if currentPosition == target {
				pathsMemory[keyInfo] = topItem.Distance
				return topItem.Distance
			}

			if visited[currentPosition] {
				continue
			}

			visited[currentPosition] = true
			for _, dir := range directions {
				nextItem := WorkItem{topItem.Position.Add(dir), topItem.Distance + 1}
				if isOutOfBounds(nextItem.Position, lines) {
					continue
				}

				char := lines[nextItem.Position.X][nextItem.Position.Y]
				if char == '#' {
					continue
				}

				// collect the key
				keyBit := 1 << uint(char-'a')
				if isKey(char) && int(char-'a') != keyInfo.Key && (keyInfo.KeysBitcode&keyBit) == 0 {
					continue
				}

				// open the door
				doorBit := 1 << uint(char-'A')
				if isDoor(char) && (keyInfo.KeysBitcode&doorBit) == 0 {
					continue
				}

				workItems = append(workItems, nextItem)
			}
		}

		pathsMemory[keyInfo] = -1
		return -1
	}

	min := math.MaxInt64

	mins := map[KeyInfo]int{}

	var search func(Vector2, int, int)
	search = func(position Vector2, keysBitcode int, pathLength int) {
		if keysBitcode == (1<<26)-1 {
			if pathLength < min {
				min = pathLength
			}

			return
		}

		if pathLength >= min {
			return
		}

		keyInfo := KeyInfo{position, keysBitcode, 0}
		if v, ok := mins[keyInfo]; ok && v <= pathLength {
			return
		}

		mins[keyInfo] = pathLength

		for nextKey := 0; nextKey < 26; nextKey++ {
			bit := 1 << uint(nextKey)
			if keysBitcode&bit != 0 {
				continue
			}

			keyInfo.Key = nextKey
			nextKeyDistance := pathFind(keyInfo)
			if nextKeyDistance == -1 {
				continue
			}

			search(keyPositions[nextKey], keysBitcode|bit, pathLength+nextKeyDistance)
		}
	}

	search(startingPosition, 0, 0)
	fmt.Println(min)
}

func findMinimumStepsForMultipleVaults() {
	lines := readLines("input.txt")
	lines = updateInput(lines)

	spcount := 0
	var startingPositions [4]Vector2
	var keyPositions [26]Vector2

	for i, line := range lines {
		for j := 0; j < len(line); j++ {
			if line[j] == '@' {
				startingPositions[spcount] = Vector2{i, j}
				spcount++
			} else if line[j] >= 'a' && line[j] <= 'z' {
				keyPositions[line[j]-'a'] = Vector2{i, j}
			}
		}
	}

	pathsMemory := map[KeyInfo]int{}

	pathFind := func(keyInfo KeyInfo) int {
		if distance, ok := pathsMemory[keyInfo]; ok {
			return distance
		}

		workItems := []WorkItem{WorkItem{keyInfo.Position, 0}}
		visited := map[Vector2]bool{}
		target := keyPositions[keyInfo.Key]

		for len(workItems) > 0 {
			topItem := workItems[0]
			workItems = workItems[1:]
			currentPosition := topItem.Position
			if currentPosition == target {
				pathsMemory[keyInfo] = topItem.Distance
				return topItem.Distance
			}

			if visited[currentPosition] {
				continue
			}

			visited[currentPosition] = true
			for _, dir := range directions {
				nextItem := WorkItem{topItem.Position.Add(dir), topItem.Distance + 1}
				if isOutOfBounds(nextItem.Position, lines) {
					continue
				}

				char := lines[nextItem.Position.X][nextItem.Position.Y]
				if char == '#' {
					continue
				}

				// collect the key
				keyBit := 1 << uint(char-'a')
				if isKey(char) && int(char-'a') != keyInfo.Key && (keyInfo.KeysBitcode&keyBit) == 0 {
					continue
				}

				// open the door
				doorBit := 1 << uint(char-'A')
				if isDoor(char) && (keyInfo.KeysBitcode&doorBit) == 0 {
					continue
				}

				workItems = append(workItems, nextItem)
			}
		}

		pathsMemory[keyInfo] = -1
		return -1
	}

	keySegs := [4]int{}
	for i := 0; i < 26; i++ {
		for j := 0; j < 4; j++ {
			if pathFind(KeyInfo{startingPositions[j], (1 << 26) - 1, i}) != -1 {
				keySegs[j] |= (1 << uint(i))
			}
		}
	}

	min := math.MaxInt64

	mins := map[[9]int]int{}

	var search func([4]Vector2, int, int)
	search = func(pos [4]Vector2, keysBitcode int, pathLength int) {
		if keysBitcode == (1<<26)-1 {
			if pathLength < min {
				min = pathLength
			}

			return
		}

		if pathLength >= min {
			return
		}

		key := [9]int{pos[0].X, pos[0].Y, pos[1].X, pos[1].Y, pos[2].X, pos[2].Y, pos[3].X, pos[3].Y, keysBitcode}
		if v, ok := mins[key]; ok && v <= pathLength {
			return
		}

		mins[key] = pathLength

		for nextKey := 0; nextKey < 26; nextKey++ {
			bit := 1 << uint(nextKey)
			if keysBitcode&bit != 0 {
				continue
			}

			var bi int
			for i := 0; i < 4; i++ {
				if keySegs[i]&bit != 0 {
					bi = i
				}
			}

			npos := pos
			npos[bi] = keyPositions[nextKey]
			nextKeyDistance := pathFind(KeyInfo{pos[bi], keysBitcode, nextKey})
			if nextKeyDistance == -1 {
				continue
			}

			search(npos, keysBitcode|bit, pathLength+nextKeyDistance)
		}
	}

	search(startingPositions, 0, 0)
	fmt.Println(min)
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

func updateInput(lines []string) []string {
	grid := make([][]byte, len(lines))
	for y, line := range lines {
		grid[y] = []byte(line)
	}

	// Find the original entrance.
	var cx, cy int
	for y, line := range lines {
		for x, char := range line {
			if char == '@' {
				cx, cy = x, y
			}
		}
	}

	grid[cy-1][cx-1], grid[cy-1][cx], grid[cy-1][cx+1] = '@', '#', '@'
	grid[cy+0][cx-1], grid[cy+0][cx], grid[cy+0][cx+1] = '#', '#', '#'
	grid[cy+1][cx-1], grid[cy+1][cx], grid[cy+1][cx+1] = '@', '#', '@'

	for y, line := range grid {
		lines[y] = string(line)
	}

	return lines
}

func isOutOfBounds(position Vector2, lines []string) bool {
	return position.X < 0 || position.Y >= len(lines) ||
		position.Y < 0 || position.Y >= len(lines[0])
}

func isKey(char byte) bool {
	return char >= 'a' && char <= 'z'
}

func isDoor(char byte) bool {
	return char >= 'A' && char <= 'Z'
}

type Vector2 struct {
	X, Y int
}

// Add returns the standard vector sum of v and ov (other vector).
func (v Vector2) Add(ov Vector2) Vector2 {
	return Vector2{
		v.X + ov.X,
		v.Y + ov.Y,
	}
}

type WorkItem struct {
	Position Vector2
	Distance int
}

type KeyInfo struct {
	Position    Vector2
	KeysBitcode int
	Key         int
}
