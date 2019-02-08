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

package vtil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestValues(t *testing.T) {
	tests := []struct {
		input func() interface{}
		want  interface{}
	}{
		{func() interface{} { return ColorRangeValues() }, _ColorRangeValues},
		{func() interface{} { return ColorSpaceValues() }, _ColorSpaceValues},
		{func() interface{} { return FieldOrderValues() }, _FieldOrderValues},
		{func() interface{} { return InterlaceTypeValues() }, _InterlaceTypeValues},
		{func() interface{} { return MediaTypeValues() }, _MediaTypeValues},
	}

	for i, test := range tests {
		got := test.input()
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("tests[%d] want %+v got %+v", i, test.want, got)
		}
	}
}

func TestEnumer(t *testing.T) {
	tests := []struct {
		input interface{}
		str   string
		isA   bool
	}{
		{ColorRangeUnspecified, "unknown", true},
		{ColorRangeMPEG, "tv", true},
		{ColorRangeJPEG, "pc", true},
		{ColorRange(1024), "ColorRange(1024)", false},
		{ColorSpaceRGB, "gbr", true},
		{ColorSpaceBT709, "bt709", true},
		{ColorSpaceUnspecified, "unknown", true},
		{ColorSpaceReserved, "reserved", true},
		{ColorSpaceFcc, "fcc", true},
		{ColorSpaceBT470Bg, "bt470bg", true},
		{ColorSpaceSMPTE170M, "smpte170m", true},
		{ColorSpaceSMPTE240M, "smpte240m", true},
		{ColorSpaceYCoCg, "ycgco", true},
		{ColorSpaceBT2020Nc, "bt2020nc", true},
		{ColorSpaceBT2020C, "bt2020c", true},
		{ColorSpaceSMTPE2085, "smpte2085", true},
		{ColorSpaceChromaDerivedNc, "chroma-derived-nc", true},
		{ColorSpaceChromaDerivedC, "chroma-derived-c", true},
		{ColorSpaceICtCp, "ictcp", true},
		{ColorSpace(1024), "ColorSpace(1024)", false},
		{FieldOrderUnknown, "unknown", true},
		{FieldOrderProgressive, "progressive", true},
		{FieldOrderTT, "tt", true},
		{FieldOrderBB, "bb", true},
		{FieldOrderTB, "tb", true},
		{FieldOrderBT, "bt", true},
		{FieldOrder(1024), "FieldOrder(1024)", false},
		{Unknown, "unknown", true},
		{Telecine, "telecine", true},
		{Interlaced, "interlaced", true},
		{InterlacedTff, "interlaced TFF", true},
		{InterlacedBff, "interlaced BFF", true},
		{Progressive, "progressive", true},
		{InterlaceType(1024), "InterlaceType(1024)", false},
		{Video, "video", true},
		{Audio, "audio", true},
		{Data, "data", true},
		{Subtitle, "subtitle", true},
		{Attachment, "attachment", true},
		{MediaType(1024), "MediaType(1024)", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			typ := reflect.TypeOf(test.input)
			want := reflect.New(typ)
			reflect.Indirect(want).Set(reflect.ValueOf(test.input))
			runTest(t, typ, want, test.str, test.isA)
		})
	}
}

func runTest(t *testing.T, typ reflect.Type, val reflect.Value, str string, isA bool) {
	// test IsA
	if meth, ok := typ.MethodByName(fmt.Sprintf("IsA%s", typ.Name())); ok {
		values := val.MethodByName(meth.Name).Call(nil)
		if values[0].Bool() != isA {
			t.Errorf("%+v %v expected %v got %v", val, str, isA, values[0].Bool())
		}
	} else {
		t.Errorf("%v does not implement IsA%s", typ.Name(), typ.Name())
	}

	// check for String()
	if _, ok := typ.MethodByName("String"); ok {
		if stringer, ok := val.Interface().(fmt.Stringer); ok {
			testString(t, stringer, str)
		} else if stringer, ok := reflect.Indirect(val).Interface().(fmt.Stringer); ok {
			testString(t, stringer, str)
		} else {
			t.Errorf("%v has String() method but does not satisfy fmt.Stringer interface", typ)
		}
	}

	// check for MarshalJSON
	if marshaler, ok := val.Interface().(json.Marshaler); ok {
		testMarshalJSON(t, marshaler, str)
	} else if marshaler, ok := reflect.Indirect(val).Interface().(json.Marshaler); ok {
		testMarshalJSON(t, marshaler, str)
	}

	// check for UnmarshalJSON
	if unmarshaler, ok := val.Interface().(json.Unmarshaler); ok {
		testUnmarshalJSON(t, unmarshaler, str, val)
	} else if unmarshaler, ok := reflect.Indirect(val).Interface().(json.Unmarshaler); ok {
		testUnmarshalJSON(t, unmarshaler, str, val)
	}
}

