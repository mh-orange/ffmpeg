// Copyright 2019 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate enumer -type=ColorRange -json=true -transform=comment
//go:generate enumer -type=ColorSpace -json=true -transform=comment
//go:generate enumer -type=FieldOrder -json=true -transform=comment
//go:generate enumer -type=InterlaceType -json=true -transform=comment
//go:generate enumer -type=MediaType -json=true -transform=comment

package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ColorRange indicates how the colors of the video stream are encoded
type ColorRange int

const (
	// ColorRangeUnspecified indicates the stream does not specify a color range
	ColorRangeUnspecified ColorRange = iota // unknown

	// ColorRangeMPEG (also known as TV) is limited range color
	ColorRangeMPEG // tv

	// ColorRangeJPEG (also known as PC) is full range color
	ColorRangeJPEG // pc
)

// ColorSpace indicates the way colors and pixels are arranged in the pixel buffer
type ColorSpace int

const (
	// ColorSpaceRGB indicates standard Red, Green, Blue color space
	ColorSpaceRGB ColorSpace = iota // gbr

	// ColorSpaceBT709 indicates high definition color
	ColorSpaceBT709 // bt709

	// ColorSpaceUnspecified indicates that the color space encoding is not specified in the stream
	ColorSpaceUnspecified // unknown

	// ColorSpaceReserved
	ColorSpaceReserved // reserved

	ColorSpaceFcc // fcc

	ColorSpaceBT470Bg // bt470bg

	ColorSpaceSMPTE170M // smpte170m

	ColorSpaceSMPTE240M // smpte240m

	ColorSpaceYCoCg // ycgco

	// ColorSpaceBT2020Nc indicates ultra high definition non-constant luminance color
	ColorSpaceBT2020Nc // bt2020nc

	// ColorSpaceBT2020C indicates ultra high definition constant luminance color
	ColorSpaceBT2020C // bt2020c

	ColorSpaceSMTPE2085 // smpte2085

	ColorSpaceChromaDerivedNc // chroma-derived-nc

	ColorSpaceChromaDerivedC // chroma-derived-c

	ColorSpaceICtCp // ictcp
)

// FieldOrder indicates whether or not a video stream is interlaced or progressive
type FieldOrder int

const (
	// FieldOrderUnknown indicates the stream field order could not be determined
	FieldOrderUnknown FieldOrder = iota // unknown

	// FieldOrderProgressive indicates progressive coded frames
	FieldOrderProgressive // progressive

	// FieldOrderTT indicates the top field was coded first and is displayed first
	FieldOrderTT // tt

	// FieldOrderBB indicates the bottom field was coded first and is displayed first
	FieldOrderBB // bb

	// FieldOrderTB indicates the top field was coded first but the bottom frame is displayed first
	FieldOrderTB // tb

	// FieldOrderBT indicates the bottom field was coded first but the top frame is displayed first
	FieldOrderBT // bt
)

// InterlaceType indicates the field order/type of interlacing (if any) that is
// detected in a video stream
type InterlaceType int

const (
	// Unknown indicates that interlacing could not be determined
	Unknown InterlaceType = iota // unknown

	// Telecine indicates the video stream appears to be encoded using telecine
	Telecine // telecine

	// Interlaced indicates the frames are interlaced, but the field order is not known
	Interlaced // interlaced

	// InterlacedTff indicates Top Frame First interlacing
	InterlacedTff // interlaced TFF

	// InterlacedBff indicates Bottom Frame First interlacing
	InterlacedBff // interlaced BFF

	// Progressive indicates no interlacing was detected
	Progressive // progressive
)

// InvalidRationalErr is returned when invalid data is encountered
// while unmarshaling a rational
type InvalidRationalErr struct {
	cause error
}

// Error returns the description of the error
func (ire *InvalidRationalErr) Error() string {
	return fmt.Sprintf("%v", ire.cause)
}

// Rational represents a rational number with a numerator and denominator
type Rational struct {
	// Numerator is the number in the top position of the rational
	Numerator int

	// Separator is the character used to separate the numerator and denoinator. This
	// is used for parsing and displaying the rational.  It is usually a '/'
	Separator string

	// Denominator is the number in the bottom position of the rational
	Denominator int
}

// MarshalJSON will return a properly formated JSON string representation of the rational number
func (r *Rational) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d%s%d"`, r.Numerator, r.Separator, r.Denominator)), nil
}

