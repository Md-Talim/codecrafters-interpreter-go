package loxerrors

type ErrorType int

const (
	LexicalError ErrorType = iota
	ParseError
	RuntimeError
)

var errorTypeNames = map[ErrorType]string{
	LexicalError: "Lexical",
	ParseError:   "Parse",
	RuntimeError: "Runtime",
}
