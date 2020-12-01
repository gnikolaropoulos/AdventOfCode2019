package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	input := readFile("input.txt")

	basepatern := []int{0, 1, 0, -1}
	var values []byte
	for _, char := range input {
		values = append(values, byte(char-'0'))
	}

	result := make([]byte, len(values))
	for p := 0; p < 100; p++ {
		for i := range values {
			sum := 0
			for j, c := range values {
				sum += int(c) * basepatern[(j+1)/(i+1)%4]
			}

			if sum < 0 {
				sum = -sum
			}

			result[i] = byte(sum % 10)
		}

		values = result
	}

	fmt.Println("--- Part One ---")
	fmt.Println(values[:8])

	signal := strings.Repeat(input, 10000)
	offset := 7
	result = []byte{}

	for _, c := range signal[offset:] {
		result = append(result, byte(c-'0'))
	}

	for p := 0; p < 100; p++ {
		sum := 0
		for i := len(result) - 1; i >= 0; i-- {
			sum += int(result[i])
			result[i] = byte(sum % 10)
		}
	}

	fmt.Println("--- Part Two ---")
	fmt.Println(result[:8])
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
