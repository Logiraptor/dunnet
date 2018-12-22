package main

import (
	"fmt"
	"io"
)

type Controller interface {
	Start() string
	Send(string) string
	Close()
	IsDead() bool
}

type controller struct {
	d     *dunnet
	read  *promptScanner
	write io.Writer
}

func (c *controller) IsDead() bool {
	return c.d.dead
}

func NewController() *controller {
	d := startDunnet()
	reader := newScanner(d.output)

	return &controller{d: d, read: reader, write: d.input}
}

func (c *controller) Start() string {
	return c.read.NextOutput()
}

func (c *controller) Close() {
	c.d.close()
}

func (c *controller) Send(msg string) string {
	_, err := fmt.Fprintln(c.write, msg)
	if err != nil {
		panic(err)
	}
	return c.read.NextOutput()
}
