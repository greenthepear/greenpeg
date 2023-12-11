package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

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
}

func globFilenames(pattern string) ([]string, error) {
	r, err := filepath.Glob(pattern)
	if len(r) == 0 {
		return nil, errors.New("no files matching pattern")
	}
	return r, err
}

func createCommandString(filename string, options ffmpegOp) (string, error) {
	beforei := ""
	needsCopy := false

	//Before -i
	if options.from != "" {
		if !regexp.MustCompile(`\d\d:\d\d:\d\d`).MatchString(options.from) {
			return "", fmt.Errorf("bad time format (not hh:mm:ss) for -from: %s", options.from)
		}
		needsCopy = true
		beforei = fmt.Sprintf("%s-ss %s ", beforei, options.from)
	}

	if options.to != "" {
		if !regexp.MustCompile(`\d\d:\d\d:\d\d`).MatchString(options.from) {
			return "", fmt.Errorf("bad time format (not hh:mm:ss) for -to: %s ", options.from)
		}
		needsCopy = true
		beforei = fmt.Sprintf("%s-to %s ", beforei, options.from)
	}

	if options.customBefore != "" {
		beforei = fmt.Sprintf("%s %s ", beforei, options.customBefore)
	}

	afteri := ""

	//After -i
	if options.vcodec != "" {
		afteri = fmt.Sprintf("%s-vcodec %s ", afteri, options.vcodec)
	}

	if options.crf != -1 {
		if options.crf < 0 || options.crf > 51 {
			return "", fmt.Errorf("bad crf option, needs to be between 0 and 51, not: %d", options.crf)
		}
		afteri = fmt.Sprintf("%s-crf %d ", afteri, options.crf)
	}
	if options.custom != "" {
		beforei = fmt.Sprintf("%s%s ", afteri, options.custom)
	}

	if needsCopy {
		beforei = fmt.Sprintf("%s-c copy ", beforei)
	}

	//output name
	ofile := filepath.Base(filename)
	ofileExt := filepath.Ext(ofile)
	ofile, _ = strings.CutSuffix(ofile, ofileExt)
	ofile = fmt.Sprintf("%s%s%s%s", options.oprefix, ofile, options.osuffix, ofileExt)
	odir := filepath.Dir(filename)
	opath := filepath.Join(odir, ofile)
	if opath == filename {
		return "", fmt.Errorf("output file name same as input, do you have empty -oprefix and -osuffinx?")
	}

	return fmt.Sprintf("ffmpeg %s-i %s %s%s", beforei, filename, afteri, opath), nil
}

func main() {
	//1. check flags
	//2. give info of files
	//3. ask to start (unless -y)
	//4. run
	//ffmpeg -i input.mp4 -vcodec libx265 -crf 28 output.mp4
	op := ffmpegOp{}

	flag.StringVar(&op.iname, "iname", "", "pattern of input files (e.g. './2023*.mp4')")
	flag.StringVar(&op.oprefix, "oprefix", "ENC_", "prefix of the name of converted files")
	flag.StringVar(&op.osuffix, "osuffix", "", "suffix of the name of converted files")
	flag.StringVar(&op.from, "from", "", "like -ss, use hh:mm:ss formatting, enables -c copy")
	flag.StringVar(&op.to, "to", "", "like -to, use hh:mm:ss formatting, enables -c copy")
	flag.StringVar(&op.vcodec, "vcodec", "", "same as ffmpeg, libx265 recommended")
	flag.IntVar(&op.crf, "crf", -1, "same as ffmpeg")
	flag.StringVar(&op.preset, "preset", "", "run custom preset, see presets with `greenpeg presets`, any overlapping flag will overwrite the preset")
	flag.StringVar(&op.custom, "custom", "", "custom options that will be appended AFTER -i")
	flag.StringVar(&op.customBefore, "custom_b", "", "custom options that will be added BEFORE -i")

	flag.Parse()
	filenames, err := globFilenames(op.iname)
	if err != nil {
		log.Fatal(err)
	}

	for i, fname := range filenames {
		commandstring, err := createCommandString(fname, op)
		if err != nil {
			log.Fatalf("error creating command for `%s`: %v", fname, err)
		}
		fmt.Printf("\t%d: for `%s`:\n$ %s\n", i, fname, commandstring)
	}
}
