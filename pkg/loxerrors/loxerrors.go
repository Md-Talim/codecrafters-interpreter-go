package loxerrors

import (
	"fmt"
)

type ErrorType int

const (
	LexicalError ErrorType = iota
	ParseError
	RuntimeError
)

type LoxError struct {
	Type    ErrorType
	Message string
	Line    int
}

func (e *LoxError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}

func NewLexicalError(line int, message string) *LoxError {
	return &LoxError{Type: LexicalError, Message: message, Line: line}
}

func NewParseError(line int, message string) *LoxError {
	return &LoxError{Type: ParseError, Message: message, Line: line}
}

func NewRuntimeError(line int, message string) *LoxError {
	return &LoxError{Type: RuntimeError, Message: message, Line: line}
}
