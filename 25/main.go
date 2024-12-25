package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

var lockHeight int

func load(f string) (locks, keys [][]int) {
	c, _ := os.ReadFile(f)
	blocks := strings.Split(strings.TrimSpace(string(c)), "\n\n")
	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		lockHeight = len(lines)
		isLock := lines[0] == strings.Repeat("#", len(lines[0]))
		if !isLock {
			slices.Reverse(lines)
		}
		curr := make([]int, len(lines[0]))
		for _, row := range lines[1:] {
			for n, r := range row {
				if r == '#' {
					curr[n]++
				}
			}
		}
		if isLock {
			locks = append(locks, curr)
		} else {
			keys = append(keys, curr)
		}
	}
	return locks, keys
}

func doesFit(lock, key []int) bool {
	for i, lv := range lock {
		kv := key[i]
		if lv+kv > lockHeight-2 {
			return false
		}
	}
	return true
}

func part1(locks, keys [][]int) int {
	r := 0
	for _, lock := range locks {
		for _, key := range keys {
			if doesFit(lock, key) {
				r++
			}
		}
	}
	return r
}

func main() {
	l, k := load("example.txt")
	fmt.Println("Example:", part1(l, k))

	l, k = load("input.txt")
	fmt.Println("Part 1:", part1(l, k))
}
