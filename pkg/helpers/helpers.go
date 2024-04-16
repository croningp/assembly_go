package helpers

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Contains returns true if i is in slice l
func Contains(l []int, i int) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

// ContainsStr returns true if i is in slice l
func ContainsStr(l []string, i string) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

// MapUpdate updates a map from an int to a slice of ints
// if the k is in the map, v is appended to the slice, otherwise k is added to the map with value {v}
func MapUpdate(k int, v int, m map[int][]int) {

	if _, ok := m[k]; ok {
		m[k] = append(m[k], v)
	} else {
		m[k] = []int{v}
	}
}

// MapUpdateAppend updates a map from an int to a slice of ints
// similar to MapUpdate, but in this case v is []int. If k is in the map, the values of v are appended to m[k]
// and if not then the m[k] is set to be a copy of v
func MapUpdateAppend(k int, v []int, m map[int][]int) {
	if _, ok := m[k]; ok {
		m[k] = append(m[k], v...)
	} else {
		m[k] = make([]int, len(v))
		copy(m[k], v)
	}
}

// CopyAppend takes a slice of int and appends a copy of it to a slice of slices
func CopyAppend(sliceOfSlices [][]int, slice []int) [][]int {
	appendSlice := make([]int, len(slice))
	copy(appendSlice, slice)
	sliceOfSlices = append(sliceOfSlices, appendSlice)
	return sliceOfSlices
}

// CopyAppendSafe calls CopyAppend with a mutex wrapped around it
func CopyAppendSafe(sliceOfSlices [][]int, slice []int) [][]int {
	var mu sync.Mutex
	mu.Lock()
	appendSlice := make([]int, len(slice))
	copy(appendSlice, slice)
	sliceOfSlices = append(sliceOfSlices, appendSlice)
	mu.Unlock()
	return sliceOfSlices
}

// CopySliceOfSlices copies a slice of slices of ints
func CopySliceOfSlices(sliceOfSlices [][]int) [][]int {
	var newSliceOfSlices [][]int
	for _, s := range sliceOfSlices{
		newSliceOfSlices = CopyAppendSafe(newSliceOfSlices, s)
	}
	return newSliceOfSlices
}

// SortInnerIntList takes a slice of slices of ints and sorts the slices
func SortInnerIntList(l [][]int) {
	for _, item := range l {
		sort.Ints(item)
	}
}

// SortSliceOfSlices sorts a slice of int slices lexographically
func SortSliceOfSlices(l [][]int) {

	// first sort the slices themselves
	SortInnerIntList(l)

	sort.Slice(l, func(i1, i2 int) bool {
		return SliceCompare(l[i1], l[i2])
	})

}

// SliceCompare returns true if slice1 < slice2 in the first element that differs
func SliceCompare(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return len(slice1) < len(slice2) // sort first by length
	} else {
		for i := 0; i < len(slice1); i++ {
			if slice1[i] != slice2[i] {
				return slice1[i] < slice2[i] // then by first element that differs
			}
		}
	}
	return false // slices are the same
}

// SliceOfSlicesOverlap returns true if there are any ints duplicated within slice of slices, e.g. {{1,2,3}, {3,4,5}} returns true due to 3
func SliceOfSlicesOverlap(sliceOfSlices [][]int) bool {
	var alreadyFound []int
	for _, slice := range sliceOfSlices {
		for _, n := range slice {
			if Contains(alreadyFound, n) {
				return true
			}
			alreadyFound = append(alreadyFound, n)
		}
	}

	return false
}

// MaxIntSlice returns the maximum value in a slice of ints
func MaxIntSlice(inputSlice []int) int {
	max := 0
	for _, i := range inputSlice {
		if i > max {
			max = i
		}
	}
	return max
}

// TimeStamp returns a formatted string timestamp to use for test data file/folder names etc
func TimeStamp() string {
	t := time.Now()
	return t.Format("20060102150405")
}

func CreateDirectoryIfNotExists(dirName string, relative bool) error {

	var dir string

	if relative {
		dir = filepath.Join(".", dirName)
	} else {
		dir = dirName
	}
	err := os.MkdirAll(dir, os.ModePerm)

	return err
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func CreateFileIfNotExists(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// IntInSliceOfSlices checks if a number is in any of a slice of slices. Returns true/false and index
// the number appears in if it does. If it appears in multiple lists, the index is that of the first
func IntInSliceOfSlices(sliceOfSlices [][]int, checkInt int) (bool, int){

	for i, slice := range sliceOfSlices{
		if Contains(slice, checkInt){
			return true, i
		}
	}

	return false, -1

}
