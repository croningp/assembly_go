package helpers

import (
	"log"
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	var tests = []struct {
		list    []int
		item    int
		desired bool
	}{
		{[]int{1, 2, 3, 4, 5}, 1, true},
		{[]int{1, 2, 3, 4, 5}, 6, false},
		{[]int{}, 1, false},
	}
	for _, test := range tests {
		output := Contains(test.list, test.item)
		if output != test.desired {
			t.Errorf("Expected Contains(%v,%v) to be %v, got %v", test.list, test.item, test.desired, output)
		}
	}
}

func TestMapUpdate(t *testing.T) {
	var testMap = map[int][]int{
		1: {},
		2: {1},
	}

	var desiredMap = map[int][]int{
		1: {2, 3},
		2: {1, 2},
		3: {1, 2},
	}

	MapUpdate(1, 2, testMap)
	MapUpdate(1, 3, testMap)
	MapUpdate(2, 2, testMap)
	MapUpdate(3, 1, testMap)
	MapUpdate(3, 2, testMap)

	eq := reflect.DeepEqual(testMap, desiredMap)

	if !eq {
		t.Error("Test maps not equal to desired map after MapUpdate")
	}

}

func TestMapUpdateAppend(t *testing.T) {
	tests := []struct {
		m       map[int][]int
		k       int
		v       []int
		desired map[int][]int
	}{
		{
			map[int][]int{
				1: {},
			},
			1,
			[]int{1, 2, 3},
			map[int][]int{
				1: {1, 2, 3},
			},
		},
		{
			map[int][]int{
				1: {},
			},
			2,
			[]int{1, 2, 3},
			map[int][]int{
				1: {},
				2: {1, 2, 3},
			},
		},
		{
			map[int][]int{
				1: {1, 2, 3},
				2: {4, 5, 6},
			},
			2,
			[]int{7, 8, 9},
			map[int][]int{
				1: {1, 2, 3},
				2: {4, 5, 6, 7, 8, 9},
			},
		},
	}

	for _, tt := range tests {
		newMap := make(map[int][]int)
		for k, v := range tt.m {
			newMap[k] = make([]int, len(v))
			copy(newMap[k], v)
		}
		MapUpdateAppend(tt.k, tt.v, newMap)
		eq := reflect.DeepEqual(newMap, tt.desired)
		if !eq {
			t.Errorf("MapUpdateAppend error: appending %v: %v to %v, expected %v got %v ",
				tt.k, tt.v, tt.m, tt.desired, newMap)
		}
	}
}

func TestSortInnerIntList(t *testing.T) {
	tests := []struct {
		l      [][]int
		sorted [][]int
	}{
		{
			[][]int{{1}, {2}, {3}, {4}, {5}},
			[][]int{{1}, {2}, {3}, {4}, {5}},
		},
		{
			[][]int{{1, 2, 3}},
			[][]int{{1, 2, 3}},
		},
		{
			[][]int{{1, 3, 2}},
			[][]int{{1, 2, 3}},
		},
		{
			[][]int{{4, 6, 5}, {1, 2, 3}},
			[][]int{{4, 5, 6}, {1, 2, 3}},
		},
		{
			[][]int{{1, 3, 3, 2}, {4, 6, 5}},
			[][]int{{1, 2, 3, 3}, {4, 5, 6}},
		},
	}

	for _, tt := range tests {
		var sorted [][]int

		for _, item := range tt.l {
			copyItem := make([]int, len(item))
			copy(copyItem, item)
			sorted = append(sorted, copyItem)
		}

		SortInnerIntList(sorted)

		eq := reflect.DeepEqual(sorted, tt.sorted)

		if !eq {
			t.Errorf("SortInnerIntList error, sorting %v, expected %v, got %v", tt.l, tt.sorted, sorted)
		}
	}
}

func TestSliceCompare(t *testing.T) {
	tests := []struct {
		slice1  []int
		slice2  []int
		compare bool
	}{
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			false,
		},
		{
			[]int{1, 2},
			[]int{1, 2, 3},
			true,
		},
		{
			[]int{1, 2, 3},
			[]int{1, 2},
			false,
		},
		{
			[]int{1, 2, 3},
			[]int{2, 3, 4},
			true,
		},
		{
			[]int{1, 3, 5},
			[]int{1, 3, 4},
			false,
		},
		{
			[]int{6, 7, 8},
			[]int{1, 3, 4},
			false,
		},
	}

	for _, tt := range tests {
		compare := SliceCompare(tt.slice1, tt.slice2)
		if compare != tt.compare {
			t.Errorf("SliceCompare error, checking if %v before %v, expected %v, got %v",
				tt.slice1, tt.slice2, tt.compare, compare)
		}
	}
}

