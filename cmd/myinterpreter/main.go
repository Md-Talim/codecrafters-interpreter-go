package main

import (
	"fmt"
	"os"
)

type TokenType string

const (
	LEFT_PAREN    TokenType = "LEFT_PAREN"
	RIGHT_PAREN   TokenType = "RIGHT_PAREN"
	LEFT_BRACE    TokenType = "LEFT_BRACE"
	RIGHT_BRACE   TokenType = "RIGHT_BRACE"
	COMMA         TokenType = "COMMA"
	DOT           TokenType = "DOT"
	MINUS         TokenType = "MINUS"
	PLUS          TokenType = "PLUS"
	SEMICOLON     TokenType = "SEMICOLON"
	STAR          TokenType = "STAR"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	SLASH         TokenType = "SLASH"
	EOF           TokenType = "EOF"
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

	tokens, lexicalError := tokenize(fileContents)

	for _, token := range tokens {
		var literalPrintValue string

		if token.literal == nil {
			literalPrintValue = "null"
		} else {
			literalPrintValue = fmt.Sprintf("%v", token.literal)
		}
		fmt.Printf("%s %s %s\n", token.tokenType, token.lexeme, literalPrintValue)
	}

	if lexicalError {
		os.Exit(65)
	}
}
