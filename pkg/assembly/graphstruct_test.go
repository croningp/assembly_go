package assembly

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestInducedSubgraph(t *testing.T) {
	tests := []struct {
		g        Graph
		vertices []int
		gNew     Graph
	}{
		{
			NewGraph([]int{0, 1, 2, 3, 4}, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}),
			[]int{1, 2, 3},
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}}),
		},
		{
			NewGraph([]int{1, 2, 3, 4, 5}, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}),
			[]int{2, 3, 4},
			NewGraph([]int{2, 3, 4}, [][2]int{{2, 3}, {3, 4}}),
		},
		{
			// disconnected graph
			NewGraph([]int{1, 2, 3, 4, 5, 6}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {4, 5}, {5, 6}, {6, 4}}),
			[]int{1, 2, 3, 4, 5},
			NewGraph([]int{1, 2, 3, 4, 5}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {4, 5}}),
		},
		{
			// disconnected graph
			NewGraph([]int{1, 2, 3, 4, 5, 6}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {4, 5}, {5, 6}, {6, 4}}),
			[]int{1, 2, 4, 5},
			NewGraph([]int{1, 2, 4, 5}, [][2]int{{1, 2}, {4, 5}}),
		},
	}

	for _, tt := range tests {

		gNew, err := InducedSubgraph(&tt.g, tt.vertices)

		if err != nil {
			t.Log("Error returned", err)
		}

		eq1 := reflect.DeepEqual(gNew.Vertices, tt.gNew.Vertices)
		eq2 := reflect.DeepEqual(gNew.Edges, tt.gNew.Edges)
		if !(eq1 && eq2) {
			t.Errorf("Induced Subgraph Error:\n"+
				"Expected %v %v"+
				"Got %v %v",
				tt.gNew.Vertices, tt.gNew.Edges, gNew.Vertices, gNew.Edges)
		}
	}
}

func TestNewGraphFromFile(t *testing.T) {
	tests := []struct {
		fileName      string
		vertices      []int
		edges         [][2]int
		vertexColours []string
		edgeColours   []string
		name          string
	}{
		{
			"testdata/graphs/triangle.txt",
			[]int{0, 1, 2},
			[][2]int{{0, 1}, {1, 2}, {2, 0}},
			[]string{},
			[]string{},
			"Triangle Graph (uncoloured)",
		},
		{
			"testdata/graphs/square_coloured.txt",
			[]int{1, 2, 3, 4},
			[][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}},
			[]string{"Red", "Blue", "Red", "Blue"},
			[]string{"A", "B", "B", "A"},
			"Square Graph (coloured)",
		},
	}

	for _, tt := range tests {
		g, name, err := NewGraphFromFile(tt.fileName)
		check(err)
		eq1 := name == tt.name
		eq2 := reflect.DeepEqual(g.Vertices, tt.vertices)
		eq3 := reflect.DeepEqual(g.Edges, tt.edges)
		eq4 := reflect.DeepEqual(g.VertexColours, tt.vertexColours)
		eq5 := reflect.DeepEqual(g.EdgeColours, tt.edgeColours)
		if !(eq1 && eq2 && eq3 && eq4 && eq5) {
			errString := "NewGraphFromFile error:\n"
			errString += fmt.Sprintf("Description expected %v, got %v - %v\n", tt.name, name, eq1)
			errString += fmt.Sprintf("Vertices expected %v, got %v - %v\n", tt.vertices, g.Vertices, eq2)
			errString += fmt.Sprintf("Edges expected %v, got %v - %v\n", tt.edges, g.Edges, eq3)
			errString += fmt.Sprintf("Vertex Colours expected %v, got %v - %v\n", tt.vertexColours, g.VertexColours, eq4)
			errString += fmt.Sprintf("Edge Colours expected %v, got %v - %v\n", tt.edgeColours, g.EdgeColours, eq5)
			t.Error(errString)
		}
	}
}

