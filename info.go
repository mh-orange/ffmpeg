// Copyright 2019 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");

// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ProgramInfo is information about programs and their streams, returned by ffprobe
type ProgramInfo struct {
}

// StreamInfo is the information from ffprobe about individual streams contained in some
// media
type StreamInfo struct {
	// Index is the stream index for this stream
	Index int `json:"index"`

	// CodecType indicates if the media is video, audio or subtitle
	CodecType MediaType `json:"codec_type"`

	// CodecName is a string representation of the codec (h264, for instance)
	CodecName string `json:"codec_name"`

	// CodecLongName is an extended version of the CodecName
	CodecLongName string `json:"codec_long_name"`

	// Profile
	Profile string `json:"profile"`

	// CodecTimeBase indicates the time base (relative to seconds) for the codec
	CodecTimeBase TimeBase `json:"codec_time_base"`

	CodecTagString string `json:"codec_tag_string"`

	CodecTag string `json:"codec_tag"`

	RFrameRate FrameRate `json:"r_frame_rate"`

	AvgFrameRate FrameRate `json:"avg_frame_rate"`

	// TimeBase indicates the time base (relative to seconds) for the format/container. TODO: Confirm this meaning
	TimeBase TimeBase `json:"time_base"`

	// StartPts indicates the Program Time Stamp (PTS) at the start of the stream
	StartPts PTS `json:"pts"`

	// Duration is the length of the stream
	Duration Time `json:"duration"`

	Disposition DispositionInfo `json:"disposition"`
}

type VideoStreamInfo struct {
	StreamInfo

	Width              int         `json:"width"`
	Height             int         `json:"height"`
	CodedWidth         int         `json:"coded_width"`
	CodedHeight        int         `json:"coded_height"`
	HasBFrames         int         `json:"has_b_frames"`
	SampleAspectRatio  AspectRatio `json:"sample_aspect_ratio"`
	DisplayAspectRatio AspectRatio `json:"display_aspect_ratio"`
	PixFmt             string      `json:"pix_fmt"`
	Level              int         `json:"level"`
	ColorRange         ColorRange  `json:"color_range"`
	ColorSpace         ColorSpace  `json:"color_space"`
	ColorTransfer      string      `json:"color_transfer"`
	ColorPrimaries     string      `json:"color_primaries"`
	ChromaLocation     string      `json:"chroma_location"`
	FieldOrder         FieldOrder  `json:"field_order"`
	Refs               int         `json:"refs"`
}

type AudioStreamInfo struct {
	StreamInfo

	SampleFmt     string `json:"sample_fmt"`
	SampleRate    string `json:"sample_rate"`
	Channels      int    `json:"channels"`
	ChannelLayout string `json:"channel_layout"`
	BitsPerSample int    `json:"bits_per_sample"`
}

type SubtitleStreamInfo struct {
	StreamInfo

	Width  int `json:"width"`
	Height int `json:"height"`
}

type DispositionInfo struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

type ChapterInfo struct {
	ID    int  `json:"id"`
	Start Time `json:"start_time"`
	End   Time `json:"end_time"`
}

// FormatInfo is everything we know about the format (container) of the file
type FormatInfo struct {
	// Filename is the URL/Filename that was passed to the ffprobe command
	Filename string `json:"filename"`

	// FormatName is a string representation of the format, such as "matroska"
	FormatName string `json:"format_name"`

	// FormatLongName is a descriptive version of the name
	FormatLongName string `json:"format_long_name"`

	// StartTime is the PTS of the first frame of the stream (in presentation order)
	StartTime Time `json:"start_time"`

	// Duration is the total length of the media
	Duration Time `json:"duration"`

	// Size is a string representation of the size of the file (size + unit)
	Size string `json:"size"`

	// BitRate is the total stream bitrate in bits/second
	BitRate string `json:"bit_rate"`

	// ProbeScore is a value between 0 and 100 that indicates how well ffprobe did when
	// trying to determine what kind of media was contained in the file
	ProbeScore int `json:"probe_score"`
}

// FileInfo is informational data about a given media file
type FileInfo struct {
	// Programs is the list of information about all the programs contained in the media
	Programs []*ProgramInfo `json:"Programs"`

	// VideoStreams is the list of information about all video streams in the media
	VideoStreams []*VideoStreamInfo

	// AudioStreams is the list of information about all audio streams in the media
	AudioStreams []*AudioStreamInfo

	// SubtitleStreams is the list of information about all the subtitles in the media
	SubtitleStreams []*SubtitleStreamInfo

	// Chapters is a list of all the chapters contained in the media
	Chapters []*ChapterInfo `json:"Chapters"`

	// Format is all the information relating to the container format
	Format FormatInfo `json:"Formats"`
}

// IsVideo determines of the FileInfo represents video media by determining
// if there is at least one video stream in the container
func (fi *FileInfo) IsVideo() bool {
	return len(fi.VideoStreams) > 0
}

// UnmarshalJSON takes the JSON string returned by ffprobe and parses it into
// the FileInfo object
func (fi *FileInfo) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Programs []*ProgramInfo
		Streams  []json.RawMessage
		Chapters []*ChapterInfo
		Format   FormatInfo
	}{}

	err := json.Unmarshal(data, &tmp)

	if err == nil {
		fi.Programs = tmp.Programs
		fi.Chapters = tmp.Chapters
		fi.Format = tmp.Format
		for _, str := range tmp.Streams {
			si := &StreamInfo{}
			err = json.Unmarshal(str, si)
			if err == nil {
				switch si.CodecType {
				case Video:
					vs := &VideoStreamInfo{StreamInfo: *si}
					err = json.Unmarshal(str, vs)
					if err == nil {
						fi.VideoStreams = append(fi.VideoStreams, vs)
					}
				case Audio:
					as := &AudioStreamInfo{StreamInfo: *si}
					err = json.Unmarshal(str, as)
					if err == nil {
						fi.AudioStreams = append(fi.AudioStreams, as)
					}
				case Subtitle:
					ss := &SubtitleStreamInfo{StreamInfo: *si}
					err = json.Unmarshal(str, ss)
					if err == nil {
						fi.SubtitleStreams = append(fi.SubtitleStreams, ss)
					}
				}
			}

			if err != nil {
				break
			}
		}
	}

	return err
}

// Stat will pass the filename to ffprobe and parse the output.  If no error occurs, then
// a FileInfo containing all the stream, program and format information is returned
func Stat(filename string) (fi *FileInfo, err error) {
	proc := ffprobe.Process()
	proc.AppendArgs(filename)
	logWriter := bytes.NewBuffer(nil)
	writer := bytes.NewBuffer(nil)
	proc.Stdout(writer)
	proc.Stderr(logWriter)
	err = proc.Start()
	if err == nil {
		err = proc.Wait()
		if err == nil {
			fi = &FileInfo{}
			err = json.Unmarshal(writer.Bytes(), fi)
		} else {
			err = fmt.Errorf("%s", string(logWriter.Bytes()))
		}
	}
	return fi, err
}

// IsVideo determines of the FileInfo represents video media by determining
// if there is at least one video stream in the container
func IsVideo(filename string) (bool, error) {
	fi, err := Stat(filename)
	if err == nil {
		return len(fi.VideoStreams) > 0, nil
	}
	return false, err
}
