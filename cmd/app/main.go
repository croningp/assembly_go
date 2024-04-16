package main

import (
	"GoAssembly/pkg/assembly"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type CommandLineOptions struct {
	inputFile *string
	molFile *bool
	logFile *string
	numWorkers *int
	bufferSize *int
	variant *string
	debug *bool
	verbose *bool
	log *bool
	pathway *bool
	tail []string
	}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// main executable will output assembly index and pathway to stdout and log file if selected in command line arguments
func main() {

	// command line arguments
	inputFile := flag.String("file", "", "the name of the input file")
	molFile := flag.Bool("molfile", true, "true if molfile, false if general graph file")
	logFile := flag.String("logfile", "log.txt", "the path to the log file")
	numWorkers := flag.Int("workers", 100, "the number of workers in the worker pool")
	bufferSize := flag.Int("buffer", 100, "the buffer size of the jobs queue")
	variant := flag.String("variant", "shortest", "the variant of the algorithm - currently only shortest implemented")
	debug := flag.Bool("debug", false, "additional logging - not currently used")
	verbose := flag.Bool("verbose", false, "stdout pathway information - if false, only assembly index output")
	log := flag.Bool("log", false, "log to file")
	pathway := flag.Bool("pathway", false, "the input file contains multiple graphs in the form of a starting pathway, e.g. an sdf file")

	flag.Parse()
	CLArgs := CommandLineOptions{
		inputFile,
		molFile,
		logFile,
		numWorkers,
		bufferSize,
		variant,
		debug,
		verbose,
		log,
		pathway,
		flag.Args(),
	}

	var logf *os.File
	var err error

	// set up log file
	if *CLArgs.log{
		logf, err = os.Create(*CLArgs.logFile)
		check(err)
		assembly.Logger.SetOutput(logf)
	}

	// get input  file
	var inFile string
	if *CLArgs.inputFile == "" {
		inFile = CLArgs.tail[0]
	} else {
		inFile = *CLArgs.inputFile
	}

	// Generate slice of Graphs. This will just contain the graph of the initial structure, unless a starting pathway is provided, in which
	// case it will contain the graphs in the pathway
	var fileGraph []assembly.Graph
	if *CLArgs.pathway{
		fileGraph = assembly.ParseSDFile(inFile, true)
	} else {
		if *CLArgs.molFile {
			fileGraph = append(fileGraph, assembly.MolColourGraph(inFile))
		} else {
			fileGraph = append(fileGraph, assembly.NewGraphOnlyFromFile(inFile))
		}
	}

	var pathways []assembly.Pathway
	start := time.Now()

	// Generate the output pathways. At present, the only variant implemented will return a single shortest pathway
	if *CLArgs.pathway{
		originalGraph, starterPathway := assembly.MolListToPathway(fileGraph, []assembly.Duplicates{})
		pathways = assembly.AssemblyPathway(originalGraph, starterPathway, *CLArgs.numWorkers, *CLArgs.bufferSize,*CLArgs.variant)
	} else {
		pathways = assembly.Assembly(fileGraph[0], *CLArgs.numWorkers, *CLArgs.bufferSize, *CLArgs.variant)
	}


	elapsed := time.Now().Sub(start)

	// calculate the assembly index from the pathways, and a string containing pathway details
	assemblyIndex := assembly.AssemblyIndex(&pathways[0], &fileGraph[0])
	assemblyString := assembly.AssemblyString(pathways, &fileGraph[0])

	// output assembly index and details to stdout
	if *CLArgs.verbose {
		fmt.Println("Running on file: ", inFile)

		fmt.Println(assemblyString)
		fmt.Println("Assembly Index: ", assemblyIndex)
		fmt.Println("Time: ", elapsed.Seconds())
	} else {
		fmt.Println(assemblyIndex)
	}

	// output assembly index and details to log file (if specified in command line arguments)
	if *CLArgs.log{
		assembly.Logger.Debug("Running on file: ", inFile)
		assembly.Logger.Debug(assemblyString)
		assembly.Logger.Debug("Assembly Index: ", assemblyIndex)
		assembly.Logger.Debug("Time: ", elapsed.Seconds())
	}
}
