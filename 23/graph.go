package main

type Node struct {
	name string
}

type Edge struct {
	Node   *Node
	Weight int
}

type Graph struct {
	Nodes []*Node
	Edges map[*Node][]*Edge
}

func (g *Graph) AddNode(n *Node) {
	g.Nodes = append(g.Nodes, n)
}

func (g *Graph) AddEdge(n1, n2 *Node, weight int) {
	if g.Edges == nil {
		g.Edges = map[*Node][]*Edge{}
	}
	e1 := Edge{n1, weight}
	e2 := Edge{n2, weight}
	g.Edges[n1] = append(g.Edges[n1], &e2)
	g.Edges[n2] = append(g.Edges[n2], &e1)
}
