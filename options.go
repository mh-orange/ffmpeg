package ffmpeg

import (
	"fmt"
	"io"
)

// TranscoderOption is passed to the Transcode function to set things like
// input, output and filters
type TranscoderOption interface {
	process(*transcodeJob) error
}

type transcoderOptionFunc func(*transcodeJob) error

func (tof transcoderOptionFunc) process(job *transcodeJob) error { return tof(job) }

// StderrOption will tee the stderr output from the underlying ffmpeg process and
// write the output to the given writer
func StderrOption(writer io.WriteCloser) TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		job.proc.Stderr(writer)
		return nil
	})
}

func LogOption(writer io.Writer) TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		pr, pw := io.Pipe()
		job.proc.Stderr(pw)
		reader := newFilterReader(pr, progPtrn, repeatPtrn)
		go func() {
			for {
				if reader.Scan() {
					if reader.Pattern() == nil {
						writer.Write(reader.Bytes())
					}
				} else {
					if reader.Err() != nil {
						break
					}
				}
			}
		}()

		return nil
	})
}

// VideoFilterOption sets a video filter chain on a transcoder
func VideoFilterOption(chaindef string) TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs("-lavfi", chaindef)
		return nil
	})
}

// DiscardOption sets the transcoder output format to null and discards the output
func DiscardOption() TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs("-f", "null", "-")
		return nil
	})
}

func MapOption(index int) TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs("-map", fmt.Sprintf("%d", index))
		return nil
	})
}

func MapMetadataOption(index int) TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs("-map_metadata", fmt.Sprintf("%d", index))
		return nil
	})
}

func DispositionOption(index int, disposition string) TranscoderOption {
	return transcoderOptionFunc(func(job *transcodeJob) error {
		job.proc.AppendArgs(fmt.Sprintf("-disposition:%d", index), disposition)
		return nil
	})
}
