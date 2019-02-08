package ffmpeg

import (
	"io"
	"testing"
)

func TestTestJobCancel(t *testing.T) {
	tj := &TestJob{}
	if tj.Canceled {
		t.Errorf("Want not canceled")
	}
	tj.Cancel()
	if !tj.Canceled {
		t.Errorf("Want canceled")
	}
}

func TestTestJobErr(t *testing.T) {
	tj := &TestJob{err: io.EOF}
	if tj.Err() != io.EOF {
		t.Errorf("Want %v got %v", io.EOF, tj.Err())
	}
}

func TestTestJobLog(t *testing.T) {
	want := "So Long, and Thanks for All the Fish"
	tj := &TestJob{log: want}
	if tj.Log() != want {
		t.Errorf("Want %v got %v", want, tj.Log())
	}
}

func TestTestJobProgress(t *testing.T) {
	tj := &TestJob{}
	ch := tj.Progress()
	select {
	case <-ch:
		_, open := <-ch
		if open {
			t.Errorf("Wanted closed channel")
		}
	default:
		t.Errorf("Want progress")
	}
}

func TestTestJobWait(t *testing.T) {
	want := io.EOF
	tj := &TestJob{err: want}
	got := tj.Wait()
	if want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestTranscode(t *testing.T) {
	tt := &TestTranscoder{}
	job, err := tt.Transcode()
	if job == nil {
		t.Errorf("Wanted non-nil job")
	} else if _, ok := job.(*TestJob); !ok {
		t.Errorf("Wanted TestJob bot %T", job)
	}

	if err == nil {
		tt.TranscodeErr = io.EOF
		_, err = tt.Transcode()
		if err != io.EOF {
			t.Errorf("Wanted %v got %v", io.EOF, err)
		}
	} else {
		t.Errorf("Unexpected error: %v", err)
	}
}
