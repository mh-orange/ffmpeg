package ffmpeg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/mh-orange/cmd"
)

func TestInterlaceRepeatedInfoUnmarshalText(t *testing.T) {
	tests := []struct {
		input       string
		wantNeither int
		wantTop     int
		wantBottom  int
		wantFrames  int
		wantErr     bool
	}{
		{"Fields: Neither: 1 Top: 2 Bottom: 3", 1, 2, 3, 6, false},
		{"Fields: neither: 1 Top: 2 Bottom: 3", 0, 0, 0, 0, true},
		{"foo bar Fields: Neither: 1 Top: 2 Bottom: 3", 1, 2, 3, 6, false},
		{"foo bar: Neither: 1 Top: 2 Bottom: 3", 0, 0, 0, 0, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			iri := &InterlaceRepeatedInfo{}
			err := iri.parse([]byte(test.input))
			if err == nil {
				if test.wantErr {
					t.Errorf("Want error got nil")
				} else {
					if test.wantNeither != iri.Neither {
						t.Errorf("Neither: want %d got %d", test.wantNeither, iri.Neither)
					}

					if test.wantTop != iri.Top {
						t.Errorf("Top: want %d got %d", test.wantTop, iri.Top)
					}

					if test.wantBottom != iri.Bottom {
						t.Errorf("Neither: want %d got %d", test.wantBottom, iri.Bottom)
					}

					if test.wantFrames != iri.Frames() {
						t.Errorf("Frames: want %d got %d", test.wantFrames, iri.Frames())
					}
				}
			} else if !test.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestInterlaceFieldInfoUnmarshalText(t *testing.T) {
	tests := []struct {
		input            string
		wantTFF          int
		wantBFF          int
		wantProgressive  int
		wantUndetermined int
		wantErr          bool
	}{
		{"detection: TFF: 1 BFF: 2 Progressive: 3 Undetermined: 4", 1, 2, 3, 4, false},
		{"detection: tff: 1 BFF: 2 Progressive: 3 Undetermined: 4", 0, 0, 0, 0, true},
		{"Detection: TFF: 1 BFF: 2 Progressive: 3 Undetermined: 4", 0, 0, 0, 0, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			iri := &InterlaceFieldInfo{}
			err := iri.parse([]byte(test.input))
			if err == nil {
				if test.wantErr {
					t.Errorf("want error got nil")
				} else {
					if test.wantTFF != iri.TFF {
						t.Errorf("TFF: want %d got %d", test.wantTFF, iri.TFF)
					}
					if test.wantBFF != iri.BFF {
						t.Errorf("BFF: want %d got %d", test.wantBFF, iri.BFF)
					}
					if test.wantProgressive != iri.Progressive {
						t.Errorf("Progressive: want %d got %d", test.wantProgressive, iri.Progressive)
					}
					if test.wantUndetermined != iri.Undetermined {
						t.Errorf("Undetermined: want %d got %d", test.wantUndetermined, iri.Undetermined)
					}
				}
			} else if !test.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func testInfo(values ...int) InterlaceInfo {
	return InterlaceInfo{
		InterlaceRepeatedInfo{values[0], values[1], values[2]},
		InterlaceFieldInfo{values[3], values[4], values[5], values[6]},
		InterlaceFieldInfo{values[7], values[8], values[9], values[10]},
	}
}

func TestInterlaceInfo(t *testing.T) {
	tests := []struct {
		input            InterlaceInfo
		wantTFF          int
		wantBFF          int
		wantInterlaced   int
		wantProgressive  int
		wantDetermined   int
		wantUndetermined int
		wantFrames       int
	}{
		{testInfo(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11), 12, 14, 26, 16, 42, 18, 6},
	}

	for i, test := range tests {
		if test.input.TFF() != test.wantTFF {
			t.Errorf("tests[%d] TFF: want %d got %d", i, test.wantTFF, test.input.TFF())
		}
		if test.input.BFF() != test.wantBFF {
			t.Errorf("tests[%d] BFF: want %d got %d", i, test.wantBFF, test.input.BFF())
		}
		if test.input.Interlaced() != test.wantInterlaced {
			t.Errorf("tests[%d] Interlaced: want %d got %d", i, test.wantInterlaced, test.input.Interlaced())
		}
		if test.input.Progressive() != test.wantProgressive {
			t.Errorf("tests[%d] Progressive: want %d got %d", i, test.wantProgressive, test.input.Progressive())
		}
		if test.input.Determined() != test.wantDetermined {
			t.Errorf("tests[%d] Determined: want %d got %d", i, test.wantDetermined, test.input.Determined())
		}
		if test.input.Undetermined() != test.wantUndetermined {
			t.Errorf("tests[%d] Undetermined: want %d got %d", i, test.wantUndetermined, test.input.Undetermined())
		}
		if test.input.Frames() != test.wantFrames {
			t.Errorf("tests[%d] Frames: want %d got %d", i, test.wantFrames, test.input.Frames())
		}
	}
}

func TestInterlaceInfoType(t *testing.T) {
	tests := []struct {
		input   InterlaceInfo
		want    InterlaceType
		wantErr error
	}{
		{testInfo(1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0), Unknown, ErrShortStream},
		{testInfo(1000, 0, 0, 0, 0, 1000, 0, 0, 0, 1000, 0), Progressive, nil},
		{testInfo(1000, 0, 0, 1000, 0, 0, 0, 1000, 0, 0, 0), InterlacedTff, nil},
		{testInfo(1000, 0, 0, 0, 1000, 0, 0, 0, 1000, 0, 0), InterlacedBff, nil},
		{testInfo(1000, 0, 0, 500, 500, 0, 0, 500, 500, 0, 0), Interlaced, nil},
	}

	for i, test := range tests {
		got, err := test.input.Type()
		if err == test.wantErr {
			if err == nil {
				if test.want != got {
					t.Errorf("tests[%d] want %v got %v", i, test.want, got)
				}
			}
		} else {
			t.Errorf("tests[%d] unexpected error: %v", i, err)
		}
	}
}

func TestInterlaceTranscode(t *testing.T) {
	names, err := filepath.Glob("testdata/interlace*.txt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, inputFile := range names {
		t.Run(inputFile, func(t *testing.T) {
			jsonFile := fmt.Sprintf("%s.json", inputFile[:len(inputFile)-len(".txt")])
			inputTxt, err := ioutil.ReadFile(inputFile)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			inputJSON, err := ioutil.ReadFile(jsonFile)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			want := InterlaceInfo{}
			err = json.Unmarshal(inputJSON, &want)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			c := &cmd.TestCmd{Stderr: inputTxt}
			Ffmpeg = c
			it := NewInterlaceTranscoder()
			got, err := it.transcode(Input())

			if err == nil {
				if want != got {
					t.Errorf("want %v got %v", want, got)
				}
			}
		})
	}
}
