package vtil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrShortStream = errors.New("stream was too short to process")
)

type InterlaceRepeatedInfo struct {
	Neither int `json:"neither"`
	Top     int `json:"top"`
	Bottom  int `json:"bottom"`
}

func (iri *InterlaceRepeatedInfo) Frames() int { return iri.Neither + iri.Top + iri.Bottom }

func (iri *InterlaceRepeatedInfo) parse(text []byte) (err error) {
	str := string(text)
	substr := "Fields: "
	if index := strings.Index(str, substr); index >= 0 {
		str = str[index+len(substr):]
		_, err = fmt.Sscanf(str, "Neither: %d Top: %d Bottom: %d", &iri.Neither, &iri.Top, &iri.Bottom)
	} else {
		err = fmt.Errorf("Input does not match pattern")
	}
	return err
}

// InterlaceFieldInfo contains the counts of each frame type that the filter
// detected
type InterlaceFieldInfo struct {
	// TFF is the number of Top Field First frames detected
	TFF int `json:"tff"`

	// BFF is the number of Bottom Field First frames detected
	BFF int `json:"bff"`

	// Progressive is the number of Progressive frames detected
	Progressive int `json:"progressive"`

	// Undetermined is the number of frames that could not be identified by the filter
	Undetermined int `json:"undetermined"`
}

func (ifi *InterlaceFieldInfo) parse(text []byte) (err error) {
	str := string(text)
	substr := "detection: "
	if index := strings.Index(str, substr); index >= 0 {
		str = str[index+len(substr):]
		_, err = fmt.Sscanf(str, "TFF: %d BFF: %d Progressive: %d Undetermined: %d", &ifi.TFF, &ifi.BFF, &ifi.Progressive, &ifi.Undetermined)
	} else {
		err = fmt.Errorf("Input does not match pattern")
	}
	return err
}

// InterlaceInfo includes the parsed information reported by the ffmpeg idet filter
type InterlaceInfo struct {
	// RepeatedFields contains the information that ffmpeg reports about repeated fields
	RepeatedFields InterlaceRepeatedInfo `json:"repeatedFields"`

	// SingleFrame contains the information reported about single-frame detection
	SingleFrame InterlaceFieldInfo `json:"singleFrame"`

	// MultieFrame contains the information reported about multi-frame detection
	MultiFrame InterlaceFieldInfo `json:"multiFrame"`
}

// TFF returns the sum of single and multi frame detected top frame first
func (ii InterlaceInfo) TFF() int { return ii.SingleFrame.TFF + ii.MultiFrame.TFF }

// BFF returns the sum of single and multi frame detected bottom frame first
func (ii InterlaceInfo) BFF() int { return ii.SingleFrame.BFF + ii.MultiFrame.BFF }

// Interlaced returns the sum of detected bottom and top frames first
func (ii InterlaceInfo) Interlaced() int { return ii.TFF() + ii.BFF() }

// Progressive returns the number of detected progressive frames
func (ii InterlaceInfo) Progressive() int {
	return ii.SingleFrame.Progressive + ii.MultiFrame.Progressive
}

// Determined returns the sum of detected interlaced and progressive frames
func (ii InterlaceInfo) Determined() int { return ii.Interlaced() + ii.Progressive() }

// Undetermined returns the number of single and multi frame detection frames that could not be determined
func (ii InterlaceInfo) Undetermined() int {
	return ii.SingleFrame.Undetermined + ii.MultiFrame.Undetermined
}

func (ii InterlaceInfo) Frames() int { return ii.RepeatedFields.Frames() }

func (ii InterlaceInfo) Type() (t InterlaceType, err error) {
	if ii.Frames() < 250 {
		err = ErrShortStream
	} else if ii.Determined() > ii.Undetermined() {
		if ii.Progressive() < ii.Interlaced()*20 {
			if ii.BFF() < ii.TFF() {
				t = InterlacedTff
			} else if ii.TFF() < ii.BFF() {
				t = InterlacedBff
			} else {
				// interlaced, not sure what order
				t = Interlaced
			}
		} else {
			t = Progressive
		}
	}
	return
}

type InterlaceTranscoder struct {
}

func NewInterlaceTranscoder() *InterlaceTranscoder {
	return &InterlaceTranscoder{}
}

func (it *InterlaceTranscoder) transcode(input TranscoderInput, options ...TranscoderOption) (info InterlaceInfo, err error) {
	r, writer := io.Pipe()
	reader := bufio.NewReader(r)
	transcoder := NewTranscoder()
	options = append([]TranscoderOption{input, StderrOption(writer)}, options...)
	job, err := transcoder.Transcode(append(options, DiscardOption())...)
	if err == nil {
		for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
			if index := bytes.Index(line, []byte("Repeated Fields:")); index >= 0 {
				info.RepeatedFields.parse(line)
			} else if index = bytes.Index(line, []byte("Single frame detection:")); index >= 0 {
				info.SingleFrame.parse(line)
			} else if index = bytes.Index(line, []byte("Multi frame detection:")); index >= 0 {
				info.MultiFrame.parse(line)
			}
		}
		err = job.Wait()
	}
	return info, err
}

func (it *InterlaceTranscoder) Deinterlace(input TranscoderInput, output TranscoderOutput) (TranscodeJob, error) {
	transcoder := NewTranscoder()
	return transcoder.Transcode(input, VideoFilterOption("yadif"), output)
}

func (it *InterlaceTranscoder) Detect(input TranscoderInput) (t InterlaceType, err error) {
	input.input().options = append(input.input().options, StartPercentOption(35), DurationOption(35*Second))
	info, err := it.transcode(input, VideoFilterOption("idet"))
	if err == nil {
		t, err = info.Type()
		if t != Unknown && t != Progressive {
			switch t {
			case InterlacedTff:
				info, err = it.transcode(input, VideoFilterOption("setfield=tff,pullup,idet"))
			case InterlacedBff:
				info, err = it.transcode(input, VideoFilterOption("setfield=bff,pullup,idet"))
			case Interlaced:
				info, err = it.transcode(input, VideoFilterOption("pullup,idet"))
			}

			if err == nil {
				t, err = info.Type()
				if t == Progressive {
					t = Telecine
				} else {
					t = Interlaced
				}
			}
		}
	}
	return t, err
}

func IsInterlaced(filename string) (isInterlaced bool, err error) {
	transcoder := NewInterlaceTranscoder()
	t, err := transcoder.Detect(Input(InputFilename(filename)))
	if t == Telecine || t == Interlaced {
		return true, err
	}
	return false, err
}
