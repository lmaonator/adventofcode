package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

type Claw struct {
	A, B, Prize Point
}

func load(f string) []Claw {
	claws := []Claw{}
	c, _ := os.ReadFile(f)

	splitLine := func(l string) Point {
		s := strings.Split(l, " ")
		i := len(s) - 2
		x, _ := strconv.Atoi(s[i][2 : len(s[i])-1])
		y, _ := strconv.Atoi(s[i+1][2:])
		return Point{x, y}
	}

	for _, s := range strings.Split(string(c), "\n\n") {
		lines := strings.Split(s, "\n")
		claw := Claw{
			splitLine(lines[0]),
			splitLine(lines[1]),
			splitLine(lines[2]),
		}
		claws = append(claws, claw)
	}
	return claws
}

func main() {
	claws := load("example.txt")
	fmt.Println("Example 1:", getTokens(claws, 0))

	claws = load("input.txt")
	fmt.Println("Part 1:", getTokens(claws, 0))
	fmt.Println("Part 2:", getTokens(claws, 10000000000000))
}

func getTokens(claws []Claw, offset int) int {
	tokens := 0
	for _, claw := range claws {
		tokens += solveForCost(claw, offset)
	}
	return tokens
}

func solveForCost(c Claw, offset int) int {
	pY, pX := c.Prize.Y+offset, c.Prize.X+offset
	div := c.A.X*c.B.Y - c.A.Y*c.B.X
	A := pX*c.B.Y - pY*c.B.X
	if A%div != 0 {
		return 0
	}
	A = A / div
	B := pY*c.A.X - pX*c.A.Y
	if B%div != 0 {
		return 0
	}
	B = B / div
	return A*3 + B
}