func TestSortSliceOfSlices(t *testing.T) {
	tests := []struct {
		sliceOfSlices [][]int
		sorted        [][]int
	}{
		{
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			[][]int{{4, 5, 6}, {1, 2, 3}},
			[][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			[][]int{{4, 6, 5}, {1, 3, 2}},
			[][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			[][]int{{4, 6, 5}, {1, 3, 2}, {1, 2}},
			[][]int{{1, 2}, {1, 2, 3}, {4, 5, 6}},
		},
	}

	for _, tt := range tests {
		var sorted [][]int

		for _, item := range tt.sliceOfSlices {
			copyItem := make([]int, len(item))
			copy(copyItem, item)
			sorted = append(sorted, copyItem)
		}

		SortSliceOfSlices(sorted)

		eq := reflect.DeepEqual(sorted, tt.sorted)
		if !eq{
			t.Errorf("SortSliceOfSlices Error, sorting %v, expected %v, got %v", tt.sliceOfSlices, tt.sorted, sorted)
		}
	}
}

func TestSliceOfSlicesOverlap(t *testing.T) {
	tests := []struct{
		sliceOfSlices [][]int
		expected bool
	}{
		{
			[][]int{},
			false,
		},
		{
			[][]int{{1, 2, 3}},
			false,
		},
		{
			[][]int{{1, 2, 3}, {4, 5, 6}},
			false,
		},
		{
			[][]int{{1, 2, 3}, {3, 5, 6}},
			true,
		},
		{
			[][]int{{1, 2, 2}, {4, 5, 6}},
			true,
		},
	}
	for i, tt := range tests{
		result := SliceOfSlicesOverlap(tt.sliceOfSlices)
		if !(result == tt.expected){
			t.Errorf("Error in test %v, input %v, expected %v, got %v", i, tt.sliceOfSlices, tt.expected, result)
		}
	}
}


func TestCreateDirectoryIfNotExists(t *testing.T) {

	err := CreateDirectoryIfNotExists("testdata/test_directory", true)
	check(err)

}

func TestCreateFileIfNotExists(t *testing.T) {
	file := CreateFileIfNotExists("testdata/tmp.txt")
	log.SetOutput(file)
	log.Println("test log")
}

func TestCopySliceOfSlices(t *testing.T) {
	tests := []struct{
		sliceOfSlices [][]int
		newSliceOfSlices [][]int
	}{

		{
			[][]int{{1}},
			[][]int{{1}},
		},
		{
			[][]int{{1, 2, 3}},
			[][]int{{1, 2, 3}},
		},
		{
			[][]int{{1, 2, 3}, {4, 5, 6}},
			[][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			[][]int{{1},{2},{3}},
			[][]int{{1},{2},{3}},
		},
	}

	for _, tt := range tests{
		newSliceOfSlices := CopySliceOfSlices(tt.sliceOfSlices)
		if !reflect.DeepEqual(tt.newSliceOfSlices, newSliceOfSlices){
			t.Errorf("CopySliceOfSlices error, slice %v, expected %v, got %v", tt.sliceOfSlices, tt.newSliceOfSlices, newSliceOfSlices)
		}
		// check modifying original
		tt.newSliceOfSlices[0] = []int{-1}
		if reflect.DeepEqual(tt.newSliceOfSlices, newSliceOfSlices){
			t.Errorf("CopySliceOfSlices error after modifying original, slice %v, modified slice %v, new slice altered %v", tt.sliceOfSlices, tt.newSliceOfSlices, newSliceOfSlices)
		}
	}
}

func TestIntInSliceOfSlices(t *testing.T) {
	tests := []struct{
		sliceOfSlices [][]int
		checkInt int
		contains bool
		index int
	}{
		{
			[][]int{{1, 2}, {3, 4, 5}},
			1,
			true,
			0,

		},
		{
			[][]int{{1, 2}, {3, 4, 5}},
			3,
			true,
			1,

		},
		{
			[][]int{{1, 2}, {3, 4, 5}},
			2,
			true,
			0,

		},
		{
			[][]int{{1, 2}, {3, 4, 5}},
			5,
			true,
			1,

		},
		{
			[][]int{{1, 2, 5}, {3, 4, 5}},
			5,
			true,
			0,

		},
		{
			[][]int{{1, 2}, {3, 4, 5}},
			6,
			false,
			-1,

		},
	}

	for _, tt := range tests{

		contains, index := IntInSliceOfSlices(tt.sliceOfSlices, tt.checkInt)
		chkContains := contains == tt.contains
		chkIndex := index == tt.index

		if !(chkContains && chkIndex){
			t.Errorf("IntInSliceOfSlices error, slices %v, int %v, expected %v, %v, got %v, %v",
				tt.sliceOfSlices, tt.checkInt, tt.contains, tt.index, contains, index)
		}


	}
}