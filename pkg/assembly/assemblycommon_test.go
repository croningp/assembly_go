package assembly

import (
	"reflect"
	"testing"
)

func TestPathwayStepsSaved(t *testing.T) {
	tests := []struct {
		pathway    Pathway
		edgeMode   bool
		stepsSaved int
	}{
		{
			// remnant and duplicated within the Pathway struct shouldn't make any difference
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{}, // not needed for this test
				[][]int{},
			},
			false,
			5,
		},
	}
	for i, tt := range tests {
		stepsSaved := PathwayStepsSaved(&tt.pathway, tt.edgeMode)
		if stepsSaved != tt.stepsSaved {
			t.Errorf("TestPathwayStepsSaved Error in test %v\nInput pathway %v\nExpected %v, Got %v",
				i, tt.pathway, tt.stepsSaved, stepsSaved)
		}
	}
}

func TestBestPathwayUpdate(t *testing.T) {
	tests := []struct {
		bestPathway Pathway
		newPathway  Pathway
		replace     bool
	}{
		{
			// Same Pathway Throughout
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			false,
		},
		{
			// Worse new pathway
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			false,
		},
		{
			// Better new pathway
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			true,
		},
	}

	_ = tests
	for i, tt := range tests {
		bestPathwayCopy := CopyPathway(&tt.bestPathway)
		BestPathwayUpdate(&tt.bestPathway, &tt.newPathway)

		eq := false
		if tt.replace {
			eq = PathwayEqual(&tt.bestPathway, &tt.newPathway)
		} else {
			eq = PathwayEqual(&tt.bestPathway, &bestPathwayCopy)
		}

		if !eq {
			t.Errorf("BestPathwayCopy error in test %v\nOriginal Best Pathway %v\nNew Pathway %v\nNew Best Pathway %v\nShould have been replaced %v",
				i, bestPathwayCopy, tt.newPathway, tt.bestPathway, tt.replace)
		}

	}
}

func TestPathwayEqual(t *testing.T) {
	tests := []struct {
		pathwayLeft  Pathway
		pathwayRight Pathway
		equal        bool
	}{
		{
			// identical pathways
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			true,
		},
		{
			// graphs differ
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			false,
		},
		{
			// duplicates differ
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			false,
		},
		{
			// remnant differs
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				[]Duplicates{},
				[][]int{},
			},
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				},
				NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				[]Duplicates{},
				[][]int{},
			},
			false,
		},
	}

	for i, tt := range tests {
		equal := PathwayEqual(&tt.pathwayLeft, &tt.pathwayRight)
		if equal != tt.equal {
			t.Errorf("PathwayEqual error in test %v\nPathway Left %v\nPathway Right %v\nExpected equality %v\nGot equality %v",
				i, tt.pathwayLeft, tt.pathwayRight, tt.equal, equal)
		}
	}

}

func TestCopyPathway(t *testing.T) {
	tests := []Pathway{
		Pathway{
			[]Graph{
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
			},
			NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
			[]Duplicates{},
			[][]int{},
		},
		Pathway{
			[]Graph{
				NewGraphOnlyFromFile("testdata/graphs/square.txt"),
				NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
			},
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			[]Duplicates{},
			[][]int{},
		},
	}

	for i, pathway := range tests {
		newPathway := CopyPathway(&pathway)
		if !PathwayEqual(&pathway, &newPathway) {
			t.Errorf("CopyPathway error in test %v\nInput Pathway %v\nCopied Pathway %v",
				i, pathway, newPathway)
		}
	}
}

func TestAssemblyIndex(t *testing.T) {
	tests := []struct {
		pathway       Pathway
		originalGraph Graph
		assemblyIndex int
	}{
		{
			Pathway{
				[]Graph{
					NewGraphOnlyFromFile("testdata/graphs/square.txt"),   // 4 edges
					NewGraphOnlyFromFile("testdata/graphs/triangle.txt"), // 3 edges
				},
				NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
				[]Duplicates{},
				[][]int{},
			},
			NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt"), // 12 edges
			6, // = (12-1) - (4-1) - (3-1) = 6
		},
	}

	for _, tt := range tests {
		assemblyIndex := AssemblyIndex(&tt.pathway, &tt.originalGraph)
		if assemblyIndex != tt.assemblyIndex {
			t.Errorf("Error in Assembly Index, Pathway: %v\nOriginal Graph %v\nExpected %v, Got %v", tt.pathway, tt.originalGraph, tt.assemblyIndex, assemblyIndex)
		}
	}
}

func TestMaxStepsSaved(t *testing.T) {
	square := NewGraphOnlyFromFile("testdata/graphs/square.txt")
	twoSquares, _ := RecombineGraphs(&square, &square)
	nineGrid := NewGraphOnlyFromFile("testdata/graphs/nine_grid.txt") // 12 edges
	doubleNine, _ := RecombineGraphs(&nineGrid, &nineGrid)            // 2 x 12 edges

	tests := []struct {
		pathway    Pathway
		stepsSaved int
	}{
		// Only the remnant is used in the function
		{
			Pathway{
				[]Graph{},
				square,
				[]Duplicates{},
				[][]int{},
			},
			1,
		},
		{
			Pathway{
				[]Graph{},
				twoSquares,
				[]Duplicates{},
				[][]int{},
			},
			2,
		},
		{
			Pathway{
				[]Graph{},
				doubleNine,
				[]Duplicates{},
				[][]int{},
			},
			16,
		},
	}

	for _, tt := range tests {
		stepSaved := MaxStepsSaved(&tt.pathway)
		if stepSaved != tt.stepsSaved {
			t.Errorf("Error in MaxStepsSaved, Pathway: %v, Expected %v, Got %v", tt.pathway, tt.stepsSaved, stepSaved)
		}
	}
}

func TestUpdateAtomEquivalents(t *testing.T) {
	testPath1 := NewPathway([]Graph{}, NewGraph([]int{}, [][2]int{}), []Duplicates{}, [][]int{{1, 2}, {4, 5}, {7}})

	tests := []struct {
		pathway        Pathway
		vertexMap      map[int]int
		newEquivalents [][]int
	}{
		{
			CopyPathway(&testPath1),
			map[int]int{
				1: 3,
			},
			[][]int{{1, 2, 3}, {4, 5}, {7}},
		},
		{
			CopyPathway(&testPath1),
			map[int]int{
				1: 3,
				7: 8,
			},
			[][]int{{1, 2, 3}, {4, 5}, {7, 8}},
		},
		{
			CopyPathway(&testPath1),
			map[int]int{
				1: 3,
				9: 10,
			},
			[][]int{{1, 2, 3}, {4, 5}, {7}, {9, 10}},
		},
	}

	for _, tt := range tests {
		UpdateAtomEquivalents(&tt.pathway, tt.vertexMap)
		if !reflect.DeepEqual(tt.newEquivalents, tt.pathway.atomEquivalents) {
			t.Errorf("UpdateAtomEquivalents error\npathway %v\nvertex map %v\nexpected %v\ngot %v",
				testPath1, tt.vertexMap, tt.newEquivalents, tt.pathway.atomEquivalents)
		}
	}
}