func testString(t *testing.T, input fmt.Stringer, want string) {
	if input.String() != want {
		t.Errorf("want %q got %q", want, input.String())
	}
}

func testMarshalJSON(t *testing.T, input json.Marshaler, want string) {
	data, err := input.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if string(data) != fmt.Sprintf(`"%s"`, want) {
		t.Errorf(`want "%s" got %s`, want, string(data))
	}
}

func testUnmarshalJSON(t *testing.T, unmarshaler json.Unmarshaler, input string, want reflect.Value) error {
	err := unmarshaler.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, input)))
	if err == nil {
		if !reflect.DeepEqual(want.Interface(), unmarshaler) {
			t.Errorf("want %+v got %+v", want.Interface(), unmarshaler)
		}

		// check for unmarshaling malformed json
		err = unmarshaler.UnmarshalJSON([]byte(input))
		if err == nil {
			t.Errorf("expected to get an error for malformed json")
		}
	}
	return err
}

func TestRationalMarshal(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		want    Rational
	}{
		{`"1/100"`, false, Rational{1, "/", 100}},
		{`"1:100"`, false, Rational{1, ":", 100}},
		{`"1-100"`, true, Rational{}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			data := []byte(test.input)
			got := &Rational{}
			err := got.UnmarshalJSON(data)
			if test.wantErr && err == nil {
				t.Errorf("Expected error got none")
			} else if err == nil {
				if test.want == *got {
					data, _ = got.MarshalJSON()
					if !bytes.Equal(data, []byte(test.input)) {
						t.Errorf("Wanted %s got %s", test.input, string(data))
					}
				} else {
					t.Errorf("Wanted %+v got %+v", test.want, *got)
				}
			} else {
				if re, ok := err.(*InvalidRationalErr); ok && test.wantErr {
					if re.Error() != re.cause.Error() {
						t.Errorf("Invalid error string: %q", re.Error())
					}
				} else {
					t.Errorf("Unexpected error %v", err)
				}
			}
		})
	}
}

func TestTimePercent(t *testing.T) {
	tests := []struct {
		input   string
		percent int
		want    string
	}{
		{"01:00:00.000000", 50, "00:30:00.000000"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t1 := Time(0)
			t1.Parse(test.input)
			t2 := t1.Percent(test.percent)
			got := t2.String()
			if test.want != got {
				t.Errorf("want %q got %q", test.want, got)
			}
		})
	}
}

func TestTimeMarshal(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		want    Time
	}{
		{`"00:25:44.530750"`, false, Time(1544530750000)},
		{`"Monday"`, true, Time(0)},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			data := []byte(test.input)
			got := Time(0)
			err := got.UnmarshalJSON(data)
			if test.wantErr && err == nil {
				t.Errorf("Expected error got none")
			} else if err == nil {
				if test.want == got {
					data, _ = got.MarshalJSON()
					if !bytes.Equal(data, []byte(test.input)) {
						t.Errorf("Wanted %s got %s", test.input, string(data))
					}
				} else {
					t.Errorf("Wanted %+v got %+v", test.want, got)
				}
			} else {
				if te, ok := err.(*InvalidTimeErr); ok && test.wantErr {
					if te.Error() != te.cause.Error() {
						t.Errorf("Invalid time error format: %q", te.Error())
					}
				} else {
					t.Errorf("Unexpected error %v", err)
				}
			}
		})
	}
}