func TestNewGraphFromString(t *testing.T) {
	tests := []struct {
		fileName      string
		vertices      []int
		edges         [][2]int
		vertexColours []string
		edgeColours   []string
		name          string
	}{
		{
			"testdata/graphs/triangle.txt",
			[]int{0, 1, 2},
			[][2]int{{0, 1}, {1, 2}, {2, 0}},
			[]string{},
			[]string{},
			"Triangle Graph (uncoloured)",
		},
		{
			"testdata/graphs/square_coloured.txt",
			[]int{1, 2, 3, 4},
			[][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}},
			[]string{"Red", "Blue", "Red", "Blue"},
			[]string{"A", "B", "B", "A"},
			"Square Graph (coloured)",
		},
	}

	for _, tt := range tests {

		// extract string from file
		graphBytes, _ := ioutil.ReadFile(tt.fileName)
		graphString := string(graphBytes)

		g, name, err := NewGraphFromString(graphString)
		check(err)
		eq1 := name == tt.name
		eq2 := reflect.DeepEqual(g.Vertices, tt.vertices)
		eq3 := reflect.DeepEqual(g.Edges, tt.edges)
		eq4 := reflect.DeepEqual(g.VertexColours, tt.vertexColours)
		eq5 := reflect.DeepEqual(g.EdgeColours, tt.edgeColours)
		if !(eq1 && eq2 && eq3 && eq4 && eq5) {
			errString := "NewGraphFromFile error:\n"
			errString += fmt.Sprintf("Description expected %v, got %v - %v\n", tt.name, name, eq1)
			errString += fmt.Sprintf("Vertices expected %v, got %v - %v\n", tt.vertices, g.Vertices, eq2)
			errString += fmt.Sprintf("Edges expected %v, got %v - %v\n", tt.edges, g.Edges, eq3)
			errString += fmt.Sprintf("Vertex Colours expected %v, got %v - %v\n", tt.vertexColours, g.VertexColours, eq4)
			errString += fmt.Sprintf("Edge Colours expected %v, got %v - %v\n", tt.edgeColours, g.EdgeColours, eq5)
			t.Error(errString)
		}
	}
}

func TestNewGraphOnlyFromString(t *testing.T) {
	tests := []struct {
		fileName      string
		vertices      []int
		edges         [][2]int
		vertexColours []string
		edgeColours   []string
	}{
		{
			"testdata/graphs/triangle.txt",
			[]int{0, 1, 2},
			[][2]int{{0, 1}, {1, 2}, {2, 0}},
			[]string{},
			[]string{},
		},
		{
			"testdata/graphs/square_coloured.txt",
			[]int{1, 2, 3, 4},
			[][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}},
			[]string{"Red", "Blue", "Red", "Blue"},
			[]string{"A", "B", "B", "A"},
		},
	}

	for _, tt := range tests {

		// extract string from file
		graphBytes, _ := ioutil.ReadFile(tt.fileName)
		graphString := string(graphBytes)

		g := NewGraphOnlyFromString(graphString)

		eq1 := reflect.DeepEqual(g.Vertices, tt.vertices)
		eq2 := reflect.DeepEqual(g.Edges, tt.edges)
		eq3 := reflect.DeepEqual(g.VertexColours, tt.vertexColours)
		eq4 := reflect.DeepEqual(g.EdgeColours, tt.edgeColours)
		if !(eq1 && eq2 && eq3 && eq4) {
			errString := "NewGraphFromFile error:\n"
			errString += fmt.Sprintf("Vertices expected %v, got %v - %v\n", tt.vertices, g.Vertices, eq1)
			errString += fmt.Sprintf("Edges expected %v, got %v - %v\n", tt.edges, g.Edges, eq2)
			errString += fmt.Sprintf("Vertex Colours expected %v, got %v - %v\n", tt.vertexColours, g.VertexColours, eq3)
			errString += fmt.Sprintf("Edge Colours expected %v, got %v - %v\n", tt.edgeColours, g.EdgeColours, eq4)
			t.Error(errString)
		}
	}
}

