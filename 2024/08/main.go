package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

// Globals lets goooooo
var antennas map[rune][]Point
var width, height int
var anodes map[Point]struct{}

func load(f string) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	antennas = map[rune][]Point{}
	width, height = 0, 0
	anodes = map[Point]struct{}{}

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		l := scanner.Text()
		for x, r := range l {
			if r != '.' {
				if _, ok := antennas[r]; !ok {
					antennas[r] = []Point{}
				}
				antennas[r] = append(antennas[r], Point{x, y})
			}
		}
		y++
		height = max(height, y)
		width = max(width, len(l))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func printAntinodeMap() {
	line := make([]byte, width)
	for y := range height {
		for x := range width {
			if _, ok := anodes[Point{x, y}]; ok {
				line[x] = '#'
			} else {
				line[x] = '.'
			}
		}
		fmt.Println(string(line))
	}
}

func main() {
	load("example.txt")
	fmt.Println("Example 1 - Antinodes:", part1())
	load("input.txt")
	fmt.Println("Part 1 - Antinodes:", part1())

	load("example.txt")
	fmt.Println("Example 2 - Antinodes:", part2())
	load("input.txt")
	fmt.Println("Part 2 - Antinodes:", part2())
}

func part1() int {
	for _, v := range antennas {
		for i, a1 := range v {
			for _, a2 := range v[i+1:] {
				findAntiNodes(a1, a2)
			}
		}
	}
	return len(anodes)
}

func findAntiNodes(a1, a2 Point) {
	dx, dy := a2.x-a1.x, a2.y-a1.y
	an := Point{a1.x - dx, a1.y - dy}
	if inBounds(an) {
		anodes[an] = struct{}{}
	}
	an = Point{a2.x + dx, a2.y + dy}
	if inBounds(an) {
		anodes[an] = struct{}{}
	}
}

func part2() int {
	for _, v := range antennas {
		for i, a1 := range v {
			for _, a2 := range v[i+1:] {
				findAntiNodesAny(a1, a2)
			}
		}
	}
	return len(anodes)
}

func findAntiNodesAny(a1, a2 Point) {
	dx, dy := a2.x-a1.x, a2.y-a1.y

	an := a1
	for inBounds(an) {
		anodes[an] = struct{}{}
		an.x -= dx
		an.y -= dy
	}

	an = a2
	for inBounds(an) {
		anodes[an] = struct{}{}
		an.x += dx
		an.y += dy
	}
}

func inBounds(p Point) bool {
	return p.x >= 0 && p.x < width && p.y >= 0 && p.y < height
}
