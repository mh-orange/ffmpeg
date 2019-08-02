package ffmpeg

import (
	"errors"
	"io"
	"strings"

	"github.com/mh-orange/cmd"
)

// Transcoder is a mechanism to take an input, optionally modify it
// with filters and produce an output.  Transcoder is useful when
// changing a video's format, resolution, encoding etc.
type Transcoder struct {
	options []TranscoderOption
}

// NewTranscoder returns a Transcoder object that will always utilize the
// given options whenever transcoding an input
func NewTranscoder(options ...TranscoderOption) *Transcoder {
	transcoder := &Transcoder{
		options: options,
	}

	return transcoder
}

// Transcode will start a new transcoding process for the specific options and return a TranscodeJob
// that can be monitored for completion.
func (transcoder *Transcoder) Transcode(options ...TranscoderOption) (TranscodeJob, error) {
	var err error
	options = append(transcoder.options, options...)

	job := &transcodeJob{
		progressCh: make(chan TranscodeInfo, 1),
	}
	job.proc = ffmpeg.Process()

	// search input for longest duration
	for _, option := range options {
		err = option.process(job)
		if err != nil {
			break
		}

		if input, ok := option.(*input); ok {
			if input.Duration > job.info.Duration {
				job.info.Duration = input.Duration
			} else if input.fi != nil && input.fi.Format.Duration > job.info.Duration {
				job.info.Duration = input.fi.Format.Duration
			}
		}
	}

	if err == nil {
		stderr, writer := io.Pipe()
		job.proc.Stderr(writer)
		err = job.proc.Start()
		if err == nil {
			cancelCh := make(chan struct{})
			job.cancelCh = cancelCh

			doneCh := make(chan struct{})
			job.doneCh = doneCh
			go job.run(cancelCh, doneCh, stderr)
		}
	}
	return job, err
}

// TranscodeJob is a transcode session that is either ready to be started
// or is currently running
type TranscodeJob interface {
	// Cancel attempts to stop/cancel a running transcode session
	Cancel()

	// Err will return any error that occurred during the transcode session
	Err() error

	// Log is the string output from the underlying transcode command.  This is useful
	// for determining if any errors occurred during transcoding
	Log() string

	// Progress returns a channel that receives TranscodeInfo objects as transcoding progresses.
	// This is useful for displaying progress and feedback to users
	Progress() <-chan TranscodeInfo

	// Wait will block until the underlying transcode process finishes.
	Wait() error
}

type transcodeJob struct {
	io.Reader
	log []string
	err error

	info TranscodeInfo
	proc cmd.Process

	progressCh chan TranscodeInfo
	cancelCh   chan<- struct{}
	doneCh     <-chan struct{}
}

func (job *transcodeJob) run(cancelCh chan struct{}, doneCh chan struct{}, stderr io.Reader) {
	defer close(doneCh)
	defer close(job.progressCh)

	reader := newFilterReader(stderr, progPtrn, statsPtrn, finalStatsPtrn, repeatPtrn)

	values := make(map[string]string)
	running := true

	for running {
		select {
		case <-cancelCh:
			job.proc.Kill()
			running = false
		default:
			if reader.Scan() {
				if reader.Pattern() == nil {
					job.log = append(job.log, reader.Text())
				} else if reader.Pattern() == progPtrn {
					tokens := strings.Split(reader.Text(), "=")
					values[strings.TrimSpace(tokens[0])] = strings.TrimSpace(tokens[1])
					if strings.TrimSpace(tokens[0]) == "progress" {
						job.info.update(values)
						select {
						case job.progressCh <- job.info:
						default:
						}
						values = make(map[string]string)
					}
				}
			} else {
				if reader.Err() != nil && reader.Err() != io.EOF {
					job.log = append(job.log, reader.Err().Error())
				}
				job.proc.Kill()
				running = false
			}
		}
	}

	job.err = job.proc.Wait()
	if job.err != nil {
		if len(job.log) >= 2 {
			job.err = errors.New(strings.Join(job.log[len(job.log)-2:], "\n"))
		} else if len(job.log) == 1 {
			job.err = errors.New(strings.TrimSpace(job.log[0]))
		}
	}
}

func (job *transcodeJob) Err() error {
	return job.err
}

func (job *transcodeJob) Log() string {
	return strings.Join(job.log, "\n")
}

func (job *transcodeJob) Progress() <-chan TranscodeInfo {
	return job.progressCh
}

func (job *transcodeJob) Cancel() {
	if job.cancelCh != nil {
		job.cancelCh <- struct{}{}
		job.cancelCh = nil
	}
}

func (job *transcodeJob) Wait() error {
	<-job.doneCh
	return job.err
}
