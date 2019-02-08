package vtil

import (
	"fmt"
	"testing"
)

func TestTranscodeInfoUpdate(t *testing.T) {
	tests := []struct {
		input   map[string]string
		want    TranscodeInfo
		wantErr bool
	}{
		{map[string]string{"frame": "1"}, TranscodeInfo{Frame: 1}, false},
		{map[string]string{"frame": "one"}, TranscodeInfo{}, true},
		{map[string]string{"fps": "2"}, TranscodeInfo{Fps: 2}, false},
		{map[string]string{"fps": "two"}, TranscodeInfo{}, true},
		{map[string]string{"bitrate": "3"}, TranscodeInfo{Bitrate: 3}, false},
		{map[string]string{"bitrate": "three"}, TranscodeInfo{}, true},
		{map[string]string{"total_size": "4"}, TranscodeInfo{TotalSize: 4}, false},
		{map[string]string{"total_size": "four"}, TranscodeInfo{}, true},
		{map[string]string{"out_time_us": "7"}, TranscodeInfo{}, false},
		{map[string]string{"out_time_ms": "6"}, TranscodeInfo{}, false},
		{map[string]string{"out_time": "00:05:03.00000"}, TranscodeInfo{Time: Time(303000000000)}, false},
		{map[string]string{"out_time": "five"}, TranscodeInfo{}, true},
		{map[string]string{"dup_frames": "8"}, TranscodeInfo{DupFrames: 8}, false},
		{map[string]string{"dup_frames": "eight"}, TranscodeInfo{}, true},
		{map[string]string{"drop_frames": "9"}, TranscodeInfo{DropFrames: 9}, false},
		{map[string]string{"drop_frames": "nine"}, TranscodeInfo{}, true},
		{map[string]string{"speed": "10"}, TranscodeInfo{Speed: 10}, false},
		{map[string]string{"speed": "ten"}, TranscodeInfo{}, true},
		{map[string]string{"frame": "N/A"}, TranscodeInfo{}, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test.input), func(t *testing.T) {
			var ti TranscodeInfo
			err := ti.update(test.input)
			if err != nil && !test.wantErr {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && test.wantErr {
				t.Errorf("wanted error got nil")
			} else {
				if test.want != ti {
					t.Errorf("wanted %+v got %+v", test.want, ti)
				}
			}
		})
	}
}
