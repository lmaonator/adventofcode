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
	bananas := 0
	bananaMap := map[[4]int]int{}
	for _, s := range secrets {
		last := s % 10
		seen := map[[4]int]struct{}{}
		mapKey := [4]int{}
		for i := range 2000 {
			s = nextNumber(s)
			price := s % 10
			change := price - last

			mapKey[0] = mapKey[1]
			mapKey[1] = mapKey[2]
			mapKey[2] = mapKey[3]
			mapKey[3] = change
			if i >= 3 {
				if _, ok := seen[mapKey]; !ok {
					seen[mapKey] = struct{}{}
					bananaMap[mapKey] += price
					bananas = max(bananas, bananaMap[mapKey])
				}
			}
			last = price
		}
	}
	return bananas
}

func Intersection(a, b map[[4]int]int) (map[[4]int]int, int) {
	highest := 0
	r := map[[4]int]int{}
	for k, va := range a {
		bananas := va + b[k]
		r[k] = bananas
		highest = max(highest, bananas)
	}
	return r, highest
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
