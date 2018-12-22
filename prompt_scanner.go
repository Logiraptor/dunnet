package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"time"
)

type promptScanner struct {
	lastData []byte
	reader   *bufio.Scanner
}

func newScanner(reader io.Reader) *promptScanner {
	promptScanner := &promptScanner{
		reader: bufio.NewScanner(reader),
	}
	promptScanner.reader.Split(promptScanner.scan)
	return promptScanner
}

var promptRegexp = regexp.MustCompile(`(>|\$|login:|[pP]assword:|ftp>|Username:|Enter it here:)`)

func (s *promptScanner) scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	s.lastData = data
	match := promptRegexp.FindIndex(data)
	if match != nil {
		i := match[1]
		return i, data[:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func (s *promptScanner) NextOutput() string {
	out := make(chan string)
	go func() {
		if !s.reader.Scan() {
			panic("EOF")
		}

		if err := s.reader.Err(); err != nil {
			panic(err)
		}

		out <- s.reader.Text()
	}()

	select {
	case result := <-out:
		return result
	case <-time.After(time.Second):
		log.Printf(`
I can't figure out where the prompt is in this data:

%s

Please update the prompt regexp.
`, s.lastData)

		return "Scanning taking too long."
	}
}
