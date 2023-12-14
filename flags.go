package main

import "flag"

func getOptionsFromFlags(op *ffmpegOp) {
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
	flag.StringVar(&op.ext, "ext", "", "file extension, aka container format, leave to keep the same as input file")
}
