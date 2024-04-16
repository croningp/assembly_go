package assembly

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestIndividualise(t *testing.T) {

	tests := []struct {
		partition      [][]int
		vertex         int
		individualised [][]int
	}{
		{
			[][]int{{1, 2, 3}, {4, 5}, {6}},
			1,
			[][]int{{1}, {2, 3}, {4, 5}, {6}},
		},
		{
			[][]int{{1}, {2, 3, 4}, {5, 6}},
			3,
			[][]int{{1}, {3}, {2, 4}, {5, 6}},
		},
	}

	for _, tt := range tests {
		individualised := Individualise(tt.partition, tt.vertex)
		if !reflect.DeepEqual(individualised, tt.individualised) {
			t.Errorf("Individualise error, partition %v, vertex %v, expected %v, got %v",
				tt.partition, tt.vertex, tt.individualised, individualised)
		}
	}

}

func TestDegreeInPart(t *testing.T) {

	tests := []struct {
		g      Graph
		vertex int
		part   []int
		degree int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			1,
			[]int{2, 3, 4},
			2,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			1,
			[]int{3},
			0,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			4,
			[]int{1, 2},
			1,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			4,
			[]int{},
			0,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			1,
			[]int{1, 2, 3, 4},
			2,
		},
	}

	for _, tt := range tests {
		degree := DegreeInPart(&tt.g, tt.vertex, tt.part)
		if degree != tt.degree {
			t.Errorf("DegreeInPart error, graph %v, vertex %v, partn %v, expected %v, got %v",
				tt.g, tt.vertex, tt.part, tt.degree, degree)
		}
	}
}

func TestShatter(t *testing.T) {
	tests := []struct {
		graph      Graph
		partLeft   []int
		partRight  []int
		shattering [][]int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{1, 2, 3},
			[]int{4},
			[][]int{{2}, {1, 3}},
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{2, 3, 4},
			[]int{1},
			[][]int{{3}, {2, 4}},
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3, 4},
			[][]int{{1, 2, 3, 4}},
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{1, 2},
			[]int{3, 4},
			[][]int{{1, 2}},
		},
	}

	for _, tt := range tests {
		shattering := Shatter(&tt.graph, tt.partLeft, tt.partRight)
		if !reflect.DeepEqual(shattering, tt.shattering) {
			t.Errorf("Shattering error, graph %v, partLeft %v, partRight %v, expected %v, got %v",
				tt.graph, tt.partLeft, tt.partRight, tt.shattering, shattering)
		}
	}
}

func TestDegree(t *testing.T) {
	tests := []struct {
		graph  Graph
		vertex int
		degree int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			1,
			2,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			2,
			2,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			3,
			2,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			4,
			2,
		},
	}

	for _, tt := range tests {
		degree := Degree(&tt.graph, tt.vertex)
		if degree != tt.degree {
			t.Errorf("Degree error, graph %v, vertex %v, expected %v, got %v",
				tt.graph, tt.vertex, tt.degree, degree)
		}
	}
}

func TestEquitableRefinement(t *testing.T) {
	tests := []struct {
		graph            Graph
		partition        [][]int
		refinedPartition [][]int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			[][]int{{1, 3, 7, 9}, {2, 4, 6, 8}, {5}},
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[][]int{{1, 3, 7, 9}, {2, 4, 6, 8}, {5}},
			[][]int{{1, 3, 7, 9}, {2, 4, 6, 8}, {5}},
		},
	}

	for _, tt := range tests {
		refinedPartition := EquitableRefinement(&tt.graph, tt.partition)
		if !reflect.DeepEqual(refinedPartition, tt.refinedPartition) {
			t.Errorf("EquitableRefinement error, graph %v, partition %v, expected %v, got %v",
				tt.graph, tt.partition, tt.refinedPartition, refinedPartition)
		}
	}
}

func TestIsEquitable(t *testing.T) {
	tests := []struct {
		graph     Graph
		partition [][]int
		equitable bool
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			false,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[][]int{{1, 3, 7, 9}, {2, 4, 6, 8}, {5}},
			true,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}},
			true,
		},
	}

	for _, tt := range tests {
		equitable := IsEquitable(&tt.graph, tt.partition)
		if equitable != tt.equitable {
			t.Errorf("IsEquitable error, graph %v, partition %v, expected %v, got %v",
				tt.graph, tt.partition, tt.equitable, equitable)
		}
	}
}

