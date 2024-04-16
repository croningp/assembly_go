package assembly

// this file contains functions common to both the serial and parallel implementations of the assembly algorithms

import (
	"GoAssembly/pkg/helpers"
	"fmt"
	"log"
	"math"
	"reflect"
	"sync"
)

// Pathway contains all information representing an assembly pathway. Pathway.pathway is a list of graphs that represent duplicated structures
// within the original graph object, and Pathway.duplicates is the related edge indices of those duplicated graphs with respect to the original graph (TODO: check this)
// Pathway.remnant is the remaining structure once duplicates have been removed / separated out. Pathway does not contain the original graph, which is also
// required for meaningful calculations, such as the assembly index.
type Pathway struct {
	pathway    []Graph
	remnant    Graph
	duplicates []Duplicates
	atomEquivalents [][]int
}

// NewStartingPathway returns a pathway with only the remnant and no duplicates etc
func NewStartingPathway(graph Graph) Pathway{
	var pathway []Graph
	var duplicates []Duplicates
	var atomEquivalents [][]int
	return NewPathway(pathway, graph, duplicates, atomEquivalents)

}

func NewPathway(pathway []Graph, remnant Graph, duplicates []Duplicates, atomEquivalents [][]int) Pathway{
	return Pathway{
		pathway,
		remnant,
		duplicates,
		atomEquivalents,
	}
}

// Duplicates details the bonds involved in pahtway duplicates. Left will be a duplicate in the pathway, whereas right will be part of the remnant
// but broken off from the rest
type Duplicates struct{
	left [][2]int
	right [][2]int
}

// check is a basic error checking function
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// PathwayStepsSaved checks the total steps saved on a pathway by looking at the size of the duplicates
// each duplicate can save the number of edges/nodes in it (depending on edgeMode) minus 1. This is because all those edges/nodes
// would otherwise need to be added individually. The -1 is because you would still need one step to join the duplicate structure.
// Currently edgemode should always be true, unless and until vertex based assembly is implemented
func PathwayStepsSaved(pathway *Pathway, edgeMode bool) int {
	stepsSaved := 0

	for i := 0; i < len(pathway.pathway); i++ {
		if edgeMode {
			stepsSaved += len(pathway.pathway[i].Edges) - 1
		} else {
			stepsSaved += len(pathway.pathway[i].Vertices) - 1
		}
	}

	return stepsSaved
}


// BestPathwayUpdate checks a pathway against the best pathway found so far. If the new pathway has more steps saved, it replaces bestPathway
func BestPathwayUpdate(bestPathway *Pathway, newPathway *Pathway) {

	var mutex = &sync.Mutex{}
	mutex.Lock()
	bestStepsSaved := PathwayStepsSaved(bestPathway, true)
	newStepsSaved := PathwayStepsSaved(newPathway, true)
	if newStepsSaved > bestStepsSaved {
		*bestPathway = *newPathway
	}
	mutex.Unlock()

}


// PathwayEqual tests if two pathways are identical. Mainly used for testing and development
// Graphs are tested for equality (automorphism) not isomorphism
// TODO: needs to be updated for new duplicates objects
func PathwayEqual(pathwayLeft *Pathway, pathwayRight *Pathway) bool {

	if len(pathwayLeft.pathway) != len(pathwayRight.pathway) || len(pathwayLeft.duplicates) != len(pathwayRight.duplicates) {
		return false
	}

	if !GraphEquals(&pathwayLeft.remnant, &pathwayRight.remnant) {
		return false
	}

	if !reflect.DeepEqual(pathwayLeft.duplicates, pathwayRight.duplicates) {
		return false
	}

	for i := 0; i < len(pathwayLeft.pathway); i++ {
		if !GraphEquals(&pathwayLeft.pathway[i], &pathwayRight.pathway[i]) {
			return false
		}
	}

	return true
}

// CopyPathway returns a full copy of a pathway
func CopyPathway(pathway *Pathway) Pathway {
	newGraphs := make([]Graph, 0)
	for _, g := range pathway.pathway {
		newGraphs = append(newGraphs, CopyGraph(&g))
	}


	newDuplicates := CopyDuplicates(pathway.duplicates)
	newRemnant := CopyGraph(&pathway.remnant)
	newAtomEquivalents := helpers.CopySliceOfSlices(pathway.atomEquivalents)

	return NewPathway(newGraphs, newRemnant, newDuplicates, newAtomEquivalents)
}


// CopyDuplicates returns a copy of a duplicates list
func CopyDuplicates(inputDuplicates []Duplicates) []Duplicates {
	var outputDuplicates []Duplicates
	for _, d := range inputDuplicates{
		newLeft := CopyEdgeList(d.left)
		newRight := CopyEdgeList(d.right)
		newDuplicate := Duplicates{newLeft, newRight}
		outputDuplicates = append(outputDuplicates, newDuplicate)
	}
	return outputDuplicates
}

