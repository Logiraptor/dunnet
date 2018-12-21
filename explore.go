package main

import (
	"bytes"
	"fmt"
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

	fmt.Fprintln(output, "Found", len(e.paths), "edges")
	for _, e := range e.paths {
		fmt.Fprintf(output, "From %6q to %6q via %q\n", e.from, e.to, e.direction)
	}
	fmt.Fprintln(output, "\n>")

	return output.String()
}

func (e *explorer) tryDirections() {
	fmt.Println(e.visited)

	e.stacker.Send("push")
	currentLocation := e.getLocation()
	e.markVisited(currentLocation)

	fmt.Println("visited", currentLocation.name)

	edges := e.findEdges()

	fmt.Println("found", len(edges), "edges", edges)

	for _, ed := range edges {
		if !e.isVisited(ed.to) {
			e.stacker.Send("push")

			e.paths = append(e.paths, ed)
			e.stacker.Send(ed.direction)
			e.tryDirections()

			e.stacker.Send("pop")
		}
	}
	e.stacker.Send("pop")
}

func (e *explorer) getLocation() location {
	return parseLocation(e.stacker.Send("l"))
}

func (e *explorer) findEdges() []edge {
	var edges []edge
	for _, dir := range directions {
		e.stacker.Send("push")
		startingPosition := e.getLocation()

		e.stacker.Send(dir)
		endingPosition := e.getLocation()

		ed := edge{
			from:      startingPosition,
			direction: dir,
			to:        endingPosition,
		}

		if ed.from != ed.to {
			edges = append(edges, ed)
			fmt.Println("non-loop", dir)
		} else {
			fmt.Println("loop", dir)
		}

		e.stacker.Send("pop")
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
	return location{
		name:     lookOutput[:firstNewline],
		fullText: lookOutput,
	}
}
