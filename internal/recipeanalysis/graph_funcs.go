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

func getOutboundVerticesForNode(g graph.Directed, n graph.Node) (outbound []graph.Node) {
	x := g.From(n.ID())

	for x.Next() {
		outbound = append(outbound, x.Node())
	}

	return outbound
}

func getInboundAndOutboundVerticesForNode(g graph.Directed, n graph.Node) (inbound, outbound []graph.Node) {
	return getInboundVerticesForNode(g, n), getOutboundVerticesForNode(g, n)
}

func nodeHasOnlyOneInboundVertex(g graph.Directed, n graph.Node) bool {
	inboundVertices, outboundVertices := getInboundAndOutboundVerticesForNode(g, n)

	inboundVerticesCount := len(inboundVertices)
	outboundVerticesCount := len(outboundVertices)

	return inboundVerticesCount == 1 && outboundVerticesCount == 0
}

func nodeHasOnlyOneOutboundVertex(g graph.Directed, n graph.Node) bool {
	inboundVertices, outboundVertices := getInboundAndOutboundVerticesForNode(g, n)

	inboundVerticesCount := len(inboundVertices)
	outboundVerticesCount := len(outboundVertices)

	return inboundVerticesCount == 0 && outboundVerticesCount == 1
}

func nodeHasMultipleInboundVertices(g graph.Directed, n graph.Node) bool {
	inboundVertices := getInboundVerticesForNode(g, n)

	inboundVerticesCount := len(inboundVertices)

	return inboundVerticesCount > 1
}

func nodeHasMultipleOutboundVertices(g graph.Directed, n graph.Node) bool {
	outboundVertices := getOutboundVerticesForNode(g, n)

	outboundVerticesCount := len(outboundVertices)

	return outboundVerticesCount > 1
}
