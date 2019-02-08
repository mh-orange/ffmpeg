package ffmpeg

type TestTranscoder struct {
	TranscodeErr error
	JobErr       error
	Log          string
}

type TestJob struct {
	Canceled bool
	log      string
	err      error
}

func (tj *TestJob) Cancel() {
	tj.Canceled = true
}
func (tj *TestJob) Err() error {
	return tj.err
}

func (tj *TestJob) Log() string {
	return tj.log
}

func (tj *TestJob) Progress() <-chan TranscodeInfo {
	ch := make(chan TranscodeInfo, 1)
	ch <- TranscodeInfo{}
	close(ch)
	return ch
}

func (tj *TestJob) Wait() error {
	return tj.err
}

func (tt *TestTranscoder) Transcode(options ...TranscoderOption) (TranscodeJob, error) {
	return &TestJob{log: tt.Log, err: tt.JobErr}, tt.TranscodeErr
}
