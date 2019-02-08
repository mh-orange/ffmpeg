package vtil

import (
	"bytes"
	"reflect"
	"runtime"
	"testing"

	"github.com/mh-orange/cmd"
)

func TestOutputOptions(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	tests := []struct {
		option OutputOption
		want   output
	}{
		{OutputFilename("test.foo"), output{filename: "test.foo"}},
		{OutputWriter(writer), output{writer: writer}},
		{CopyAudioOption(), output{aCodec: "copy"}},
		{CopyOutput(), output{aCodec: "copy", vCodec: "copy"}},
		{OutputFormat("matroska"), output{format: "matroska"}},
		{DefaultH264(), output{vCodec: "libx264", vCodecOptions: []string{"-preset", "medium", "-tune", "film"}}},
		{DefaultMatroska(), output{format: "matroska", formatOptions: []string{"-map_chapters", "0"}}},
	}

	for _, test := range tests {
		name := runtime.FuncForPC(reflect.ValueOf(test.option).Pointer()).Name()
		t.Run(name, func(t *testing.T) {
			got := Output(test.option).output()
			got.process(&transcodeJob{proc: (&cmd.TestCmd{}).Process()})
			got.options = nil
			if !reflect.DeepEqual(test.want, *got) {
				t.Errorf("Want %v got %v", test.want, *got)
			}
		})
	}
}
