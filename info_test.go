package ffmpeg

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mh-orange/cmd"
)

func TestInfoUnmarshalJSON(t *testing.T) {
	oldFfprobe := Ffprobe

	names, err := filepath.Glob("testdata/info*.json")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			input, err := ioutil.ReadFile(name)
			if err == nil {
				c := &cmd.TestCmd{}
				c.Stdout = input
				Ffprobe = c
				fi, err := Stat(name)
				if strings.HasSuffix(name, "_video.json") {
					if !fi.IsVideo() {
						t.Errorf("Expected the information to indicate that the media is a video")
					}
				} else if strings.HasSuffix(name, "_err.json") {
					if err == nil {
						t.Errorf("Expected parse error")
					}
					err = nil
				} else if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					if len(fi.VideoStreams) > 0 && !fi.IsVideo() {
						t.Errorf("want true got false")
					}
				}
			} else {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
	Ffprobe = oldFfprobe
}
