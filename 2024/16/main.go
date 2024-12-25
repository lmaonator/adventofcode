package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"slices"
	"strings"
)

func main() {
	maze := load("example.txt")
	score, paths := shortestPaths(maze, false)
	fmt.Println("Example Part 1:", score)
	score, paths = shortestPaths(maze, true)
	fmt.Println("Example Part 2:", numTilesInPaths(maze, paths))

	maze = load("input.txt")
	score, paths = shortestPaths(maze, false)
	fmt.Println("Part 1:", score)
	score, paths = shortestPaths(maze, true)
	fmt.Println("Part 2:", numTilesInPaths(maze, paths))
}

type Maze struct {
	grid       map[image.Point]bool
	start, end image.Point
}

func load(f string) Maze {
	c, _ := os.ReadFile(f)
	maze := Maze{
		grid: map[image.Point]bool{},
	}
	for y, line := range strings.Split(string(c), "\n") {
		if len(line) == 0 {
			break
		}
		for x, r := range line {
			p := image.Point{x, y}
			switch r {
			case '.':
				maze.grid[p] = true
			case 'S':
				maze.grid[p] = true
				maze.start = p
			case 'E':
				maze.grid[p] = true
				maze.end = p
			}
		}
	}
	return maze
}

var dirs = []image.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

type SeenKey struct {
	p   image.Point
	dir int
}
type Tile struct {
	p          image.Point
	dir, score int
	history    []*Tile
}

func (t Tile) Key() SeenKey {
	return SeenKey{t.p, t.dir}
}

type shortPath struct {
	score int
	path  []*Tile
}

func shortestPaths(m Maze, relaxed bool) (int, []shortPath) {
	offset := 0
	if relaxed {
		offset = 1000
	}
	paths := []shortPath{}
	queue := []Tile{{m.start, 0, 0, []*Tile{}}}
	seen := map[SeenKey]int{{m.start, 0}: 0}
	for len(queue) > 0 {
		tile := queue[0]
		queue = queue[1:]
		if tile.p == m.end {
			paths = append(paths, shortPath{tile.score, tile.history})
			continue
		}
		if v, ok := seen[tile.Key()]; ok && (v+offset < tile.score) {
			continue
		}
		for _, i := range []int{0, 1, 3} {
			nt := Tile{}
			nt.dir = (tile.dir + i) % 4
			nt.p = tile.p.Add(dirs[nt.dir])
			if !m.grid[nt.p] {
				continue
			}
			nt.score = tile.score
			if nt.dir == tile.dir {
				nt.score += 1
			} else {
				nt.score += 1000
				nt.p = tile.p
			}
			if v, ok := seen[nt.Key()]; !ok || (ok && nt.score < v+offset) {
				nt.history = append(slices.Clone(tile.history), &tile)
				seen[nt.Key()] = nt.score
				queue = append(queue, nt)
			}
		}
	}

	score := math.MaxInt64
	for i := range dirs {
		if v, ok := seen[SeenKey{m.end, i}]; ok {
			score = min(score, v)
		}
	}
	return score, paths
}

func numTilesInPaths(m Maze, paths []shortPath) int {
	tiles := map[image.Point]bool{}
	shortest := math.MaxInt64
	for _, path := range paths {
		shortest = min(shortest, path.score)
	}
	for _, path := range paths {
		if path.score > shortest {
			continue
		}
		for _, tile := range path.path {
			tiles[tile.p] = true
		}
	}
	tiles[m.end] = true
	return len(tiles)
}