func TestListPairSort(t *testing.T) {
	tests := []struct {
		input  [][2]int
		output [][2]int
	}{
		{
			[][2]int{{1, 2}, {3, 4}, {5, 6}},
			[][2]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			[][2]int{{1, 2}, {5, 6}, {3, 4}},
			[][2]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			[][2]int{{2, 1}, {6, 5}, {4, 3}},
			[][2]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			[][2]int{{6, 5}, {1, 2}, {4, 3}},
			[][2]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			[][2]int{{1, 8}, {1, 7}},
			[][2]int{{1, 7}, {1, 8}},
		},
	}

	for _, tt := range tests {
		output := ListPairSort(tt.input)
		eq := reflect.DeepEqual(output, tt.output)
		if !eq {
			t.Errorf("LisPairSort error, input %v, expected %v, got %v",
				tt.input, tt.output, output)
		}
	}
}

func TestGraphEquals(t *testing.T) {
	tests := []struct {
		g1 Graph
		g2 Graph
		eq bool
	}{
		{ //0
			NewGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 4}}),
			NewGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 4}}),
			true,
		},
		{ //1
			NewGraph([]int{3, 2, 1}, [][2]int{{2, 3}, {2, 1}, {1, 3}}),
			NewGraph([]int{1, 2, 3}, [][2]int{{1, 2}, {2, 3}, {3, 1}}),
			true,
		},
		{ //2
			NewGraph([]int{3, 2, 1, 4}, [][2]int{{2, 3}, {2, 1}, {1, 3}, {4, 3}}),
			NewGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}}),
			true,
		},
		{ //3
			NewGraph([]int{3, 2, 1, 4}, [][2]int{{2, 3}, {2, 1}, {1, 3}, {1, 4}}),
			NewGraph([]int{1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}}),
			false,
		},
		{ //4
			NewColourGraph([]int{3, 2, 1, 4},
				[][2]int{{2, 3}, {2, 1}, {1, 3}, {1, 4}},
				[]string{"A", "B", "A", "B"},
				[]string{"A", "B", "A", "B"}),
			NewColourGraph([]int{1, 2, 3, 4},
				[][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}},
				[]string{"A", "B", "A", "B"},
				[]string{"A", "B", "A", "B"}),
			false,
		},
		{ //5
			NewColourGraph([]int{3, 2, 1, 4},
				[][2]int{{2, 3}, {2, 1}, {1, 3}, {1, 4}},
				[]string{"A", "B", "A", "B"},
				[]string{"A", "B", "A", "B"}),
			NewColourGraph([]int{1, 2, 3, 4},
				[][2]int{{1, 2}, {2, 3}, {3, 1}, {1, 4}},
				[]string{"A", "B", "A", "B"},
				[]string{"B", "A", "A", "B"}),
			true,
		},
		{ //6
			NewColourGraph([]int{3, 2, 1, 4},
				[][2]int{{2, 3}, {2, 1}, {1, 3}, {1, 4}},
				[]string{"A", "B", "A", "B"},
				[]string{"A", "B", "A", "B"}),
			NewColourGraph([]int{1, 2, 3, 4},
				[][2]int{{1, 2}, {2, 3}, {3, 1}, {1, 4}},
				[]string{"A", "B", "B", "A"}, // not equal
				[]string{"A", "B", "A", "B"}),
			false,
		},
		{ //7
			NewColourGraph([]int{3, 2, 1, 4},
				[][2]int{{2, 3}, {2, 1}, {1, 3}, {1, 4}}, // [1 2][1 3][1 4][2 3]
				[]string{"A", "B", "A", "B"},             // 1:A, 2:B, 3:A, 4:B
				[]string{"A", "B", "A", "B"}),            // [1 2]:B, [1 3]:A, [1 4]:B, [2 3]:A,
			NewColourGraph([]int{1, 2, 3, 4},
				[][2]int{{1, 2}, {2, 3}, {3, 1}, {1, 4}}, // [1 2][1 3][1 4][2 3]
				[]string{"A", "B", "A", "B"},             // 1:A, 2:B, 3:A, 4:B
				[]string{"A", "B", "B", "A"}),            // [1 2]:A, [1 3]:B, [1 4]:A, [2 3]:B
			false,
		},
		{ //8
			NewColourGraph([]int{3, 2, 1, 4},
				[][2]int{{2, 3}, {2, 1}, {1, 3}, {1, 4}},
				[]string{"A", "B", "A", "B"},
				[]string{"A", "B", "A", "B"}),
			NewColourGraph([]int{1, 2, 3, 4},
				[][2]int{{1, 2}, {2, 3}, {3, 1}, {1, 4}},
				[]string{"A", "B", "B", "A"},  // not equal
				[]string{"A", "B", "B", "A"}), // not equal
			false,
		},
	}

	for i, tt := range tests {
		eq := GraphEquals(&tt.g1, &tt.g2)
		if eq != tt.eq {
			errStr := fmt.Sprintf("GraphEquals error on test %v:\n", i)
			errStr += fmt.Sprintf(GraphPrint(&tt.g1) + "\n")
			errStr += fmt.Sprintf(GraphPrint(&tt.g2) + "\n")
			errStr += fmt.Sprintf("Expected %v got %v", tt.eq, eq)
			t.Error(errStr)
		}
	}
}

