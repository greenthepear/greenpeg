package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const displayCommandsLimit = 5

type ffmpegOp struct {
	iname        string
	oprefix      string
	osuffix      string
	from, to     string
	vcodec       string
	crf          int
	preset       string
	custom       string
	customBefore string
	ext          string
}

func globFilenames(pattern string) ([]string, error) {
	r, err := filepath.Glob(pattern)
	if len(r) == 0 {
		return nil, fmt.Errorf("no files matching pattern")
	}
	return r, err
}

func main() {
	//1. check flags
	//2. give info of files
	//3. ask to start (unless -y)
	//4. run
	//ffmpeg -i input.mp4 -vcodec libx265 -crf 28 output.mp4

	op := ffmpegOp{}

	getOptionsFromFlags(&op)

	flag.Parse()
	if op.iname == "" {
		log.Fatalf("No input name (`-iname`).")
	}
	filenames, err := globFilenames(op.iname)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d files from iname: %v\nGenerating command strings...\n", len(filenames), filenames)

	commandStrings := createCommandStrings(filenames, op, false)

	fmt.Printf("\nThe following %d commands will be executed:\n", len(commandStrings))
	for i, c := range commandStrings {
		if i > displayCommandsLimit {
			fmt.Printf("... and %d more.\n", len(commandStrings)-displayCommandsLimit)
			break
		}
		fmt.Printf("%s\n", c)
	}

	confirm := ""
	for {
		fmt.Printf("Do you wish to execute the commands? [y/n]\n")
		_, err = fmt.Scan(&confirm)
		if err != nil {
			fmt.Printf("Input error: %v\nt", err)
			continue
		}
		if confirm == "N" || confirm == "n" {
			fmt.Printf("Aborting...\n")
			os.Exit(0)
		}
		if confirm == "Y" || confirm == "y" {
			break
		}
		fmt.Printf("Invalid input.\n")
	}
}
