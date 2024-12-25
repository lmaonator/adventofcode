package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
)

func main() {
	example1()
	example2()

	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	part1(input)
	part2(input)
}

func example1() {
	input, err := os.ReadFile("example1.txt")
	if err != nil {
		panic(err)
	}
	grid := newGrid(input)
	count := find(grid, []byte("XMAS"))
	fmt.Println("Example 1 result:", count)
}

func example2() {
	input, err := os.ReadFile("example2.txt")
	if err != nil {
		panic(err)
	}
	grid := newGrid(input)
	count := findMASCross(grid)
	fmt.Println("Example 2 result:", count)
}

func part1(input []byte) {
	grid := newGrid(input)
	count := find(grid, []byte("XMAS"))
	fmt.Println("Part 1 result:", count)
}

func part2(input []byte) {
	grid := newGrid(input)
	count := findMASCross(grid)
	fmt.Println("Part 2 result:", count)
}

type grid struct {
	width  int
	height int
	grid   []byte
}

func newGrid(input []byte) grid {
	g := grid{}
	g.width = bytes.IndexByte(input, '\n')
	for _, b := range input {
		if b != '\n' {
			g.grid = append(g.grid, b)
		}
	}
	g.height = len(g.grid) / g.width
	return g
}

func find(g grid, needle []byte) int {
	needleR := slices.Clone(needle)
	slices.Reverse(needleR)
	size := len(needle)
	count := 0

	part := make([]byte, size)

	for i := range g.height {
		for n := range g.width {
			pos := i*g.width + n
			if g.grid[pos] != needle[0] && g.grid[pos] != needleR[0] {
				continue
			}

			// horizontal
			if n <= g.width-size {
				part := g.grid[pos : pos+size]
				if bytes.Equal(part, needle) || bytes.Equal(part, needleR) {
					count++
				}
			}
			// vertical
			last := pos + (size-1)*g.width
			if last < len(g.grid) {
				for x := range size {
					part[x] = g.grid[pos+x*g.width]
				}
				if bytes.Equal(part, needle) || bytes.Equal(part, needleR) {
					count++
				}
			}
			// diagonal forward
			if i <= g.height-size && n <= g.width-size {
				for x := range size {
					part[x] = g.grid[pos+x*g.width+x]
				}
				if bytes.Equal(part, needle) || bytes.Equal(part, needleR) {
					count++
				}
			}
			// diagonal backwards
			if i <= g.height-size && n >= size-1 {
				for x := range size {
					part[x] = g.grid[pos+x*g.width-x]
				}
				if bytes.Equal(part, needle) || bytes.Equal(part, needleR) {
					count++
				}
			}
		}
	}
	return count
}

func findMASCross(g grid) int {
	isMAS := func(b []byte) bool {
		return bytes.Equal([]byte("MAS"), b) || bytes.Equal([]byte("SAM"), b)
	}
	size := 3
	count := 0

	part := make([]byte, size)

	for i := range g.height - size + 1 {
		for n := range g.width {
			if n > g.width-size {
				continue
			}
			pos := i*g.width + n
			// diagonal forward
			for x := range size {
				part[x] = g.grid[pos+x*g.width+x]
			}
			if isMAS(part) {
				// diagonal backwards
				for x := range size {
					part[x] = g.grid[pos+x*g.width-x+size-1]
				}
				if isMAS(part) {
					count++
				}
			}
		}
	}
	return count
}
