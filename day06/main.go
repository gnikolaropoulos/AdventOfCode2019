package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	lines := readLines("input.txt")

	orbits := make(map[string]string)

	for _, line := range lines {
		parts := strings.Split(line, ")")
		orbits[parts[1]] = parts[0]
	}
	fmt.Println("--- Part One ---")

	total := 0
	for object := range orbits {
		for {
			parent, ok := orbits[object]
			if !ok {
				break
			}
			object = parent
			total++
		}
	}

	fmt.Println(total)

	fmt.Println("--- Part Two ---")

	path := make(map[string]int)

	object, distance := orbits["YOU"], 0
	for {
		path[object] = distance
		parent, ok := orbits[object]
		if !ok {
			break
		}
		object = parent
		distance++
	}

	object, distance = orbits["SAN"], 0
	for {
		pathDistance, ok := path[object]
		if ok {
			distance += pathDistance
			break
		}
		parent, ok := orbits[object]
		if !ok {
			panic("YOU and SAN are not connected")
		}
		object = parent
		distance++
	}

	fmt.Println(distance)
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
