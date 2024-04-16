package assembly

import (
	"GoAssembly/pkg/helpers"
	"reflect"
	"sort"
	"testing"
)

func TestBreakGraphOnEdges(t *testing.T) {
	tests := []struct {
		g               Graph
		edges           []int
		raisesException bool
		breakGraph      Graph
		remnantGraph    Graph
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]int{0, 1},
			false,
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}}),
			NewGraph([]int{3, 4, 1}, [][2]int{{3, 4}, {4, 1}}),
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
			[]int{0},
			false,
			NewGraph([]int{0, 1}, [][2]int{{0, 1}}),
			NewGraph([]int{0, 1, 2}, [][2]int{{1, 2}, {2, 0}}),
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
			[]int{0, 1, 2},
			false,
			NewColourGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 4}}, []string{"Red", "Blue", "Red", "Blue"}, []string{"A", "B", "B"}),
			NewColourGraph([]int{1, 4}, [][2]int{{1, 4}}, []string{"Red", "Blue"}, []string{"A"}),
		},
	}

	for _, tt := range tests {
		breakGraph, remnantGraph := BreakGraphOnEdges(&tt.g, tt.edges)

		eq1 := GraphEquals(&breakGraph, &tt.breakGraph)
		eq2 := GraphEquals(&remnantGraph, &tt.remnantGraph)
		if !(eq1 && eq2) {
			t.Errorf("BreakGraphOnEdges error\ninput graph %v\nbreak on edges %v\nexpected\n%v\n%v\ngot\n%v\n%v",
				tt.g, tt.edges, tt.breakGraph, tt.remnantGraph, breakGraph, remnantGraph)
		}
	}

}


func TestRecombineGraphs(t *testing.T) {
	tests := []struct {
		gLeft   Graph
		gRight  Graph
		gOutput Graph
	}{
		{
			NewColourGraph([]int{1, 2, 4}, [][2]int{{1, 2}, {1, 4}}, []string{"B", "B", "A"}, []string{"Y", "X"}),
			NewColourGraph([]int{2, 4, 3, 5}, [][2]int{{2, 3}, {3, 4}, {3, 5}}, []string{"A", "B", "A", "A"}, []string{"X", "X", "Y"}),
			NewColourGraph([]int{1, 2, 4, 6, 7, 3, 5}, [][2]int{{1, 2}, {1, 4}, {6, 3}, {3, 7}, {3, 5}}, []string{"B", "B", "A", "A", "B", "A", "A"}, []string{"Y", "X", "X", "X", "Y"}),
		},
	}

	for _, tt := range tests {
		gOutput, _ := RecombineGraphs(&tt.gLeft, &tt.gRight)
		if !GraphEquals(&gOutput, &tt.gOutput) {
			t.Errorf("Error in Recombine Graphs\nG1: %v\nG2: %v\nExpected %v\nGot %v", tt.gLeft, tt.gRight, tt.gOutput, gOutput)
		}
	}
}

func TestConnectedComponent(t *testing.T) {
	tests := []struct {
		graph     Graph
		edge      int
		component []int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			0,
			[]int{0, 1, 2, 3},
		},
		{
			NewGraph([]int{0, 1, 2, 3}, [][2]int{{0, 1}, {2, 3}}),
			0,
			[]int{0},
		},
		{
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 1}, {4, 5}}),
			1,
			[]int{0, 1, 2, 3},
		},
		{
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 1}, {4, 5}}),
			4,
			[]int{4},
		},
		{
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {2, 3}, {4, 5}, {3, 1}}),
			4,
			[]int{0, 1, 2, 4},
		},
	}

	for _, tt := range tests {
		edgeAdj := tt.graph.EdgeAdjacencies()
		component := ConnectedComponent(tt.edge, edgeAdj)
		sort.Ints(component)
		if !reflect.DeepEqual(component, tt.component) {
			t.Errorf("Error in ConnectedComponent, Graph: %v, edge: %v, Expected %v, Got %v", tt.graph, tt.edge, tt.component, component)
		}
	}
}

func TestConnectedComponentEdges(t *testing.T) {
	tests := []struct {
		graph      Graph
		components [][]int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[][]int{{0, 1, 2, 3}},
		},
		{
			NewGraph([]int{0, 1, 2, 3}, [][2]int{{0, 1}, {2, 3}}),
			[][]int{{0}, {1}},
		},
		{
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 1}, {4, 5}}),
			[][]int{{0, 1, 2, 3}, {4}},
		},
		{
			NewGraph([]int{0, 1, 2, 3, 4, 5}, [][2]int{{0, 1}, {1, 2}, {2, 3}, {4, 5}, {3, 1}}),
			[][]int{{0, 1, 2, 4}, {3}},
		},
	}

	for _, tt := range tests{
		components := ConnectedComponentEdges(&tt.graph)
		helpers.SortSliceOfSlices(components)
		helpers.SortSliceOfSlices(tt.components)
		if !reflect.DeepEqual(components, tt.components){
			t.Errorf("Error in ConnectedComponent, Graph: %v, Expected %v, Got %v", tt.graph, tt.components, components)
		}
	}
}
