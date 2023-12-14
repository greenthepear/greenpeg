**greenpeg** is a simple CLI wrapper for FFmpeg made to make mass conversion easy on any platform. Using Go flags and presets it makes using FFmpeg less about googling commands and more about getting conversions done.

# Build
If you have [Go](https://go.dev/doc/install), [FFmpeg](https://www.ffmpeg.org/download.html) and git:

    git clone https://github.com/greenthepear/greenpeg
    cd greenpeg
    go build

# Running
Check flags with `greenpeg -h`.

For example, running:

    greenpeg -iname="/home/marcel/anaxi*.mp4" -oprefix="" -osuffix="_encoded" -vcodec="libx265"

will generate the following commands

    ffmpeg -i /home/marcel/anaxi1.mp4 -vcodec libx265 /home/marcel/anaxi1_encoded.mp4
    ffmpeg -i /home/marcel/anaxidev.mp4 -vcodec libx265 /home/marcel/anaxidev_encoded.mp4
    ffmpeg -i /home/marcel/anaxifyne.mp4 -vcodec libx265 /home/marcel/anaxifyne_encoded.mp4
    ...

and so on for any file matching the pattern from `iname`.

# Releases
There are no releases yet as it is in a very early stage. Handle with care, as using FFmpeg improperly in general can cause filling up storage space, especially when doing it on a mass scale.