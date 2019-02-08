package ffmpeg

import (
	"io"
)

type TranscoderOutput interface {
	TranscoderOption
	output() *output
}

type OutputOption func(*output) error

type output struct {
	filename string
	writer   io.Writer

	aCodec        string
	aCodecOptions []string

	vCodec        string
	vCodecOptions []string

	format        string
	formatOptions []string

	options []OutputOption
}

func (out *output) output() *output {
	return out
}

func (out *output) process(job *transcodeJob) error {
	for _, option := range out.options {
		option(out)
	}

	if out.vCodec != "" {
		job.proc.AppendArgs("-c:v", out.vCodec)
		job.proc.AppendArgs(out.vCodecOptions...)
	}

	if out.aCodec != "" {
		job.proc.AppendArgs("-c:a", out.aCodec)
		job.proc.AppendArgs(out.aCodecOptions...)
	}

	if out.format != "" {
		job.proc.AppendArgs("-f", out.format)
		job.proc.AppendArgs(out.formatOptions...)
	}

	if out.filename != "" {
		job.proc.AppendArgs("-y", out.filename)
	} else if out.writer != nil {
		job.proc.AppendArgs("-")
		job.proc.Stdout(out.writer)
	}
	return nil
}

func Output(options ...OutputOption) TranscoderOutput {
	return &output{options: options}
}

func DefaultH264() OutputOption {
	return func(output *output) error {
		output.vCodec = "libx264"
		output.vCodecOptions = []string{"-preset", "medium", "-tune", "film"}
		return nil
	}
}

func DefaultMatroska() OutputOption {
	return func(output *output) error {
		output.format = "matroska"
		output.formatOptions = []string{"-map_chapters", "0"}
		return nil
	}
}

func OutputFilename(filename string) OutputOption {
	return func(output *output) error {
		output.filename = filename
		return nil
	}
}

func OutputWriter(writer io.Writer) OutputOption {
	return func(output *output) error {
		output.writer = writer
		return nil
	}
}

func CopyAudioOption() OutputOption {
	return func(output *output) error {
		output.aCodec = "copy"
		return nil
	}
}

func CopyOutput() OutputOption {
	return func(output *output) error {
		output.aCodec = "copy"
		output.vCodec = "copy"
		return nil
	}
}

func OutputFormat(format string) OutputOption {
	return func(output *output) error {
		output.format = format
		return nil
	}
}
