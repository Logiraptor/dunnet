package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/profile"
	"io"
	"os"
	"os/exec"
	"time"
)

type dunnet struct {
	cmd *exec.Cmd
}

func startDunnet(output io.Writer, input io.Reader) *dunnet {
	command := exec.Command("emacs", "-batch", "-l", "dunnet")
	d := &dunnet{
		cmd: command,
	}

	out, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}
	in, err := command.StdinPipe()
	if err != nil {
		panic(err)
	}
	err = command.Start()
	if err != nil {
		panic(err)
	}

	go io.Copy(output, out)
	go func() {
		defer d.cmd.Wait()
		io.Copy(in, input)
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

func main() {
	go func() {
		start := profile.Start(profile.CPUProfile)
		time.Sleep(15 * time.Second)
		start.Stop()
	}()

	controller := NewExplorer(NewStacker(NewRewinder(func() Controller {
		return NewController()
	})))

	runInteractive(controller)
}
