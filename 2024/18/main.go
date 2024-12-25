package main

import (
	"container/heap"
	"fmt"
	"image"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func load(f string) []image.Point {
	c, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(c)), "\n")
	points := make([]image.Point, len(lines))
	for i, line := range lines {
		s := strings.Split(line, ",")
		x, _ := strconv.Atoi(s[0])
		y, _ := strconv.Atoi(s[1])
		points[i] = image.Point{x, y}
	}
	return points
}

type MemSpace struct {
	grid        map[image.Point]bool
	falling     []image.Point
	fallIndex   int
	wh          int
	start, goal image.Point
	bounds      image.Rectangle
}

func (m *MemSpace) Drop(count int) image.Point {
	for range count {
		if m.fallIndex >= len(m.falling) {
			panic("No more falling bytes")
		}
		m.grid[m.falling[m.fallIndex]] = false
		m.fallIndex++
	}
	return m.falling[m.fallIndex-1]
}

var dirs = []image.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func (m MemSpace) Neighbours(point image.Point) []image.Point {
	r := []image.Point{}
	for _, dir := range dirs {
		p := point.Add(dir)
		if p.In(m.bounds) && m.grid[p] {
			r = append(r, p)
		}
	}
	return r
}

func (m MemSpace) NewGraph() Graph {
	g := Graph{
		Nodes: []*Item{},
		Edges: map[*Item][]*Edge{},
	}
	nodes := map[image.Point]*Item{}
	for y := range m.wh {
		for x := range m.wh {
			src := image.Pt(x, y)
			if !m.grid[src] {
				continue
			}
			if _, ok := nodes[src]; !ok {
				n := Item{value: src}
				nodes[src] = &n
				g.AddNode(&n)
			}
			for _, nb := range m.Neighbours(src) {
				if _, ok := nodes[nb]; !ok {
					n := Item{value: nb}
					nodes[nb] = &n
					g.AddNode(&n)
				}
				g.AddEdge(nodes[src], nodes[nb], 1)
			}
		}
	}
	return g
}

func NewMemSpace(wh int, falling []image.Point) MemSpace {
	m := MemSpace{
		grid:      map[image.Point]bool{},
		falling:   falling,
		fallIndex: 0,
		wh:        wh,
		start:     image.Pt(0, 0),
		goal:      image.Pt(wh-1, wh-1),
		bounds:    image.Rectangle{image.Pt(0, 0), image.Pt(wh, wh)},
	}
	for y := range wh {
		for x := range wh {
			m.grid[image.Pt(x, y)] = true
		}
	}
	return m
}

func (m MemSpace) getShortestPath() (length int, path []image.Point) {
	graph := m.NewGraph()
	dist := map[image.Point]int{}
	prev := map[image.Point]image.Point{}
	pq := PriorityQueue{}

	for _, v := range graph.Nodes {
		if v.value == m.start {
			dist[m.start] = 0
			v.priority = 0
			pq.Push(v)
		} else {
			dist[v.value] = math.MaxInt
			v.priority = math.MaxInt
			pq.Push(v)
		}
	}
	heap.Init(&pq)

	for len(pq) > 0 {
		u := heap.Pop(&pq).(*Item)

		if u.value == m.goal {
			break
		}

		for _, v := range graph.Edges[u] {
			alt := dist[u.value] + v.Weight
			if alt < dist[v.Node.value] && alt >= 0 {
				prev[v.Node.value] = u.value
				dist[v.Node.value] = alt
				pq.updatePriority(v.Node, alt)
			}
		}
	}

	if u, ok := prev[m.goal]; ok || m.goal == m.start {
		path := make([]image.Point, 0, dist[m.goal])
		for ok {
			path = slices.Insert(path, 0, u)
			u, ok = prev[u]
		}
		return dist[m.goal], path
	}
	return -1, nil
}

func part2(m MemSpace, path []image.Point) image.Point {
	for {
		fallen := m.Drop(1)
		index := slices.Index(path, fallen)
		if index < 0 {
			continue
		}
		length, newPath := m.getShortestPath()
		if length < 0 {
			return fallen
		}
		path = newPath
	}
}

func main() {
	falling := load("example.txt")
	m := NewMemSpace(7, falling)
	m.Drop(12)
	length, path := m.getShortestPath()
	fmt.Println("Example 1:", length)
	coord := part2(m, path)
	fmt.Printf("Example 2: %v,%v\n", coord.X, coord.Y)

	falling = load("input.txt")
	m = NewMemSpace(71, falling)
	m.Drop(1024)
	length, _ = m.getShortestPath()
	fmt.Println("Part 1:", length)
	coord = part2(m, path)
	fmt.Printf("Part 2: %v,%v\n", coord.X, coord.Y)
}
