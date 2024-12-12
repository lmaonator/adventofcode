package main

import (
	"fmt"
	"os"
	"strings"
)

type Plot struct {
	name rune
	x, y int
}

type Field struct {
	plots [][]Plot
	h, w  int
}

type Region struct {
	value                    rune
	area, perimeter, corners int
}

func (r Region) Price() int {
	return r.area * r.perimeter
}

func (r Region) PriceBulkDiscount() int {
	return r.area * r.corners
}

func (r Region) String() string {
	return fmt.Sprint("{ Region ", string(r.value), " Area ", r.area, " Perimeter ", r.perimeter, " Corners ", r.corners, " }")
}

func load(f string) Field {
	d, _ := os.ReadFile(f)
	s := strings.Split(strings.TrimSpace(string(d)), "\n")
	h, w := len(s), len(s[0])
	plots := make([][]Plot, h)
	for y, l := range s {
		plots[y] = make([]Plot, w)
		for x, r := range l {
			p := Plot{r, x, y}
			plots[y][x] = p
		}
	}
	return Field{plots, h, w}
}

func main() {
	field := load("example.txt")
	regions := findRegions(field)
	fmt.Println("Example 1:", sumRegionPrice(regions))

	field = load("example2.txt")
	regions = findRegions(field)
	fmt.Println("Example 2:", sumRegionPriceBulkDiscount(regions))

	field = load("input.txt")
	regions = findRegions(field)
	fmt.Println("Part 1:", sumRegionPrice(regions))
	fmt.Println("Part 2:", sumRegionPriceBulkDiscount(regions))
}

func sumRegionPrice(regions []Region) int {
	price := 0
	for _, r := range regions {
		price += r.Price()
	}
	return price
}

func sumRegionPriceBulkDiscount(regions []Region) int {
	price := 0
	for _, r := range regions {
		price += r.PriceBulkDiscount()
	}
	return price
}

type point struct{ x, y int }

var dirs = [4]point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
var innerCorners = [4]point{{1, 1}, {-1, 1}, {-1, -1}, {1, -1}}

func findRegions(field Field) []Region {
	regions := []Region{}
	inRegion := map[Plot]struct{}{}
	for y := range field.plots {
		for x := range field.plots {
			origin := field.plots[y][x]
			if _, ok := inRegion[origin]; ok {
				continue
			}

			inRegion[origin] = struct{}{}
			region := Region{origin.name, 1, 0, 0}

			seen := map[point]struct{}{}
			seen[point{origin.x, origin.y}] = struct{}{}
			stack := []point{{x, y}}
			for len(stack) > 0 {
				curr := stack[0]
				stack = stack[1:]
				dirsOutsideRegion := [4]bool{}
				for i, dir := range dirs {
					currPoint := point{curr.x + dir.x, curr.y + dir.y}
					if _, ok := seen[currPoint]; ok {
						continue
					}
					if !inBounds(field.h, field.w, currPoint) {
						dirsOutsideRegion[i] = true
						region.perimeter++
						continue
					}
					currPlot := field.plots[currPoint.y][currPoint.x]
					if currPlot.name == origin.name {
						stack = append(stack, currPoint)
						region.area++
						inRegion[currPlot] = struct{}{}
						seen[currPoint] = struct{}{}
					} else {
						dirsOutsideRegion[i] = true
						region.perimeter++
					}
				}
				for i := range 4 {
					if dirsOutsideRegion[i] && dirsOutsideRegion[(i+1)%4] {
						region.corners++
					} else if !dirsOutsideRegion[i] && !dirsOutsideRegion[(i+1)%4] {
						currPoint := point{curr.x + innerCorners[i].x, curr.y + innerCorners[i].y}
						if !inBounds(field.h, field.w, currPoint) {
							region.corners++
						} else if field.plots[currPoint.y][currPoint.x].name != origin.name {
							region.corners++
						}
					}
				}

			}
			regions = append(regions, region)
		}
	}
	return regions
}

func inBounds(h, w int, p point) bool {
	return p.x >= 0 && p.y >= 0 && p.x < w && p.y < h
}
