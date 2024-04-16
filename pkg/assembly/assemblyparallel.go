package assembly

import (
	"GoAssembly/pkg/helpers"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sort"
	"sync"
	"time"
)

// This file contains functions specific to the main parallel implementation of the assembly algorithm, the main one being Assembly

// Worker takes pathways from the jobs queue and extends them, placing the results back in the jobs queue
func Worker(jobs chan Pathway, graph *Graph, bestPathways *[]Pathway, activeWorkers *WorkerCounter, variant string, done chan bool) {

	// Initially wait for a job (a Pathway)
	currentPathway := <-jobs

	for {

		// Extend the pathway, putting any results back in the jobs queue for other workers to pick up
		ExtendPathway(&currentPathway, bestPathways, graph, variant, jobs, activeWorkers)

		// TODO: rename, since activeWorkers is now really active jobs
		if activeWorkers.NumWorkers() == 0{
			done <- true
		}
		currentPathway = <-jobs

	}
}

// WorkerCounter is a struct to track the number of active workers. Zero active workers and an empty queue means the process is complete
// WorkerCounter.activeWorkers is initialised as -1, then set to 1 by the first call of Increment, to avoid the process terminating
// on startup
// TODO: Rename, as this now counts active jobs rather than workers
type WorkerCounter struct {
	activeWorkers int64
	mu            sync.Mutex
}
func (workerCounter *WorkerCounter) Increment() {
	workerCounter.mu.Lock()
	if workerCounter.activeWorkers == -1 {
		workerCounter.activeWorkers = 1
	} else {
		workerCounter.activeWorkers++
	}

	workerCounter.mu.Unlock()
}
func (workerCounter *WorkerCounter) Decrement() {
	workerCounter.mu.Lock()
	workerCounter.activeWorkers--
	workerCounter.mu.Unlock()
}

func (workerCounter *WorkerCounter) NumWorkers() int64 {
	var numWorkers int64
	workerCounter.mu.Lock()
	numWorkers = workerCounter.activeWorkers
	workerCounter.mu.Unlock()
	return numWorkers
}

// AssemblyFromMultiMolString take a set of graphs and use as starting pathway. TODO: include duplicates also
// the original graph is the first one, then a pathway with the final residue at the end
func AssemblyFromMultiMolString(mols string, numWorkers int, chanBufferSize int, variant string) []Pathway {
	graphs := ParseMultiMolString(mols, true)

	originalGraph := graphs[0]
	startingPathway := Pathway{
		graphs[1:len(graphs)-1],
		graphs[len(graphs)-1],
		[]Duplicates{},
		[][]int{},
	}

	outputPathway := AssemblyPathway(originalGraph, startingPathway, numWorkers, chanBufferSize, variant)
	return outputPathway
}

// AssemblyToString takes an input molecule in the form of a mol block and returns the
// assembly index and pathway as a string. This is for use in calling from assembly calculator API
// rather than from the executable
func AssemblyToString(molBlock string, numWorkers int, chanBufferSize int, variant string) string {

	graph := MolBlockColourGraph(molBlock)
	start := time.Now()
	pathways := Assembly(graph, numWorkers, chanBufferSize, variant)
	elapsed := time.Now().Sub(start)

	assemblyIndex := AssemblyIndex(&pathways[0], &graph)

	outString := AssemblyString(pathways, &graph)
	outString += fmt.Sprintf("Assembly Index: %v\n", assemblyIndex)
	outString += fmt.Sprintf("Time (seconds):  %v\n", elapsed.Seconds())

	return outString
}

// AssemblySDFBlock is similar to AssemblyToString, but takes an sdf block as input rather than a
// single mol block. The SDF block is interpreted as the first mol being the original molecule,
// the last being the remnant, and the intermediates being the duplicates. The main assembly algorithm
// works on the remnant, and then extends the duplicates. There's no check that the input pathway is sensible,
// e.g. there is no check that the remnant or duplicates exist within the target molecule.
func AssemblySDFBlock(sdfBlock string, numWorkers int, chanBufferSize int, variant string) string {
	graphs := ParseMultiMolString(sdfBlock, true)
	originalGraph, starterPathway := MolListToPathway(graphs, []Duplicates{})

	start := time.Now()
	pathways := AssemblyPathway(originalGraph, starterPathway, numWorkers, chanBufferSize, variant)
	elapsed := time.Now().Sub(start)

	assemblyIndex := AssemblyIndex(&pathways[0], &originalGraph)

	outString := AssemblyString(pathways, &originalGraph)
	outString += fmt.Sprintf("Assembly Index: %v\n", assemblyIndex)
	outString += fmt.Sprintf("Time (seconds):  %v\n", elapsed.Seconds())

	return outString
}


