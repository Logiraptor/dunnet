package main

import (
	"os"
	"io"
	"bufio"
	"fmt"
	"io/ioutil"
	"bytes"
)

func main() {
	fileName := os.Args[1]
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	rewrite(bytes.NewBuffer(buf), out)
}

func rewrite(input io.Reader, output io.Writer) {
	lines := bufio.NewScanner(input)

	for lines.Scan() {
		switch lines.Text() {
		case "push":
			discardUntilPop(lines, 0)
		case "l", "i":
		default:
			_, err := fmt.Fprintln(output, lines.Text())
			if err != nil {
				panic(err)
			}
		}
	}

	if err := lines.Err(); err != nil {
		panic(err)
	}
}

func discardUntilPop(scanner *bufio.Scanner, depth int) {
	for scanner.Scan() {
		switch scanner.Text() {
		case "push":
			discardUntilPop(scanner, depth+1)
		case "pop":
			return
		}
	}
}