func TestCoarsestEquitableColourings(t *testing.T) {
	g := NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt")
	CoarsestEquitableColourings(&g, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}})
	CoarsestEquitableColourings(&g, [][]int{{1, 3, 7, 9}, {2, 4, 6, 8}, {5}})

	//fmt.Println(EquitableRefinement(&g, [][]int{{1}, {3, 7, 9}, {2, 4, 6, 8}, {5}}))

}

func TestIsDiscrete(t *testing.T) {
	tests := []struct {
		colouring [][]int
		discrete  bool
	}{
		{
			[][]int{{1}, {2}, {3}, {4}, {5}},
			true,
		},
		{
			[][]int{{1, 2}, {3}, {4}, {5}},
			false,
		},
		{
			[][]int{{1}, {2}, {3}, {4}, {5, 6}},
			false,
		},
		{
			[][]int{},
			true,
		},
	}

	for _, tt := range tests {
		discrete := IsDiscrete(tt.colouring)
		if discrete != tt.discrete {
			t.Errorf("IsDiscrete error, colouring %v, expected %v, got %v",
				tt.colouring, tt.discrete, discrete)
		}
	}
}

func TestDiscreteColouringToIntSlice(t *testing.T) {
	tests := []struct {
		colouring [][]int
		intSlice  []int
	}{
		{
			[][]int{{1}},
			[]int{1},
		},
		{
			[][]int{{1}, {2}, {3}, {4}},
			[]int{1, 2, 3, 4},
		},
		{
			[][]int{{9}, {7}, {6}, {8}, {5}},
			[]int{9, 7, 6, 8, 5},
		},
	}

	for _, tt := range tests {
		intSlice := DiscreteColouringToIntSlice(tt.colouring)
		if !reflect.DeepEqual(intSlice, tt.intSlice) {
			t.Errorf("DiscreteColouringToIntSlice error, colouring %v, expected %v, got %v",
				tt.colouring, tt.intSlice, intSlice)
		}
	}
}

func TestPermuteGraph(t *testing.T) {

	squareGraph := NewGraphOnlyFromFile("testdata/graphs/square.txt")

	PermuteGraph(&squareGraph, []int{3, 1, 2, 4})

	tests := []struct {
		graph       Graph
		permutation []int
		permuted    Graph
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{2, 1, 3, 4},
			NewColourGraph([]int{2, 1, 3, 4}, [][2]int{{2, 1}, {1, 3}, {3, 4}, {4, 2}}, []string{}, []string{}),
		},

		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{3, 2, 4, 1},
			NewColourGraph([]int{4, 2, 1, 3}, [][2]int{{4, 2}, {2, 1}, {1, 3}, {3, 4}}, []string{}, []string{}),
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			[]int{3, 2, 4, 1},
			NewColourGraph([]int{4, 2, 1, 3}, [][2]int{{4, 2}, {2, 1}, {1, 3}, {3, 4}}, []string{"Red", "Blue", "Red", "Blue"}, []string{"A", "B", "B", "A"}),
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[]int{2, 4, 1, 3, 7, 9, 8, 6, 5},
			NewColourGraph([]int{3, 1, 4, 2, 9, 8, 5, 7, 6}, [][2]int{{3, 1}, {1, 4}, {2, 9}, {9, 8}, {5, 7}, {7, 6}, {3, 2}, {2, 5}, {1, 9}, {9, 7}, {4, 8}, {8, 6}}, []string{}, []string{}),
		},
	}

	for _, tt := range tests {
		permutedGraph := PermuteGraph(&tt.graph, tt.permutation)
		if !GraphEquals(&tt.permuted, &permutedGraph) {
			t.Errorf("PermuteGraph error\ngraph %v\npermutatiopn %v\nexpected %v\ngot %v",
				tt.graph, tt.permutation, tt.permuted, permutedGraph)
		}
	}

}

