package bazooka

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	reader *bufio.Reader
	line   []byte
	err    error
	dead   bool
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(reader)}
}

func (s *Scanner) Scan() bool {
	if s.dead {
		return false
	}
	line, err := s.reader.ReadBytes('\n')
	line = bytes.TrimSuffix(line, []byte{'\n'})

	// ReadBytes can return full or partial output even when it failed.
	// e.g. it can return a full entry and EOF.
	if err == nil || len(line) > 0 {
		s.line = line
		s.err = nil
		return true
	}

	if err != nil {
		s.line = nil
		s.dead = true
		if err != io.EOF {
			s.err = err
		}
	}
	return false
}

func (s *Scanner) Text() string {
	return string(s.line)
}

func (s *Scanner) Err() error {
	return s.err
}
