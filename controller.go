package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

	reader.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		i := bytes.IndexByte(data, '>')
		if i != -1 {
			return i + 1, data[:i+1], nil
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
	// Sleeping here is necessary, because saving doesn't happen fast enough.
	time.Sleep(time.Millisecond)
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
		return "dead"
	}
}
