package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	left := []int{}
	right := []int{}

	splitInput := strings.Split(string(input), "\n")
	for _, line := range splitInput {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			break
		}
		splitLine := strings.Split(line, " ")
		l, _ := strconv.Atoi(splitLine[0])
		left = append(left, l)
		r, _ := strconv.Atoi(splitLine[len(splitLine)-1])
		right = append(right, r)
	}

	fmt.Println("Part 1:", part1(left, right))
	fmt.Println("Part 2:", part2(left, right))
}

func part1(inLeft, inRight []int) int {
	left := slices.Clone(inLeft)
	right := slices.Clone(inRight)
	slices.Sort(left)
	slices.Sort(right)

	total := 0
	for i := range len(left) {
		distance := left[i] - right[i]
		// add generic Abs() already..
		if distance < 0 {
			distance = -distance
		}
		total += distance
	}
	return total
}

func part2(left, right []int) int {
	counts := map[int]int{}
	for _, value := range right {
		if count, exists := counts[value]; exists {
			counts[value] = count + 1
		} else {
			counts[value] = 1
		}
	}

	score := 0
	for _, value := range left {
		score += value * counts[value]
	}
	return score
}
