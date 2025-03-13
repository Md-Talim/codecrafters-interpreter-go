package lox

import (
	"fmt"
	"os"

	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/parser"
	"codecrafters-interpreter-go/internal/scanner"
)

type Lox struct {
	hadError bool
}

func (l *Lox) Tokenize(source string) {
	scanner := scanner.NewScanner(source, l.error)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func (l *Lox) Parse(source string) {
	scanner := scanner.NewScanner(source, l.error)
	tokens := scanner.ScanTokens()
	parser := parser.NewParser(tokens)

	expr, err := parser.Parse()
	if err != nil {
		l.parseError(err.Token(), err.Error())
		return
	}

	printer := &ast.AstPrinter{}
	fmt.Println(printer.Print(expr))
}

func (l *Lox) HadError() bool {
	return l.hadError
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) parseError(token ast.Token, message string) {
	if token.Type == ast.EofToken {
		l.report(token.Line, " at end", message)
	} else {
		l.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}
