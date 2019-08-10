package ffmpeg

type CopyTranscoder struct {
	transcoder *Transcoder
}

func NewCopyTranscoder() *CopyTranscoder {
	return &CopyTranscoder{NewTranscoder()}
}

func (ct *CopyTranscoder) Transcode(input TranscoderInput, output TranscoderOutput) (TranscodeJob, error) {
	return ct.transcoder.Transcode(input, output)
}

func Copy(input TranscoderInput, output TranscoderOutput) (TranscodeJob, error) {
	transcoder := NewCopyTranscoder()
	output.output().options = append(output.output().options, CopyOutput())
	return transcoder.Transcode(input, output)
}
