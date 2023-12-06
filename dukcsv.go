package dukcsv

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Reader provides a random-access pattern reader to a CSV file on disk.
type Reader struct {
	loc     string
	offsets []int64
	file    *os.File
	header  []string
}

// NewReader creates a new reader for the file at the given location. The
// hasHeader flag will determine whether the first line will be considered a
// header row with column names. The header row is not included in the line
// count or line indexing.
func NewReader(loc string, hasHeader bool) (*Reader, error) {
	rdr := &Reader{
		loc:     loc,
		offsets: make([]int64, 0),
	}

	var err error
	rdr.file, err = os.Open(rdr.loc)
	if err != nil {
		return nil, err
	}

	// initialize offset values per line
	buf := make([]byte, 1024*1024)
	var q bool
	var c int64
	for {
		l, err := rdr.file.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			break
		}

		for i := 0; i < l; i++ {
			c++
			switch rune(buf[i]) {
			case '"':
				q = !q
			case '\n':
				if !q {
					rdr.offsets = append(rdr.offsets, c)
				}
			}
		}
	}

	// adds the last line, ignores if empty
	if rdr.offsets[len(rdr.offsets)-1] != c {
		rdr.offsets = append(rdr.offsets, c)
	}

	firstLine, err := rdr.read(0)
	if err != nil {
		return nil, err
	}
	if hasHeader {
		rdr.header = firstLine
	}

	return rdr, nil
}

// Close closes the associated file.
func (rdr *Reader) Close() error {
	return rdr.file.Close()
}

// HasHeader indicates whether the CSV was read in with a header row. Note:
// the number of rows and the line indexes do not include the header row.
func (rdr *Reader) HasHeader() bool {
	return rdr.header != nil
}

// Read reads the given from the CSV file and returns the values mapped by
// column name; if there was no header, the column name will be the number of
// the column, surrounded by square brackets.
func (rdr *Reader) Read(line int64) (map[string]string, error) {
	v, err := rdr.read(line)
	if err != nil {
		return nil, err
	}

	nCols := len(v)
	header := rdr.header
	if header == nil {
		header = make([]string, nCols)
		for i := range header {
			header[i] = fmt.Sprintf("[%d]", i+1)
		}
	}

	if nCols != len(header) {
		return nil, fmt.Errorf(
			"no. of cols on line %d doesn't match no. of cols in header: %d vs %d",
			line,
			nCols,
			len(v),
		)
	}

	cols := make(map[string]string, nCols)
	for i, h := range header {
		cols[h] = v[i]
	}

	return cols, nil
}

func (rdr *Reader) read(line int64) ([]string, error) {
	idx := line
	nLines := len(rdr.offsets) - 1
	if !rdr.HasHeader() {
		idx--
		nLines++
	}

	if line >= int64(nLines) {
		return nil, fmt.Errorf(
			"line %d is past the capacity (%d) of the CSV file",
			line,
			nLines,
		)
	}

	var offset int64
	if idx >= 0 {
		offset = rdr.offsets[idx]
	}
	lb := make([]byte, rdr.offsets[idx+1]-offset)
	if _, err := rdr.file.ReadAt(lb, offset); err != nil {
		return nil, err
	}

	return SplitCSVLine(lb), nil
}

// Count returns the number of lines in the CSV file, not including the header,
// if this file has one.
func (rdr *Reader) Count() int64 {
	nLines := len(rdr.offsets) - 1
	if !rdr.HasHeader() {
		nLines++
	}

	return int64(nLines)
}

// SplitCSVLine splits the single line into separate values, respecting quotes.
// CRLF at the end of lines are discarded.
func SplitCSVLine(line []byte) []string {
	items := make([]string, 0)
	buf := make([]byte, 0, 32*1024)
	var q bool
	for _, r := range line {
		switch r {
		case '"':
			q = !q
			if !q {
				buf = append(buf, r)
			}
		case ',':
			if !q {
				items = append(items, string(buf))
				buf = buf[:0]
			}
		case '\r', '\n':
			if q {
				buf = append(buf, r)
			}
		default:
			buf = append(buf, r)
		}
	}

	items = append(items, string(buf))

	return items
}
