package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func load(f string) *Graph {
	b, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	g := Graph{}
	nodes := map[string]*Node{}
	for _, line := range lines {
		s := strings.Split(line, "-")
		n1, ok := nodes[s[0]]
		if !ok {
			n1 = &Node{s[0]}
			nodes[s[0]] = n1
			g.AddNode(n1)
		}
		n2, ok := nodes[s[1]]
		if !ok {
			n2 = &Node{s[1]}
			nodes[s[1]] = n2
			g.AddNode(n2)
		}
		g.AddEdge(n1, n2, 1)
	}
	return &g
}

func getTriangles(g *Graph) [][]*Node {
	r := [][]*Node{}
	stack := [][]*Node{}
	for _, node := range g.Nodes {
		stack = append(stack, []*Node{node})
		for len(stack) > 0 {
			path := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			curr := path[len(path)-1]
			for _, edge := range g.Edges[curr] {
				if len(path) == 3 {
					if edge.Node == node {
						r = append(r, path)
					}
				} else {
					if strings.Compare(curr.name, edge.Node.name) < 0 {
						np := append(slices.Clone(path), edge.Node)
						stack = append(stack, np)
					}
				}
			}
		}
		clear(stack)
	}
	return r
}

func getLargestClique(g *Graph) []*Node {
	r := []*Node{}
	stack := [][]*Node{}
	for _, node := range g.Nodes {
		stack = append(stack, []*Node{node})
		for len(stack) > 0 {
			path := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(path) > len(r) {
				r = path
			}
			curr := path[len(path)-1]
			for _, edge := range g.Edges[curr] {
				if strings.Compare(curr.name, edge.Node.name) < 0 {
					if connectsToAllInPath(g, edge.Node, path) {
						np := append(slices.Clone(path), edge.Node)
						stack = append(stack, np)
					}
				}
			}
		}
		clear(stack)
	}
	return r
}

func connectsToAllInPath(g *Graph, node *Node, path []*Node) bool {
outer:
	for _, pathNode := range path {
		for _, newEdge := range g.Edges[node] {
			if pathNode == newEdge.Node {
				continue outer
			}
		}
		return false
	}
	return true
}

func part1(g *Graph, nameStart byte) int {
	r := 0
	for _, group := range getTriangles(g) {
		for _, n := range group {
			if n.name[0] == byte(nameStart) {
				r++
				break
			}
		}
	}
	return r
}

func part2(g *Graph) string {
	c := getLargestClique(g)
	return cliqueToStr(c)
}

func cliqueToStr(clique []*Node) string {
	str := ""
	for i, n := range clique {
		if i > 0 {
			str += ","
		}
		str += n.name
	}
	return str
}

func main() {
	g := load("example.txt")
	fmt.Println("Example 1:", part1(g, 't'))
	fmt.Println("Example 2:", part2(g))

	g = load("input.txt")
	fmt.Println("Part 1:", part1(g, 't'))
	fmt.Println("Part 2:", part2(g))
}
