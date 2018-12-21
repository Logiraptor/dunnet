package main

import (
	"bufio"
	"fmt"
	"io"
)

type Controller interface {
	Start() string
	Send(string) string
	Close()
}

type controller struct {
	d     *dunnet
	read  *bufio.Reader
	write io.Writer
}

func NewController() *controller {
	inputPipe, input := io.Pipe()
	output, outputPipe := io.Pipe()
	d := startDunnet(outputPipe, inputPipe)
	reader := bufio.NewReader(output)

	return &controller{d: d, read: reader, write: input}
}

func (c *controller) Start() string {
	result, err := c.read.ReadString('>')
	if err != nil {
		panic(err)
	}
	return result
}

func (c *controller) Close() {
	c.d.close()
}

func (c *controller) Send(msg string) string {
	_, err := fmt.Fprintln(c.write, msg)
	if err != nil {
		panic(err)
	}
	result, err := c.read.ReadString('>')
	if err != nil {
		panic(err)
	}
	return result
}
