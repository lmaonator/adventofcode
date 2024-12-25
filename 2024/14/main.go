package main

import (
	"fmt"
	"image"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Robot struct {
	pos, vel image.Point
}

func (r *Robot) Move(w, h, seconds int) {
	move := r.vel.Mul(seconds)
	pos := r.pos.Add(move)
	pos.X %= w
	pos.Y %= h
	if pos.X < 0 {
		pos.X += w
	}
	if pos.Y < 0 {
		pos.Y += h
	}
	r.pos = pos
}

var loadRe = regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

func load(f string) []Robot {
	c, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(c)), "\n")
	robots := make([]Robot, len(lines))
	for i, line := range lines {
		m := loadRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		robot := Robot{}
		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		robot.pos.X = x
		robot.pos.Y = y
		x, _ = strconv.Atoi(m[3])
		y, _ = strconv.Atoi(m[4])
		robot.vel.X = x
		robot.vel.Y = y
		robots[i] = robot
	}
	return robots
}

func main() {
	robots := load("example.txt")
	fmt.Println("Example 1:", part1(robots, 11, 7, 100))

	robots = load("input.txt")
	fmt.Println("Part 1:", part1(slices.Clone(robots), 101, 103, 100))

	seconds, robots := part2(robots, 101, 103)
	printGrid(robots, 101, 103)
	fmt.Println("Part 2:", seconds)
}

func part1(robots []Robot, w, h, seconds int) int {
	quadrants := [4]image.Rectangle{
		{image.Point{0, 0}, image.Point{w / 2, h / 2}},     // top left
		{image.Point{w/2 + 1, 0}, image.Point{w, h / 2}},   // top right
		{image.Point{0, h/2 + 1}, image.Point{w / 2, h}},   // bottom left
		{image.Point{w/2 + 1, h/2 + 1}, image.Point{w, h}}, // bottom right
	}
	qcount := [4]int{}
	for i, robot := range robots {
		robot.Move(w, h, seconds)
		robots[i] = robot
		for n := range 4 {
			if robot.pos.In(quadrants[n]) {
				qcount[n]++
			}
		}
	}
	return qcount[0] * qcount[1] * qcount[2] * qcount[3]
}

// Find the tree based on lowest danger level, hopefully..
func part2(robots []Robot, w, h int) (int, []Robot) {
	lowest := part1(robots, w, h, 1)
	lowestSec := 0
	lowestRobots := slices.Clone(robots)
	for sec := 2; sec < w*h; sec++ {
		curr := part1(robots, w, h, 1)
		if curr < lowest {
			lowest = curr
			lowestSec = sec
			lowestRobots = slices.Clone(robots)
		}
	}
	return lowestSec, lowestRobots
}

func printGrid(robots []Robot, w, h int) {
	rcount := map[image.Point]int{}
	for _, r := range robots {
		rcount[r.pos]++
	}
	for y := range h {
		for x := range w {
			if count := rcount[image.Point{x, y}]; count > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
