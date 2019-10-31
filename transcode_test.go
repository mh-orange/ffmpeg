package ffmpeg

import (
	"io"
	"testing"

	"github.com/mh-orange/cmd"
)

func TestTranscoderTranscode(t *testing.T) {
	oldFfmpeg := Ffmpeg

	tests := []struct {
		name     string
		startErr error
	}{
		{"no error", nil},
		{"io error", io.EOF},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			optionCalled := false
			option := transcoderOptionFunc(func(job *transcodeJob) error {
				optionCalled = true
				return nil
			})

			tc := &cmd.TestCmd{}
			tc.StartErr = test.startErr
			Ffmpeg = tc
			job, err := NewTranscoder().Transcode(option)
			job.Cancel()
			if err != test.startErr {
				t.Errorf("wanted %v got %v", test.startErr, err)
			}

			if !optionCalled {
				t.Errorf("expected option to be called by transcode")
			}
		})
	}

	Ffmpeg = oldFfmpeg
}

func TestTranscoderRun(t *testing.T) {

}
