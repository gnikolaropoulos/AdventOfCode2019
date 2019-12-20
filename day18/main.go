// package main

// // package main

// import (
// 	"fmt"
// 	"hash/fnv"
// 	"io/ioutil"
// 	"math"
// 	"strings"
// 	"unicode"
// )

// // Point ...
// type Point struct {
// 	row, col int
// }

// // Key ...
// type Key struct {
// 	Point
// 	steps int
// }

// var initGrid [][]rune
// var startingPos Point

// func main() {
// 	startingPos = readInput("input.txt")

// 	fmt.Println(shortestPath(initGrid, startingPos, make(map[int64]int)))
// }

// func shortestPath(grid [][]rune, pos Point, memo map[int64]int) int {
// 	keys := findReachableKeys(grid, pos)
// 	stateHash := serializeState(grid)
// 	if _, ok := memo[stateHash]; ok {
// 		return memo[stateHash]
// 	}

// 	if len(keys) == 0 {
// 		return 0
// 	}

// 	minPath := math.MaxInt32
// 	cnt := 0
// 	for _, key := range keys {
// 		if pos == startingPos {
// 			fmt.Println("STEP:", cnt)
// 		}

// 		g, p := removeKey(grid, pos, key)
// 		minPath = intMin(minPath, key.steps+shortestPath(g, p, memo))

// 		cnt++
// 	}

// 	memo[stateHash] = minPath

// 	return minPath
// }

// func findReachableKeys(grid [][]rune, pos Point) []Key {
// 	return bfs(grid, pos)
// }

// func bfs(grid [][]rune, pos Point) []Key {
// 	type Step struct {
// 		Point
// 		steps int
// 	}

// 	visited := make([][]bool, len(grid))
// 	for row := range grid {
// 		visited[row] = make([]bool, len(grid[row]))
// 	}

// 	queue := []Step{}
// 	queue = append(queue, Step{pos, 0})

// 	keys := []Key{}
// 	for len(queue) > 0 {
// 		step := queue[0]
// 		queue = queue[1:]

// 		p := step.Point
// 		if p.row < 0 || p.row >= len(grid) || p.col < 0 || p.col >= len(grid[p.row]) {
// 			continue
// 		} else if grid[p.row][p.col] == '#' || unicode.IsUpper(grid[p.row][p.col]) {
// 			continue
// 		} else if visited[p.row][p.col] {
// 			continue
// 		}

// 		if unicode.IsLower(grid[p.row][p.col]) {
// 			keys = append(keys, Key{p, step.steps})
// 		}

// 		visited[p.row][p.col] = true
// 		queue = append(queue, Step{Point{p.row - 1, p.col}, step.steps + 1})
// 		queue = append(queue, Step{Point{p.row, p.col + 1}, step.steps + 1})
// 		queue = append(queue, Step{Point{p.row + 1, p.col}, step.steps + 1})
// 		queue = append(queue, Step{Point{p.row, p.col - 1}, step.steps + 1})

// 		// diagonals
// 		queue = append(queue, Step{Point{p.row - 1, p.col - 1}, step.steps + 1})
// 		queue = append(queue, Step{Point{p.row - 1, p.col + 1}, step.steps + 1})
// 		queue = append(queue, Step{Point{p.row + 1, p.col - 1}, step.steps + 1})
// 		queue = append(queue, Step{Point{p.row + 1, p.col + 1}, step.steps + 1})
// 	}

// 	return keys
// }

// func removeKey(grid [][]rune, pos Point, key Key) ([][]rune, Point) {
// 	newGrid := make([][]rune, len(grid))
// 	for row := range grid {
// 		newGrid[row] = make([]rune, len(grid[row]))
// 		copy(newGrid[row], grid[row])
// 	}

// 	keyLetter := grid[key.row][key.col]
// 	for row := range grid {
// 		for col := range grid[row] {
// 			if newGrid[row][col] == unicode.ToUpper(keyLetter) {
// 				newGrid[row][col] = '.'
// 			}
// 		}
// 	}

