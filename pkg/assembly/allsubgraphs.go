package assembly

import (
	"GoAssembly/pkg/helpers"
)

// NOTE: These are mot used in the main assembly algorithm. They are retained for use in future functionality.
// Path tracing algorithm implementation for all subgraphs.
// from Automatic Enumeration of All Connected Subgraphs, Rucker &  Rucker, 2000
// These functions are not present in the main assembly algorithm, as the algorithms are written into the
// architecture of the assembly functions. They can be used to check subgraphs though.

// AllSubgraphsToChan is the main subraph enumeration algorithm used, returning the subgraphs into a channel
// returns the count of the number of subgraphs, and an error
// countMode determines if the individual subgraphs are returned or not
func AllSubgraphsToChan(g *Graph, subgraphChan chan []int, subCountChan chan int, countMode bool) {
	var chans []chan [][]int
	var intChans []chan int
	defer close(subgraphChan)
	defer close(subCountChan)

	// initialise list of channels and run for all subgraphs containing the given edge, and those with higher indices
	for i := range g.Edges {
		chans = append(chans, make(chan [][]int, 1000))
		intChans = append(intChans, make(chan int, 1000))
		go AllSubsOnEdge(g, i, chans[i], intChans[i], countMode)
	}

	for i := 0; i < len(chans); i++ {

		// put all the subgraphs in the output channel if we are not in countMode
		edgeSubs := <-chans[i]
		if !countMode {
			for _, sub := range edgeSubs {
				subgraphChan <- sub
			}
		}

		// put the subgraph count for this edge in the subcount channel
		thisSubCount := <-intChans[i]
		subCountChan <- thisSubCount
	}

}

// SubgraphCount returns the number of subgraphs of a graph g as an integer
func SubgraphCount(g *Graph) int {
	_, subCount, _  := AllSubgraphs(g, true)
	return subCount
}

// MolSubgraphCount returns a count of all subgraphs of a molfile or mol block
func MolSubgraphCount(mol string, molBlock bool) int {
	var g Graph
	if molBlock{
		g = MolBlockColourGraph(mol)
	} else {
		g = MolColourGraph(mol)
	}
	return SubgraphCount(&g)
}

// AllSubgraphs returns all subgraphs as edge lists. If countMode is true, this will only return the count of the number of subgraphs,
// and will return nil for the list of subgraphs
func AllSubgraphs(g *Graph, countMode bool) ([][]int, int, error) {

	var edgeSubgraphs [][]int = nil
	subCount := 0
	var chans []chan [][]int
	var intChans []chan int

	for i := range g.Edges {
		chans = append(chans, make(chan [][]int))
		intChans = append(intChans, make(chan int))
		go AllSubsOnEdge(g, i, chans[i], intChans[i], countMode)
	}

	for i := 0; i < len(chans); i++ {
		eSub := <-chans[i]
		thisSubCount := <-intChans[i]
		edgeSubgraphs = append(edgeSubgraphs, eSub...)
		subCount += thisSubCount
	}

	return edgeSubgraphs, subCount, nil

}

// AllSubsOnEdge is called from AllSubgraphs to return all subgraphs that include a particular edge, and
// edges with higher indices
func AllSubsOnEdge(g *Graph, e int, c chan [][]int, cInt chan int, countMode bool) {

	edgeAdjacencies, sub, subCount, edgeSubgraphs, forbidden, forbiddenSize := InitialiseSubsOnEdge(g, e)

	for len(sub) != 0 {
		edgeSubgraphs, subCount, sub, forbidden, forbiddenSize = NextSubgraph(edgeSubgraphs, subCount, sub, forbidden, forbiddenSize, edgeAdjacencies, countMode)
	}

	if !countMode {
		c <- edgeSubgraphs
	} else {
		c <- nil
	}

	cInt <- subCount
}

// InitialiseSubsOnEdge sets up initial values for the AllSubsOnEdge function
func InitialiseSubsOnEdge(g *Graph, e int) (map[int][]int, []int, int, [][]int, map[int]bool, map[int]int) {
	edgeAdjacencies := g.EdgeAdjacencies() // map of which edges are adjacent, maps edge index to slice of edge indices
	forbidden := make(map[int]bool)        // map for whether a edge is forbidden
	forbiddenSize := make(map[int]int)     // map of the size of the list an edge was forbidden from
	var sub []int                          // the current subgraph under construction
	var edgeSubgraphs [][]int

	modifyEdgeAdjacencies(edgeAdjacencies, e) // remove edges with index <= e

	sub = []int{e}                                             // sub starts with just the current edge
	subCount := 1                                              // Initially just the one sub (current edge)
	edgeSubgraphs = helpers.CopyAppendSafe(edgeSubgraphs, sub) // current edge only is a valid subgraph, so add to list

	return edgeAdjacencies, sub, subCount, edgeSubgraphs, forbidden, forbiddenSize
}

// NextSubgraph takes the current subgraphs, subgraph count, forbidden list, forbidden size list,
// and returns the updated version for the next step in the path
func NextSubgraph(edgeSubgraphs [][]int, subCount int, sub []int, forbidden map[int]bool,
	forbiddenSize map[int]int, edgeAdjacencies map[int][]int, countMode bool) ([][]int, int, []int, map[int]bool, map[int]int) {
	// find a non-forbidden neighbour
	neighbour, found := nonForbiddenNeighbour(sub, edgeAdjacencies, forbidden)

	if found {

		// if a neighbour is available add it to sub
		subCount++
		sub = append(sub, neighbour)
		if !countMode {
			edgeSubgraphs = helpers.CopyAppendSafe(edgeSubgraphs, sub)
		}

	} else {
		// forbid the last item added, remove it, and update forbidden lists
		thisForbidSize := len(sub)
		thisForbid := sub[len(sub)-1]                                      // item to forbid
		sub = sub[:len(sub)-1]                                             // pop the forbidden element out
		forbidUpdate(thisForbid, thisForbidSize, forbidden, forbiddenSize) // update forbidden lists

	}

	return edgeSubgraphs, subCount, sub, forbidden, forbiddenSize
}

// strips out any edges in edgeAdj that are lower in index than v
func modifyEdgeAdjacencies(edgeAdj map[int][]int, e int) {
	for k := range edgeAdj {
		if k < e {
			delete(edgeAdj, k)
		} else {
			for i := 0; i < len(edgeAdj[k]); i++ {
				if edgeAdj[k][i] < e {
					edgeAdj[k] = append(edgeAdj[k][:i], edgeAdj[k][i+1:]...)
					i--
				}
			}
		}
	}
}
