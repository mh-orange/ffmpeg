package vtil

type CheckTranscoder struct {
}

func NewCheckTranscoder() *CheckTranscoder {
	return &CheckTranscoder{}
}

func (ct *CheckTranscoder) Check(input TranscoderInput) (TranscodeJob, error) {
	transcoder := NewTranscoder()
	options := append([]TranscoderOption{input}, DiscardOption())
	return transcoder.Transcode(options...)
}

func Check(input TranscoderInput) (string, error) {
	transcoder := NewCheckTranscoder()
	job, err := transcoder.Check(input)
	if err == nil {
		err = job.Wait()
	}
	return job.Log(), err
}
