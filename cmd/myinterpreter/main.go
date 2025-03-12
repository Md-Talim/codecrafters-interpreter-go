package main

import (
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

var lox = Lox{}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	lox.tokenize(string(fileContents))

	if lox.hadError {
		os.Exit(65)
	}
}

func (l *Lox) tokenize(source string) {
	scanner := NewScanner(source)
	scanner.scanTokens()

	for _, token := range scanner.tokens {
		fmt.Println(token)
	}
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}
