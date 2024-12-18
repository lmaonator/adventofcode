package main

type Edge struct {
	Node   *Item
	Weight int
}

type Graph struct {
	Nodes []*Item
	Edges map[*Item][]*Edge
}

func (g *Graph) AddNode(n *Item) {
	g.Nodes = append(g.Nodes, n)
}

func (g *Graph) AddEdge(n1, n2 *Item, weight int) {
	if g.Edges == nil {
		g.Edges = map[*Item][]*Edge{}
	}
	e1 := Edge{n1, weight}
	e2 := Edge{n2, weight}
	g.Edges[n1] = append(g.Edges[n1], &e2)
	g.Edges[n2] = append(g.Edges[n2], &e1)
}
