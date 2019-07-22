package ffmpeg

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
)

var (
	progPtrn       = regexp.MustCompile(`^([^=]+)=\s*([^\s]+)$`)
	statsPtrn      = regexp.MustCompile(`^frame=\s*[^\s]+\s+fps=\s*[^\s]+\s+q=\s*[^\s]+\s+L?size=\s*[^\s]+\s+time=\s*[^\s]+\s+bitrate=\s*[^\s]+\s+speed=\s*[^\s]+$`)
	finalStatsPtrn = regexp.MustCompile(`^video:[^\s]+\s+audio:[^\s]+\s+subtitle:[^\s]+\s+other\s+streams:[^\s]+\s+global\s+headers:[^\s]+\s+muxing\s+overhead:\s+[^\s]+$`)
	repeatPtrn     = regexp.MustCompile(`^Last message repeated`)
)

type filterReader struct {
	scanner  *bufio.Scanner
	patterns []*regexp.Regexp
	text     string
	pattern  *regexp.Regexp
	err      error
}

func newFilterReader(reader io.Reader, patterns ...*regexp.Regexp) *filterReader {
	fr := &filterReader{
		scanner:  bufio.NewScanner(reader),
		patterns: patterns,
	}
	fr.scanner.Split(scanLines)
	return fr
}

func (fr *filterReader) Scan() bool {
	ok := fr.scanner.Scan()
	fr.pattern = nil
	fr.text = strings.TrimSpace(fr.scanner.Text())
	for _, pattern := range fr.patterns {
		if pattern.MatchString(fr.text) {
			fr.pattern = pattern
		}
	}
	fr.err = fr.scanner.Err()
	if !ok && fr.err == nil {
		fr.err = io.EOF
	}
	return ok
}

func (fr *filterReader) Pattern() *regexp.Regexp {
	return fr.pattern
}

func (fr *filterReader) Text() string {
	return fr.text
}

func (fr *filterReader) Bytes() []byte {
	return []byte(fr.text)
}

func (fr *filterReader) Err() error {
	return fr.err
}

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		if i < len(data)-2 {
			return i + 1, data[0:i], nil
		}
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		if data[i-1] == '\r' {
			return i + 1, data[0 : i-1], nil
		}
		return i + 1, data[0:i], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		if data[len(data)-1] == '\r' {
			return len(data), data[0 : len(data)-1], nil
		}
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}
