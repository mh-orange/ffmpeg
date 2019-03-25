package ffmpeg

import (
	"io"
)

// TranscoderOutput is an Option that can be used for the output of the transcoder
type TranscoderOutput interface {
	TranscoderOption
	output() *output
}

// OutputOption is an option that is applied to a TranscoderOutput these are useful
// for setting things like output format, output codec, etc
type OutputOption func(*output) error

type output struct {
	filename string
	writer   io.Writer

	aCodec        string
	aCodecOptions []string

	vCodec        string
	vCodecOptions []string

	sCodec string

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

	if out.sCodec != "" {
		job.proc.AppendArgs("-c:s", out.sCodec)
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

// Output returns a TranscoderOutput with the given output options
func Output(options ...OutputOption) TranscoderOutput {
	return &output{options: options}
}

// DefaultH264 sets the Output video codec to libx264 using the medium preset and film tuning
func DefaultH264() OutputOption {
	return func(output *output) error {
		output.vCodec = "libx264"
		output.vCodecOptions = []string{"-preset", "medium", "-tune", "film"}
		return nil
	}
}

// DefaultMatroska sets the output to use the matroska format
func DefaultMatroska() OutputOption {
	return func(output *output) error {
		output.format = "matroska"
		output.formatOptions = []string{"-map_chapters", "0"}
		return nil
	}
}

// OutputFilename sets the output to write to a file named by the filename string
func OutputFilename(filename string) OutputOption {
	return func(output *output) error {
		output.filename = filename
		return nil
	}
}

// OutputWriter will send the transcoder output to the given io.Writer
func OutputWriter(writer io.Writer) OutputOption {
	return func(output *output) error {
		output.writer = writer
		return nil
	}
}

// CopyAudioOption will set the audio codec to "copy"
func CopyAudioOption() OutputOption {
	return func(output *output) error {
		output.aCodec = "copy"
		return nil
	}
}

// CopySubtitlesOption will set the subtitle codec to copy
func CopySubtitlesOption() OutputOption {
	return func(output *output) error {
		output.sCodec = "copy"
		return nil
	}
}

// CopyOutput sets both the audio and video codecs to copy
func CopyOutput() OutputOption {
	return func(output *output) error {
		output.aCodec = "copy"
		output.vCodec = "copy"
		return nil
	}
}

// OutputFormat sets the output format to the format string.  No checking
// is done to make sure the format string is valid
func OutputFormat(format string) OutputOption {
	return func(output *output) error {
		output.format = format
		return nil
	}
}
