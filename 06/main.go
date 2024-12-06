package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
)

var printGrid bool

func main() {
	printGrid = len(os.Args) > 1 && os.Args[1] == "print"

	m := load("example1.txt")
	part1(m)
	m = load("input.txt")
	part1(m)
	m = load("input.txt")
	part2(m)
}

func load(filename string) Map {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	return ParseMap(file)
}

func part1(m Map) {
	m.MoveUntilDone()
	if printGrid {
		fmt.Println(m)
	}
	fmt.Println("Part 1 - Distinct positions visited:", len(m.Visited))
}

func part2(m Map) {
	startPos := m.pos
	startDir := m.dir
	for !m.MoveOnce() {
	}

	visited := maps.Clone(m.Visited)

	type Loc struct {
		coord Coord
		dir   Direction
	}

	count := 0

	for curr := range visited {
		m.pos = startPos
		m.dir = startDir
		m.grid[curr.y][curr.x] = Obstruction

		history := map[Loc]struct{}{}
		for !m.MoveOnce() {
			loc := Loc{m.pos, m.dir}
			if _, exists := history[loc]; exists {
				count++
				break
			}
			history[loc] = struct{}{}
		}

		m.grid[curr.y][curr.x] = Empty
	}

	fmt.Println("Part 2 - Number of positions for obstructions:", count)
}

type Coord struct {
	x, y int
}

type Tile int

const (
	Empty Tile = iota
	Obstruction
	Visited
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Map struct {
	grid    [][]Tile
	pos     Coord
	dir     Direction
	Visited map[Coord]struct{}
}

func (m Map) String() string {
	str := make([]byte, 0, len(m.grid)*len(m.grid[0]))
	for y, row := range m.grid {
		line := make([]byte, len(row)+1)
		for x, tile := range row {
			if m.pos.y == y && m.pos.x == x {
				switch m.dir {
				case Up:
					line[x] = '^'
				case Right:
					line[x] = '>'
				case Down:
					line[x] = 'v'
				case Left:
					line[x] = '<'
				}
				continue
			}
			switch tile {
			case Empty:
				line[x] = '.'
			case Obstruction:
				line[x] = '#'
			case Visited:
				line[x] = 'X'
			}
		}
		if y < len(m.grid)-1 {
			line[len(line)-1] = '\n'
		}
		str = append(str, line...)
	}
	return string(str)
}

func (m *Map) MoveUntilDone() {
	for !m.MoveOnce() {
	}
}

func (m *Map) MoveOnce() (done bool) {
	g := m.grid
	x, y := m.pos.x, m.pos.y
	m.Visited[m.pos] = struct{}{}
	g[y][x] = Visited
	switch m.dir {
	case Up:
		if y == 0 {
			return true
		}
		y = y - 1
		if g[y][x] == Obstruction {
			m.dir = Right
		} else {
			m.pos.y = y
		}
	case Right:
		if x == len(g[y])-1 {
			return true
		}
		x = x + 1
		if g[y][x] == Obstruction {
			m.dir = Down
		} else {
			m.pos.x = x
		}
	case Down:
		if y == len(g)-1 {
			return true
		}
		y = y + 1
		if g[y][x] == Obstruction {
			m.dir = Left
		} else {
			m.pos.y = y
		}
	case Left:
		if x == 0 {
			return true
		}
		x = x - 1
		if g[y][x] == Obstruction {
			m.dir = Up
		} else {
			m.pos.x = x
		}
	}
	return false
}

func ParseMap(input io.Reader) Map {
	m := Map{}
	m.Visited = map[Coord]struct{}{}
	scanner := bufio.NewScanner(input)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		m.grid = append(m.grid, []Tile{})
		for x, r := range line {
			switch r {
			case '.':
				m.grid[y] = append(m.grid[y], Empty)
			case '#':
				m.grid[y] = append(m.grid[y], Obstruction)
			case '^':
				m.grid[y] = append(m.grid[y], Empty)
				m.dir = Up
				m.pos = Coord{x, y}
			case '>':
				m.grid[y] = append(m.grid[y], Empty)
				m.dir = Right
				m.pos = Coord{x, y}
			case 'v':
				m.grid[y] = append(m.grid[y], Empty)
				m.dir = Down
				m.pos = Coord{x, y}
			case '<':
				m.grid[y] = append(m.grid[y], Empty)
				m.dir = Left
				m.pos = Coord{x, y}
			default:
				log.Fatalln("Invalid tile:", r)
			}
		}
		y++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return m
}
