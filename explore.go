package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

var directions = [...]string{
	"n", "ne", "nw", "w", "e", "se", "sw", "s", "u", "d",
}

type location struct {
	name     string
	fullText string
}

func (l location) String() string {
	return l.name
}

type edge struct {
	from      location
	direction string
	to        location
}

func (e edge) String() string {
	return fmt.Sprintf("%s %s %s.", e.from, e.direction, e.to)
}

type explorer struct {
	stacker *stacker
	paths   []edge
	visited map[location]struct{}
}

func (e *explorer) Start() string {
	return e.stacker.Start()
}

func (e *explorer) Send(msg string) string {
	if msg == "explore" {
		return e.explore()
	}
	return e.stacker.Send(msg)
}

func (e *explorer) Close() {
	e.stacker.Close()
}

func NewExplorer(stacker *stacker) *explorer {
	return &explorer{
		stacker: stacker,
		visited: make(map[location]struct{}),
	}
}

func (e *explorer) explore() string {
	e.visited = make(map[location]struct{})
	e.tryDirections()

	output := new(bytes.Buffer)

	for loc := range e.visited {
		fmt.Fprintln(output, "Visited", loc.name)
	}
	fmt.Fprintln(output, "Found", len(e.paths), "edges")

	outFile, err := ioutil.TempFile("dunnet", "*.dot")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	fmt.Fprint(outFile, "digraph {\n")
	for _, e := range e.paths {
		fmt.Fprintf(output, "From %q to %q via %q\n", e.from, e.to, e.direction)
		fmt.Fprintf(outFile, "\t%q -> %q [label=%q];\n", e.from, e.to, e.direction)
	}
	fmt.Fprint(outFile, "}\n")

	fmt.Fprintf(output, "dot file generated at: %s\n", "./dunnet/map.png")
	fmt.Fprint(output, "\n>")

	exec.Command("dot", "-Tpng", outFile.Name(), "-o", "dunnet/map.png").Run()

	return output.String()
}

func (e *explorer) tryDirections() {

	e.Send("push")
	currentLocation := e.getLocation()
	e.markVisited(currentLocation)

	edges := e.findEdges()

	for _, ed := range edges {
		e.paths = append(e.paths, ed)
		if !e.isVisited(ed.to) {
			e.Send("push")

			e.Send(ed.direction)
			e.tryDirections()

			e.Send("pop")
		}
	}
	e.Send("pop")
}

func (e *explorer) getLocation() location {
	return parseLocation(e.Send("l"))
}

func (e *explorer) findEdges() []edge {
	var edges []edge
	for _, dir := range directions {

		e.Send("push")

		startingPosition := e.getLocation()
		e.Send(dir)
		endingPosition := e.getLocation()

		e.Send("pop")

		if startingPosition == endingPosition {
			continue
		}

		ed := edge{
			from:      startingPosition,
			direction: dir,
			to:        endingPosition,
		}

		edges = append(edges, ed)
	}

	return edges
}

func (e *explorer) markVisited(l location) {
	e.visited[l] = struct{}{}
}

func (e *explorer) isVisited(l location) bool {
	_, ok := e.visited[l]
	return ok
}

func parseLocation(lookOutput string) location {
	firstNewline := strings.Index(lookOutput, "\n")
	if firstNewline == -1 {
		return location{
			name:     lookOutput,
			fullText: "",
		}
	}
	return location{
		name:     lookOutput[:firstNewline],
		fullText: lookOutput[firstNewline+1:],
	}
}
