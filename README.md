Alien

This implementation is done with directional maps. Each city is treated as a node with a maximum of 4 bi directional paths in/out, we track the paths traversed to facilitate in destroying cities as well. 

Dependencies:
Make sure Go version 1.9+ is installed

Build:

make build

Run:

$ ./alien --map=<INPUT-MAP> --out=<OUTPUT-FILE-NAME> --n=<ALIENS-ON-MAP>

A mock test file has been included 
"testmap.txt"

main.go
Houses the creation and writing of a map as well as generating a simulation (Simulation/sim.go)