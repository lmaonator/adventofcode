package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

var lookup map[string]bool
var lookupAll map[string]int

func load(f string) ([]string, []string) {
	c, _ := os.ReadFile(f)
	s := strings.Split(string(c), "\n\n")
	patterns := []string{}
	for _, pattern := range strings.Split(s[0], ", ") {
		patterns = append(patterns, pattern)
	}
	slices.SortFunc(patterns, func(a, b string) int { return len(a) - len(b) })
	designs := []string{}
	for _, design := range strings.Split(strings.TrimSpace(s[1]), "\n") {
		designs = append(designs, design)
	}
	lookup = map[string]bool{}
	lookupAll = map[string]int{}
	return patterns, designs
}

func matchPatterns(patterns []string, design string) bool {
	if len(design) == 0 {
		return true
	}
	if v, ok := lookup[design]; ok {
		return v
	}
	for _, p := range patterns {
		if strings.HasPrefix(design, p) {
			if matchPatterns(patterns, design[len(p):]) {
				lookup[design] = true
				return true
			}
		}
	}
	lookup[design] = false
	return false
}

func part1(patterns []string, designs []string) int {
	r := 0
	for _, design := range designs {
		if matchPatterns(patterns, design) {
			r++
		}
	}
	return r
}

func matchPatternsAll(patterns []string, design string) int {
	if v, ok := lookup[design]; ok && !v {
		return 0
	}
	if v, ok := lookupAll[design]; ok {
		return v
	}
	r := 0
	for _, p := range patterns {
		if strings.HasPrefix(design, p) {
			if len(design[len(p):]) == 0 {
				r++
			} else {
				r += matchPatternsAll(patterns, design[len(p):])
			}
			lookupAll[design] = r
		}
	}
	lookupAll[design] = r
	return r
}

func part2(patterns []string, designs []string) int {
	r := 0
	for _, design := range designs {
		r += matchPatternsAll(patterns, design)
	}
	return r
}

func main() {
	p, d := load("example.txt")
	fmt.Println("Example 1:", part1(p, d))
	fmt.Println("Example 2:", part2(p, d))

	p, d = load("input.txt")
	fmt.Println("Part 1:", part1(p, d))
	fmt.Println("Part 2:", part2(p, d))
}
