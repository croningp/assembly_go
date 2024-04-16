# Go Assembly Code

Go package to calculate Assembly Numbers for Molecules. 

# Build

Make sure `go` is installed, the instructions for that can be found on the (`golang docs`)[https://go.dev/doc/install]

Confirm you've installed `go` correctly by running `go version` (a version later than 1.15 should be shown)

after that building the command line tool is easy just run `go build cmd/app/main.go -o assembly`

# Usage

Currently basic usage is as follows - note flags (preceded by a dash) can be in any order

filename as argument at the end of the command, outputs assembly index only:

`./assembly my_mol.mol`

verbose flag will result in more detailed pathway output. In this usage, my_mol.mol must be at the end:

`./assembly -verbose my_mol.mol`

filename can also be passed as a flag - this takes priority over putting the filename at the end if both are done

`./assembly -file=my_mol.mol -verbose`

number of workers in the worker pool and the buffer size of the jobs queue. Currently defaults both to 100 (yet to test how optimal the default is).

`./assembly -file=my_mol.mol -workers=500 -buffer=500`

The -molfile flag defaults to true to read input as a mol file. Switch to false to read input as custom
graph txt file. This is a basic format that was used in testing and development as is simply 5 lines, 
those being name, list of vertex indices, list of associated edges (will be read in pairs, 
e.g. 1 2 2 3 is edges {1, 2} and {2, 3}), vertex colours, edge colours. Vertex and edge colours interpreted
as strings, delineated by spaces. If there are no vertex or edge colours, replace the corresponding
line with an exclamation mark "!"

```
Square graph (name - this line can be anything)
1 2 3 4 5
1 2 2 3 3 4 4 5
A B A B B
Red Blue Red Blue
```

`./assembly -file=my_graph.txt -molfile=false`

The `-log` flag is a boolean, and if present will log the pathway output to a file (default log.txt)

To specify a log file use e.g. `-logfile my_log_file.txt` (must also have log flag to do anything)

## Example
Here's an example with aspirin:

`./assembly file=aspirin.mol -verbose`

Which produces this output:

```
Running on file:  aspirin.mol
ORIGINAL GRAPH
+++++++++++++++
Vertices [0 1 2 3 4 5 6 7 8 9 10 11 12]
Edges [[0 1] [2 0] [0 3] [1 4] [5 2] [2 6] [3 10] [4 7] [7 5] [6 8] [6 9] [10 11] [10 12]]
VertexColours [C C C O C C C C O O C O C]
EdgeColours [double single single single double single single double single double single double single]
+++++++++++++++
PATHWAY
Pathway Graphs
======
Vertices [2 0 5]
Edges [[2 0] [5 2]]
VertexColours [C C C]
EdgeColours [single double]
======
======
Vertices [3 10 11 12]
Edges [[3 10] [10 11] [10 12]]
VertexColours [O C O C]
EdgeColours [single double single]
======
======
Vertices [13 1 14]
Edges [[13 1] [1 14]]
VertexColours [C C C]
EdgeColours [double single]
======
----------
Remnant Graph
Vertices [0 3 2 6 8 9 4 7 5]
Edges [[0 3] [2 6] [6 8] [6 9] [4 7] [7 5]]
VertexColours [C O C C O O C C C]
EdgeColours [single single double single double single]
----------
Duplicated Edges
[0 2]
[1 4 5]
[1 2]
+++++++++++++++

Assembly Index:  8
Time:  0.0449225
```

Sample build command from repo root

go build -o bin\assembly.exe GoAssembly/cmd/app