func TestVertexColourMap(t *testing.T) {
	tests := []struct {
		g Graph
		m map[int]string
	}{
		{NewColourGraph([]int{1, 2, 3, 4},
			[][2]int{{1, 2}, {2, 3}, {3, 4}},
			[]string{"A", "A", "A", "B"},
			[]string{}),
			map[int]string{
				1: "A",
				2: "A",
				3: "A",
				4: "B",
			},
		},
	}

	for _, tt := range tests {
		m := VertexColourMap(&tt.g)
		eq := reflect.DeepEqual(tt.m, m)
		if !eq {
			t.Errorf("VertexColourMap Error\nGraph %v %v %v %v\nExpected %v\nGot %v",
				tt.g.Vertices, tt.g.Edges, tt.g.VertexColours, tt.g.EdgeColours, tt.m, m)
		}
	}
}

func TestEdgeColourMap(t *testing.T) {
	tests := []struct {
		g Graph
		m map[[2]int]string
	}{
		{NewColourGraph([]int{1, 2, 3, 4},
			[][2]int{{1, 2}, {2, 3}, {3, 4}},
			[]string{},
			[]string{"A", "A", "B"}),
			map[[2]int]string{
				[2]int{1, 2}: "A",
				[2]int{2, 3}: "A",
				[2]int{3, 4}: "B",
			},
		},
	}

	for _, tt := range tests {
		m := EdgeColourMap(&tt.g)
		eq := reflect.DeepEqual(tt.m, m)
		if !eq {
			t.Errorf("VertexColourMap Error\nGraph %v %v %v %v\nExpected %v\nGot %v",
				tt.g.Vertices, tt.g.Edges, tt.g.VertexColours, tt.g.EdgeColours, tt.m, m)
		}
	}
}

func TestCopyGraph(t *testing.T) {
	tests := []Graph{
		NewGraphOnlyFromFile("testdata/graphs/fish_graph.txt"),
		NewGraphOnlyFromFile("testdata/graphs/square.txt"),
		NewGraphOnlyFromFile("testdata/graphs/square_coloured.txt"),
		NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
	}
	for i, graph := range tests {
		newGraph := CopyGraph(&graph)
		if !GraphEquals(&graph, &newGraph) {
			t.Errorf("CopyGraph error, test %v\nInput graph %v\nCopied Graph %v", i, graph, newGraph)
		}
	}
}
