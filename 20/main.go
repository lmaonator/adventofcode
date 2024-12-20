package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"slices"
	"strings"
)

type Track struct {
	grid        map[image.Point]int
	w, h        int
	start, end  image.Point
	regularTime int
}

func load(f string) Track {
	c, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(c)), "\n")
	t := Track{
		grid: map[image.Point]int{},
		h:    len(lines),
	}
	for y, line := range lines {
		t.w = len(line)
		for x, r := range line {
			p := image.Pt(x, y)
			switch r {
			case '.':
				t.grid[p] = math.MaxInt
			case 'S':
				t.grid[p] = 0
				t.start = p
			case 'E':
				t.grid[p] = math.MaxInt
				t.end = p
			}
		}
	}
	t.setTimes()
	return t
}

var dirs = []image.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func (t Track) nextPoint(p image.Point) image.Point {
	for _, dir := range dirs {
		n := p.Add(dir)
		if v, ok := t.grid[n]; ok && v > t.grid[p] {
			return n
		}
	}
	panic("already reached the end..")
}

func (t *Track) setTimes() {
	time := 0
	t.grid[t.start] = 0
	next := t.start
	for next != t.end {
		time++
		next = t.nextPoint(next)
		t.grid[next] = time
	}
	t.regularTime = time
}

func (t Track) validShortcuts(p image.Point) []image.Point {
	r := []image.Point{}
	for _, dir := range dirs {
		n := p.Add(dir.Mul(2))
		if v, ok := t.grid[n]; ok && v > t.grid[p] {
			r = append(r, n)
		}
	}
	return r
}

type Cheat [2]image.Point
type Cheats map[Cheat]int

func (cheats Cheats) Print() {
	x := map[int]int{}
	for _, v := range cheats {
		x[v]++
	}
	keys := make([]int, 0, len(cheats))
	for k := range x {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		if x[k] == 1 {
			fmt.Print("There is one cheat that saves ")
		} else {
			fmt.Print("There are ", x[k], " cheats that save ")
		}
		fmt.Println(k, "picoseconds.")
	}
}

func (t Track) findShortcuts(minSaved int) Cheats {
	r := Cheats{}
	curr := t.start
	for curr != t.end {
		scs := t.validShortcuts(curr)
		for _, sc := range scs {
			saved := t.grid[sc] - t.grid[curr] - 2
			if saved >= minSaved {
				r[Cheat{curr, sc}] = saved
			}
		}
		curr = t.nextPoint(curr)
	}
	return r
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func taxicab(a, b image.Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

var taxicabDirs []image.Point

func init() {
	c := image.Pt(0, 0)
	for y := -20; y < 21; y++ {
		for x := -20; x < 21; x++ {
			if x == 0 && y == 0 {
				continue
			}
			n := image.Pt(x, y)
			if taxicab(c, n) <= 20 {
				taxicabDirs = append(taxicabDirs, n)
			}
		}
	}
}

func (t Track) findLongShortcuts(minSaved int) Cheats {
	r := Cheats{}
	curr := t.start
	for curr != t.end {
		for _, dir := range taxicabDirs {
			e := curr.Add(dir)
			if v, ok := t.grid[e]; !ok || v <= t.grid[curr] {
				continue
			}
			saved := t.grid[e] - t.grid[curr] - taxicab(curr, e)
			if saved >= minSaved {
				r[Cheat{curr, e}] = saved
			}
		}
		curr = t.nextPoint(curr)
	}
	return r
}

func main() {
	t := load("example.txt")
	scs := t.findShortcuts(1)
	//scs.Print()
	fmt.Println("Example 1:", len(scs))
	scs = t.findLongShortcuts(50)
	//scs.Print()
	fmt.Println("Example 2:", len(scs))

	t = load("input.txt")
	fmt.Println("Part 1:", len(t.findShortcuts(100)))
	fmt.Println("Part 2:", len(t.findLongShortcuts(100)))
}
