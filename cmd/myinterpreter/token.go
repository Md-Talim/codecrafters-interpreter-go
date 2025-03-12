package main

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	if t.Literal == nil {
		return fmt.Sprintf("%v %v null", tokenTypeNames[t.Type], t.Lexeme)
	}
	return fmt.Sprintf("%v %v %v", tokenTypeNames[t.Type], t.Lexeme, t.Literal)
}
