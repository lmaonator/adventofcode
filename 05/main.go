package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	example()

	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	manuals := parseInput(file)

	r := addCorrectMiddlePages(manuals)
	fmt.Println("Part 1 result:", r)
	r = fixAndAddInccorectPages(manuals)
	fmt.Println("Part 2 result:", r)
}

func example() {
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	manuals := parseInput(file)
	r := addCorrectMiddlePages(manuals)
	fmt.Println("Example 1 result:", r)
	r = fixAndAddInccorectPages(manuals)
	fmt.Println("Example 2 result:", r)
}

type manuals struct {
	rules   map[string][]string
	updates [][]string
}

func newManuals() manuals {
	return manuals{
		map[string][]string{},
		[][]string{},
	}
}

func parseInput(r io.Reader) manuals {
	scanner := bufio.NewScanner(r)
	m := newManuals()
	var readUpdates bool
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			readUpdates = true
			continue
		}
		if !readUpdates {
			split := strings.Split(line, "|")
			m.rules[split[0]] = append(m.rules[split[0]], split[1])
		} else {
			split := strings.Split(line, ",")
			m.updates = append(m.updates, split)
		}
	}
	return m
}

func addCorrectMiddlePages(m manuals) int {
	updates := filterUpdates(m, false)
	return addMiddlePages(updates)
}

func fixAndAddInccorectPages(m manuals) int {
	updates := filterUpdates(m, true)

	for uIndex, u := range updates {
		pageIndexes := map[string]int{}
		currentPageIndex := -1
		for {
			currentPageIndex++
			if currentPageIndex >= len(u) {
				break
			}
			currentPage := u[currentPageIndex]

			for i, p := range u {
				pageIndexes[p] = i
			}

			rules, exists := m.rules[currentPage]
			if !exists {
				continue
			}
			moveTo := -1
			for _, ruleNum := range rules {
				index, exists := pageIndexes[ruleNum]
				if !exists {
					continue
				}
				if currentPageIndex < index {
					continue
				}
				moveTo = index
			}
			if moveTo > -1 {
				u[moveTo], u[currentPageIndex] = u[currentPageIndex], u[moveTo]
				currentPageIndex = -1
				continue
			}
		}
		updates[uIndex] = u
	}

	return addMiddlePages(updates)
}

func addMiddlePages(updates [][]string) int {
	r := 0
	for _, u := range updates {
		v, err := strconv.Atoi(u[len(u)/2])
		if err != nil {
			panic(err)
		}
		r += v
	}
	return r
}

func filterUpdates(m manuals, incorrect bool) [][]string {
	result := [][]string{}
outer:
	for _, update := range m.updates {
		seenPages := map[string]struct{}{}
		for _, p := range update {
			if x, exists := m.rules[p]; exists {
				for _, y := range x {
					if _, exists := seenPages[y]; exists {
						if incorrect {
							result = append(result, update)
						}
						continue outer
					}
				}
			}
			seenPages[p] = struct{}{}
		}
		if !incorrect {
			result = append(result, update)
		}
	}
	return result
}
