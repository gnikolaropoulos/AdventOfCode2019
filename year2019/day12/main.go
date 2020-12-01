package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Moon struct {
	pos, vel Vector3
}

func main() {
	lines := readLines("input.txt")

	var input []Moon

	regex := regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)
	for _, line := range lines {
		match := regex.FindStringSubmatch(line)
		x, y, z := toInt(match[1]), toInt(match[2]), toInt(match[3])
		moon := Moon{
			pos: Vector3{x, y, z},
			vel: Vector3{0, 0, 0},
		}

		input = append(input, moon)
	}

	fmt.Println("--- Part One ---")

	moons := make([]Moon, len(input))
	copy(moons, input)

	for step := 0; step < 1000; step++ {
		simulate(moons)
	}

	totalEnergy := 0
	for _, moon := range moons {
		potentialEnergy := moon.pos.ManhattenLength()
		kineticEnergy := moon.vel.ManhattenLength()
		totalEnergy += potentialEnergy * kineticEnergy
	}

	fmt.Println(totalEnergy)

	fmt.Println("--- Part Two ---")
	moons = make([]Moon, len(input))
	copy(moons, input)

	xSteps, ySteps, zSteps := 0, 0, 0
	for steps := 1; xSteps == 0 || ySteps == 0 || zSteps == 0; steps++ {
		simulate(moons)

		if xSteps == 0 {
			found := true
			for i, moon := range moons {
				if moon.pos.X != input[i].pos.X || moon.vel.X != input[i].vel.X {
					found = false
					break
				}
			}
			if found {
				xSteps = steps
			}
		}

		if ySteps == 0 {
			found := true
			for i, moon := range moons {
				if moon.pos.Y != input[i].pos.Y || moon.vel.Y != input[i].vel.Y {
					found = false
					break
				}
			}
			if found {
				ySteps = steps
			}
		}

		if zSteps == 0 {
			found := true
			for i, moon := range moons {
				if moon.pos.Z != input[i].pos.Z || moon.vel.Z != input[i].vel.Z {
					found = false
					break
				}
			}
			if found {
				zSteps = steps
			}
		}
	}

	tempResult := LCM(xSteps, ySteps)
	result := LCM(tempResult, zSteps)
	fmt.Println(result)
	
}

func simulate(moons []Moon) {
	for ai, a := range moons {
		for bi, b := range moons {
			if bi == ai {
				continue
			}

			a.vel = a.vel.Add(b.pos.Sub(a.pos).Sign())
		}
		moons[ai] = a
	}

	for index, moon := range moons {
		moons[index].pos = moon.pos.Add(moon.vel)
	}
}

type Vector3 struct {
	X, Y, Z int
}

func (v Vector3) Add(ov Vector3) Vector3 {
	return Vector3{
		v.X + ov.X,
		v.Y + ov.Y,
		v.Z + ov.Z,
	}
}

func (v Vector3) Sub(ov Vector3) Vector3 {
	return Vector3{
		v.X - ov.X,
		v.Y - ov.Y,
		v.Z - ov.Z,
	}
}

func (v Vector3) Sign() Vector3 {
	return Vector3{
		sign(v.X),
		sign(v.Y),
		sign(v.Z),
	}
}

func (v Vector3) ManhattenLength() int {
	return abs(v.X) + abs(v.Y) + abs(v.Z)
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

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
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

func GCD(a, b int) int {
	if b == 0 {
		return a
	}

	return GCD(b, a%b)
}

func LCM(a, b int) int {
	return a / GCD(a, b) * b
}
