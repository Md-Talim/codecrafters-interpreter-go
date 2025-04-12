package lox

import (
	"fmt"
	"os"

	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/interpreter"
	"codecrafters-interpreter-go/internal/parser"
	"codecrafters-interpreter-go/internal/scanner"
)

func Tokenize(source string) {
	scanner := scanner.NewScanner(source)
	tokens, hadError := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}

	if hadError {
		os.Exit(65)
	}
}

func Parse(source string) {
	parser := parser.NewParser(source)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	printer := ast.NewAstPrinter()
	expr.Accept(printer)
}

func Evaluate(source string) {
	interpreter := &interpreter.Interpreter{}
	value, err := interpreter.Interpret(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
	}
	fmt.Println(value)
}
