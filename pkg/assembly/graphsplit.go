package assembly

import (
	"GoAssembly/pkg/helpers"
	"errors"
)

// Code related to splitting graphs up etc. This could probably be moved to some other file


// BreakGraphOnEdges returns two graph, one comprising the edges specified, and the other the remaining part
func BreakGraphOnEdges(g *Graph, edges []int) (Graph, Graph) {

	breakGraph := NewColourGraph(
		[]int{},
		[][2]int{},
		[]string{},
		[]string{},
	)

	remnantGraph := NewColourGraph(
		[]int{},
		[][2]int{},
		[]string{},
		[]string{},
	)

	// distribute the edges across the two graphs
	for i, _ := range g.Edges {
		if helpers.Contains(edges, i) {
			err := CopyGraphEdge(g, &breakGraph, i)
			check(err)
		} else {
			err := CopyGraphEdge(g, &remnantGraph, i)
			check(err)
		}
	}

	return breakGraph, remnantGraph

}

// CopyGraphEdge copies an edge from one graph to another
func CopyGraphEdge(oldGraph *Graph, newGraph *Graph, edgeIndex int) error {

	// copy the edge to the new graph
	newGraph.Edges = append(newGraph.Edges, oldGraph.Edges[edgeIndex])

	// copy the edge colours to the new graph, if there are any
	if len(oldGraph.EdgeColours) != 0 {
		if len(oldGraph.EdgeColours) != len(oldGraph.Edges) {
			return errors.New("graph has Edge Colours specified, but the number of edge colours does not equal number of edges")
		} else {
			newGraph.EdgeColours = append(newGraph.EdgeColours, oldGraph.EdgeColours[edgeIndex])
		}
	}

	// copy the associated vertices
	err := CopyGraphVerticesFromEdge(oldGraph, newGraph, edgeIndex)
	check(err)

	return nil
}

// CopyGraphVerticesFromEdge copies vertices associated with a particular edge from one graph to another
func CopyGraphVerticesFromEdge(oldGraph *Graph, newGraph *Graph, edgeIndex int) error {

	for _, v := range oldGraph.Edges[edgeIndex] {
		if !helpers.Contains(newGraph.Vertices, v) {

			// add vertex to new graph
			newGraph.Vertices = append(newGraph.Vertices, v)

			// get position of v in old graph
			vPosition := -1
			for i, posV := range oldGraph.Vertices {
				if posV == v {
					vPosition = i
				}
			}
			if vPosition == -1 {
				return errors.New("there is a vertex in the edge set of the input graph that does not appear in the vertex list of the input graph")
			}

			// add vertex colour to new graph if vertex colours were specified
			if len(oldGraph.VertexColours) != 0 {
				if len(oldGraph.VertexColours) != len(oldGraph.Vertices) {
					return errors.New("graph has Vertex Colours specified, but the number of vertex colours does not equal number of vertices")
				} else {
					newGraph.VertexColours = append(newGraph.VertexColours, oldGraph.VertexColours[vPosition])
				}
			}

		}

	}
	return nil
}

// RecombineGraphs takes a pair of input graphs and puts them into a single graph object, relabeling the vertices of graphRight
// No edges are added between the two graphs in the new object
func RecombineGraphs(graphLeft *Graph, graphRight *Graph) (Graph, map[int]int) {
	var outputEdges [][2]int
	var outputVertices []int
	var outputEdgeColours []string
	var outputVertexColours []string

	// copy left edges
	for i, edge := range graphLeft.Edges {
		outputEdges = append(outputEdges, edge)
		if GraphIsEdgeColoured(graphLeft) {
			outputEdgeColours = append(outputEdgeColours, graphLeft.EdgeColours[i])
		}
	}

	// copy left vertices
	for i, vertex := range graphLeft.Vertices {
		outputVertices = append(outputVertices, vertex)
		if GraphIsVertexColoured(graphLeft) {
			outputVertexColours = append(outputVertexColours, graphLeft.VertexColours[i])
		}
	}

	maxVertexLeft := helpers.MaxIntSlice(graphLeft.Vertices)
	maxVertexRight := helpers.MaxIntSlice(graphRight.Vertices)
	nextVertex := helpers.MaxIntSlice([]int{maxVertexLeft, maxVertexRight}) + 1

	// copy right vertices
	vertexMap := make(map[int]int)
	for i, vertex := range graphRight.Vertices {
		newVertex := vertex
		if helpers.Contains(outputVertices, vertex) {
			newVertex = nextVertex
			nextVertex++
		}
		outputVertices = append(outputVertices, newVertex)
		vertexMap[vertex] = newVertex
		if GraphIsVertexColoured(graphRight) {
			outputVertexColours = append(outputVertexColours, graphRight.VertexColours[i])
		}
	}

	// copy right edges
	for i, edge := range graphRight.Edges {
		outputEdges = append(outputEdges, [2]int{vertexMap[edge[0]], vertexMap[edge[1]]})
		if GraphIsEdgeColoured(graphRight) {
			outputEdgeColours = append(outputEdgeColours, graphRight.EdgeColours[i])
		}
	}

	return NewColourGraph(outputVertices, outputEdges, outputVertexColours, outputEdgeColours), vertexMap

}



// ConnectedComponentEdges finds sets of edges corresponding to all connected components in the graph
func ConnectedComponentEdges(g *Graph) [][]int {
	edgeAdj := g.EdgeAdjacencies()

	var edgesUsed []int
	var edgeSets [][]int
	for i := range g.Edges{
		if !helpers.Contains(edgesUsed, i){
			component := ConnectedComponent(i, edgeAdj)
			edgeSets = append(edgeSets, component)
			for _, j := range component{
				edgesUsed = append(edgesUsed, j)
			}
		}
	}
	return edgeSets
}

// ConnectedComponent returns the connected component of a graph that contains a given edge
func ConnectedComponent(edge int, edgeAdj map[int][]int) []int {
	component := []int{edge}
	for i := 0; i < len(component); i++{
		for _, j := range edgeAdj[component[i]]{
			if !helpers.Contains(component, j){
				component = append(component, j)
			}
		}
	}
	return component
}