package main

import (
	"fmt"
	"image"
	"os"
	"strings"
)

type Tile int

const (
	Empty Tile = iota
	Wall
	Box
	DBoxL
	DBoxR
)

type Map struct {
	Grid        [][]Tile
	Robot       image.Point
	Moves       []image.Point
	CurrentMove int
}

func (m Map) PrintGrid() {
	for y, row := range m.Grid {
		for x, p := range row {
			if m.Robot.Eq(image.Point{x, y}) {
				fmt.Print("@")
				continue
			}
			switch p {
			case Empty:
				fmt.Print(".")
			case Wall:
				fmt.Print("#")
			case Box:
				fmt.Print("O")
			case DBoxL:
				fmt.Print("[")
			case DBoxR:
				fmt.Print("]")
			}
		}
		fmt.Println()
	}
}

func (m *Map) MoveRobot() bool {
	if m.CurrentMove == len(m.Moves) {
		return false
	}
	dir := m.Moves[m.CurrentMove]
	m.CurrentMove++
	newPos := m.Robot.Add(dir)
	switch m.Grid[newPos.Y][newPos.X] {
	case Empty:
		m.Robot = newPos
	case Wall:
		// do nothing
	case Box:
		if m.MoveBox(newPos, dir) {
			m.Robot = newPos
		}
	case DBoxL, DBoxR:
		if clr, set := m.MoveDBox(newPos, dir); clr != nil && set != nil {
			for _, p := range clr {
				m.Grid[p.Y][p.X] = Empty
				m.Grid[p.Y][p.X+1] = Empty
			}
			for _, p := range set {
				m.Grid[p.Y][p.X] = DBoxL
				m.Grid[p.Y][p.X+1] = DBoxR
			}
			m.Robot = newPos
		}
	default:
		panic("Invalid tile")
	}
	return true
}

func (m *Map) MoveBox(pos, dir image.Point) bool {
	newPos := pos.Add(dir)
	switch m.Grid[newPos.Y][newPos.X] {
	case Empty:
		m.Grid[pos.Y][pos.X] = Empty
		m.Grid[newPos.Y][newPos.X] = Box
		return true
	case Wall:
		return false
	case Box:
		if m.MoveBox(newPos, dir) {
			m.Grid[pos.Y][pos.X] = Empty
			m.Grid[newPos.Y][newPos.X] = Box
			return true
		}
		return false
	default:
		panic("Invalid tile")
	}
}

func (m *Map) MoveDBox(pos, dir image.Point) ([]image.Point, []image.Point) {
	var posL, posR image.Point
	if m.Grid[pos.Y][pos.X] == DBoxL {
		posL = pos
		posR = pos
		posR.X++
	} else {
		posR = pos
		posL = pos
		posL.X--
	}

	newPL, newPR := posL.Add(dir), posR.Add(dir)
	gridL, gridR := m.Grid[newPL.Y][newPL.X], m.Grid[newPR.Y][newPR.X]

	clr := []image.Point{posL}
	set := []image.Point{newPL}

	// horizontal
	if dir.Y == 0 {
		if (dir.X == -1 && gridL == Empty) || (dir.X == 1 && gridR == Empty) {
			return clr, set
		}
		if dir.X == -1 && gridL == DBoxR {
			if nclr, nset := m.MoveDBox(newPL, dir); nclr != nil && nset != nil {
				clr = append(clr, nclr...)
				set = append(set, nset...)
				return clr, set
			}
		}
		if dir.X == 1 && gridR == DBoxL {
			if nclr, nset := m.MoveDBox(newPR, dir); nclr != nil && nset != nil {
				clr = append(clr, nclr...)
				set = append(set, nset...)
				return clr, set
			}
		}
		return nil, nil
	}

	// rest vertical

	if gridL == Empty && gridR == Empty {
		return clr, set
	}

	if gridL == Wall || gridR == Wall {
		return nil, nil
	}

	// only single box has to be moved
	if gridL == DBoxL || (gridL == DBoxR && gridR == Empty) || (gridR == DBoxL && gridL == Empty) {
		var pos image.Point
		if gridL == DBoxL || gridL == DBoxR {
			pos = newPL
		} else {
			pos = newPR
		}
		if nclr, nset := m.MoveDBox(pos, dir); nclr != nil && nset != nil {
			clr = append(clr, nclr...)
			set = append(set, nset...)
			return clr, set
		}
		return nil, nil
	}

	// 2 boxes have to be moved
	if gridL == DBoxR && gridR == DBoxL {
		if lclr, lset := m.MoveDBox(newPL, dir); lclr != nil && lset != nil {
			if rclr, rset := m.MoveDBox(newPR, dir); rclr != nil && rset != nil {
				clr = append(clr, lclr...)
				clr = append(clr, rclr...)
				set = append(set, lset...)
				set = append(set, rset...)
				return clr, set
			}
		}
		return nil, nil
	}

	panic("Unreachable")
}

func (m *Map) MoveUntilEnd() {
	for m.MoveRobot() {
	}
}

func (m Map) BoxGPS(pos image.Point) int {
	return 100*pos.Y + pos.X
}

func (m Map) BoxGPSSum() int {
	r := 0
	for y, row := range m.Grid {
		for x, p := range row {
			if p == Box || p == DBoxL {
				r += m.BoxGPS(image.Point{x, y})
			}
		}
	}
	return r
}

func load(f string, doubleWidth bool) Map {
	d, _ := os.ReadFile(f)
	s := strings.Split(string(d), "\n\n")
	grid := strings.TrimSpace(s[0])
	moves := strings.TrimSpace(s[1])
	m := Map{}

	s = strings.Split(grid, "\n")
	m.Grid = make([][]Tile, len(s))
	for y, line := range s {
		var row []Tile
		if doubleWidth {
			row = make([]Tile, len(line)*2)
		} else {
			row = make([]Tile, len(line))
		}
		x := 0
		for _, r := range line {
			switch r {
			case '.':
				row[x] = Empty
				if doubleWidth {
					row[x+1] = Empty
				}
			case '#':
				row[x] = Wall
				if doubleWidth {
					row[x+1] = Wall
				}
			case 'O':
				if doubleWidth {
					row[x] = DBoxL
					row[x+1] = DBoxR
				} else {
					row[x] = Box
				}
			case '@':
				m.Robot = image.Point{x, y}
			default:
				panic("Invalid map")
			}
			if doubleWidth {
				x += 2
			} else {
				x++
			}
		}
		m.Grid[y] = row
	}

	m.Moves = make([]image.Point, 0, len(moves))
	for _, r := range moves {
		if r == '\n' {
			continue
		}
		var dir image.Point
		switch r {
		case '<':
			dir.X, dir.Y = -1, 0
		case '>':
			dir.X, dir.Y = 1, 0
		case '^':
			dir.X, dir.Y = 0, -1
		case 'v':
			dir.X, dir.Y = 0, 1
		default:
			panic("Invalid move")
		}
		m.Moves = append(m.Moves, dir)
	}

	return m
}

func main() {
	m := load("input.txt", false)
	m.MoveUntilEnd()
	fmt.Println("Part 1:", m.BoxGPSSum())

	m = load("input.txt", true)
	m.MoveUntilEnd()
	fmt.Println("Part 2:", m.BoxGPSSum())
}
