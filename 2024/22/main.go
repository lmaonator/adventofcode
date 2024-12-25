package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func load(f string) []int {
	d, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(d)), "\n")
	secrets := make([]int, 0, len(lines))
	for _, line := range lines {
		num, _ := strconv.Atoi(line)
		secrets = append(secrets, num)
	}
	return secrets
}

var pruneNum = 16777216

func nextNumber(secret int) int {
	secret ^= secret * 64
	secret %= pruneNum

	secret ^= secret / 32
	secret %= pruneNum

	secret ^= secret * 2048
	secret %= pruneNum
	return secret
}

func part1(secrets []int) int {
	r := 0
	for _, s := range secrets {
		for range 2000 {
			s = nextNumber(s)
		}
		r += s
	}
	return r
}

func part2(secrets []int) int {
	bananas := int16(0)
	bananaMap := make([]int16, 19*19*19*19)
	seen := make([]bool, 19*19*19*19)
	for _, s := range secrets {
		last := s % 10
		key := 0
		for i := range 2000 {
			s = nextNumber(s)
			price := s % 10
			change := price - last
			key = (change+9)*19*19*19 + key/19
			if i >= 3 && !seen[key] {
				seen[key] = true
				bananaMap[key] += int16(price)
				bananas = max(bananas, bananaMap[key])
			}
			last = price
		}
		clear(seen)
	}
	return int(bananas)
}

func main() {
	secrets := load("example.txt")
	fmt.Println("Example 1:", part1(secrets))
	secrets = load("example2.txt")
	fmt.Println("Example 2:", part2(secrets))

	secrets = load("input.txt")
	fmt.Println("Part 1:", part1(secrets))
	fmt.Println("Part 2:", part2(secrets))
}
