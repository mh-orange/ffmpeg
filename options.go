package vtil

import (
	"io"
)

type TranscoderOption interface {
	process(*transcodeJob) error
}

type TranscoderOptionFunc func(*transcodeJob) error

func (tof TranscoderOptionFunc) process(job *transcodeJob) error { return tof(job) }

func StderrOption(writer io.WriteCloser) TranscoderOption {
	return TranscoderOptionFunc(func(job *transcodeJob) error {
		job.proc.Stderr(writer)
		return nil
	})
}

func VideoFilterOption(chaindef string) TranscoderOption {
	return TranscoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs("-lavfi", chaindef)
		return nil
	})
}

func DiscardOption() TranscoderOption {
	return TranscoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs("-f", "null", "-")
		return nil
	})
}
