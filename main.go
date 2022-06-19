package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"alien/alienworld"
	"alien/simulation"
)

func main() {
	var (
		mapgenfile string
		outpf      string
		aliens     uint
	)

	flag.StringVar(&mapgenfile, "map", "", "map definition file")
	
	flag.StringVar(&outpf, "out", "", "output file")
	
	flag.UintVar(&aliens, "n", 0, "number of aliens")

	flag.Parse()

	if aliens == 0 {
	
		cmdErrorMsg("number of aliens should beatleast 1")
	
	} else if len(mapgenfile) == 0 {
	
		cmdErrorMsg("unspecified map file")
	
	} else if len(outpf) == 0 {
	
		cmdErrorMsg("unspecified output file")
	} 

	alienmap, err := buildalienmap(mapgenfile)
	
	if err != nil {
	
		log.Fatalf("failed to build map from: %v", err)
	
	}

	// Seed map iwth aliens
	alienmap.AddAliens(aliens)
	
	alienmap.ExecuteFights()

	sim := simulation.Newsim(alienmap)

	if err := sim.Run(); err != nil {
	
		log.Fatalf("failed to execute simulation: %v", err)
	
	}

	log.Println("sim complete")

	if err := writemap(alienmap, outpf); err != nil {
	
		log.Fatalf("failed to write to file: %v", err)
	
	}
}

func cmdErrorMsg(errMsg string) {
	
	fmt.Printf("%s\n\n", errMsg)
	
	fmt.Println("usage:")
	
	flag.PrintDefaults()
	
	os.Exit(1)
}

// Build the map
func buildalienmap(mapgenfile string) (*alienworld.Map, error) {
	
	file, err := os.Open(mapgenfile)
	
	if err != nil {
	
		return nil, err
	
	}
	
	defer file.Close()

	alienmap := alienworld.NewMap()

	// read input line by line
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
	
		line := scanner.Text()
	
		tokens := strings.Split(line, " ")

		if len(tokens) == 0 {

			return nil, errors.New("invalid line in map definition")

		}

		cityName := tokens[0]

		if len(tokens) > 1 {

			for _, link := range tokens[1:] {

				linkTokens := strings.Split(link, "=")

				if len(linkTokens) != 2 {

					return nil, errors.New("invalid line in map definition")

				}

				alienmap.Addpaths(cityName, linkTokens[0], linkTokens[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {

		return nil, err

	}

	return alienmap, nil
}

// writemap writes a given alienworld map to the file at path 'outPath'.
func writemap(alienmap *alienworld.Map, outPath string) error {

	fileHandle, err := os.Create(outPath)

	if err != nil {

		return err

	}

	defer fileHandle.Close()

	writer := bufio.NewWriter(fileHandle)

	defer fileHandle.Close()

	for _, city := range alienmap.Cities() {
		
		s := city.String()

		if len(s) != 0 {
		
			fmt.Fprintln(writer, s)
		
			writer.Flush()
		
		}
	}

	return nil
}
