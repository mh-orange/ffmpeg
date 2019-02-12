package ffmpeg

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

// TranscoderInput is a TranscoderOption used to indicate how to access the input
// media.  It could be a filename string, an io.ReadWriteSeeker or a URL.  The
// TranscoderInput also has options for processing the input, such as seeking to
// a start position or only processing a given duration
type TranscoderInput interface {
	TranscoderOption
	input() *input
}

// InputOption is an option passed to a TranscoderInput to alter the input in some
// way (such as specifying a starting position or duration)
type InputOption func(*input) error

type input struct {
	URL      *url.URL
	Start    Time
	Duration Time

	fi      *FileInfo
	file    io.ReadWriteSeeker
	args    []string
	options []InputOption
}

func (in *input) input() *input {
	return in
}

func (in *input) process(job *transcodeJob) (err error) {
	if len(in.args) == 0 {
		for _, option := range in.options {
			err = option(in)
			if err != nil {
				break
			}
		}

		if err == nil {
			if in.Start != 0 {
				in.args = append(in.args, "-ss", in.Start.String())
			}

			if in.Duration != 0 {
				in.args = append(in.args, "-t", in.Duration.String())
			}

			if in.URL != nil {
				in.args = append(in.args, "-i", in.URL.String())
			} else if in.fi != nil {
				in.args = append(in.args, "-i", in.fi.Format.Filename)
			}
		}
	}
	job.proc.AppendArgs(in.args...)
	return
}

// Input creates a TranscoderInput and applies the options
func Input(options ...InputOption) TranscoderInput {
	return &input{options: options}
}

// InputFilename creates an InputOption that will pass the filename on
// to the ffmpeg process
func InputFilename(filename string) InputOption {
	return func(input *input) (err error) {
		input.fi, err = Stat(filename)
		return err
	}
}

// InputURL creates an InputOption that will pass the URL on to the
// underlying ffmpeg process
func InputURL(url *url.URL) InputOption {
	return func(input *input) (err error) {
		input.URL = url
		return nil
	}
}

// InputFile will create an InputOption that reads the input file
// and sends the data to the ffmpeg process using STDIN
func InputFile(file *os.File) InputOption {
	return func(input *input) (err error) {
		input.fi, err = Stat(file.Name())
		input.file = file
		return err
	}
}

// StartOption sets the starting position that ffmpeg will attempt
// to seek.  This sets the -ss option on the input on the ffmpeg command
// line
func StartOption(start Time) InputOption {
	return func(input *input) error {
		input.Start = start
		return nil
	}
}

// StartPercentOption can only be used if input media can be read by
// the Stat function.  If the media can be read by the Stat function then
// the -ss argument passed to ffmpeg will be a percentage of the input
// streams duration
func StartPercentOption(percent int) InputOption {
	return func(input *input) error {
		if input.fi == nil {
			return fmt.Errorf("StartPercent can only be used on a File input")
		}

		input.Start = input.fi.Format.Duration.Percent(percent)
		return nil
	}
}

// DurationOption will tell ffmpeg to process only "duration" frames
// from the input file.  This uses the ffmpeg -t option
func DurationOption(duration Time) InputOption {
	return func(input *input) error {
		input.Duration = duration
		return nil
	}
}
