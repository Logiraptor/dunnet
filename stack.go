package main

import (
	"fmt"
)

type stacker struct {
	controller *controller
	stackLen   int
	cons       func() *controller
}

func NewStacker(cons func() *controller) *stacker {
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
		if s.controller.d.dead {
			return "<DEAD>\n>"
		}
		return s.controller.Send(msg)
	}
}

func (s *stacker) pop() string {
	if s.controller.d.dead {
		s.Close()
		s.Start()
	}

	if s.stackLen == 0 {
		return "nothing to pop\n>"
	}
	s.stackLen--
	stackName := fmt.Sprintf("restore dunnet/stack-%d", s.stackLen)
	s.controller.Send(stackName)
	return stackName + "\n>"
}

func (s *stacker) push() string {
	stackName := fmt.Sprintf("save dunnet/stack-%d", s.stackLen)
	s.Send(stackName)
	s.stackLen++
	return stackName + "\n>"
}

func (s *stacker) Close() {
	s.controller.Close()
}
