package positionalgames

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

type colorEdges map[string][]string

func (g *Game) Dot() string {
	dotAst, _ := gographviz.ParseString("digraph G {}")
	graph := gographviz.NewGraph()
	_ = gographviz.Analyse(dotAst, graph)

	edges := make(colorEdges)

	colorsList := []string{"red", "blue", "green", "purple", "orange"}

	for i, payoff := range g.r.v {
		col := colorsList[i%len(colorsList)]
		highlightPath(g.r, payoff, col, edges)
	}

	addNodesAndEdges(graph, g.r, edges)
	return graph.String()
}

func highlightPath(n *node, payoff tuple, color string, edges colorEdges) {
	for _, child := range n.children {
		if containsPayoff(child.v, payoff) {
			key := fmt.Sprintf("%d->%d", n.id, child.id)
			edges[key] = append(edges[key], color)
			highlightPath(child, payoff, color, edges)
		}
	}
}

func addNodesAndEdges(gv *gographviz.Graph, n *node, edgeColors colorEdges) {
	label := fmt.Sprintf("Player %d\n%d", n.playerNum+1, n.id)
	if len(n.v) > 0 {
		label += "\\n("
		for i, payoff := range n.v {
			label += tupleToString(payoff)
			if i < len(n.v)-1 {
				label += "), ("
			}
		}
		label += ")"
	}
	quotedLabel := fmt.Sprintf("\"%s\"", label)
	nodeID := strconv.Itoa(n.id)

	gv.AddNode("G", nodeID, map[string]string{"label": quotedLabel})

	for _, child := range n.children {
		key := fmt.Sprintf("%d->%d", n.id, child.id)
		childID := strconv.Itoa(child.id)

		if colorList, ok := edgeColors[key]; !ok || len(colorList) == 0 {
			gv.AddEdge(nodeID, childID, true, nil)
		} else {
			for _, col := range colorList {
				attrs := map[string]string{
					"color":    col,
					"penwidth": "1",
				}

				gv.AddEdge(nodeID, childID, true, attrs)
			}
		}

		addNodesAndEdges(gv, child, edgeColors)
	}
}

func containsPayoff(v []tuple, payoff tuple) bool {
	for _, x := range v {
		if slices.Equal(x, payoff) {
			return true
		}
	}
	return false
}

func tupleToString(t tuple) string {
	s := ""
	for i, val := range t {
		s += strconv.Itoa(val)
		if i < len(t)-1 {
			s += ","
		}
	}
	return s
}
