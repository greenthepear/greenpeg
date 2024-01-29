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
	iname    string
	oprefix  string
	osuffix  string
	from, to string
	vcodec   string
	crf      int
	preset   string
	ext      string
}

func globFilenames(pattern string) ([]string, error) {
	r, err := filepath.Glob(pattern)
	if len(r) == 0 {
		return nil, fmt.Errorf("no files matching pattern")
	}
	return r, err
}

func main() {
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

	fmt.Printf("Found %d files from iname: %v\n\nGenerating commands...\n", len(filenames), filenames)

	commands, err := createCommands(filenames, op, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The following %d commands will be executed:\n", len(commands))
	for i, c := range commands {
		if i >= displayCommandsLimit {
			fmt.Printf("... and %d more.\n", len(commands)-displayCommandsLimit)
			break
		}
		fmt.Printf("%s\n", c.String())
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

	for i, c := range commands {
		fmt.Printf("████ [%d/%d] ████\nRunning `%s`...\n", i+1, len(commands), c.String())
		out, err := c.CombinedOutput()
		if err != nil {
			log.Fatalf("Error while running command: %s\n%s", err, out)
		}
		fmt.Printf("▒▒▒▒ Finished successfully. Output:\n%s\n", out)
	}
}
