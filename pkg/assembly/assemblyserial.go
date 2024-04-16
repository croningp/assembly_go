package assembly

import (
	"fmt"
	"math"
)

// this is a version of the assembly algorithm that runs without any concurrency
// These functions are mainly for testing, and are not implemented in the main functions of the program
// You can probably ignore them, unless you have some specific reason to be interested in them

func GraphAssemblySerial(g Graph) Pathway {


	var bestPathway Pathway
	//initPathway := Pathway{
	//	[]Graph{},
	//	g,
	//	[][]int{},
	//}
	initPathway := NewStartingPathway(g)

	//GraphAssemblySerialInner(&initPathway, &bestPathway, &g)
	GraphAssemblySerialInnerDG(&initPathway, &bestPathway, &g, 0)
	fmt.Println("Complete, assembly index: ", AssemblyIndex(&bestPathway, &g))

	return bestPathway
}

// GraphAssemblySerialInnerDG Attemps to implement Daniel's improvement into the algorithm
func GraphAssemblySerialInnerDG(currentPathway *Pathway, bestPathway *Pathway, originalGraph *Graph, level int) {

	if AssemblyIndex(bestPathway, originalGraph) < BestAssemblyIndex(originalGraph, currentPathway) {
		return
	}

	remnantEdges := len(currentPathway.remnant.Edges)
	sizesToCheck := int(math.Floor(float64(remnantEdges) / 2))
	// fmt.Println("RemnantEdges: ", remnantEdges)
	// fmt.Println("Sized to Check: ", sizesToCheck)
	BestPathwayUpdate(bestPathway, currentPathway)

	edgeAdjacencies := currentPathway.remnant.EdgeAdjacencies()
	forbidden := make(map[int]bool)              // map for whether a vertex is forbidden
	forbiddenSize := make(map[int]int)           // map of the size of the list a vertex was forbidden from
	var sub []int

	for i := 0; i < len(currentPathway.remnant.Edges); i++{
		sub = []int{i}  // subgraph starts with just the current edge
		for{

			neighbour, found := nonForbiddenNeighbour(sub, edgeAdjacencies, forbidden)

			if found && (len(sub) <= sizesToCheck){
				sub = append(sub, neighbour)
				// if level == 0 {fmt.Println("sub: ", sub)}
				subgraph, remnant := BreakGraphOnEdges(&currentPathway.remnant, sub)
				match := AllSubgraphsMatch(currentPathway, bestPathway, originalGraph, &subgraph, &remnant, level)
				if match{
					// if level == 0 {fmt.Println("MATCH")}
					continue

				}
				continue // cancel DG bit...
			}

			// backtrack
			thisForbidSize := len(sub)
			thisForbid := sub[len(sub)-1]                                      // item to forbid
			sub = sub[:len(sub)-1]                                             // pop the forbidden element out
			forbidUpdate(thisForbid, thisForbidSize, forbidden, forbiddenSize) // update forbidden lists


			if len(sub) == 0{
				break
			}
		}
	}

}

func AllSubgraphsMatch(currentPathway *Pathway, bestPathway *Pathway, originalGraph *Graph, subgraph *Graph, remnant *Graph, level int) bool {

	k := len(subgraph.Edges) // size of the subgraphs to search for

	//var edgeSubgraphs [][]int
	edgeAdjacencies := remnant.EdgeAdjacencies() // map of which edges are adjacent, maps edge index to slice of edge indices
	forbidden := make(map[int]bool)              // map for whether a vertex is forbidden
	forbiddenSize := make(map[int]int)           // map of the size of the list a vertex was forbidden from
	var sub []int                                // the current subgraph under construction
	match := false


	for i := 0; i < len(remnant.Edges); i++ {

		sub = []int{i}

		for {

			neighbour, found := nonForbiddenNeighbour(sub, edgeAdjacencies, forbidden)
			if found && (len(sub) <= k) {
				// if a neighbour is available add it to sub
				sub = append(sub, neighbour)


				if len(sub) == k {
					//edgeSubgraphs = helpers.CopyAppend(edgeSubgraphs, sub)
					possibleDuplicate, newRemnant := BreakGraphOnEdges(remnant, sub)
					if GraphsIsomorphic(subgraph, &possibleDuplicate) {
						match = true
						//fmt.Println("match: ", sub, match)
						newPathway := CopyPathway(currentPathway)
						newPathway.pathway = append(newPathway.pathway, CopyGraph(&possibleDuplicate))

						dupLeft := CopyEdgeList(subgraph.Edges)
						dupRight := CopyEdgeList(possibleDuplicate.Edges)
						newDuplicate := Duplicates{dupLeft, dupRight}
						newPathway.duplicates = append(newPathway.duplicates, newDuplicate)


						newPathway.duplicates = append(newPathway.duplicates, newDuplicate)
						newGraph, vertexMap := RecombineGraphs(&newRemnant, &possibleDuplicate)
						newPathway.remnant = CopyGraph(&newGraph)
						_ = vertexMap
						GraphAssemblySerialInnerDG(&newPathway, bestPathway, originalGraph, level + 1)

						//extend = true // continue building subgraph only if there is a match

					}
				}
				continue
			}

			// forbid the last item added, remove it, and update forbidden lists
			thisForbidSize := len(sub)
			thisForbid := sub[len(sub)-1]                                      // item to forbid
			sub = sub[:len(sub)-1]                                             // pop the forbidden element out
			forbidUpdate(thisForbid, thisForbidSize, forbidden, forbiddenSize) // update forbidden lists

			// break out of the loop if the subgraph is empty, and move on to the next starting edge
			if len(sub) == 0 {
				break
			}

		}

	}

	return match
}

