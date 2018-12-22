package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

type Controller interface {
	Start() string
	Send(string) string
	Close()
}

type controller struct {
	d     *dunnet
	read  *bufio.Scanner
	write io.Writer
}

func NewController() *controller {
	d := startDunnet()
	reader := bufio.NewScanner(d.output)

	promptRegexp := regexp.MustCompile(`(>|\$|login:|[pP]assword:|ftp>)`)

	reader.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		match := promptRegexp.FindIndex(data)
		if match != nil {
			i := match[1]
			return i, data[:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	return &controller{d: d, read: reader, write: d.input}
}

func (c *controller) Start() string {
	return c.nextOutput()
}

func (c *controller) Close() {
	c.d.close()
}

func (c *controller) Send(msg string) string {
	_, err := fmt.Fprintln(c.write, msg)
	if err != nil {
		panic(err)
	}
	return c.nextOutput()
}

func (c *controller) nextOutput() string {
	out := make(chan string)
	go func() {
		if !c.read.Scan() {
			panic("EOF")
		}

		if err := c.read.Err(); err != nil {
			panic(err)
		}

		out <- c.read.Text()
	}()

	select {
	case result := <-out:
		return result
	case <-time.After(time.Second):
		return "ERROR: could not scan response after 1 second"
	}
}