// Assembly takes an input graph and returns assembly pathways, either a shortest pathway
// or all shortest pathways depending on the variant (shortest, all_shortest). More variants may
// be added. all_shortest may have some duplication at the moment, currently only shortest is supported. TODO: test all_shortest
// The process spawns a number of worker goroutines, that take pathways from the jobs queue
// and extend them by a step in all possible ways through finding duplicates. The resultant
// extended pathways are placed back into the jobs queue to be extended further. The jobs queue (a channel)
// is buffered. If full, the goroutine will process the job in a depth first manner, until there is space in the queue.
// numWorkers is the number of worker threads, and chanBufferSize is the buffer size of the queue.
func Assembly(graph Graph, numWorkers int, chanBufferSize int, variant string) []Pathway{

	initPathway := NewStartingPathway(graph)
	return AssemblyPathway(graph, initPathway, numWorkers, chanBufferSize, variant)
}

// AssemblyPathway is called by Assembly to generate pathways based on an initial graph. This can also be used as an entry
// point if starting with a pathway, e.g. to specify a duplicate that must be used.
func AssemblyPathway(graph Graph, initPathway Pathway, numWorkers int, chanBufferSize int, variant string) []Pathway {

	// will return shortest pathway, or all shortest pathways depending on the variant
	// could be extended to all pathways
	ValidateVariants(variant)

	bestPathways := []Pathway{initPathway}
	jobs := make(chan Pathway, chanBufferSize)
	done := make(chan bool, 1)

	jobs <- initPathway

	// Listener for keyboard interrupt. Sends done signal on interrupt, exiting and outputing best pathway found.
	cInt := make(chan os.Signal, 1)
	signal.Notify(cInt, os.Interrupt)
	go func(){
		for sig := range cInt {
			fmt.Printf("Captured %v - exiting with best found pathway\n", sig)
			done <- true
		}
	}()


	activeWorkers := WorkerCounter{
		1,
		sync.Mutex{},
	}

	// var workerMu sync.Mutex

	for i := 0; i < numWorkers; i++ {
		go Worker(jobs, &graph, &bestPathways, &activeWorkers, variant, done)
	}


	<-done // block until something sent to done

	return bestPathways
}


// ExtendPathway takes an input pathway and checks if it can be extended by matching duplicate subgraphs
// New pathways are placed in the jobs queue if it is not full. If the jobs queue is full, the thread proceeds to extend pathways
// In a depth-first fashion. The matching process cycles through subgraphs within the remnant graph based on the path trace algorithm
// described in "Automatic Enumeration of All Connected Subgraphs, Rucker &  Rucker, 2000". For each of those subgraphs, CheckSubgraphMatches
// is called which uses a similar subgraph search on the remaining part of the remnant, checking for matches. There is some
// pruning within the process also.
func ExtendPathway(currentPathway *Pathway, bestPathways *[]Pathway, originalGraph *Graph, variant string, jobs chan Pathway, activeWorkers *WorkerCounter) {

	// If this pathway cannot in principle be extended to a better pathway than the best found so far, then return
	// TODO: If implementing output of all pathways, this will need to be disabled
	if AssemblyIndex(&(*bestPathways)[0], originalGraph) < BestAssemblyIndex(originalGraph, currentPathway) {

		// activeWorkers is set to 1 at the start of the program for the first job, then is incremented when
		// new jobs are added to the jobs pool.
		activeWorkers.Decrement()

		return
	}

	// We only need to consider subgraphs up to half the size of the main graph when checking for duplicates
	sizesToCheck := int(math.Floor(float64(len(currentPathway.remnant.Edges)) / 2))

	// Update the best pathway list if this pathway is better than the ones found so far
	BestPathwayListUpdate(bestPathways, currentPathway, variant)

	// Initialisation for the path tracing algorithm to find all subgraphs
	edgeAdjacencies := currentPathway.remnant.EdgeAdjacencies()
	forbidden := make(map[int]bool)    // map for whether a vertex is forbidden
	forbiddenSize := make(map[int]int) // map of the size of the list a vertex was forbidden from
	var sub []int

	// for each edge
	for i := 0; i < len(currentPathway.remnant.Edges); i++ {
		sub = []int{i} // subgraph starts with just the current edge
		for {

			neighbour, found := nonForbiddenNeighbour(sub, edgeAdjacencies, forbidden)

			// grow the subgraph if a valid neighbour is found
			if found && (len(sub) <= sizesToCheck) {
				sub = append(sub, neighbour)

				// break out this subgraph from the main graph
				subgraph, remnant := BreakGraphOnEdges(&currentPathway.remnant, sub)

				// the subgraph and remnant are sent into CheckSubgraphMatches, which will look for the subgraph being contained within the rest
				// of the remnant. The matches that are found are used to construct new pathways that are placed into the jobs queue.
				// CheckSubgraphMatches returns true if any matches are found (there might be multiple matches)
				match := true
				if len(sub) > 1 {
					match = CheckSubgraphMatches(currentPathway, bestPathways, originalGraph, &subgraph, &remnant, variant, jobs, activeWorkers)
				}

				// if we have found matches of the current subgraph, or if the subgraph is of size 1, then we continue and keep trying to grow the
				// subgraph. If no matches, there will be no matches to larger subgraphs and we can continue to the backtracking steps
				if match{
					continue
				}
			}

			// backtracking steps
			thisForbidSize := len(sub)
			thisForbid := sub[len(sub)-1]                                      // item to forbid
			sub = sub[:len(sub)-1]                                             // pop the forbidden element out
			forbidUpdate(thisForbid, thisForbidSize, forbidden, forbiddenSize) // update forbidden lists

			// backtracking from the first edge means we are done with this edge
			if len(sub) == 0 {
				break
			}

		}
	}

	// activeWorkers is set to 1 at the start of the program for the first job, then is incremented when
	// new jobs are added to the jobs pool.
	activeWorkers.Decrement()

}

