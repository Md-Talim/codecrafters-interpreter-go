package main

import (
	"fmt"
	"os"
)

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	STAR        TokenType = "STAR"
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
	lineNumber := 1
	lexicalError := false

	for _, char := range fileContents {
		switch char {
		case '(':
			tokens = append(tokens, Token{LEFT_PAREN, "(", nil})
		case ')':
			tokens = append(tokens, Token{RIGHT_PAREN, ")", nil})
		case '{':
			tokens = append(tokens, Token{LEFT_BRACE, "{", nil})
		case '}':
			tokens = append(tokens, Token{RIGHT_BRACE, "}", nil})
		case ',':
			tokens = append(tokens, Token{COMMA, ",", nil})
		case '.':
			tokens = append(tokens, Token{DOT, ".", nil})
		case '-':
			tokens = append(tokens, Token{MINUS, "-", nil})
		case '+':
			tokens = append(tokens, Token{PLUS, "+", nil})
		case ';':
			tokens = append(tokens, Token{SEMICOLON, ";", nil})
		case '*':
			tokens = append(tokens, Token{STAR, "*", nil})
		case '\n':
			lineNumber++
		default:
			fmt.Fprintf(os.Stderr, "[line %d]: Unexpected character: %c\n", lineNumber, char)
			lexicalError = true
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

	if lexicalError {
		os.Exit(65)
	}
}
