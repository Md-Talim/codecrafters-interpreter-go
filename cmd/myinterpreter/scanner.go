package main

import (
	"fmt"
	"os"
)

func tokenize(fileContents []byte) ([]Token, bool) {
	tokens := []Token{}
	lineNumber := 1
	lexicalError := false
	length := len(fileContents)

	for i := 0; i < length; i++ {
		char := fileContents[i]

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
		case '\t':
			continue
		case ' ':
			continue
		case '=':
			if i+1 < length && fileContents[i+1] == '=' {
				tokens = append(tokens, Token{EQUAL_EQUAL, "==", nil})
				i++
			} else {
				tokens = append(tokens, Token{EQUAL, "=", nil})
			}
		case '!':
			if i+1 < length && fileContents[i+1] == '=' {
				tokens = append(tokens, Token{BANG_EQUAL, "!=", nil})
				i++
			} else {
				tokens = append(tokens, Token{BANG, "!", nil})
			}
		case '<':
			if i+1 < length && fileContents[i+1] == '=' {
				tokens = append(tokens, Token{LESS_EQUAL, "<=", nil})
				i++
			} else {
				tokens = append(tokens, Token{LESS, "<", nil})
			}
		case '>':
			if i+1 < length && fileContents[i+1] == '=' {
				tokens = append(tokens, Token{GREATER_EQUAL, ">=", nil})
				i++
			} else {
				tokens = append(tokens, Token{GREATER, ">", nil})
			}
		case '/':
			if i+1 < length && fileContents[i+1] == '/' {
				i++
				for i < length && fileContents[i] != '\n' {
					i++
				}
				lineNumber++
			} else {
				tokens = append(tokens, Token{SLASH, "/", nil})
			}
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", lineNumber, char)
			lexicalError = true
		}
	}

	tokens = append(tokens, Token{EOF, "", nil})

	return tokens, lexicalError
}