// UnmarshalJSON will parse a JSON string and assign the numerator and denominator to the
// rational.  UnmarshalJSON looks for either a colon (':') or slash ('/') as a separator.
// If neither a colon or slash is found, then an InvalidRationalErr is returned.
func (r *Rational) UnmarshalJSON(data []byte) (err error) {
	if index := bytes.Index(data, []byte(":")); index > 0 {
		r.Separator = ":"
		_, err = fmt.Sscanf(string(data), `"%d:%d"`, &r.Numerator, &r.Denominator)
	} else if index = bytes.Index(data, []byte("/")); index > 0 {
		r.Separator = "/"
		_, err = fmt.Sscanf(string(data), `"%d/%d"`, &r.Numerator, &r.Denominator)
	} else {
		err = fmt.Errorf("Unknown format for %q", string(data))
	}

	if err != nil {
		err = &InvalidRationalErr{err}
	}
	return err
}

// AspectRatio indicates the ratio of width and height of the video stream
type AspectRatio struct{ Rational }

// TimeBase indicates the base rate for the clock of the video stream.  In
// essence, this is a fraction of a second that represents on tick of the
// clock for the video encoder or decoder
type TimeBase struct{ Rational }

// FrameRate indicates the number of frames per second contained in the
// video stream
type FrameRate struct{ Rational }

// InvalidTimeErr indicates that the time value (string, json, etc) could
// not be parsed due to a wrong format
type InvalidTimeErr struct {
	cause error
}

// Error provides the reason for the parse error
func (ite *InvalidTimeErr) Error() string {
	return fmt.Sprintf("%v", ite.cause)
}

// Time is the relative time in nanoseconds. It is used to represent things
// like chapter start and end times
type Time uint64

const (
	// Nanosecond is one billionth of a second.  A nanosecond is the basic unit
	// of time
	Nanosecond Time = 1

	// Microsecond is 1000 nanoseconds or one millionth of a second
	Microsecond = 1000 * Nanosecond

	// Millisecond is 1000 microseconds or one thousandth of a second
	Millisecond = 1000 * Microsecond

	// Second is 1000 milliseconds
	Second = 1000 * Millisecond

	// Minute is 60 seeconds
	Minute = 60 * Second

	// Hour is 60 minutes
	Hour = 60 * Minute
)

// Percent computes the correct percentage of time.  For instance, 10 percent of
// 1 hour returns a time representing 6 minutes
func (t Time) Percent(percent int) Time {
	tt := uint64(t) / 100
	tt *= uint64(percent)
	return Time(tt)
}

// Parse will take a time string in the form of "HH:MM:SS.mmmmmm" and assign
// the values (hour, minute, second, milliseconds) to the Time receiver
func (t *Time) Parse(str string) error {
	hr, min, sec, ms := Time(0), Time(0), Time(0), Time(0)
	_, err := fmt.Sscanf(str, `%d:%d:%d.%d`, &hr, &min, &sec, &ms)
	if err == nil {
		*t = (hr * Hour) + (min * Minute) + (sec * Second) + (ms * Microsecond)
	} else {
		err = &InvalidTimeErr{err}
	}
	return err
}

// UnmarshalJSON takes a JSON string and parses it into the correct
// Time
func (t *Time) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err == nil {
		err = t.Parse(str)
	}
	return err
}

// String returns a Time string in the form of "HH:MM:SS.mmmmmm" where H is hour,
// M is minute, S is second and m is milliseconds
func (t Time) String() string {
	val := t
	hr := val / Hour
	val -= (hr * Hour)
	min := val / Minute
	val -= (min * Minute)
	sec := val / Second
	val -= (sec * Second)
	ms := val / Microsecond
	return fmt.Sprintf("%02d:%02d:%02d.%06d", hr, min, sec, ms)
}

// MarshalJSON returns a properly formated JSON string representing the time
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// PTS is the Presentation Time Stamp
type PTS uint64

// MediaType is used to indicate the type of media a stream contains
type MediaType int

const (
	// Video includes any Video codec
	Video MediaType = iota // video

	// Audio includes any audio codec
	Audio // audio

	// Data are data streams
	Data // data

	// Subtitle matches subtitles in the container
	Subtitle // subtitle

	// Attachment are attachment streams
	Attachment // attachment
)