// 	newGrid[pos.row][pos.col] = '.'
// 	newGrid[key.row][key.col] = '@'

// 	return newGrid, Point{key.row, key.col}
// }

// func serializeState(grid [][]rune) int64 {
// 	h := fnv.New32a()
// 	for row := range grid {
// 		for col := range grid[row] {
// 			h.Write([]byte{byte(grid[row][col])})
// 		}
// 	}

// 	return int64(h.Sum32())
// }

// func intMin(a, b int) int {
// 	if a < b {
// 		return a
// 	}

// 	return b
// }

// func readInput(filename string) Point {
// 	bs, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	lines := strings.Split(string(bs), "\n")
// 	lines = lines[:len(lines)-1]

// 	var startingPos Point
// 	initGrid = make([][]rune, len(lines))
// 	for row := range lines {
// 		initGrid[row] = make([]rune, len(lines[row]))
// 		for col := range lines[row] {
// 			initGrid[row][col] = rune(lines[row][col])
// 			if initGrid[row][col] == '@' {
// 				startingPos = Point{row, col}
// 			}
// 		}
// 	}

// 	return startingPos
// }

// func printGrid(grid [][]rune) {
// 	for row := range grid {
// 		for col := range grid[row] {
// 			fmt.Print(string(grid[row][col]))
// 		}

// 		fmt.Println()
// 	}
// }

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	// fmt.Println("--- Part One ---")
	// findMinimumSteps()

	fmt.Println("--- Part Two ---")
	findMinimumStepsForMultipleVaults()
}

func findMinimumSteps() {
	lines := readLines("input.txt")

	var startPos [2]int
	var doorPos [26][2]int
	var keyPos [26][2]int

	for i, line := range lines {
		for j := 0; j < len(line); j++ {
			if line[j] == '@' {
				startPos = [2]int{i, j}
			} else if line[j] >= 'a' && line[j] <= 'z' {
				keyPos[line[j]-'a'] = [2]int{i, j}
			} else if line[j] >= 'A' && line[j] <= 'Z' {
				doorPos[line[j]-'A'] = [2]int{i, j}
			}
		}
	}

	pmemo := map[[4]int]int{}

	pathFind := func(pos [2]int, keys int, key int) int {
		pkey := [4]int{pos[0], pos[1], keys, key}
		if v, ok := pmemo[pkey]; ok {
			return v
		}

		work := [][3]int{{pos[0], pos[1], 0}}
		seen := map[[2]int]bool{}
		target := keyPos[key]
		dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
		for len(work) > 0 {
			w := work[0]
			work = work[1:]
			k := [2]int{w[0], w[1]}
			if k == target {
				pmemo[pkey] = w[2]
				return w[2]
			}
			if seen[k] {
				continue
			}
			seen[k] = true
			for _, dir := range dirs {
				w2 := [3]int{w[0] + dir[0], w[1] + dir[1], w[2] + 1}
				if w2[0] < 0 || w2[0] >= len(lines) || w2[1] < 0 || w2[1] >= len(lines[0]) {
					continue
				}
				c := lines[w2[0]][w2[1]]
				if c == '#' {
					continue
				}
				b := 1 << uint(c-'a')
				if c >= 'a' && c <= 'z' && int(c-'a') != key && (keys&b) == 0 {
					continue
				}
				b = 1 << uint(c-'A')
				if c >= 'A' && c <= 'Z' && (keys&b) == 0 {
					continue
				}
				work = append(work, w2)
			}
		}
		pmemo[pkey] = -1
		return -1
	}

	min := math.MaxInt64

	mins := map[[3]int]int{}

	var search func([2]int, int, int)
	search = func(pos [2]int, keys int, plen int) {
		if keys == (1<<26)-1 {
			if plen < min {
				min = plen
			}
			return
		}
		if plen >= min {
			return
		}

		key := [3]int{pos[0], pos[1], keys}
		if v, ok := mins[key]; ok && v <= plen {
			return
		}
		mins[key] = plen

		for nextKey := 0; nextKey < 26; nextKey++ {
			bit := 1 << uint(nextKey)
			if keys&bit != 0 {
				continue
			}
			aplen := pathFind(pos, keys, nextKey)
			if aplen == -1 {
				continue
			}
			search(keyPos[nextKey], keys|bit, plen+aplen)
		}
	}
	search(startPos, 0, 0)
	fmt.Println(min)
}

