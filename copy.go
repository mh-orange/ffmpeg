package vtil

type CopyTranscoder struct {
	transcoder *Transcoder
}

func NewCopyTranscoder() *CopyTranscoder {
	return &CopyTranscoder{NewTranscoder()}
}

func (ct *CopyTranscoder) Transcode(input TranscoderInput, output TranscoderOutput) (TranscodeJob, error) {
	return ct.transcoder.Transcode(input, output)
}

func Copy(input TranscoderInput, output TranscoderOutput) error {
	transcoder := NewCopyTranscoder()
	job, err := transcoder.Transcode(input, output)
	if err == nil {
		err = job.Wait()
	}
	return err
}
