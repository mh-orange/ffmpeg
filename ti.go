package ffmpeg

import (
	"strconv"
	"strings"
)

// TranscodeInfo contains information periodically emitted by a transcode job.  The TranscodeInfo
// contains information that is useful when monitorring a running transcode job.
type TranscodeInfo struct {
	// Duration is the length of the video stream being transcoded
	Duration Time

	// Frame is the current frame being processed
	Frame int

	// Fps indicates the number of frames being processed every second
	Fps float64

	// Bitrate is the current bitrate that the file is being processed
	Bitrate float64

	// TotalSize is the total current size (in bytes) of the output file
	TotalSize int64

	// Time indicates the time offset currently being processed in the stream
	Time Time

	// DupFrames is the number of duplicate frames encoudntered during processing
	DupFrames int

	// DropFrames is the number of frames dropped during processing
	DropFrames int

	// Speed is the processing speed, relative to real/time.  For instance if the file is
	// being processed twice as fast as it would be played then Speed is 2.0.  Likewise, if
	// it is taking twice as long to process as to play, then it will be 0.5
	Speed float64
}

func (ti *TranscodeInfo) update(values map[string]string) (err error) {
	for k, v := range values {
		v = strings.TrimSpace(v)
		if v == "N/A" {
			continue
		}

		switch k {
		case "frame":
			ti.Frame, err = strconv.Atoi(v)
		case "fps":
			ti.Fps, err = strconv.ParseFloat(v, 64)
		case "bitrate":
			ti.Bitrate, err = strconv.ParseFloat(strings.TrimSuffix(v, "kbits/s"), 64)
		case "total_size":
			ti.TotalSize, err = strconv.ParseInt(v, 10, 64)
		case "out_time_us":
			continue
		case "out_time_ms":
			continue
		case "out_time":
			err = ti.Time.Parse(v)
		case "dup_frames":
			ti.DupFrames, err = strconv.Atoi(v)
		case "drop_frames":
			ti.DropFrames, err = strconv.Atoi(v)
		case "speed":
			ti.Speed, err = strconv.ParseFloat(strings.TrimSuffix(v, "x"), 64)
		case "progress":
		}
	}
	return
}
