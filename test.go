package ffmpeg

// TestJob is a dummy transcode job returned by TestTranscoder.Transcode.  The
// TestJob performs just like a normal transcode job, but will return error and log
// based on those values set in the TestTranscoder
type TestJob struct {
	// Canceled is set by the Cancel() method and is useful to test if a job
	// was successfully canceled
	Canceled bool
	log      string
	err      error
}

func (tj *TestJob) Inspect() string {
	return "test job"
}

// Cancel will set the Canceled property true
func (tj *TestJob) Cancel() {
	tj.Canceled = true
}

// Err returns the JobErr set in the TestTranscoder
func (tj *TestJob) Err() error {
	return tj.err
}

// Log returns the Log value set in TestTranscoder
func (tj *TestJob) Log() string {
	return tj.log
}

// Progress returns a channel that will be populated with exactly one
// TranscodeInfo object and then immediately closed.
func (tj *TestJob) Progress() <-chan TranscodeInfo {
	ch := make(chan TranscodeInfo, 1)
	ch <- TranscodeInfo{}
	close(ch)
	return ch
}

// Wait returns the JobErr set in the TestTranscoder
func (tj *TestJob) Wait() error {
	return tj.err
}

// TestTranscoder will not actually call the ffmpeg command line tool
// and will simply return a TestJob when Transcode is called.  This
// is useful when writing unit tests for consumers of the ffmpeg
// package
type TestTranscoder struct {
	// TranscodeErr will be returned immediately by the Transcode function. If
	// this is nil, no error is returned and a TestJob is returned instead
	TranscodeErr error

	// JobErr is the error returned by TestJob.Wait and TestJob.Err
	JobErr error

	// Log is the string returned by TestJob.Log
	Log string
}

// Transcode will return a new TestJob with the log and err values set to the corresponding
// TestTranscoder values, as well as any error set on TranscodeErr
func (tt *TestTranscoder) Transcode(options ...TranscoderOption) (TranscodeJob, error) {
	return &TestJob{log: tt.Log, err: tt.JobErr}, tt.TranscodeErr
}