func TestSliceGreaterThan(t *testing.T) {
	tests := []struct {
		sliceLeft   []int
		sliceRight  []int
		greaterThan bool
	}{
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			false,
		},
		{
			[]int{1, 2, 3},
			[]int{1, 2},
			true,
		},
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3, 4},
			false,
		},
		{
			[]int{2, 3},
			[]int{1, 2, 3, 4},
			true,
		},
		{
			[]int{1, 2, 3, 4, 5},
			[]int{4},
			false,
		},
		{
			[]int{2},
			[]int{1, 2, 3, 4, 5, 6, 7},
			true,
		},
	}

	for _, tt := range tests {
		greaterThan := SliceGreaterThan(tt.sliceLeft, tt.sliceRight)
		if greaterThan != tt.greaterThan {
			t.Errorf("SliceGreaterThan error, sliceLeft %v, sliceRight %v, expected %v, got %v",
				tt.sliceLeft, tt.sliceRight, tt.greaterThan, greaterThan)
		}
	}
}

func TestSliceEqual(t *testing.T) {
	tests := []struct {
		sliceLeft  []int
		sliceRight []int
		equal      bool
	}{
		{
			[]int{},
			[]int{},
			true,
		},
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			true,
		},
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3, 4},
			false,
		},
		{
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3},
			false,
		},
		{
			[]int{4, 5, 6},
			[]int{1, 2, 3},
			false,
		},
	}

	for _, tt := range tests {
		equal := SliceEqual(tt.sliceLeft, tt.sliceRight)
		if equal != tt.equal {
			t.Errorf("SliceEqual error, sliceLeft %v, sliceRight %v, expected %v, got %v",
				tt.sliceLeft, tt.sliceRight, tt.equal, equal)
		}
	}
}

func TestFlattenEdgeList(t *testing.T) {
	tests := []struct {
		edgeList [][2]int
		flatList []int
	}{
		{
			[][2]int{{1, 2}, {3, 4}},
			[]int{1, 2, 3, 4},
		},
		{
			[][2]int{{2, 3}, {3, 4}, {4, 1}},
			[]int{2, 3, 3, 4, 4, 1},
		},
		{
			[][2]int{},
			nil,
		},
	}

	for _, tt := range tests {
		flatList := FlattenEdgeList(tt.edgeList)
		if !reflect.DeepEqual(flatList, tt.flatList) {
			t.Errorf("FlattenEdgeList error, edgeList %v, expected %v, got %v",
				tt.edgeList, tt.flatList, flatList)
		}
	}
}

func TestGraphGreaterThan(t *testing.T) {
	tests := []struct {
		graphLeft   Graph
		graphRight  Graph
		greaterThan bool
	}{
		{
			NewGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}}),
			NewGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {1, 4}}),
			true,
		},
		{
			NewGraph([]int{1, 2, 3, 4, 5}, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}, {1, 5}}),
			NewGraph([]int{1, 2, 3, 4, 5}, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}, {4, 5}}),
			false,
		},
	}

	for _, tt := range tests {
		greaterThan := GraphGreaterThan(&tt.graphLeft, &tt.graphRight)
		if greaterThan != tt.greaterThan {
			t.Errorf("GraphGreaterThan error, graphLeft %v, graphRight %v, expected %v, got %v",
				tt.graphLeft, tt.graphRight, tt.greaterThan, greaterThan)
		}
	}
}

func TestSearchTree(t *testing.T) {
	nineGrid := NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt")
	initialColouring := [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}}
	SearchTree(&nineGrid, initialColouring, true)
}

func TestRandomPermutation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	inputPerm := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < 10; i++ {
		fmt.Println(RandomPermutation(inputPerm))
	}
}

func TestRandomPermutationList(t *testing.T) {
	inputList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println(RandomPermutationList(inputList, 10))

}

func TestCanonicalGraphTest(t *testing.T) {
	nineGrid := NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt")
	pass := CanonicalGraphTest(&nineGrid, [][]int{{1, 3}, {7, 9}, {2, 4, 6}, {8}, {5}}, 100)
	fmt.Println("all canonicals equal", pass)
}

func TestRandomGraph(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	RandomGraph(5, 3, []string{"A", "B"})
}

func TestGraphColourPartition(t *testing.T) {
	tests := []struct {
		graph           Graph
		colourPartition [][]int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			[][]int{{2, 4}, {1, 3}},
		},
		{
			NewColourGraph([]int{1, 2, 3, 4, 5, 6}, [][2]int{{1, 2},{2, 3},{3,4},{4,5},{5,6}}, []string{"A", "B", "B", "C", "B", "A"}, []string{}),
			[][]int{{1, 6}, {2, 3, 5}, {4}},
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"),
			[][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		},
	}

	for _, tt := range tests{
		colourPartition := GraphColourPartition(&tt.graph)
		if !reflect.DeepEqual(colourPartition, tt.colourPartition){
			t.Errorf("GraphColourPartition error, graph %v, expected %v, got %v",
				tt.graph, tt.colourPartition, colourPartition)
		}
	}
}

