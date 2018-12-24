package main

import (
	"regexp"
	"fmt"
)

var comboRegexp = regexp.MustCompile(`The combination is (\d+)\.`)

type comboCracker struct {
	inner Controller
	combo string
}

func NewComboCracker(inner Controller) *comboCracker {
	return &comboCracker{
		inner: inner,
	}
}

func (c *comboCracker) Start() string {
	return c.inner.Start()
}

func (c *comboCracker) Send(msg string) string {
	switch msg {
	case "crack combo":
		output := c.inner.Send("type foo.txt")
		c.combo = comboRegexp.FindStringSubmatch(output)[1]
		return fmt.Sprintf("%s\n => %q\n>", output, c.combo)
	case "enter combo":
		return c.inner.Send(c.combo)
	default:
		return c.inner.Send(msg)
	}
}

func (c *comboCracker) Close() {
	c.inner.Close()
}

func (c *comboCracker) IsDead() bool {
	return c.inner.IsDead()
}