// CheckSubgraphMatches takes the takes a remnant graph from a pathway, and a subgraph of that graph, and looks for matches within the remaining part of the remnant.
// It does this by searching through subgraphs of the remnant in a similar way to how the input subgraph was found in ExtendPathway
func CheckSubgraphMatches(currentPathway *Pathway, bestPathways *[]Pathway, originalGraph *Graph, subgraph *Graph, remnant *Graph,
	variant string, jobs chan Pathway, activeWorkers *WorkerCounter) bool {

	// We are only looking for subgraphs of size k to match
	k := len(subgraph.Edges)

	// Initialisation for the path tracing algorithm to find all subgraphs
	edgeAdjacencies := remnant.EdgeAdjacencies() // map of which edges are adjacent, maps edge index to slice of edge indices
	forbidden := make(map[int]bool)              // map for whether a vertex is forbidden
	forbiddenSize := make(map[int]int)           // map of the size of the list a vertex was forbidden from
	var sub []int                                // the current subgraph under construction
	match := false


	for i := 0; i < len(remnant.Edges); i++ {

		sub = []int{i}

		for {

			neighbour, found := nonForbiddenNeighbour(sub, edgeAdjacencies, forbidden)

			// not building subgraphs bigger than k, which is the size of the subgraph we are matching
			if found && (len(sub) <= k) {

				// if a neighbour is available add it to sub
				sub = append(sub, neighbour)

				// check all subgraphs of length k for a match
				if len(sub) == k {

					possibleDuplicate, newRemnant := BreakGraphOnEdges(remnant, sub)

					// Check if the subgraph and possible duplicate are isomorphic. First SubgraphEdgeCompare checks that the
					// sorted edge list of the subgraph is less than that of possibleDuplicate. This is to prevent duplication as otherwise
					// all matchine pairs of subgraphs would be investigated twice
					if SubgraphEdgeCompare(subgraph.Edges, possibleDuplicate.Edges) && GraphsIsomorphic(subgraph, &possibleDuplicate) {
						match = true


						newPathway := CopyPathway(currentPathway)
						newPathway.pathway = append(newPathway.pathway, CopyGraph(subgraph))

						// create bond lists from duplicates
						dupLeft := CopyEdgeList(subgraph.Edges)
						dupRight := CopyEdgeList(possibleDuplicate.Edges)
						newDuplicate := Duplicates{dupLeft, dupRight}
						newPathway.duplicates = append(newPathway.duplicates, newDuplicate)

						// the remnant to use in the new pathway is the a graph comprised of newRemnant and possibleDuplicte
						// but no longer connected. TODO: refactor variable names - I have inadvertenly made them quite confusing
						// As an example, if the graph was A-B-C-D (with A, B, C, D being subgraphs), then we found A was the same as B
						// The new pathway would have duplicate A and Remnant B C-D (i.e. with B not connected to C-D)
						newGraph, vertexMap := RecombineGraphs(&newRemnant, &possibleDuplicate)
						newPathway.remnant = CopyGraph(&newGraph)
						UpdateAtomEquivalents(&newPathway, vertexMap)

						// The newPathway will be added to the jobs queue, if there is space in the queue
						// If there is not, then this goroutine will also extend the pathway, i.e. proceed in a depth first way
						// This is done as workers are likely to write to the jobs channel more than they
						// read from it, and they may all be blocked if the channel is buffered
						// TODO: rename as activeWorkers is now more like active jobs
						select {
						case jobs <- newPathway:
							activeWorkers.Increment()
						default:
							activeWorkers.Increment()
							ExtendPathway(&newPathway, bestPathways, originalGraph, variant, jobs, activeWorkers)

						}

					}


				}
				continue
			}

			// backtracking steps
			thisForbidSize := len(sub)
			thisForbid := sub[len(sub)-1]                                      // item to forbid
			sub = sub[:len(sub)-1]                                             // pop the forbidden element out
			forbidUpdate(thisForbid, thisForbidSize, forbidden, forbiddenSize) // update forbidden lists

			// backtracking from the first edge means we are done with this edge
			if len(sub) == 0 {
				break
			}

		}

	}

	return match
}

