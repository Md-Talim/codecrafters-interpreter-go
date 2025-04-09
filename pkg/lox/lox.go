package lox

import (
	"fmt"
	"os"

	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/interpreter"
	"codecrafters-interpreter-go/internal/parser"
	"codecrafters-interpreter-go/internal/scanner"
	"codecrafters-interpreter-go/pkg/loxerrors"
)

func Tokenize(source string) {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	for _, token := range tokens {
		fmt.Println(token)
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
		if loxErr, ok := err.(*loxerrors.LoxError); ok {
			fmt.Fprintln(os.Stderr, loxErr)
			if loxErr.Type == loxerrors.RuntimeError {
				os.Exit(70)
			}
			os.Exit(65)
		}
	}
	fmt.Println(value)
}
