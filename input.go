package ffmpeg

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

type TranscoderInput interface {
	TranscoderOption
	input() *input
}

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

func Input(options ...InputOption) TranscoderInput {
	return &input{options: options}
}

func InputFilename(filename string) InputOption {
	return func(input *input) (err error) {
		input.fi, err = Stat(filename)
		return err
	}
}

func InputURL(url *url.URL) InputOption {
	return func(input *input) (err error) {
		input.URL = url
		return nil
	}
}

func InputFile(file *os.File) InputOption {
	return func(input *input) (err error) {
		input.fi, err = Stat(file.Name())
		input.file = file
		return err
	}
}

func StartOption(start Time) InputOption {
	return func(input *input) error {
		input.Start = start
		return nil
	}
}

func StartPercentOption(percent int) InputOption {
	return func(input *input) error {
		if input.fi == nil {
			return fmt.Errorf("StartPercent can only be used on a File input")
		}

		input.Start = input.fi.Format.Duration.Percent(percent)
		return nil
	}
}

func DurationOption(duration Time) InputOption {
	return func(input *input) error {
		input.Duration = duration
		return nil
	}
}
