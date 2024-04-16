package assembly

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestAssembly(t *testing.T) {

	// This testing covers ExtendPathway and CheckSubgraphMatches due to the recursive nature of the whole thing



	tests := []struct{
		graph Graph
		assemblyIndex int
	}{
		{
			NewGraphOnlyFromFile("testdata/graphs/square.txt"),
			2,
		},
		{
			NewGraphOnlyFromFile("testdata/graphs/triangle.txt"),
			2,
		},
		{
			MolColourGraph("testdata/aspirin.mol"),
			8,
		},

	}

	workers := 100
	buf := 100
	for _, tt := range tests{
		testPathway := Assembly(tt.graph, workers, buf, "shortest")[0]
		pathwayString := PathwayString(&testPathway)
		assemblyIndex := AssemblyIndex(&testPathway, &tt.graph)
		if assemblyIndex != tt.assemblyIndex{
			t.Errorf("Assembly error in graph %v\nWorkers, Buffer %v %v\nExpected %v got %v\n%v",
				tt.graph, workers, buf, tt.assemblyIndex, assemblyIndex, pathwayString)
		}
	}

}


func AssemblyTimeTest(graph *Graph, workers int, bufSize int, variant string) (int, time.Duration){

	start := time.Now()
	pathways := Assembly(*graph, 1000, 100000, variant)
	elapsed := time.Now().Sub(start)

	index := AssemblyIndex(&pathways[0], graph)

	return index, elapsed


}

func TestAssemblyToString(t *testing.T) {
	molBytes, _ := ioutil.ReadFile("testdata/aspirin.mol")
	molBlock := string(molBytes)

	resultString := AssemblyToString(molBlock, 100, 100, "shortest")
	fmt.Println(resultString)
}
func TestAssemblyFromMultiMolString(t *testing.T) {
	fileName := "testdata/dual_ring_test.sdf"
	molBytes, _ := ioutil.ReadFile(fileName)
	molString := string(molBytes)
	fmt.Println("MOL STRING: ", molString)

	// strings.Split not working propery here...

	graphs := ParseMultiMolString(molString, true)

	fmt.Println(graphs)


	originalGraph := graphs[0]

	pathways := AssemblyFromMultiMolString(molString, 100, 500, "shortest")

	fmt.Println(AssemblyString(pathways, &originalGraph))
}

func TestAssemblySDFBlock(t *testing.T) {
	fileName := "testdata/dual_ring_test.sdf"
	molBytes, _ := ioutil.ReadFile(fileName)
	molString := string(molBytes)
	//fmt.Println("MOL STRING: ", molString)

	fmt.Println(AssemblySDFBlock(molString, 100, 500, "shortest"))
}

func TestDGImprovements(t *testing.T){
	fmt.Println("test test")

	aspirin := MolColourGraph("testdata/aspirin.mol")
	_= aspirin

	pathways := Assembly(aspirin, 100, 500, "shortest")
	index := AssemblyIndex(&pathways[0], &aspirin)
	fmt.Println("Index: ", index)
	fmt.Println(AssemblyString(pathways, &aspirin))
}


func TestRaceCondition(t *testing.T){
	fmt.Println("Test Race Condition")

	mol := MolColourGraph("testdata/inconsistency.mol")

	counts := make(map[int]int)
	times := 1000
	for i:=0; i < times ; i++{
		pathways := Assembly(mol, 1, 500, "shortest")
		index := AssemblyIndex(&pathways[0], &mol)
		// fmt.Println("Index: ", index)

		if _, ok := counts[index]; ok{
			counts[index] += 1
		} else {
			counts[index] =1
		}

	}
	fmt.Println(counts)
	//fmt.Println(AssemblyString(pathways, &mol))

}

func TestEdgeSort(t *testing.T) {
	testSlice := [][2]int{{3, 4}, {2, 1}, {6, 5}}
	EdgeSort(testSlice)
}

func TestFlatten(t *testing.T) {
	testSlice := [][]int{{3, 4}, {2, 1}, {6, 5}}
	flatSlice := Flatten(testSlice)
	fmt.Println(flatSlice)
}

func TestSubgraphEdgeCompare(t *testing.T) {
	testL := [][2]int{{1, 2}, {3, 4}, {5, 7}}
	testR := [][2]int{{6, 5}, {3, 4}, {2, 1}}
	fmt.Println(SubgraphEdgeCompare(testL, testR))

}