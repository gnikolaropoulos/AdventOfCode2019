package main

import (
	"bufio"
	"fmt"
	"math/big"
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

type BigShuffle struct {
	shuffleType ShuffleType
	arg         int64
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

	fmt.Println("--- Part Two ---")

	const bigCount = 119315717514047
	const iterations = 101741582076661

	var bigShuffles []BigShuffle
	newShuffles := readBigShuffles(lines)
	factor := compact(newShuffles, bigCount)

	for iterationsLeft := bigCount - iterations - 1; iterationsLeft != 0; iterationsLeft /= 2 {
		if iterationsLeft%2 == 1 {
			bigShuffles = append(bigShuffles, factor...)
			bigShuffles = compact(bigShuffles, bigCount)
		}

		factor = append(factor, factor...)
		factor = compact(factor, bigCount)
	}

	pos := big.NewInt(2020)
	for _, shuffle := range bigShuffles {
		if shuffle.shuffleType == ShuffleTypeDeal {
			increment := shuffle.arg
			pos.Mul(pos, big.NewInt(increment))
			pos.Mod(pos, big.NewInt(bigCount))

		} else if shuffle.shuffleType == ShuffleTypeStack {
			pos.Sub(big.NewInt(bigCount-1), pos)

		} else if shuffle.shuffleType == ShuffleTypeCut {
			cut := shuffle.arg

			if pos.Int64() < cut {
				pos.Add(pos, big.NewInt(bigCount-cut))
			} else {
				pos.Sub(pos, big.NewInt(cut))
			}
		}
	}

	fmt.Println(pos.Int64())

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

func readBigShuffles(lines []string) []BigShuffle {
	result := make([]BigShuffle, len(lines))
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		shuffle := strings.Join(tokens[:len(tokens)-1], " ")
		switch shuffle {
		case "deal with increment":
			incr, err := strconv.Atoi(tokens[len(tokens)-1])
			if err != nil {
				panic(err)
			}

			result = append(result, BigShuffle{ShuffleTypeDeal, int64(incr)})

		case "cut":
			cut, err := strconv.Atoi(tokens[len(tokens)-1])
			if err != nil {
				panic(err)
			}

			result = append(result, BigShuffle{ShuffleTypeCut, int64(cut)})

		case "deal into new":
			result = append(result, BigShuffle{ShuffleTypeStack, 0})

		default:
			panic("unknown command")
		}
	}

	return result
}

func compact(input []BigShuffle, count int64) []BigShuffle {
	{
		compacted := make([]BigShuffle, 0, len(input))
		reverse := false
		for _, shuffle := range input {
			if shuffle.shuffleType == ShuffleTypeStack {
				reverse = !reverse
				continue
			}
			if !reverse {
				compacted = append(compacted, shuffle)
				continue
			}
			switch shuffle.shuffleType {
			case ShuffleTypeDeal:
				compacted = append(compacted, shuffle)
				compacted = append(compacted, BigShuffle{ShuffleTypeCut, count + 1 - shuffle.arg})

			case ShuffleTypeCut:
				cut := (shuffle.arg + count) % count
				cut = count - cut
				compacted = append(compacted, BigShuffle{ShuffleTypeCut, cut})
			}
		}

		if reverse {
			compacted = append(compacted, BigShuffle{ShuffleTypeStack, 0})
		}

		input = compacted
	}

	{
		compacted := make([]BigShuffle, 0, len(input))
		cut := big.NewInt(0)
		for _, shuffle := range input {
			switch shuffle.shuffleType {
			case ShuffleTypeStack:
				if value := cut.Int64(); value != 0 {
					compacted = append(compacted, BigShuffle{ShuffleTypeCut, value})
					cut.SetInt64(0)
				}

				compacted = append(compacted, shuffle)

			case ShuffleTypeDeal:
				compacted = append(compacted, shuffle)
				cut.Mul(cut, big.NewInt(shuffle.arg))
				cut.Mod(cut, big.NewInt(count))

			case ShuffleTypeCut:
				cut.Add(cut, big.NewInt(shuffle.arg))
				cut.Mod(cut, big.NewInt(count))
			}
		}

		if value := cut.Int64(); value != 0 {
			compacted = append(compacted, BigShuffle{ShuffleTypeCut, value})
			cut.SetInt64(0)
		}

		input = compacted
	}

	{
		compacted := make([]BigShuffle, 0, len(input))
		increment := big.NewInt(1)
		for _, shuffle := range input {
			switch shuffle.shuffleType {
			case ShuffleTypeDeal:
				increment.Mul(increment, big.NewInt(shuffle.arg))
				increment.Mod(increment, big.NewInt(count))

			default:
				if value := increment.Int64(); value != 1 {
					compacted = append(compacted, BigShuffle{ShuffleTypeDeal, value})
					increment.SetInt64(1)
				}
				compacted = append(compacted, shuffle)
			}
		}
		if value := increment.Int64(); value != 1 {
			compacted = append(compacted, BigShuffle{ShuffleTypeDeal, value})
			increment.SetInt64(1)
		}

		input = compacted
	}

	return input
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