// BestPathwayListUpdate replaces the pathways in bestPathways if newPathway is shorter
// and appends to bestPathway is newPathway is equal in length to bestPathway and we are using
// all_shortest variant (note: all_shortest not yet fully implemented/tested)
func BestPathwayListUpdate(bestPathways *[]Pathway, newPathway *Pathway, variant string) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	bestStepsSaved := PathwayStepsSaved(&(*bestPathways)[0], true)
	newStepsSaved := PathwayStepsSaved(newPathway, true)

	// if more steps are saved, replace the contents of bestStepsSaved
	// if all the shortest paths are requred and the same number of steps is saved, then append the new pathway to the best pathway list
	// TODO: can incorporate check to see if saved pathway already exists using canonicalisation check on duplicates and remnant
	if newStepsSaved > bestStepsSaved {
		*bestPathways = []Pathway{*newPathway}
	} else if newStepsSaved == bestStepsSaved && variant == "all_shortest" {
		*bestPathways = append(*bestPathways, *newPathway)
	}

	mutex.Unlock()
}

// ValidateVariants checks that the variant (e.g. shortest, all_shortest) is valid
func ValidateVariants(variant string) {
	allowedVariants := []string{"shortest", "all_shortest", "all"}
	if helpers.ContainsStr(allowedVariants, variant) {
		return
	}
	errString := "Invalid variant: " + variant
	check(errors.New(errString))
}


// SubgraphEdgeCompare compares two subgraphs to see if the sorted edge list of one is less than the sorted edge list of the other
func SubgraphEdgeCompare(left [][2]int, right [][2]int) bool {

	sortedLeft := Flatten(EdgeSort(left))
	sortedRight := Flatten(EdgeSort(right))

	return helpers.SliceCompare(sortedLeft, sortedRight)
}

// EdgeSort sorts edges by their vertices, and then sorts the whole edge list by the first element. It returns a slice of slice, rather than a
// slice of [2]int, since it's simpler and all that is needed for this purpose
func EdgeSort(edges [][2]int) [][]int {

	// TODO: complete testing

	var newEdgeList [][]int
	for _, e := range edges{
		newEdgeList = append(newEdgeList, []int{e[0], e[1]} )
	}


	for _, e := range newEdgeList{
		sort.Ints(e)
	}


	sort.Slice(newEdgeList, func(i, j int) bool {
		// edge cases
		if len(newEdgeList[i]) == 0 && len(newEdgeList[j]) == 0 {
			return false // two empty slices - so one is not less than other i.e. false
		}
		if len(newEdgeList[i]) == 0 || len(newEdgeList[j]) == 0 {
			return len(newEdgeList[i]) == 0 // empty slice listed "first" (change to != 0 to put them last)
		}

		// both slices len() > 0, so can test this now:
		return newEdgeList[i][0] < newEdgeList[j][0]
	})


	return newEdgeList
}

// Flatten flattens a [][]int into a []int, e.g. {{1, 2}, {3, 4}} -> {1, 2, 3, 4}
func Flatten(s [][]int) []int {
	var outputSlice []int
	for _, outer := range s{
		for _, element := range outer{
			outputSlice = append(outputSlice, element)
		}
	}
	return outputSlice
}