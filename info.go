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

type ProgramInfo struct {
}

type StreamInfo struct {
	Index          int               `json:"index"`
	CodecType      MediaType         `json:"codec_type"`
	CodecName      string            `json:"codec_name"`
	CodecLongName  string            `json:"codec_long_name"`
	Profile        string            `json:"profile"`
	CodecTimeBase  TimeBase          `json:"codec_time_base"`
	CodecTagString string            `json:"codec_tag_string"`
	CodecTag       string            `json:"codec_tag"`
	RFrameRate     FrameRate         `json:"r_frame_rate"`
	AvgFrameRate   FrameRate         `json:"avg_frame_rate"`
	TimeBase       TimeBase          `json:"time_base"`
	StartPts       PTS               `json:"pts"`
	Duration       Time              `json:"duration"`
	Disposition    DispositionInfo   `json:"disposition"`
	Tags           map[string]string `json:"tags"`
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
	ID    int               `json:"id"`
	Start Time              `json:"start_time"`
	End   Time              `json:"end_time"`
	Tags  map[string]string `json:"tags"`
}

type FormatInfo struct {
	Filename       string            `json:"filename"`
	NbStreams      int               `json:"nb_streams"`
	NbPrograms     int               `json:"nb_programs"`
	FormatName     string            `json:"format_name"`
	FormatLongName string            `json:"format_long_name"`
	StartTime      Time              `json:"start_time"`
	Duration       Time              `json:"duration"`
	Size           string            `json:"size"`
	BitRate        string            `json:"bit_rate"`
	ProbeScore     int               `json:"probe_score"`
	Tags           map[string]string `json:"tags"`
}

type FileInfo struct {
	Programs        []*ProgramInfo `json:"Programs"`
	VideoStreams    []*VideoStreamInfo
	AudioStreams    []*AudioStreamInfo
	SubtitleStreams []*SubtitleStreamInfo
	Chapters        []*ChapterInfo `json:"Chapters"`
	Format          FormatInfo     `json:"Formats"`
}

func (fi *FileInfo) IsVideo() bool {
	return len(fi.VideoStreams) > 0
}

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
