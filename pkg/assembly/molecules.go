package assembly

import (
	"GoAssembly/pkg/helpers"
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Code relating specifically to parsing molecules from mol files / mol blocks


func MolColourGraph(molFile string) Graph {
	atomTypes, bonds, bondTypes, atomIndices := ParseMolFile(molFile, true)
	bondTypeString := make([]string, len(bondTypes))
	for i, bondType := range bondTypes{
		switch bondType {
		case 1:
			bondTypeString[i] = "single"
		case 2:
			bondTypeString[i] = "double"
		case 3:
			bondTypeString[i] = "triple"
		case 4:
			bondTypeString[i] = "aromatic"
		default:
			bondTypeString[i] = "error"
		}

	}
	outGraph := NewColourGraph(atomIndices, bonds, atomTypes, bondTypeString)
	return outGraph
}


func MolBlockColourGraph(molBlock string) Graph {
	// TODO: refactor code reuse with MolColourGraph
	// TODO: ParseMolString should not always be true
	atomTypes, bonds, bondTypes, atomIndices := ParseMolString(molBlock, true)
	bondTypeString := make([]string, len(bondTypes))
	for i, bondType := range bondTypes{
		switch bondType {
		case 1:
			bondTypeString[i] = "single"
		case 2:
			bondTypeString[i] = "double"
		case 3:
			bondTypeString[i] = "triple"
		case 4:
			bondTypeString[i] = "aromatic"
		default:
			bondTypeString[i] = "error"
		}

	}
	outGraph := NewColourGraph(atomIndices, bonds, atomTypes, bondTypeString)
	return outGraph
}
// ParseMultiMolString parses string input that is in the form of an sdfile, i.e. a sequence of mol blocks with $$$$ as delimiter
func ParseMultiMolString(multiMolString string, stripH bool) []Graph {
	multiMolString = strings.ReplaceAll(multiMolString, "\r\n", "\n")  // deal with windows insertion of carriage return
	mols := strings.Split(multiMolString, "$$$$\n")
	var molGraphs []Graph
	for _, mol := range mols{
		molGraph := MolBlockColourGraph(mol)
		if len(molGraph.Vertices) != 0 {
			molGraphs = append(molGraphs, molGraph)
		}
	}
	return molGraphs
}


func MolListToPathway(mols []Graph, duplicates []Duplicates) (Graph, Pathway){
	// TODO: validate inputs
	originalGraph := mols[0]
	pathway := Pathway{
			mols[1:len(mols)-1],
			mols[len(mols)-1],
			duplicates,
		[][]int{},
	}

	return originalGraph, pathway

}

func ParseSDFile(filePath string, stripH bool) []Graph {
	fileBytes, _ := ioutil.ReadFile(filePath)
	fileString := string(fileBytes)
	return ParseMultiMolString(fileString, stripH)
}

func ParseMolScanner(scanner *bufio.Scanner, stripH bool)([]string, [][2]int, []int, []int){
	var atoms []string
	var atomIndices []int
	var bonds [][2]int
	var bondTypes []int

	i := 0
	atNum := 0
	var atomEnd, bondEnd int
	for scanner.Scan() {
		if i == 3 {
			line3 := scanner.Text()

			atomString := strings.ReplaceAll(line3[:3], " ", "")
			bondString := strings.ReplaceAll(line3[3:6], " ", "")

			atoms, err := strconv.Atoi(atomString)
			check(err)
			bonds, err := strconv.Atoi(bondString)
			check(err)

			atomEnd = 4 + atoms
			bondEnd = atomEnd + bonds

		}

		// atom block
		if i >= 4 && i < atomEnd {
			line := strings.Fields(scanner.Text())
			atoms = append(atoms, line[3])

			atomIndices = append(atomIndices, atNum)
			atNum++
		}

		// bond block
		if i >= atomEnd && i < bondEnd {

			bondLine := scanner.Text()
			at1String := strings.ReplaceAll(bondLine[:3], " ", "")
			at2String := strings.ReplaceAll(bondLine[3:6], " ", "")
			typeString := strings.ReplaceAll(bondLine[6:9], " ", "")

			// line := strings.Fields(scanner.Text())
			at1, err := strconv.Atoi(at1String)
			check(err)
			at2, err := strconv.Atoi(at2String)
			check(err)
			bondType, err := strconv.Atoi(typeString)
			check(err)
			bonds = append(bonds, [2]int{at1 - 1, at2 - 1}) // -1 as changing to zero indexing
			bondTypes = append(bondTypes, bondType)

		}

		i++
	}


	if stripH{
		return stripHAtoms(atoms, bonds, bondTypes, atomIndices)
	} else {
		return atoms, bonds, bondTypes, atomIndices
	}
}

// ParseMolFile extracts lists of atoms, bonds, bond types from a mol file
func ParseMolFile(filePath string, stripH bool) ([]string, [][2]int, []int, []int) {
	f, err := os.Open(filePath)
	check(err)
	scanner := bufio.NewScanner(f)
	atoms, bonds, bondTypes, atomIndices := ParseMolScanner(scanner, stripH)

	cErr := f.Close()
	check(cErr)

	return atoms, bonds, bondTypes, atomIndices

}

// ParseMolString extracts lists of atoms, bonds, bond types from a string of a mol block
func ParseMolString(molString string, stripH bool)([]string, [][2]int, []int, []int){
	scanner := bufio.NewScanner(strings.NewReader(molString))
	atoms, bonds, bondTypes, atomIndices := ParseMolScanner(scanner, stripH)
	return atoms, bonds, bondTypes, atomIndices
}

// stripHAtoms takes out all the H atoms, while maintaining the correct connectivity etc
func stripHAtoms(atoms []string, bonds [][2] int, bondTypes []int, atomIndices []int)([]string, [][2]int, []int, []int){
		atomMap := make(map[int]int)
		var newAtoms []string
		var newBonds [][2]int
		var newBondTypes []int
		var newAtomIndices []int
		var HIndices []int

		// update atoms
		newInd := 0
		for i := 0; i < len(atoms); i++{
			if atoms[i] != "H" {
				atomMap[atomIndices[i]] = newInd
				newAtomIndices = append(newAtomIndices, newInd)
				newAtoms = append(newAtoms, atoms[i])
				newInd++
			} else {
				HIndices = append(HIndices, atomIndices[i])
			}
		}

		// update bonds
		for i, b := range bonds{
			at1 := b[0]
			at2 := b[1]
			if !helpers.Contains(HIndices, at1) && !helpers.Contains(HIndices, at2){
				newBonds = append(newBonds, [2]int{atomMap[at1], atomMap[at2]})
				newBondTypes = append(newBondTypes, bondTypes[i])
			}
		}

		return newAtoms, newBonds, newBondTypes, newAtomIndices
}


