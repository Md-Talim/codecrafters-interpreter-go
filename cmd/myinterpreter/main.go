package main

import (
	"fmt"
	"os"
)

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	EOF         TokenType = "EOF"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
}

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

	tokens := []Token{}

	for _, char := range fileContents {
		switch char {
		case '(':
			tokens = append(tokens, Token{LEFT_PAREN, "(", nil})
		case ')':
			tokens = append(tokens, Token{RIGHT_PAREN, ")", nil})
		}
	}

	tokens = append(tokens, Token{EOF, "", nil})

	for _, token := range tokens {
		var literalPrintValue string

		if token.literal == nil {
			literalPrintValue = "null"
		} else {
			literalPrintValue = fmt.Sprintf("%v", token.literal)
		}
		fmt.Printf("%s %s %s\n", token.tokenType, token.lexeme, literalPrintValue)
	}
}
