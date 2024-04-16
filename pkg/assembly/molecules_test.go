package assembly

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseMolFile(t *testing.T) {

	var tests = []struct {
		fileName             string
		atomsH, atomsNoH     []string
		bondsH, bondsNoH     [][2]int
		atomIndH, atomIndNoH []int
		bondTH, bondTNoH     []int
	}{
		{
			"testdata/formic_acid_with_H.mol",
			[]string{"C", "O", "O", "H", "H"},
			[]string{"C", "O", "O"},
			[][2]int{{0, 1}, {0, 2}, {0, 3}, {2, 4}},
			[][2]int{{0, 1}, {0, 2}},
			[]int{0, 1, 2, 3, 4},
			[]int{0, 1, 2},
			[]int{2, 1, 1, 1},
			[]int{2, 1},
		},
		{
			"testdata/big_mol_test.mol",
			[]string{"C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","Br","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H","H"},
			[]string{"C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","C","Br"},
			[][2]int{{0, 1},{1, 2},{2, 3},{3, 4},{4, 5},{5, 6},{6, 7},{7, 8},{8, 9},{9, 10},{10, 11},{11, 12},{12, 13},{13, 14},{14, 15},{15, 16},{16, 17},{17, 18},{18, 19},{19, 20},{20, 21},{21, 22},{22, 23},{23, 24},{24, 25},{25, 26},{26, 27},{27, 28},{28, 29},{29, 30},{30, 31},{31, 32},{32, 33},{33, 34},{34, 35},{35, 36},{36, 37},{37, 38},{38, 39},{39, 40},{0, 41},{0, 42},{1, 43},{2, 44},{2, 45},{3, 46},{3, 47},{4, 48},{4, 49},{5, 50},{5, 51},{6, 52},{6, 53},{7, 54},{7, 55},{8, 56},{8, 57},{9, 58},{9, 59},{10, 60},{10, 61},{11, 62},{11, 63},{12, 64},{12, 65},{13, 66},{13, 67},{14, 68},{14, 69},{15, 70},{15, 71},{16, 72},{16, 73},{17, 74},{17, 75},{18, 76},{18, 77},{19, 78},{19, 79},{20, 80},{20, 81},{21, 82},{21, 83},{22, 84},{22, 85},{23, 86},{23, 87},{24, 88},{24, 89},{25, 90},{25, 91},{26, 92},{26, 93},{27, 94},{27, 95},{28, 96},{28, 97},{29, 98},{29, 99},{30, 100},{30, 101},{31, 102},{31, 103},{32, 104},{32, 105},{33, 106},{33, 107},{34, 108},{34, 109},{35, 110},{35, 111},{36, 112},{36, 113},{37, 114},{37, 115},{38, 116},{38, 117},{39, 118},{39, 119}},
			[][2]int{{0, 1},{1, 2},{2, 3},{3, 4},{4, 5},{5, 6},{6, 7},{7, 8},{8, 9},{9, 10},{10, 11},{11, 12},{12, 13},{13, 14},{14, 15},{15, 16},{16, 17},{17, 18},{18, 19},{19, 20},{20, 21},{21, 22},{22, 23},{23, 24},{24, 25},{25, 26},{26, 27},{27, 28},{28, 29},{29, 30},{30, 31},{31, 32},{32, 33},{33, 34},{34, 35},{35, 36},{36, 37},{37, 38},{38, 39},{39, 40}},
			[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,101,102,103,104,105,106,107,108,109,110,111,112,113,114,115,116,117,118,119},
			[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40},
			[]int{2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
			[]int{2,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
		},
	}

	for _, tt := range tests {
		atomsH, bondsH, bondTH, atomIndH := ParseMolFile(tt.fileName, false)
		atomsNoH, bondsNoH, bondTNoH, atomIndNoH := ParseMolFile(tt.fileName, true)

		var err = false
		var errString string
		if !reflect.DeepEqual(atomsH, tt.atomsH){
			err = true
			errString += fmt.Sprintf("Atoms with H expecte %v got %v\n", tt.atomsH, atomsH)
		}
		if !reflect.DeepEqual(bondsH, tt.bondsH){
			err = true
			errString += fmt.Sprintf("Bonds with H expected %v got %v\n", tt.bondsH, bondsH)
		}
		if !reflect.DeepEqual(atomIndH, tt.atomIndH){
			err = true
			errString += fmt.Sprintf("Atom Indices with H expected %v got %v\n", tt.atomIndH, atomIndH)
		}
		if !reflect.DeepEqual(bondTH, tt.bondTH){
			err = true
			errString += fmt.Sprintf("Bond Types with H expected %v got %v\n", tt.bondTH, bondTH)
		}
		if !reflect.DeepEqual(atomsNoH, tt.atomsNoH){
			err = true
			errString += fmt.Sprintf("Atoms without H expected %v got %v\n", tt.atomsNoH, atomsNoH)
		}
		if !reflect.DeepEqual(bondsNoH, tt.bondsNoH){
			err = true
			errString += fmt.Sprintf("Bonds without H expected %v got %v\n", tt.bondsNoH, bondsNoH)
		}
		if !reflect.DeepEqual(atomIndNoH, tt.atomIndNoH){
			err = true
			errString += fmt.Sprintf("Atom Indices without H expected %v got %v\n", tt.atomIndNoH, atomIndNoH)
		}
		if !reflect.DeepEqual(bondTNoH, tt.bondTNoH){
			err = true
			errString += fmt.Sprintf("Bond Types without H expected %v got %v\n", tt.bondTNoH, bondTNoH)
		}
		if err {
			t.Error(errString)
		}
	}

}

func TestParseMolString(t *testing.T) {

	var tests = []struct {
		fileName             string
		atomsH, atomsNoH     []string
		bondsH, bondsNoH     [][2]int
		atomIndH, atomIndNoH []int
		bondTH, bondTNoH     []int
	}{
		{
			"testdata/formic_acid_with_H.mol",
			[]string{"C", "O", "O", "H", "H"},
			[]string{"C", "O", "O"},
			[][2]int{{0, 1}, {0, 2}, {0, 3}, {2, 4}},
			[][2]int{{0, 1}, {0, 2}},
			[]int{0, 1, 2, 3, 4},
			[]int{0, 1, 2},
			[]int{2, 1, 1, 1},
			[]int{2, 1},
		},
	}

	for _, tt := range tests {

		molBytes, _ := ioutil.ReadFile(tt.fileName)
		molString := string(molBytes)

		atomsH, bondsH, bondTH, atomIndH := ParseMolString(molString, false)
		atomsNoH, bondsNoH, bondTNoH, atomIndNoH := ParseMolString(molString, true)

		var err = false
		var errString string
		if !reflect.DeepEqual(atomsH, tt.atomsH){
			err = true
			errString += fmt.Sprintf("Atoms with H expecte %v got %v\n", tt.atomsH, atomsH)
		}
		if !reflect.DeepEqual(bondsH, tt.bondsH){
			err = true
			errString += fmt.Sprintf("Bonds with H expected %v got %v\n", tt.bondsH, bondsH)
		}
		if !reflect.DeepEqual(atomIndH, tt.atomIndH){
			err = true
			errString += fmt.Sprintf("Atom Indices with H expected %v got %v\n", tt.atomIndH, atomIndH)
		}
		if !reflect.DeepEqual(bondTH, tt.bondTH){
			err = true
			errString += fmt.Sprintf("Bond Types with H expected %v got %v\n", tt.bondTH, bondTH)
		}
		if !reflect.DeepEqual(atomsNoH, tt.atomsNoH){
			err = true
			errString += fmt.Sprintf("Atoms without H expected %v got %v\n", tt.atomsNoH, atomsNoH)
		}
		if !reflect.DeepEqual(bondsNoH, tt.bondsNoH){
			err = true
			errString += fmt.Sprintf("Bonds without H expected %v got %v\n", tt.bondsNoH, bondsNoH)
		}
		if !reflect.DeepEqual(atomIndNoH, tt.atomIndNoH){
			err = true
			errString += fmt.Sprintf("Atom Indices without H expected %v got %v\n", tt.atomIndNoH, atomIndNoH)
		}
		if !reflect.DeepEqual(bondTNoH, tt.bondTNoH){
			err = true
			errString += fmt.Sprintf("Bond Types without H expected %v got %v\n", tt.bondTNoH, bondTNoH)
		}
		if err {
			t.Error(errString)
		}
	}

}

func TestParseMultiMolString(t *testing.T) {
	fileName := "testdata/dual_ring_test.sdf"
	molBytes, _ := ioutil.ReadFile(fileName)
	molString := string(molBytes)

	molGraphs := ParseMultiMolString(molString, true)

	for _, g := range molGraphs{
		fmt.Println(g)
	}
}

func TestParseSDFile(t *testing.T) {
	fileName := "testdata/dual_ring_test.sdf"
	molGraphs:= ParseSDFile(fileName, true)
	for _, g := range molGraphs{
		fmt.Println(g)
	}

}


func TestMolListToPathway(t *testing.T) {
	fileName := "testdata/taxol_test.sdf"
	molGraphs:= ParseSDFile(fileName, true)
	originalGraph, pathway := MolListToPathway(molGraphs,[]Duplicates{})

	fmt.Println("Original Graph")
	fmt.Println(originalGraph)
	fmt.Println("Pathway")
	fmt.Println(PathwayString(&pathway))


	fmt.Println("**********************")

	assemblyPathway := AssemblyPathway(originalGraph, pathway, 100, 500, "shortest")

	fmt.Println("RESULTING PATHWAY")
	fmt.Println(PathwayString(&assemblyPathway[0]))
	fmt.Println("Assembly Index")
	fmt.Println(AssemblyIndex(&assemblyPathway[0], &originalGraph))
}