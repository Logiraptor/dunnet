package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type replay struct {
	inner    Controller
	commands []string
}

func (r *replay) IsDead() bool {
	return r.inner.IsDead()
}

func NewReplay(inner Controller) *replay {
	return &replay{inner: inner}
}

func (r *replay) Start() string {
	return r.inner.Start()
}

func (r *replay) Send(msg string) string {
	if strings.HasPrefix(msg, "record") {
		return r.saveFile(msg)
	} else if strings.HasPrefix(msg, "replay") {
		return r.loadFile(msg)
	}
	r.commands = append(r.commands, msg)
	return r.inner.Send(msg)
}

func (r *replay) Close() {
	r.inner.Close()
}

func (r *replay) saveFile(msg string) string {
	fileName, err := extractArg(msg)
	if err != nil {
		return err.Error()
	}

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, cmd := range r.commands {
		fmt.Fprintln(f, cmd)
	}
	return fmt.Sprintf("Recorded %d commands to %s\n>", len(r.commands), fileName)
}

func (r *replay) loadFile(msg string) string {
	fileName, err := extractArg(msg)
	if err != nil {
		return err.Error()
	}

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		r.Send(scanner.Text())
	}

	return fmt.Sprintf("Loaded %d commands from %s\n>", len(r.commands), fileName)
}