func findMinimumStepsForMultipleVaults() {
	lines := readLines("input.txt")

	spcount := 0
	var startPoss [4][2]int
	var doorPos [26][2]int
	var keyPos [26][2]int

	for i, line := range lines {
		for j := 0; j < len(line); j++ {
			if line[j] == '@' {
				startPoss[spcount] = [2]int{i, j}
				spcount++
			} else if line[j] >= 'a' && line[j] <= 'z' {
				keyPos[line[j]-'a'] = [2]int{i, j}
			} else if line[j] >= 'A' && line[j] <= 'Z' {
				doorPos[line[j]-'A'] = [2]int{i, j}
			}
		}
	}

	pmemo := map[[4]int]int{}

	pathFind := func(pos [2]int, keys int, key int) int {
		pkey := [4]int{pos[0], pos[1], keys, key}
		if v, ok := pmemo[pkey]; ok {
			return v
		}

		work := [][3]int{{pos[0], pos[1], 0}}
		seen := map[[2]int]bool{}
		target := keyPos[key]
		dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
		for len(work) > 0 {
			w := work[0]
			work = work[1:]
			k := [2]int{w[0], w[1]}
			if k == target {
				pmemo[pkey] = w[2]
				return w[2]
			}
			if seen[k] {
				continue
			}
			seen[k] = true
			for _, dir := range dirs {
				w2 := [3]int{w[0] + dir[0], w[1] + dir[1], w[2] + 1}
				if w2[0] < 0 || w2[0] >= len(lines) || w2[1] < 0 || w2[1] >= len(lines[0]) {
					continue
				}
				c := lines[w2[0]][w2[1]]
				if c == '#' {
					continue
				}
				b := 1 << uint(c-'a')
				if c >= 'a' && c <= 'z' && int(c-'a') != key && (keys&b) == 0 {
					continue
				}
				b = 1 << uint(c-'A')
				if c >= 'A' && c <= 'Z' && (keys&b) == 0 {
					continue
				}
				work = append(work, w2)
			}
		}
		pmemo[pkey] = -1
		return -1
	}

	keySegs := [4]int{}
	for i := 0; i < 26; i++ {
		for j := 0; j < 4; j++ {
			if pathFind(startPoss[j], (1<<26)-1, i) != -1 {
				keySegs[j] |= (1 << uint(i))
			}
		}
	}

	min := math.MaxInt64

	mins := map[[9]int]int{}

	var search func([4][2]int, int, int)
	search = func(pos [4][2]int, keys int, plen int) {
		if keys == (1<<26)-1 {
			if plen < min {
				min = plen
			}
			return
		}
		if plen >= min {
			return
		}

		key := [9]int{pos[0][0], pos[0][1], pos[1][0], pos[1][1], pos[2][0], pos[2][1], pos[3][0], pos[3][1], keys}
		if v, ok := mins[key]; ok && v <= plen {
			return
		}
		mins[key] = plen

		for nextKey := 0; nextKey < 26; nextKey++ {
			bit := 1 << uint(nextKey)
			if keys&bit != 0 {
				continue
			}
			var bi int
			for i := 0; i < 4; i++ {
				if keySegs[i]&bit != 0 {
					bi = i
				}
			}
			npos := pos
			npos[bi] = keyPos[nextKey]
			aplen := pathFind(pos[bi], keys, nextKey)
			if aplen == -1 {
				continue
			}
			search(npos, keys|bit, plen+aplen)
		}
	}
	search(startPoss, 0, 0)
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
