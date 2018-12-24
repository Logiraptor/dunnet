package main

import "fmt"

func (e *explorer) navigate(msg string) string {
	destinationName, err := extractArg(msg)
	if err != nil {
		return err.Error()
	}

	var destination location
	for visit := range e.visited {
		if visit.name == destinationName {
			destination = visit
		}
	}

	if result, ok := e.navigateTo(e.getLocation(), destination, nil); ok {
		for i := len(result) - 1; i >= 0 ; i-- {
			e.stacker.Send(result[i])
		}
	} else {
		fmt.Println("Couldn't find a path.️")
	}

	return e.stacker.Send("l")
}

func (e *explorer) navigateTo(from, to location, visited []location) ([]string, bool) {
	if from == to {
		return nil, true
	}

pathLoop:
	for _, path := range e.paths {
		if path.from == from {
			for _, visit := range visited {
				if visit == path.to {
					continue pathLoop
				}
			}

			if rest, ok := e.navigateTo(path.to, to, append(visited, path.to)); ok {
				return append(rest, path.direction), true
			}
		}
	}

	return nil, false
}
