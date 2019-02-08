package ffmpeg

import (
	"io/ioutil"
	"net/url"
	"reflect"
	"runtime"
	"testing"

	"github.com/mh-orange/cmd"
)

func TestInput(t *testing.T) {
	tests := []struct {
		option  InputOption
		fi      *FileInfo
		want    []string
		wantErr bool
	}{
		{StartOption(1 * Hour), nil, []string{"-ss", "01:00:00.000000"}, false},
		{StartPercentOption(10), nil, nil, true},
		{StartPercentOption(10), &FileInfo{Format: FormatInfo{Filename: "foo.mkv", Duration: 1 * Hour}}, []string{"-ss", "00:06:00.000000", "-i", "foo.mkv"}, false},
		{DurationOption(42 * Minute), nil, []string{"-t", "00:42:00.000000"}, false},
		{InputURL(&url.URL{"http", "", nil, "video.net", "/foo", "", false, "", ""}), nil, []string{"-i", "http://video.net/foo"}, false},
	}

	for _, test := range tests {
		name := runtime.FuncForPC(reflect.ValueOf(test.option).Pointer()).Name()
		t.Run(name, func(t *testing.T) {
			in := Input(test.option).input()
			if test.fi != nil {
				in.fi = test.fi
			}
			err := in.process(&transcodeJob{proc: (&cmd.TestCmd{}).Process()})
			if err == nil {
				if test.wantErr {
					t.Errorf("Expected error, got nil")
				} else {
					got := in.args
					if !reflect.DeepEqual(test.want, got) {
						t.Errorf("Want %v got %v", test.want, got)
					}
				}
			} else if !test.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestInputFilename(t *testing.T) {
	in := &input{}
	inputTxt, err := ioutil.ReadFile("testdata/info1.json")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	c := &cmd.TestCmd{Stdout: inputTxt}
	ffprobe = c

	err = InputFilename("test.mkv")(in)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if in.fi == nil {
		t.Errorf("Expected InputFilename to set fi")
	}
}
