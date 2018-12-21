package main

import (
	"strconv"
	"strings"
)

type rewinder struct {
	commands []string
	current  Controller
	cons     func() Controller
}

func NewRewinder(cons func() Controller) *rewinder {
	return &rewinder{commands: nil, current: cons(), cons: cons}
}

func (r *rewinder) Start() string {
	r.current = r.cons()
	return r.current.Start()
}

func (r *rewinder) Send(msg string) string {
	if strings.HasPrefix(msg, "rewind") {
		parts := strings.SplitN(msg, " ", 2)

		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return err.Error() + "\n"
		}

		r.retry(n)
		return "Backing up n steps\n>"
	}
	r.commands = append(r.commands, msg)
	return r.current.Send(msg)
}

func (r *rewinder) Close() {
	r.current.Close()
}

func (r *rewinder) retry(n int) {
	r.commands = r.commands[:len(r.commands)-n]
	r.Close()
	r.Start()
	for _, cmd := range r.commands {
		r.Send(cmd)
	}
}
