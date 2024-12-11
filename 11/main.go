package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func load(f string) map[int]int {
	c, _ := os.ReadFile(f)
	s := strings.Split(strings.TrimSpace(string(c)), " ")
	r := map[int]int{}
	for _, numStr := range s {
		num, _ := strconv.Atoi(numStr)
		r[num]++
	}
	return r
}

func main() {
	stones := load("example.txt")
	fmt.Println("Example 1 result:", blink(stones, 6))

	stones = load("input.txt")
	fmt.Println("Part 1 result:", blink(stones, 25))
	fmt.Println("Part 2 result:", blink(stones, 75))
}

func blink(stones map[int]int, times int) int {
	for range times {
		next := map[int]int{}
		for stone, count := range stones {
			if stone == 0 {
				next[1] += count
			} else if numDigits(stone)%2 == 0 {
				numStr := strconv.Itoa(stone)
				left, _ := strconv.Atoi(numStr[:len(numStr)/2])
				right, _ := strconv.Atoi(numStr[len(numStr)/2:])
				next[left] += count
				next[right] += count
			} else {
				next[stone*2024] += count
			}
		}
		stones = next
	}
	r := 0
	for _, count := range stones {
		r += count
	}
	return r
}

func numDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return count
}
