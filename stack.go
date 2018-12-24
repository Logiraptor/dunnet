package main

import (
	"fmt"
	"time"
)

type stacker struct {
	controller Controller
	stackLen   int
	cons       func() Controller
}

func (s *stacker) IsDead() bool {
	return s.controller.IsDead()
}

func NewStacker(cons func() Controller) *stacker {
	return &stacker{cons: cons}
}

func (s *stacker) Start() string {
	s.controller = s.cons()
	return s.controller.Start()
}

func (s *stacker) Send(msg string) string {
	if msg == "push" {
		return s.push()
	} else if msg == "pop" {
		return s.pop()
	} else {
		if s.controller.IsDead() {
			return "<DEAD>\n>"
		}
		return s.controller.Send(msg)
	}
}

func (s *stacker) pop() string {
	// Sleeping here is necessary, because saving doesn't always happen fast enough.
	time.Sleep(time.Millisecond * 10)

	if s.controller.IsDead() {
		s.Close()
		s.Start()
	}

	if s.stackLen == 0 {
		return "nothing to pop\n>"
	}
	s.stackLen--
	stackName := fmt.Sprintf("restore generated/stack-%d", s.stackLen)
	s.controller.Send(stackName)
	return stackName + "\n>"
}

func (s *stacker) push() string {
	stackName := fmt.Sprintf("save generated/stack-%d", s.stackLen)
	s.Send(stackName)
	s.stackLen++
	return stackName + "\n>"
}

func (s *stacker) Close() {
	s.controller.Close()
}
