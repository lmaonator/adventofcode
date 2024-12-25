package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

var m [][]int
var w, h int
var trailheads []point
var seenPeaks map[point]struct{}

func load(f string) {
	d, _ := os.ReadFile(f)
	s := strings.Split(strings.TrimSpace(string(d)), "\n")
	h = len(s)
	m = make([][]int, h)
	trailheads = []point{}
	for i, l := range s {
		m[i] = make([]int, len(l))
		for n, r := range l {
			num, err := strconv.Atoi(string(r))
			// for examples with .
			if err != nil {
				num = -1
			}
			m[i][n] = num
			if num == 0 {
				trailheads = append(trailheads, point{n, i})
			}
		}
	}
	w = len(m[0])
}

func main() {
	load("example.txt")
	fmt.Println("Example result:", countValidTrails(false))

	load("input.txt")
	fmt.Println("Part 1 result:", countValidTrails(false))
	fmt.Println("Part 2 result:", countValidTrails(true))
}

var dirs = []point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func countValidTrails(allRoutes bool) int {
	r := 0
	for _, th := range trailheads {
		seenPeaks = map[point]struct{}{}
		r += traverse(th, allRoutes)
	}
	return r
}

func traverse(curr point, allRoutes bool) int {
	count := 0
	for _, dir := range dirs {
		next := point{curr.x + dir.x, curr.y + dir.y}
		if next.x < 0 || next.x >= w || next.y < 0 || next.y >= h {
			continue
		}
		if m[next.y][next.x] == m[curr.y][curr.x]+1 {
			if m[next.y][next.x] == 9 {
				if allRoutes {
					count++
					continue
				} else if _, seen := seenPeaks[next]; !seen {
					seenPeaks[next] = struct{}{}
					count++
					continue
				}
			}
			count += traverse(next, allRoutes)
		}
	}
	return count
}
