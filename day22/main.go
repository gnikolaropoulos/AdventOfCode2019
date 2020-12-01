package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ShuffleType int

const (
	ShuffleTypeCut ShuffleType = iota
	ShuffleTypeDeal
	ShuffleTypeStack
)

type Shuffle struct {
	shuffleType ShuffleType
	arg         int
}

func main() {
	lines := readLines("input.txt")
	shuffles := readShuffles(lines)
	const count = 10007
	deck := createDeck(count)
	for _, shuffle := range shuffles {
		switch shuffle.shuffleType {
		case ShuffleTypeCut:
			cutShuffle(deck, shuffle.arg)
		case ShuffleTypeDeal:
			dealShuffle(deck, shuffle.arg)
		case ShuffleTypeStack:
			stackShuffle(deck)
		default:
			panic("unknown shuffle type")
		}
	}

	fmt.Println("--- Part One ---")
	for index, card := range deck {
		if card == 2019 {
			fmt.Println(index)
			break
		}
	}

}

func cutShuffle(deck []int, n int) {
	var newDeck []int
	if n > 0 {
		newDeck = append(deck[n:], deck[:n]...)
	} else {
		n = abs(n)
		newDeck = append(deck[len(deck)-n:], deck[:len(deck)-n]...)
	}

	copy(deck, newDeck)
}

func dealShuffle(deck []int, n int) {
	newDeck := make([]int, len(deck))
	increment := 0
	for i := 0; i < len(deck); i++ {
		if newDeck[increment] != 0 {
			panic("dealing with increment overrides a position")
		}

		newDeck[increment] = deck[i]
		increment = (increment + n) % len(deck)
	}

	copy(deck, newDeck)
}

func stackShuffle(deck []int) {
	newStack := make([]int, len(deck))
	for i := 0; i < len(deck); i++ {
		newStack[len(newStack)-i-1] = deck[i]
	}

	copy(deck, newStack)
}

func readShuffles(lines []string) []Shuffle {
	result := make([]Shuffle, len(lines))
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		shuffle := strings.Join(tokens[:len(tokens)-1], " ")
		switch shuffle {
		case "deal with increment":
			incr, err := strconv.Atoi(tokens[len(tokens)-1])
			if err != nil {
				panic(err)
			}

			result = append(result, Shuffle{ShuffleTypeDeal, incr})
		case "cut":
			cut, err := strconv.Atoi(tokens[len(tokens)-1])
			if err != nil {
				panic(err)
			}

			result = append(result, Shuffle{ShuffleTypeCut, cut})
		case "deal into new":
			result = append(result, Shuffle{ShuffleTypeStack, 0})
		default:
			panic("unknown command")
		}
	}

	return result
}

func createDeck(count int) []int {
	deck := []int{}
	for i := 0; i < count; i++ {
		deck = append(deck, i)
	}

	return deck
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
