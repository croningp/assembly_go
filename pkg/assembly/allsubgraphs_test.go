package assembly

import (
	"GoAssembly/pkg/helpers"
	"fmt"
	"reflect"
	"testing"
)


func TestAllSubsPTParallel(t *testing.T) {
	tests := []struct {
		g         Graph
		countMode bool
		subGraphs [][]int
		subCount  int
	}{
		{
			// graph with two vertices and one edge between
			NewGraph([]int{1, 2}, [][2]int{{1, 2}}),
			false,
			[][]int{{0}},
			1,
		},
		{
			// graph with two vertices and one edge, countMode so should return no subgraphs
			NewGraph([]int{1, 2}, [][2]int{{1, 2}}),
			true,
			nil,
			1,
		},
		{
			// graph with two vertices and one edge between
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}, {3, 1}}),
			false,
			[][]int{{0}, {1}, {2}, {0, 1}, {1, 2}, {0, 2}, {0, 1, 2}},
			7,
		},
		{
			// triangle graph
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}, {3, 1}}),
			true,
			nil,
			7,
		},
		{
			// disconnected graph 0-1-2 3-4
			NewGraph([]int{0, 1, 2, 3, 4}, [][2]int{{0, 1}, {1, 2}, {3, 4}}),
			false,
			[][]int{{0}, {1}, {2}, {0, 1}},
			4,
		},
		{
			// disconnected graph 0-1-2 triangle plus 3-4-5 linear
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {0, 2}, {3, 4}, {4, 5}}),
			false,
			[][]int{{0}, {1}, {2}, {3}, {4}, {0, 1}, {1, 2}, {2, 0}, {0, 1, 2}, {3, 4}},
			10,
		},
		{
			// aspirin has 579 subgraphs
			MolColourGraph("testdata/aspirin.mol"),
			true,
			nil,
			579,
		},
	}

	for _, tt := range tests {
		subGraphs, subCount, err := AllSubgraphs(&tt.g, tt.countMode)
		check(err)
		helpers.SortSliceOfSlices(subGraphs)
		helpers.SortSliceOfSlices(tt.subGraphs)
		eq1 := reflect.DeepEqual(subGraphs, tt.subGraphs)
		eq2 := subCount == tt.subCount
		if !(eq1 && eq2) {
			errMessage := "Error in AllSubgraphs\n"
			errMessage += fmt.Sprintf("Subgraphs expected %v got %v\n", tt.subGraphs, subGraphs)
			errMessage += fmt.Sprintf("Subgraph count expected %v got %v\n", tt.subCount, subCount)
			t.Error(errMessage)
		}
	}
}


func TestGetAllSubgraphs(t *testing.T) {
	tests := []struct {
		g         Graph
		countMode bool
		subGraphs [][]int
		subCount  int
	}{
		{
			// graph with two vertices and one edge between
			NewGraph([]int{1, 2}, [][2]int{{1, 2}}),
			false,
			[][]int{{0}},
			1,
		},
		{
			// graph with two vertices and one edge, countMode so should return no subgraphs
			NewGraph([]int{1, 2}, [][2]int{{1, 2}}),
			true,
			nil,
			1,
		},
		{
			// graph with two vertices and one edge between
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}, {3, 1}}),
			false,
			[][]int{{0}, {1}, {2}, {0, 1}, {1, 2}, {0, 2}, {0, 1, 2}},
			7,
		},
		{
			// triangle graph
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}, {3, 1}}),
			true,
			nil,
			7,
		},
		{
			// disconnected graph 0-1-2 3-4
			NewGraph([]int{0, 1, 2, 3, 4}, [][2]int{{0, 1}, {1, 2}, {3, 4}}),
			false,
			[][]int{{0}, {1}, {2}, {0, 1}},
			4,
		},
		{
			// disconnected graph 0-1-2 triangle plus 3-4-5 linear
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {0, 2}, {3, 4}, {4, 5}}),
			false,
			[][]int{{0}, {1}, {2}, {3}, {4}, {0, 1}, {1, 2}, {2, 0}, {0, 1, 2}, {3, 4}},
			10,
		},
		{
			// aspirin has 579 subgraphs
			MolColourGraph("testdata/aspirin.mol"),
			true,
			nil,
			579,
		},
	}

	for _, tt := range tests {

		subgraphChan := make(chan []int)
		subCountChan := make(chan int)
		var subGraphs [][]int
		subCount := 0


		go AllSubgraphsToChan(&tt.g, subgraphChan, subCountChan, tt.countMode)


		for{
			select {
				case sub, ok := <- subgraphChan:
					if ok{
						subGraphs = append(subGraphs, sub)
					} else {
						subgraphChan = nil
					}
				case count, ok := <- subCountChan:
					if ok{
						subCount += count
					} else {
						subCountChan = nil
					}
			}

			if subgraphChan == nil && subCountChan == nil{
				break
			}

		}





		// subGraphs, subCount, err := AllSubgraphs(&tt.g, tt.countMode)
		helpers.SortSliceOfSlices(subGraphs)
		helpers.SortSliceOfSlices(tt.subGraphs)
		eq1 := reflect.DeepEqual(subGraphs, tt.subGraphs)
		eq2 := subCount == tt.subCount
		if !(eq1 && eq2) {
			errMessage := "Error in AllSubgraphs\n"
			errMessage += fmt.Sprintf("Subgraphs expected %v got %v\n", tt.subGraphs, subGraphs)
			errMessage += fmt.Sprintf("Subgraph count expected %v got %v\n", tt.subCount, subCount)
			t.Error(errMessage)
		}
	}
}

func TestSubgraphCount(t *testing.T) {
	tests := []struct{
		g Graph
		subCount int
	}{
		{
			MolColourGraph("testdata/aspirin.mol"),
			579,
		},
	}

	for _, tt := range tests{
		subCount := SubgraphCount(&tt.g)
		if subCount != tt.subCount{
			t.Errorf("Subgraph count expected %v got %v\n", tt.subCount, subCount)
		}
	}
}

func TestMolSubgraphCount(t *testing.T) {
	tests := []struct{
		mol string
		molBlock bool
		subCount int
	}{
		{
			"testdata/aspirin.mol",
			false,
			579,
		},
	}

	for _, tt := range tests{
		subCount := MolSubgraphCount(tt.mol, tt.molBlock)
		if subCount != tt.subCount{
			t.Errorf("MolSubgraphCount expected %v got %v\n", tt.subCount, subCount)
		}
	}
}

func TestAspirinSubs(t *testing.T){
	aspirin := MolColourGraph("testdata/aspirin.mol")
	subs, _, _ := AllSubgraphs(&aspirin, false)
	for _, s := range subs{
		fmt.Println(s)
	}
}