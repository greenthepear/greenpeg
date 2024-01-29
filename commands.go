package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

func createCommandString(filename string, options ffmpegOp) (string, error) {
	beforei := " "
	needsCopy := false

	//Before -i
	if options.from != "" {
		if !regexp.MustCompile(`\d\d:\d\d:\d\d`).MatchString(options.from) {
			return "", fmt.Errorf("bad time format (not hh:mm:ss) for -from: %s", options.from)
		}
		needsCopy = true
		//beforei = fmt.Sprintf("%s-ss %s ", beforei, options.from)
		beforei = beforei + "-ss " + options.from
	}

	if options.to != "" {
		if !regexp.MustCompile(`\d\d:\d\d:\d\d`).MatchString(options.to) {
			return "", fmt.Errorf("bad time format (not hh:mm:ss) for -to: %s ", options.to)
		}
		needsCopy = true
		//beforei = fmt.Sprintf("%s-to %s ", beforei, options.from)
		beforei = beforei + "-to " + options.to
	}

	if options.customBefore != "" {
		beforei += " " + options.customBefore
	}

	if beforei != " " {
		beforei = beforei + " "
	}

	afteri := ""

	//After -i
	if options.vcodec != "" {
		//afteri = fmt.Sprintf("%s-vcodec %s ", afteri, options.vcodec)
		afteri = afteri + " -vcodec " + options.vcodec
	}

	if options.crf != -1 {
		if options.crf < 0 || options.crf > 51 {
			return "", fmt.Errorf("bad crf option, needs to be between 0 and 51, not: %d", options.crf)
		}
		afteri = fmt.Sprintf("%s-crf %d ", afteri, options.crf)
		afteri = afteri + " -crf " + fmt.Sprint(options.crf)
	}

	if options.custom != "" {
		afteri = afteri + " " + options.custom
	}

	if needsCopy {
		//beforei = fmt.Sprintf("%s-c copy ", beforei)
		afteri = afteri + " -c copy"
	}

	//output name
	ofile := filepath.Base(filename)
	ofileExt := filepath.Ext(ofile)
	ofile, _ = strings.CutSuffix(ofile, ofileExt)
	ofile = options.oprefix + ofile + options.osuffix + ofileExt
	odir := filepath.Dir(filename)
	opath := filepath.Join(odir, ofile)
	if opath == filename {
		return "", fmt.Errorf("output file name same as input, do you have empty -oprefix and -osuffinx?")
	}

	return fmt.Sprintf("ffmpeg%s-i %s%s %s", beforei, filename, afteri, opath), nil
}

func createCommandStrings(filenames []string, op ffmpegOp, debugInfo bool) []string {
	commandStrings := make([]string, 0)
	for i, fname := range filenames {
		commandstring, err := createCommandString(fname, op)
		if err != nil {
			log.Fatalf("Error creating command for `%s`: %v", fname, err)
		}
		if debugInfo {
			fmt.Printf("\t%d: for `%s`:\n$ %s\n", i, fname, commandstring)
		}
		commandStrings = append(commandStrings, commandstring)
	}
	return commandStrings
}
