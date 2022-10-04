package recipeanalysis

import (
	"gonum.org/v1/gonum/graph"
)

func getInboundVerticesForNode(g graph.Directed, n graph.Node) (inbound []graph.Node) {
	x := g.To(n.ID())

	for x.Next() {
		inbound = append(inbound, x.Node())
	}

	return inbound
}

func nodeHasMultipleInboundVertices(g graph.Directed, n graph.Node) bool {
	inboundVertices := getInboundVerticesForNode(g, n)

	inboundVerticesCount := len(inboundVertices)

	return inboundVerticesCount > 1
}
