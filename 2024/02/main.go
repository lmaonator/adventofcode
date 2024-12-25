package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reports := [][]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		split := strings.Split(line, " ")
		levels := make([]int, len(split))
		for i, level := range split {
			levelInt, err := strconv.Atoi(level)
			if err != nil {
				log.Fatalln(err)
			}
			levels[i] = levelInt
		}
		reports = append(reports, levels)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	part1(reports)
	part2(reports)
}

func isSafe(levels []int) bool {
	increasing := levels[0] < levels[1]
	for i := 0; i < len(levels)-1; i++ {
		curr := levels[i]
		next := levels[i+1]
		if (increasing && curr > next) || (!increasing && curr < next) {
			return false
		}
		diff := curr - next
		if diff < 0 {
			diff = -diff
		}
		if diff < 1 || diff > 3 {
			return false
		}
	}
	return true
}

func part1(reports [][]int) {
	safeReports := 0

	for _, report := range reports {
		if isSafe(report) {
			safeReports++
		}
	}

	fmt.Println("Part 1:", safeReports, "safe reports")
}

func part2(reports [][]int) {
	safeReports := 0

	for _, report := range reports {
		if isSafe(report) {
			safeReports++
			continue
		}

		for i := 0; i < len(report); i++ {
			dampened := make([]int, 0, len(report)-1)
			for n, val := range report {
				if n == i {
					continue
				}
				dampened = append(dampened, val)
			}
			if isSafe(dampened) {
				safeReports++
				break
			}
		}
	}

	fmt.Println("Part 2:", safeReports, "safe reports")
}
