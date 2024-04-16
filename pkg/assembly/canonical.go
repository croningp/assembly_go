package assembly

// This file contains code that implements a version of the nauty graph canonicalisation algorithm within Go (https://pallini.di.uniroma1.it/).
// References:
// Practical Graph Isomorhipsm II, McKay & Piperno
// McKay's Canonical Graph Labeling Algorithm, Hartke & Radcliffe
// TODO: get rid of some of the testing code
// I implemented the algorithm in Go as I think calling the nauty C code might interfere with my goroutines. I don't know this for sure though.
// Might switch to the nauty C code at some point, which is likely much faster.

import (
	"GoAssembly/pkg/helpers"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type AutomorphismData struct {
	automorphisms      []*Graph
	automorphismLevels [][]int
}

// Individualise individualises a vertex within a partition, placing it in its own part of the partition in front of the rest of its original partition
// For example, individualise 3 in [[1], [2, 3, 4], [5, 6]] would result in [[1], [3], [2, 4], [5, 6]]
func Individualise(partition [][]int, vertex int) [][]int {
	var individualised [][]int

	for _, part := range partition {
		if helpers.Contains(part, vertex) {

			// add individualised vertex in its own part
			individualised = append(individualised, []int{vertex})

			// add the rest of the vertices in order in a part
			var remainder []int
			for _, v := range part {
				if v != vertex {
					remainder = append(remainder, v)
				}
			}
			if len(remainder) != 0 {
				individualised = append(individualised, remainder)
			}

		} else {
			individualised = append(individualised, part)
		}
	}

	return individualised
}

// Degree returns the degree of a vertex within a graph. It simply calls DegreeInPart, but with the Part being all the vertices in the graph
func Degree(graph *Graph, vertex int) int {
	return DegreeInPart(graph, vertex, graph.Vertices)

}

// DegreeInPart returns the degree of a vertex with respect to a set of other vertices (a part of a partition)
// e.g. if vertex is 4 and part is [1, 2, 3] then this returns the total number of edges between 4 and [1, 2, 3].
// The part can contain the vertex in question.
func DegreeInPart(graph *Graph, vertex int, part []int) int {

	degreeInPart := 0

	for _, edge := range graph.Edges {
		// if one side of the edge is vertex, and the other is in part
		if (vertex == edge[0] && helpers.Contains(part, edge[1])) || (vertex == edge[1] && helpers.Contains(part, edge[0])) {
			degreeInPart++
		}
	}

	return degreeInPart

}

// Shatter returns the shattering of partLeft by partRight, splitting partLeft into parts ordered by degree into partRight
// for example, take the square graph with edges (1,2), (2,3), (3, 4), (4, 1) with partLeft being (1, 2, 3) and partRight being (4).
// then 2 has degree 0 with respect to (4) and both 1 and 3 have degree 1 with respect to 4, so the output is ((2), (1, 3))
func Shatter(graph *Graph, partLeft []int, partRight []int) [][]int {
	var shattering [][]int

	degreeMap := make(map[int][]int)

	for _, v := range partLeft {
		degreeInPart := DegreeInPart(graph, v, partRight)
		helpers.MapUpdate(degreeInPart, v, degreeMap)
	}

	var degrees []int
	for d := range degreeMap {
		degrees = append(degrees, d)
	}
	sort.Ints(degrees)

	for _, d := range degrees {
		shattering = append(shattering, degreeMap[d])
	}

	return shattering
}

// EquitableRefinement returns a more refined partition (where each part of the new partition is a subset of - including possibly being equal to - a part of the original partition)
// The new partition is equitable, meaning that all the vertices within any given part of the partition have the same degree into any other given part. E.g. for the partition
// [[1, 2, 3], [4, 5, 6], [7, 8 ,9]], vertices 1, 2 and 3 will all have the same degree into [1, 2, 3], and the same degree into [4, 5, 6] and [7, 8, 9]. E.g. 1, 2, 3 all have
// degree d1 into [4, 5, 6] and all have degree d2 into [7, 8, 9] but d1 and d2 do not have to be equal.
// If the input partition is already equitable, it is returned unchanged
func EquitableRefinement(graph *Graph, partition [][]int) [][]int {

	refinedPartition := helpers.CopySliceOfSlices(partition)

	for {

		done := true

		for i := 0; i < len(refinedPartition); i++ {
			breakOut := false
			for j := 0; j < len(refinedPartition); j++ {
				shattering := Shatter(graph, refinedPartition[i], refinedPartition[j])

				if len(shattering) > 1 {
					newPartition := helpers.CopySliceOfSlices(refinedPartition[:i])
					newPartition = append(newPartition, shattering...)
					if len(refinedPartition) > i {
						newPartition = append(newPartition, refinedPartition[i+1:]...)
					}
					refinedPartition = helpers.CopySliceOfSlices(newPartition)
					done = false
					breakOut = true
					break
				}

			}
			if breakOut {
				break
			}
		}

		if done {
			break
		}
	}

	return refinedPartition
}

// IsEquitable returns true if a partition is equitable, i.e. every vertex in any given part has the same degree into each other given part
func IsEquitable(graph *Graph, partition [][]int) bool {
	for _, partLeft := range partition {
		for _, partRight := range partition {
			shattering := Shatter(graph, partLeft, partRight)
			if len(shattering) > 1 {
				return false
			}
		}

	}
	return true

}

// CoarsestEquitableColourings returns the set of coarsest equitable colourings of a given partition. If the partition is not equitable, it refines the
// partition using EquitableRefinement. If it is equitable, it chooses the first non-trivial part of the partition (length > 1) and returns the list of equitable
// refinements after individualising each vertex in that part, as well as the individualised vertices
// Colouring and partition mean the same thing, essentially. I should probably make the names consistent at some point.
// TODO: tests
func CoarsestEquitableColourings(graph *Graph, partition [][]int) ([][][]int, []int) {
	var outputColourings [][][]int
	var individualisedVertices []int

	if !IsEquitable(graph, partition) {
		refinement := EquitableRefinement(graph, partition)
		return [][][]int{refinement}, []int{}
	} else {
		for _, part := range partition {
			if len(part) > 1 {
				for _, v := range part {
					individualised := Individualise(partition, v)
					outputColourings = append(outputColourings, EquitableRefinement(graph, individualised))
					individualisedVertices = append(individualisedVertices, v)
				}
				return outputColourings, individualisedVertices
			}
		}
	}

	return nil, nil
}

// SearchTree is the main canonical graph algorithm, which returns a canonical version of the input graph
// An initial colouring can be specified (a partition) if the input graph is already coloured
// TODO: tests
func SearchTree(graph *Graph, initialColouring [][]int, prune bool) Graph {

	canonicalGraphContainer := make([]Graph, 1)
	var bestInvariant [][]int
	var v []int
	treeLevel := []int{0}
	backtrack := []bool{false}
	backtrackLevel := []int{-1}
	automorphisms := AutomorphismData{
		[]*Graph{},
		[][]int{},
	}

	SearchTreeInner(graph, initialColouring, v, canonicalGraphContainer, bestInvariant, &automorphisms, treeLevel, prune, backtrack, backtrackLevel)

	return canonicalGraphContainer[0]
}

// SearchTreeInner is the inner recursive part of the SearchTree algorithm
func SearchTreeInner(graph *Graph, colouring [][]int, v []int, canonicalGraphContainer []Graph, bestInvariant [][]int, automorphisms *AutomorphismData, treeLevel []int, prune bool, backtrack []bool, backtrackLevel []int) {

	if prune && backtrack[0] {
		if len(treeLevel) == backtrackLevel[0] {
			backtrack[0] = false
			backtrackLevel[0] = -1 // prob not necessary
		} else {
			return
		}
	}

	if IsDiscrete(colouring) {

		canonicalCandidate := PermuteGraph(graph, DiscreteColouringToIntSlice(colouring))

		autoFound, level := AutomorphismCheck(&canonicalCandidate, automorphisms)
		_, _ = autoFound, level

		if autoFound {
			// start backtracking to
			// fmt.Println("automorphism found: ", automorphisms.automorphismLevels[level] ," this level ", treeLevel)
			backtrack[0] = true
			backtrackLevel[0] = MaxEqualLevel(automorphisms.automorphismLevels[level], treeLevel)
			return
		} else {
			automorphisms.automorphisms = append(automorphisms.automorphisms, &canonicalCandidate)
			automorphisms.automorphismLevels = append(automorphisms.automorphismLevels, treeLevel)
		}

		if len(canonicalGraphContainer[0].Vertices) == 0 || GraphGreaterThan(&canonicalCandidate, &canonicalGraphContainer[0]) {
			canonicalGraphContainer[0] = canonicalCandidate
		}

	} else {
		equitableColourings, vertices := CoarsestEquitableColourings(graph, colouring)

		for i, newColouring := range equitableColourings {
			newV := make([]int, len(v))
			copy(newV, v)

			if len(vertices) != 0 {
				newV = append(newV, vertices[i])
			}

			newTreeLevel := append(treeLevel, i)

			SearchTreeInner(graph, newColouring, newV, canonicalGraphContainer, bestInvariant, automorphisms, newTreeLevel, prune, backtrack, backtrackLevel)
		}

	}

}

// AutomorphismCheck checks the canonical graph candidate against a list of graphs for automosphisms, and returns whether the graph is found, as well as the
// position in the list
func AutomorphismCheck(canonicalCandidate *Graph, automorphismData *AutomorphismData) (bool, int) {
	for i, g := range automorphismData.automorphisms {
		if GraphEquals(canonicalCandidate, g) {
			return true, i
		}
	}
	return false, -1
}

// IsDiscrete returns true if a colouring is discrete, i.e. every part of the colouring has only 1 vertex
// e.g. [[1],[2],[3],[4],[5]] is discrete, [[1,2],[3],[4],[5]] is not
func IsDiscrete(colouring [][]int) bool {
	for _, part := range colouring {
		if len(part) > 1 {
			return false
		}

	}
	return true
}

// DiscreteColouringToIntSlice takes a discrete colouring, which is a slice of slices each having a single member,
// e.g. {{1}, {2}, {3}, {4}} and flattens it to an int slice, e.g. {1, 2, 3, 4}
func DiscreteColouringToIntSlice(colouring [][]int) []int {
	if !IsDiscrete(colouring) {
		check(errors.New("DiscreteColouringToIntSlice error - colouring is not discrete"))
	}

	var intSlice []int
	for _, part := range colouring {
		if len(part) == 1 {
			intSlice = append(intSlice, part[0])
		} else {
			log.Fatal(errors.New("DiscreteColouringToIntSlice error - colouring is not discrete"))
		}
	}
	return intSlice

}

// PermuteGraph relabels graph nodes based on input permutation. The permutation [4, 1, 3, 2, 5]  on [1, 2, 3, 4, 5]
// should move 1 to where 4 was, 2 to where 1 was etc. The permuted graph is returned.
func PermuteGraph(graph *Graph, permutation []int) Graph {


	if len(permutation) != len(graph.Vertices) {
		log.Fatal(errors.New("PermuteGraph error - permutation must have the same length as the number of vertices in the input graph"))
	}

	// handle coloured graphs - edge colours should be the same
	vertexColoured := len(graph.VertexColours) == len(graph.Vertices)
	edgeColoured := len(graph.EdgeColours) == len(graph.Edges)
	var newVertexColours []string
	var newEdgeColours []string
	var vertexColourMap map[int]string

	if vertexColoured {
		newVertexColours = make([]string, len(graph.Vertices))
		vertexColourMap = VertexColourMap(graph)
	} else {
		newVertexColours = make([]string, 0)
	}


	sortedVertices := make([]int, len(graph.Vertices))
	for i, v := range graph.Vertices {
		sortedVertices[i] = v
	}
	sort.Ints(sortedVertices)

	permutationMap := make(map[int]int)
	for i, vertex := range sortedVertices {
		permutationMap[permutation[i]] = vertex
	}


	newVertices := make([]int, len(graph.Vertices))
	for i, _ := range newVertices {
		newVertices[i] = permutationMap[sortedVertices[i]]
		if vertexColoured {
			newVertexColours[i] = vertexColourMap[sortedVertices[i]]
		}
	}

	var newEdges [][2]int
	for i, e := range graph.Edges {
		newEdges = append(newEdges, [2]int{permutationMap[e[0]], permutationMap[e[1]]})
		if edgeColoured {
			newEdgeColours = append(newEdgeColours, graph.EdgeColours[i])
		}

	}


	return NewColourGraph(newVertices, newEdges, newVertexColours, newEdgeColours)
}

// PermuteColouring returns a permuted colouring (partition). This was just used for testing, not part of the canonical algorithm
func PermuteColouring(graph *Graph, colouring [][]int, permutation []int) [][]int {
	var newColouring [][]int

	sortedVertices := make([]int, len(graph.Vertices))
	for i, v := range graph.Vertices {
		sortedVertices[i] = v
	}
	sort.Ints(sortedVertices)

	permutationMap := make(map[int]int)
	for i, vertex := range sortedVertices {
		permutationMap[permutation[i]] = vertex
	}

	for _, part := range colouring {
		var newPart []int
		for _, v := range part {
			newPart = append(newPart, permutationMap[v])
		}
		newColouring = append(newColouring, newPart)
	}

	return newColouring
}

// FlattenEdgeList converts a list of edges, e.g. {{1, 2}, {3, 4}} to a flattened slice of ints e.g. {1, 2, 3, 4}
// This flat list is used to order the graphs
func FlattenEdgeList(edgeList [][2]int) []int {
	var flatList []int

	for _, e := range edgeList {
		flatList = append(flatList, e[0])
		flatList = append(flatList, e[1])
	}

	return flatList
}

// SliceGreaterThan returns true if the left slice is greater than the right in lexographic order, i.e. the first element that differs between the slices is greater
// or the left slice is longer than the right, where the right is a prefix of the left
func SliceGreaterThan(intSliceLeft []int, intSliceRight []int) bool {

	lenLeft := len(intSliceLeft)
	lenRight := len(intSliceRight)

	i := 0
	for {

		if intSliceLeft[i] != intSliceRight[i] {
			return intSliceLeft[i] > intSliceRight[i]
		}

		// reached end of left fragment with every element equal so far, left cannot be greater
		if i == lenLeft-1 {
			return false
		}

		// reached end of right fragment but not left with every element equal so far, left  is greater
		if i == lenRight-1 {
			return true
		}

		i++
	}

}

// SliceEqual returns true if two slices of ints are equal at each position
func SliceEqual(intSliceLeft []int, intSliceRight []int) bool {

	if len(intSliceLeft) != len(intSliceRight) {
		return false
	}

	for i, _ := range intSliceLeft {
		if intSliceLeft[i] != intSliceRight[i] {
			return false
		}
	}

	return true

}

// GraphGreaterThan returns true if graphLeft is greater than graphRight, as determined by the lexographic ordering of the sorted edge list
func GraphGreaterThan(graphLeft *Graph, graphRight *Graph) bool {
	edgeListLeft := ListPairSort(graphLeft.Edges)
	flatEdgeListLeft := FlattenEdgeList(edgeListLeft)

	edgeListRight := ListPairSort(graphRight.Edges)
	flatEdgeListRight := FlattenEdgeList(edgeListRight)

	return SliceGreaterThan(flatEdgeListLeft, flatEdgeListRight)
}

// GraphColourPartition creates an initial partition for a graph based on colours. If there are no graph colours the initial partition just contains one
// part with all of the vertices. If there are colours, then the initial partition is based on those, in order of colour.
func GraphColourPartition(graph *Graph) [][]int {

	var outputPartition [][]int
	if len(graph.VertexColours) == 0 {
		fullPart := make([]int, len(graph.Vertices))
		copy(fullPart, graph.Vertices)
		outputPartition = append(outputPartition, fullPart)
		return outputPartition
	}

	// make list of all the graph colours
	var colours []string
	for _, c := range graph.VertexColours {
		if !helpers.ContainsStr(colours, c) {
			colours = append(colours, c)
		}
	}
	sort.Strings(colours)  // colours must be in order

	// generate partition based on ordered colours
	outputPartition = make([][]int, len(colours))
	for i, c := range colours {
		for j, v := range graph.Vertices {
			if graph.VertexColours[j] == c {
				outputPartition[i] = append(outputPartition[i], v)
			}
		}
	}

	return outputPartition

}

// MaxEqualLevel returns the final position in a pair of lists where the lists are equal. -1 if not equal at any point.
func MaxEqualLevel(sliceLeft []int, sliceRight []int) int {

	levelMatch := -1

	for i := 0; i < len(sliceLeft); i++ {
		if sliceLeft[i] == sliceRight[i] {
			levelMatch = i
		} else {
			break
		}
	}

	return levelMatch
}

// RandomPermutation randomlu permutes an input slice. For testing.
func RandomPermutation(inputList []int) []int {
	permutedList := make([]int, len(inputList))
	copy(permutedList, inputList)
	rand.Shuffle(len(permutedList), func(i, j int) { permutedList[i], permutedList[j] = permutedList[j], permutedList[i] })
	return permutedList
}

// RandomPermutationList returns a number of random permutations of an input slice. For testing
func RandomPermutationList(inputList []int, numberOfPermutations int) [][]int {

	rand.Seed(time.Now().UnixNano())
	var outputPermutations [][]int
	for i := 0; i < numberOfPermutations; i++ {
		outputPermutations = append(outputPermutations, RandomPermutation(inputList))

	}

	return outputPermutations
}

// CanonicalGraphTest is for some additional canonical testing
func CanonicalGraphTest(graph *Graph, inputColouring [][]int, numberOfPermutations int) bool {
	pass := true

	permutations := RandomPermutationList(graph.Vertices, numberOfPermutations)

	var colouring [][]int
	if len(inputColouring) == 0 {
		colouring = make([][]int, 1)
		vertexList := make([]int, len(graph.Vertices))
		copy(vertexList, graph.Vertices)
		colouring[0] = vertexList
	} else {
		colouring = inputColouring
	}

	var permutedGraphs []Graph
	var permutedColourings [][][]int
	for _, p := range permutations {
		permutedGraphs = append(permutedGraphs, PermuteGraph(graph, p))
		permutedColourings = append(permutedColourings, PermuteColouring(graph, colouring, p))
	}

	var canonicals []Graph
	for i, g := range permutedGraphs {
		canonicalGraph := SearchTree(&g, permutedColourings[i], true)
		canonicals = append(canonicals, canonicalGraph)
	}

	for _, c := range canonicals {
		if !GraphEquals(&c, &canonicals[0]) {
			fmt.Println("not equal", c.Edges, canonicals[0].Edges)
			pass = false
		}
	}


	return pass
}

// RandomGraph builds a random graph for testing by stringing together vertices and then adding some extra edges
// Colours will be assigned to vertices randomly
func RandomGraph(numVertices int, numExtraEdges int, colours []string) Graph {

	var vertices []int
	var edges [][2]int
	var vertexColours []string
	_ = vertexColours

	for i := 1; i <= numVertices; i++ {
		vertices = append(vertices, i)
	}

	vertices = RandomPermutation(vertices)

	// add edges to make a skeleton
	includedVertices := []int{vertices[0]}
	for i := 1; i < len(vertices); i++ {
		vertexToJoin := includedVertices[rand.Intn(len(includedVertices))]
		edges = append(edges, [2]int{vertices[i], vertexToJoin})
		includedVertices = append(includedVertices, vertices[i])

	}


	for i := 0; i < numExtraEdges; i++ {
		randomVertices := RandomPermutation(vertices)
		edges = append(edges, [2]int{randomVertices[0], randomVertices[1]})
	}

	if len(colours) != 0 {
		for i := 0; i < len(vertices); i++ {
			vertexColours = append(vertexColours, colours[rand.Intn(len(colours))])
		}
	}

	return NewColourGraph(vertices, edges, vertexColours, []string{})
}

// RandomGraphCanonicalTest tests canonicalisation across a specified number of random graphs
func RandomGraphCanonicalTest(numGraphs int, numPermutations int, numVerticesRange [2]int, numExtraEdgesRange [2]int, vertexColours []string) bool {

	time.Sleep(1 * time.Second) // in case multiple calls are so fast that the seed is the same
	rand.Seed(time.Now().UnixNano())

	pass := true

	numVertices := rand.Intn(numVerticesRange[1]-numVerticesRange[0]) + numVerticesRange[0]
	numExtraEdges := rand.Intn(numExtraEdgesRange[1]-numExtraEdgesRange[0]) + numExtraEdgesRange[0]

	for i := 0; i < numGraphs; i++ {
		testGraph := RandomGraph(numVertices, numExtraEdges, vertexColours)
		colouring := GraphColourPartition(&testGraph)
		result := CanonicalGraphTest(&testGraph, colouring, numPermutations)
		// fmt.Println(result)
		pass = pass && result
	}

	return pass

}



func GraphsIsomorphic(graphLeft *Graph, graphRight *Graph) bool {


	if !CanonicalInitialCheck(graphLeft, graphRight){

		return false
	}

	var checkGraphLeft Graph
	var checkGraphRight Graph

	// vertex index mapping
	if reflect.DeepEqual(graphLeft.Vertices, graphRight.Vertices){
		checkGraphLeft = *graphLeft
		checkGraphRight = *graphRight
	} else {
		// relabel graphRight to have the same vertex labels as graphLeft
		checkGraphLeft = *graphLeft
		checkGraphRight = GraphVertexRelabel(graphRight, graphLeft.Vertices)
	}


	if len(graphLeft.Edges) == len(graphLeft.EdgeColours){
		checkGraphLeft = EdgeColourConversion(&checkGraphLeft)
		checkGraphRight = EdgeColourConversion(&checkGraphRight)
	}

	canonicalLeft := SearchTree(&checkGraphLeft, GraphColourPartition(&checkGraphLeft), true)
	canonicalRight := SearchTree(&checkGraphRight, GraphColourPartition(&checkGraphRight), true)

	return GraphEquals(&canonicalLeft, &canonicalRight)
}

// GraphVertexRelabel returns a relabeled version of the input graph with the vertices and edges relabeled according to the given labeling
// this is used to ensure the same set of vertex labels is used for the canonicalisation check
func GraphVertexRelabel(graph *Graph, labeling []int) Graph {

	if len(graph.Vertices) != len(labeling){
		log.Fatal("GraphVertexRelabel error - number of vertices does not equal size of labeling")
	}

	// map of the current vertex labels to the new ones
	labelMap := make(map[int]int)
	for i, v := range graph.Vertices{
		labelMap[v] = labeling[i]
	}

	// use label map to generate new edges with the appropriate new vertex labels
	newEdges := make ([][2]int, len(graph.Edges))
	for i, e := range graph.Edges{
		newEdges[i] = [2]int{labelMap[e[0]],labelMap[e[1]]}
	}

	// vertex and edge colours are unchanged
	newVertexColours := make([]string, len(graph.VertexColours))
	newEdgeColours := make([]string, len(graph.EdgeColours))
	copy(newVertexColours, graph.VertexColours)
	copy(newEdgeColours, graph.EdgeColours)

	return NewColourGraph(labeling, newEdges,newVertexColours, newEdgeColours)

}

func CanonicalInitialCheck(graphLeft *Graph, graphRight *Graph) bool{

	// Some basic canonicalisation checks to ensure that the partition of colours is the same
	possibleMatch := len(graphLeft.Vertices) == len(graphRight.Vertices)
	possibleMatch = possibleMatch &&  (len(graphLeft.Edges) == len(graphRight.Edges))
	possibleMatch = possibleMatch &&  (len(graphLeft.VertexColours) == len(graphRight.VertexColours))
	possibleMatch = possibleMatch &&  (len(graphLeft.EdgeColours) == len(graphRight.EdgeColours))

	if !possibleMatch{
		return false
	}

	leftEdgeColours := make([]string, len(graphLeft.EdgeColours))
	copy(leftEdgeColours, graphLeft.EdgeColours)
	sort.Strings(leftEdgeColours)
	rightEdgeColours := make([]string, len(graphRight.EdgeColours))
	copy(rightEdgeColours, graphRight.EdgeColours)
	sort.Strings(rightEdgeColours)
	if !reflect.DeepEqual(leftEdgeColours, rightEdgeColours){
		return false
	}

	leftVertexColours := make([]string, len(graphLeft.VertexColours))
	copy(leftVertexColours, graphLeft.VertexColours)
	sort.Strings(leftVertexColours)
	rightVertexColours := make([]string, len(graphRight.VertexColours))
	copy(rightVertexColours, graphRight.VertexColours)
	sort.Strings(rightVertexColours)
	if !reflect.DeepEqual(leftVertexColours, rightVertexColours){
		return false
	}

	// TODO: test, and add check for map of vertex colours to edge colours

	return true

}

// EdgeColourConversion takes a graph with edge colours and converts it to a layered graph with a different layer for each edge colour. E.g. if there are 3 edge colours
// 1, 2, 3 and vertex colours a, b, c, then each vertex of colour a is split into a linear graph a1-a2-a3, with edges of type 1 going between vertices of type 1, and so on.
// There are more space efficient ways of doing this, e.g. type 1 is layer 1, type 2 is layer 2, type 3 is both layers 1 and 2. Not implemented for now though.
func EdgeColourConversion(graph *Graph) Graph {

	// don't need to do anything if no distinct edge colours
	if len(graph.EdgeColours) < 2 {
		return CopyGraph(graph)
	}

	outGraph := NewColourGraph([]int{}, [][2]int{}, []string{}, []string{})

	vertexMap := make(map[int][]int)
	nextVertex := helpers.MaxIntSlice(graph.Vertices) + 1 // the next vertex index to use

	var edgeColours []string
	for _, c := range graph.EdgeColours {
		if !helpers.ContainsStr(edgeColours, c) {
			edgeColours = append(edgeColours, c)
		}
	}
	sort.Strings(edgeColours)

	numLayers := len(edgeColours)

	// used to map the edge colours to the correct layers
	edgeColourMap := make(map[string]int)
	for i, c := range edgeColours {
		edgeColourMap[c] = i
	}

	// add new vertices to the graph and link the vertices in each layer
	for k, v := range graph.Vertices {
		vertexMap[v] = make([]int, numLayers)
		for i := 0; i < numLayers; i++ {
			if i == 0 {
				outGraph.Vertices = append(outGraph.Vertices, v) // layer 1 is the initial vertices
				vertexMap[v][i] = v
			} else {
				outGraph.Vertices = append(outGraph.Vertices, nextVertex)
				vertexMap[v][i] = nextVertex
				outGraph.Edges = append(outGraph.Edges, [2]int{vertexMap[v][i-1], vertexMap[v][i]}) // link the layers with edges
				nextVertex++
			}
			outGraph.VertexColours = append(outGraph.VertexColours, graph.VertexColours[k]+strconv.Itoa(i)) // the vertices on each layer have a different colour
		}
	}

	// add edges from original graph into the correct layer
	for i, e := range graph.Edges {
		edgeLayerIndex := edgeColourMap[graph.EdgeColours[i]]
		vertLeft := vertexMap[e[0]][edgeLayerIndex]
		vertRight := vertexMap[e[1]][edgeLayerIndex]
		outGraph.Edges = append(outGraph.Edges, [2]int{vertLeft, vertRight})
	}

	return outGraph

}

// EdgeColourRandomGraph returns a random edge coloured graph, for testing
func EdgeColourRandomGraph(numVertices int, numExtraEdges int, vertexColours []string, edgeColours []string) Graph {

	randomGraph := RandomGraph(numVertices, numExtraEdges, vertexColours)

	for i:=0;i<len(randomGraph.Edges);i++ {
		randomGraph.EdgeColours = append(randomGraph.EdgeColours, edgeColours[rand.Intn(len(edgeColours))])
	}

	return randomGraph
}

// EdgeColourRandomCanonicalTest contains some code to test Edge Coloured graph canonicalisation
func EdgeColourRandomCanonicalTest(numPermutations int, numGraphs int, numVerticesRange [2]int, numExtraEdgesRange [2]int, vertexColours []string, edgeColours []string)bool{


	time.Sleep(1 * time.Second) // in case multiple calls are so fast that the seed is the same
	rand.Seed(time.Now().UnixNano())

	pass := true

	numVertices := rand.Intn(numVerticesRange[1]-numVerticesRange[0]) + numVerticesRange[0]
	numExtraEdges := rand.Intn(numExtraEdgesRange[1]-numExtraEdgesRange[0]) + numExtraEdgesRange[0]


	for i := 0; i < numGraphs; i++ {
		canonicals := make([]Graph, 0)
		// canonicals = []Graph{}
		testGraph := EdgeColourRandomGraph(numVertices, numExtraEdges, vertexColours, edgeColours)
		permutations := RandomPermutationList(testGraph.Vertices, numPermutations)

		for _, p := range permutations {
			permutedGraph := PermuteGraph(&testGraph, p)
			edgeGraph := EdgeColourConversion(&permutedGraph)
			edgeGraphColouring := GraphColourPartition(&edgeGraph)


			thisCanonical := SearchTree(&edgeGraph, edgeGraphColouring, true)
			canonicals = append(canonicals, thisCanonical)
		}

		passed := 0
		failed := 0
		for _, c := range canonicals {
			if !GraphEquals(&c, &canonicals[0]) {
				// fmt.Println("not equal", c.Edges, canonicals[0].Edges)
				pass = false
				failed += 1
			} else {
				// fmt.Println("EQUALS")
				passed += 1
			}
		}
		// fmt.Println("Passed: ", passed, ", Failed: ", failed)
	}



	return pass

}