func TestRandomGraphCanonicalTest(t *testing.T) {
	result := RandomGraphCanonicalTest(50, 100, [2]int{3, 20}, [2]int{3, 10}, []string{"A", "B", "X", "F", "yy"})
	print("Random Graph Canonical Test Result: ", result)
}



func TestMaxEqualLevel(t *testing.T) {
	tests := []struct{
		sliceLeft []int
		sliceRight []int
		maxEqualLevel int
	}{
		{
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 5, 5},
			2,
		},
		{
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
			4,
		},
		{
			[]int{0, 1, 2 ,3, 4},
			[]int{0, 9, 9, 9, 9},
			0,
		},
		{
			[]int{0, 1, 2 ,3, 4},
			[]int{9, 9, 9, 9, 9},
			-1,
		},
	}

	for _, tt := range tests{
		maxEqualLevel := MaxEqualLevel(tt.sliceLeft, tt.sliceRight)
		if maxEqualLevel != tt.maxEqualLevel{
			t.Errorf("MaxEqualLevel error, sliceLeft %v, sliceRight %v, expected %v, got %v",
				tt.sliceLeft, tt.sliceRight, tt.maxEqualLevel, maxEqualLevel)
		}
	}
}


func TestEdgeColourConversion(t *testing.T) {
	squareColoured := NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt")

	fmt.Println("Square Coloured:\n", squareColoured)

	edgeGraph := EdgeColourConversion(&squareColoured)
	fmt.Println("Square Coloured Edge Converted:\n", edgeGraph)

}

func TestEdgeColourRandomGraph(t *testing.T) {
	fmt.Println(EdgeColourRandomGraph(6, 3, []string{"A", "B"}, []string{"X", "Y", "Z"}))
}

func TestEdgeColourRandomCanonicalTest(t *testing.T) {
	result := EdgeColourRandomCanonicalTest(20, 20, [2]int{3, 8}, [2]int{1, 3}, []string{"A", "B", "C"}, []string{"X","Y","Z"})
	fmt.Println("Edge Colour Random Test Result: ", result)
	if !result{
		t.Error("TestEdgeColourRandomCanonicalTest fail - creation of random edge coloured graph, conversion to layered graph, permutation and canonicalisation check. Some or all canonical graphs not equal")
	}
}

func TestGraphsIsomorphic(t *testing.T) {
	tests := []struct{
		graphLeft Graph
		graphRight Graph
		isomorphic bool
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			NewGraphOnlyFromFile("testdata/graphs/square_isomorph.txt"),
			true,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			NewGraphOnlyFromFile("testdata/graphs/square_isomorph.txt"),
			false,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			NewGraphOnlyFromFile("testdata/graphs/not_square.txt"),
			false,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			NewGraphOnlyFromFile("testdata/graphs/square_coloured_isomorphic.txt"),
			true,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			NewGraphOnlyFromFile("testdata/graphs/square_coloured_relabeled.txt"),
			true,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/a_test_1.txt"),
			NewGraphOnlyFromFile("testdata/graphs/a_test_2.txt"),
			true,
		},
	}

	for _, tt := range tests{
		isomorphic := GraphsIsomorphic(&tt.graphLeft, &tt.graphRight)
		if isomorphic != tt.isomorphic{
			t.Errorf("GraphsIsomorphic error, graphLeft %v, graphRight %v, expected %v, got %v",
				tt.graphLeft, tt.graphRight, tt.isomorphic, isomorphic)
		}
	}
}

func TestGraphVertexRelabel(t *testing.T) {
	tests := []struct {
		graph Graph
		labeling []int
		newGraph Graph
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			[]int{0, 1, 2, 3},
			NewGraphOnlyFromFile("testdata/graphs/square_coloured_relabeled.txt"),
		},
	}

	for _, tt := range tests{
		newGraph := GraphVertexRelabel(&tt.graph, tt.labeling)
		if !GraphEquals(&newGraph, &tt.newGraph){
			t.Errorf("GraphVertexRelabel error, graph %v, labeling %v, expected new graph %v, got %v",
				tt.graph, tt.labeling, tt.newGraph, newGraph)
		}
	}
}
