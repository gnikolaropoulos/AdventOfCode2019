package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lines := readLines("input.txt")
	asteroids := make([]Vector2, 0)
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				asteroids = append(asteroids, Vector2{x, y})
			}
		}
	}

	var mostVisibleAsteroids int
	var bestAsteroidLocation Vector2
	for _, location := range asteroids {
		visibleAsteroids := len(findAsteroidsInSight(location, asteroids))
		if visibleAsteroids > mostVisibleAsteroids {
			mostVisibleAsteroids = visibleAsteroids
			bestAsteroidLocation = location
		}
	}

	fmt.Println("--- Part One ---")
	fmt.Println(mostVisibleAsteroids)
	fmt.Println(bestAsteroidLocation)

	// var vaporizationOrder []Vector2
	// for len(asteroids) > 1 {
	// 	// Find the asteroids that will be vaporized during this rotation of the laser.
	// 	list := findVisibleAsteroids(bestAsteroidLocation, asteroids)

	// 	// Sort them by angle to find the exact order they will be vaproized in.
	// 	calculateAngle := func(asteroid Vector2) float64 {
	// 		dist := asteroid.Minus(bestAsteroidLocation)
	// 		return 2.0*math.Pi - (math.Atan2(float64(dist.X), float64(dist.Y)) + math.Pi)
	// 	}
	// 	sort.Slice(list, func(i, j int) bool {
	// 		return calculateAngle(list[i]) < calculateAngle(list[j])
	// 	})

	// 	// Add them to the global list and remove them from the asteroid field.
	// 	vaporizationOrder = append(vaporizationOrder, list...)
	// 	for _, asteroid := range list {
	// 		delete(asteroids, asteroid)
	// 	}
	// }

	// {
	// 	fmt.Println("--- Part Two ---")
	// 	target := vaporizationOrder[199]
	// 	fmt.Println(target.X*100 + target.Y)
	// }
}

func findAsteroidsInSight(location Vector2, asteroids []Vector2) []Vector2 {
	visible := make(map[Vector2]Vector2)

	for _, asteroid := range asteroids {
		if asteroid == location {
			continue
		}

		// Calculate the location vector and find the direction by normalizing it.
		dist := asteroid.Sub(location)
		dir := dist.Normalize()

		if blockingAsteroid, ok := visible[dir]; ok {
			blockingAsteroidDist := blockingAsteroid.Sub(location)
			if (dir.X != 0 && blockingAsteroidDist.X/dir.X < dist.X/dir.X) ||
				(dir.Y != 0 && blockingAsteroidDist.Y/dir.Y < dist.Y/dir.Y) {
				continue
			}
		}

		visible[dir] = asteroid
	}

	var result []Vector2
	for _, asteroid := range visible {
		result = append(result, asteroid)
	}

	return result
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

type Vector2 struct {
	X, Y int
}

// Sub returns the standard vector difference of v and ov (other vector).
func (v Vector2) Sub(ov Vector2) Vector2 {
	return Vector2{
		X: v.X - ov.X,
		Y: v.Y - ov.Y,
	}
}

// Normalize returns a unit vector in the same direction as v.
func (v Vector2) Normalize() Vector2 {
	gcd := GCD(abs(v.X), abs(v.Y))
	return Vector2{v.X / gcd, v.Y / gcd}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GCD(a, b int) int {
	if b == 0 {
		return a
	}

	return GCD(b, a%b)
}