// CopyEdgeList returns a copy of an input edge list
func CopyEdgeList(inputList [][2]int) [][2]int{
	var outputList [][2]int
	for _, edge := range inputList{
		p := [2]int{}
		p[0] = edge[0]
		p[1] = edge[1]
		outputList = append(outputList, p)
	}
	return outputList
}


// PathwayPrint prints full pathway information to stdout
func PathwayPrint(pathway *Pathway) {
	fmt.Println(PathwayString(pathway))
}

// PathwayPrintLog outputs full pathway information to the logger
func PathwayPrintLog(pathway *Pathway) {
	Logger.Debug(PathwayString(pathway))
}

// PathwayString outputs pathway information as a string
func PathwayString(pathway *Pathway) string {
	outString := "Pathway Graphs\n"
	for _, g := range pathway.pathway {
		outString += "======\n"
		outString += GraphPrint(&g) + "\n"
		outString += "======\n"
	}

	outString += "----------\n"
	outString += "Remnant Graph\n"
	outString += GraphPrint(&pathway.remnant) + "\n"
	outString += "----------\n"

	outString += "Duplicated Edges\n"
	for _, e := range pathway.duplicates {
		outString += fmt.Sprintf("%v\n", e)
	}
	outString +="+++++++++++++++\n"
	outString +="###############\n"
	outString += "Atom Equivalents\n"
	for _, e := range pathway.atomEquivalents{
		outString += fmt.Sprintf("%v\n", e)
	}
	outString +="###############\n"

	return outString
}

// AssemblyString returns a string with details of the original graph and pathway
func AssemblyString(pathways []Pathway, originalGraph *Graph) string {
	outString := "ORIGINAL GRAPH\n"
	outString += "+++++++++++++++\n"
	outString += GraphPrint(originalGraph) + "\n"
	outString += "+++++++++++++++\n"
	outString += "PATHWAY\n"

	for _, p := range pathways {
		outString += PathwayString(&p)
	}

	return outString
}

// AssemblyIndex returns the integer assembly index of the pathway, which is the max assembly index for originalGraph (number of edges - 1)
// minus the total steps saved through all the duplicates
func AssemblyIndex(pathway *Pathway, originalGraph *Graph) int {

	index := len(originalGraph.Edges) - 1
	index -= PathwayStepsSaved(pathway, true)
	return index

}

// forbidUpdate adds thisForbid to the forbidden list, updates forbiddenSize, and un-forbids any edges
// that have forbiddenSize > thisForbidSize. This function is part of the process of finding all subgraphs
func forbidUpdate(thisForbid int, thisForbidSize int, forbidden map[int]bool, forbiddenSize map[int]int) {
	for k, v := range forbiddenSize {
		if v > thisForbidSize {
			forbidden[k] = false
		}
	}
	forbidden[thisForbid] = true
	forbiddenSize[thisForbid] = thisForbidSize
}

// nonForbiddenNeighbour returns a non-forbidden neighbour of a subgraph (where a subgraph is a slice of int).
// This is part of the process of finding all subgraphs
func nonForbiddenNeighbour(sub []int, edgeAdjacencies map[int][]int, forbidden map[int]bool) (int, bool) {
	for _, e := range sub {
		for _, adj := range edgeAdjacencies[e] {
			if !forbidden[adj] && !helpers.Contains(sub, adj) {
				return adj, true
			}
		}
	}
	return -1, false
}

// BestAssemblyIndex returns the best possible assembly index of a pathway, based on the maximum possible
// steps saved on the remnant graph. This is used to bound the assembly process
func BestAssemblyIndex(g *Graph, pathway *Pathway) int {
	return AssemblyIndex(pathway, g) - MaxStepsSaved(pathway)
}

// MaxStepsSaved returns the maximum possible additional steps that could be saved within the remnant portion of a
// pathway. This is based on the assuming that each connected component remaining in the remnant can be constructed in the
// shortest possible way, which is bounded by the log base 2 of the number of edges (i.e. repeatedly duplicating an structure,
// from 1 edge, to 2, to 4, to 8 etc.
func MaxStepsSaved(pathway *Pathway) int {
	connectedComponentEdges := ConnectedComponentEdges(&pathway.remnant)
	maxStepsSaved := 0

	for _, c := range connectedComponentEdges {
		numEdges := len(c)
		naiveMA := numEdges - 1
		bestMA := int(math.Floor(math.Log2(float64(numEdges))))
		bestStepsSaved := naiveMA - bestMA
		maxStepsSaved += bestStepsSaved
	}

	return maxStepsSaved

}

// UpdateAtomEquivalents uses the vertexMap created in RecombineGraphs to update the pathway list of atom equivalents
func UpdateAtomEquivalents(pathway *Pathway, vertexMap map[int]int){
	for originalInt, newInt := range vertexMap{
		if newInt != originalInt {
			contains, index := helpers.IntInSliceOfSlices(pathway.atomEquivalents, originalInt)

			if contains {
				pathway.atomEquivalents[index] = append(pathway.atomEquivalents[index], newInt)
			} else {
				pathway.atomEquivalents = append(pathway.atomEquivalents, []int{originalInt, newInt})
			}
		}
	}
}

