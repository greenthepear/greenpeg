package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func createCommandFlags(filename string, options ffmpegOp) ([]string, error) {
	needsCopy := false
	flags := make([]string, 0)

	//Before -i
	if options.from != "" {
		if !regexp.MustCompile(`\d\d:\d\d:\d\d`).MatchString(options.from) {
			return []string{},
				fmt.Errorf("bad time format (not hh:mm:ss) for -from: %s", options.from)
		}
		needsCopy = true
		flags = append(flags, "-ss", options.from)
	}

	if options.to != "" {
		if !regexp.MustCompile(`\d\d:\d\d:\d\d`).MatchString(options.to) {
			return []string{},
				fmt.Errorf("bad time format (not hh:mm:ss) for -to: %s ", options.to)
		}
		needsCopy = true
		flags = append(flags, "-to", options.to)
	}

	flags = append(flags, "-i", filename)

	//After -i
	if options.vcodec != "" {
		flags = append(flags, "-vcodec", options.vcodec)
	}

	if options.crf != -1 {
		if options.crf < 0 || options.crf > 51 {
			return []string{},
				fmt.Errorf("bad crf option, needs to be between 0 and 51, not: %d", options.crf)
		}
		flags = append(flags, "-crf", fmt.Sprint(options.crf))
	}

	if needsCopy {
		flags = append(flags, "-c", "copy")
	}

	ofile := filepath.Base(filename)
	ofileExt := filepath.Ext(ofile)
	ofile, _ = strings.CutSuffix(ofile, ofileExt)
	ofile = options.oprefix + ofile + options.osuffix + ofileExt
	odir := filepath.Dir(filename)
	opath := filepath.Join(odir, ofile)
	if opath == filename {
		return []string{}, fmt.Errorf("output file name same as input, do you have empty -oprefix and -osuffinx?")
	}
	flags = append(flags, opath, "-hide_banner", "-loglevel", "error")
	return flags, nil
}

func createCommands(filenames []string, op ffmpegOp, debugInfo bool) ([]*exec.Cmd, error) {
	commands := make([]*exec.Cmd, 0)
	for i, fname := range filenames {
		commandflags, err := createCommandFlags(fname, op)
		if err != nil {
			return nil, fmt.Errorf("error creating command flags for `%s`: %v", fname, err)
		}
		newCommand := exec.Command("ffmpeg", commandflags...)
		if debugInfo {
			fmt.Printf("\t%d: for `%s`:\n$ %s\n", i, fname, newCommand.String())
		}
		commands = append(commands, newCommand)
	}
	return commands, nil
}
