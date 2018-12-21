package main

import "fmt"

type stacker struct {
	r     *rewinder
	stack []int
}

func NewStacker(r *rewinder) *stacker {
	return &stacker{r: r}
}

func (s *stacker) Start() string {
	return s.r.Start()
}

func (s *stacker) Send(msg string) string {
	if msg == "push" {
		s.stack = append(s.stack, 0)
		return "starting new stack frame\n>"
	} else if msg == "pop" {
		if len(s.stack) == 0 {
			return "nothing to pop\n>"
		}
		lastFrame := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		s.r.Send(fmt.Sprintf("rewind %d", lastFrame))
		return fmt.Sprintf("popped %d commands\n>", lastFrame)
	} else {
		if len(s.stack) > 0 {
			s.stack[len(s.stack)-1]++
		}
		return s.r.Send(msg)
	}
}

func (s *stacker) Close() {
	s.r.Close()
}
