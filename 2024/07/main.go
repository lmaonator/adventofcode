package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type equation struct {
	result  int
	numbers []int
}

var equations []equation

func load(f string) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	equations = []equation{}
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ": ")
		r, err := strconv.Atoi(s[0])
		if err != nil {
			panic(err)
		}
		e := equation{result: r}
		s = strings.Split(s[1], " ")
		e.numbers = make([]int, 0, len(s))
		for _, nStr := range s {
			n, err := strconv.Atoi(nStr)
			if err != nil {
				panic(err)
			}
			e.numbers = append(e.numbers, n)
		}
		equations = append(equations, e)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	load("example.txt")
	fmt.Println("Example 1:", totalCalibrationResult(false))
	fmt.Println("Example 2:", totalCalibrationResult(true))

	load("input.txt")
	fmt.Println("Part 1:", totalCalibrationResult(false))
	fmt.Println("Part 2:", totalCalibrationResult(true))
}

func totalCalibrationResult(concat bool) int {
	result := 0
	for _, e := range equations {
		if solve(e.result, e.numbers[0], e.numbers[1:], concat) {
			result += e.result
		}
	}
	return result
}

func solve(result, curr int, remaining []int, concat bool) bool {
	if len(remaining) == 0 {
		return result == curr
	}
	next := remaining[0]

	temp := curr * next
	if temp <= result && solve(result, temp, remaining[1:], concat) {
		return true
	}

	temp = curr + next
	if temp <= result && solve(result, temp, remaining[1:], concat) {
		return true
	}

	if concat && len(remaining) >= 1 {
		curr = int(math.Pow(10, 1+math.Floor(math.Log10(float64(next)))))*curr + next
		if curr <= result && solve(result, curr, remaining[1:], concat) {
			return true
		}
	}
	return false
}
