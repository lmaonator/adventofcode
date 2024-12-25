package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Panicln(err)
	}
	dataStr := string(data)

	part1(dataStr)
	part2(dataStr)
}

func part1(data string) {
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	result := 0

	for _, match := range matches {
		a, err := strconv.Atoi(match[1])
		if err != nil {
			log.Panicln(err)
		}
		b, err := strconv.Atoi(match[2])
		if err != nil {
			log.Panicln(err)
		}

		result += a * b
	}

	fmt.Println("Part 1 - Result:", result)
}

func part2(data string) {
	re := regexp.MustCompile(`(mul|do|don't)\(((\d+),(\d+)|)\)`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	result := 0
	enabled := true

	for _, match := range matches {
		switch match[1] {
		case "do":
			enabled = true
		case "don't":
			enabled = false
		case "mul":
			if !enabled {
				break
			}
			a, err := strconv.Atoi(match[3])
			if err != nil {
				log.Panicln(err)
			}
			b, err := strconv.Atoi(match[4])
			if err != nil {
				log.Panicln(err)
			}
			result += a * b
		}
	}

	fmt.Println("Part 2 - Result:", result)
}
