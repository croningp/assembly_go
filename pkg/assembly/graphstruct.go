package assembly

import (
	"GoAssembly/pkg/helpers"
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Code in this file relates mainly to the Graph struct

type Graph struct {
	Vertices      []int
	Edges         [][2]int
	// Adjacencies   map[int][]int
	VertexColours []string
	EdgeColours   []string
}

// NewColourGraph constructs a new Graph based on input vertices,edged, vertex and edge colours
func NewColourGraph(vertices []int, edges [][2]int, vColours []string, eColours []string) Graph {
	g := Graph{
		Vertices:      vertices,
		Edges:         edges,
		VertexColours: vColours,
		EdgeColours:   eColours,
	}

	// g.CalculateAdjacencies()
	return g

}

// NewGraph creates a new graph from vertices and edges. It just calls NewColourGraph, but with
// blank vertex and edge colours
func NewGraph(vertices []int, edges [][2]int) Graph {

	return NewColourGraph(vertices, edges, []string{}, []string{})
}

// NewGraphOnlyFromFile returns a graph from a graph file, but without a graph name or error
func NewGraphOnlyFromFile(filePath string) Graph {
	g, _, _ := NewGraphFromFile(filePath)
	return g
}


// NewGraphFromScanner returns a graph based from a scanner object, either from NewGraphFromFile or
// NewGraphFromString . The graph text should contain 5 lines:
// 1. A brief description of the graph; 2. Vertex indices separated by spaces, e.g. "0 1 2 3"
// 3. A list of edges as vertex indices separated with spaces, e.g. 0 1 1 2 2 3 (must have even length)
// 4. A list of vertex colour as strings (length = length of vertex list), or single "!" if not coloured
// 5. A list of edge colours as strings (length = length of edge list), or single "!" if not coloured
// If the file only has 3 lines, graph is assumed to have no vertex or edge colours
func NewGraphFromScanner(scanner *bufio.Scanner) (Graph, string, error){
	var graphName string
	var vertices []int
	var edges [][2]int
	vertexColours := make([]string, 0)
	edgeColours := make([]string, 0)
	i := 0

	for scanner.Scan() {

		if i == 0 {
			// the name of the graph is on line 0 - can be anything, useful for testing
			graphName = scanner.Text()
		}

		if i == 1 {
			// fill the vertices from line 1
			splitLine := strings.Fields(scanner.Text())
			for _, s := range splitLine {
				n, err := strconv.Atoi(s)
				check(err)
				vertices = append(vertices, n)
			}
		}
		if i == 2 {
			// fill the edges from line 2
			splitLine := strings.Fields(scanner.Text())
			if len(splitLine)%2 != 0 {
				return NewGraph([]int{}, [][2]int{}), "", errors.New("edges line in file must contain even number of digits")
			}
			var v1, v2 int
			for i, vertex := range splitLine {
				n, err := strconv.Atoi(vertex)
				check(err)
				if i%2 == 0 {
					v1 = n
				} else {
					v2 = n
					edges = append(edges, [2]int{v1, v2})
				}
			}

		}
		if i == 3 {
			lineText := scanner.Text()
			if lineText != "!" {
				vertexColours = strings.Fields(lineText)
				if len(vertexColours) != len(vertices) {
					return NewGraph([]int{}, [][2]int{}), "", errors.New("if vertex colours are specified, must be the same number as vertices")
				}
			}
		}
		if i == 4 {
			lineText := scanner.Text()
			if lineText != "!" {
				edgeColours = strings.Fields(lineText)
				if len(edgeColours) != len(edges) {
					return NewGraph([]int{}, [][2]int{}), "", errors.New("if edge colours are specified, must be the same number as edges")
				}
			}
			break // anything after the 5th line of the file should be ignored
		}

		i++
	}

	return NewColourGraph(vertices, edges, vertexColours, edgeColours), graphName, nil
}

// NewGraphFromFile returns a graph from text file input. See NewGraphFromScanner comments for required graph format
func NewGraphFromFile(filePath string) (Graph, string, error) {
	f, err := os.Open(filePath)
	check(err)
	scanner := bufio.NewScanner(f)

	graph, name, graphError := NewGraphFromScanner(scanner)
	closeErr := f.Close()
	check(closeErr)

	return graph, name, graphError
}

// NewGraphFromString returns a graph from text file input. See NewGraphFromScanner comments for required graph format
func NewGraphFromString(graphString string)(Graph, string, error){
	scanner := bufio.NewScanner(strings.NewReader(graphString))
	return NewGraphFromScanner(scanner)
}

// NewGraphOnlyFromString returns a graph from a string, but without a graph name or error
func NewGraphOnlyFromString(graphString string) Graph {
	g, _, _ := NewGraphFromString(graphString)
	return g
}

// InducedSubgraph returns a Graph object which is the subgraph of Graph g induced by vertices in subVertices
// Currently not used
func InducedSubgraph(g *Graph, subVertices []int) (Graph, error) {
	var outGraph Graph
	for _, v := range subVertices {
		if !helpers.Contains(g.Vertices, v) {
			fmt.Println(g.Vertices, v)
			return outGraph, errors.New("subvertices Contains vertices that are not in the Graph g")
		}
	}

	var newEdges [][2]int
	for _, e := range g.Edges {
		if helpers.Contains(subVertices, e[0]) && helpers.Contains(subVertices, e[1]) {
			newEdges = append(newEdges, e)
		}
	}

	outGraph = NewGraph(subVertices, newEdges)
	return outGraph, nil

}


// EdgeAdjacencies returns a map of indexed edge adjacencies, i.e. which edges are related to a given edge
// through sharing a vertex
// TODO: doesn't have a test
func (g *Graph) EdgeAdjacencies() map[int][]int {

	outMap := make(map[int][]int)

	for i, e1 := range g.Edges {
		for j, e2 := range g.Edges {
			if i != j {
				if e1[0] == e2[0] || e1[0] == e2[1] || e1[1] == e2[0] || e1[1] == e2[1] {
					helpers.MapUpdate(i, j, outMap)
				}
			}
		}
	}
	return outMap
}



// ListPairSort takes an [][2]int slice and sorts it, sorting the internal [2]int and then sorting the whole list
// by the first item in the pair
func ListPairSort(l [][2]int) [][2]int {
	var outputSlice [][2]int

	// order the pairs in the list
	for _, pair := range l {
		if pair[0] < pair[1] {
			outputSlice = append(outputSlice, [2]int{pair[0], pair[1]})
		} else {
			outputSlice = append(outputSlice, [2]int{pair[1], pair[0]})
		}
	}

	// order the list by the first of the pairs, or the 2nd if they are equal
	sort.Slice(outputSlice, func(i, j int) bool {
		if outputSlice[i][0] != outputSlice[j][0] {
			return outputSlice[i][0] < outputSlice[j][0]
		} else {

			return outputSlice[i][1] < outputSlice[j][1]
		}
	})

	return outputSlice
}

// GraphEquals checks for equality between two graphs, in that it contains the same list of vertices, and the
// same edges between them. It does not care about the order of vertices or edges, but it does not test for
// isomorphisms, only the same labels and edges between them.
func GraphEquals(g1 *Graph, g2 *Graph) bool {

	// copy vertices to new array and sort
	vert1 := make([]int, len(g1.Vertices))
	vert2 := make([]int, len(g2.Vertices))
	copy(vert1, g1.Vertices)
	copy(vert2, g2.Vertices)
	sort.Ints(vert1)
	sort.Ints(vert2)
	eqVert := reflect.DeepEqual(vert1, vert2)

	// check if vertex and edge colours are equal
	if !checkColoursEqual(g1, g2) {
		return false
	}

	// copy edges and sort
	edge1 := ListPairSort(g1.Edges)
	edge2 := ListPairSort(g2.Edges)
	eqEdge := reflect.DeepEqual(edge1, edge2)

	return eqVert && eqEdge
}

// checkColoursEqual checks if the vertex and edge colours of two graphs are the same
func checkColoursEqual(g1 *Graph, g2 *Graph) bool {
	g1VColoured := len(g1.Vertices) == len(g1.VertexColours)
	g1EColoured := len(g1.Edges) == len(g1.EdgeColours)
	g2VColoured := len(g2.Vertices) == len(g2.VertexColours)
	g2EColoured := len(g2.Edges) == len(g2.EdgeColours)

	// false if the are not both coloured or uncoloured
	if !(g1VColoured == g2VColoured && g1EColoured == g2EColoured) {
		return false
	}

	// Check if vertex colours equal
	if g1VColoured && !reflect.DeepEqual(VertexColourMap(g1), VertexColourMap(g2)) {
		return false
	}

	// check if edge colours equal
	if g1EColoured {
		e1ColMap := orderEdgeColourMap(EdgeColourMap(g1))
		e2ColMap := orderEdgeColourMap(EdgeColourMap(g2))
		if !reflect.DeepEqual(e1ColMap, e2ColMap) {
			return false
		}
	}

	return true

}

// VertexColourMap returns a map of vertices to vertex colours in a graph
func VertexColourMap(g *Graph) map[int]string {
	outMap := make(map[int]string)
	for i, v := range g.Vertices {
		outMap[v] = g.VertexColours[i]
	}
	return outMap
}

// EdgeColourMap returns a map of edges to edge colours in a graph
func EdgeColourMap(g *Graph) map[[2]int]string {
	outMap := make(map[[2]int]string)
	for i, e := range g.Edges {
		outMap[e] = g.EdgeColours[i]
	}
	return outMap
}

// orderEdgeColourMap takes a map of edge colours and orders the keys (edges) so that the edges vertices are in order
// e.g. keys {1, 2} and {2, 1}, which represent the same edge, both become {1, 2}
func orderEdgeColourMap(edgeMap map[[2]int]string) map[[2]int]string {
	outMap := make(map[[2]int]string)
	for k, v := range edgeMap {
		if k[0] < k[1] {
			outMap[k] = v
		} else {
			outMap[[2]int{k[1], k[0]}] = v
		}
	}
	return outMap
}

// GraphPrint returns a string containing some graph information
func GraphPrint(g *Graph) string {
	return fmt.Sprintf("Vertices %v\nEdges %v\nVertexColours %v\nEdgeColours %v ", g.Vertices, g.Edges, g.VertexColours, g.EdgeColours)
}

// GraphIsVertexColoured returns true if the graph is vertex coloured
func GraphIsVertexColoured(g *Graph) bool {
	if len(g.Vertices) == len(g.VertexColours){
		return true
	}
	return false
}

// GraphIsEdgeColoured returns true if the graph is edge coloured
func GraphIsEdgeColoured(g *Graph) bool {
	if len(g.Edges) == len(g.EdgeColours){
		return true
	}
	return false
}

// CopyGraph returns a copy of a graph
func CopyGraph(g *Graph) Graph {
	newVertices := make([]int, len(g.Vertices))
	newEdges := make([][2]int, len(g.Edges))
	newVertexColours := make([]string, len(g.VertexColours))
	newEdgeColours := make([]string, len(g.EdgeColours))

	copy(newVertices, g.Vertices)
	copy(newVertexColours, g.VertexColours)
	copy(newEdgeColours, g.EdgeColours)

	for i, e := range g.Edges {
		newEdges[i] = e
	}

	return NewColourGraph(newVertices, newEdges, newVertexColours, newEdgeColours)


}