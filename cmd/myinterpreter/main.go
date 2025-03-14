package main

import (
	"fmt"
	"os"

	"codecrafters-interpreter-go/pkg/lox"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	lox := lox.Lox{}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		lox.Tokenize(string(fileContents))
	case "parse":
		lox.Parse(string(fileContents))
	case "evaluate":
		lox.Evaluate(string(fileContents))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	if lox.HadError() {
		os.Exit(65)
	}
}
