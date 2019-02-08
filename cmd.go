package ffmpeg

import (
	"log"
	"os/exec"

	"github.com/mh-orange/cmd"
)

var (
	ffmpeg  = cmd.New("ffmpeg", "-hide_banner", "-nostdin", "-nostats", "-progress", "/dev/stderr")
	ffprobe = cmd.New("ffprobe", "-hide_banner", "-v", "error", "-print_format", "json", "-sexagesimal", "-show_format", "-show_streams", "-show_chapters", "-show_programs")
)

func init() {
	for _, c := range []cmd.Command{ffmpeg, ffprobe} {
		path, err := exec.LookPath(c.Path())
		if err != nil {
			log.Printf("%v no such file or directory", c.Path())
		}
		c.SetPath(path)
	}
}
