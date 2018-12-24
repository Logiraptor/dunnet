package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"errors"
)

type dunnet struct {
	cmd    *exec.Cmd
	output io.Reader
	input  io.Writer
	dead   bool
}

func startDunnet() *dunnet {
	command := exec.Command("emacs", "-batch", "-l", "dunnet")
	out, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}
	in, err := command.StdinPipe()
	if err != nil {
		panic(err)
	}

	d := &dunnet{
		cmd:    command,
		output: out,
		input:  in,
	}

	err = command.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		d.cmd.Wait()
		d.dead = true
	}()

	return d
}

func (d *dunnet) close() {
	_ = d.cmd.Process.Kill()
}

func runInteractive(c Controller) {
	scanner := bufio.NewScanner(os.Stdin)
	output := os.Stdout

	fmt.Fprint(output, c.Start())

	for scanner.Scan() {
		fmt.Fprint(output, c.Send(scanner.Text()))
	}

	c.Close()
}

func extractArg(msg string) (string, error) {
	parts := strings.SplitN(msg, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("Syntax: (record|replay) FILENAME\n>")
	}
	return parts[1], nil
}

func main() {
	//go func() {
	//	start := profile.Start(profile.CPUProfile)
	//	time.Sleep(15 * time.Second)
	//	start.Stop()
	//}()

	controller := NewExplorer(NewReplay(NewComboCracker(NewStacker(func() Controller {
		return NewController()
	}))))

	runInteractive(controller)
}
