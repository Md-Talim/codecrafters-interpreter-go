package interpreter

import "codecrafters-interpreter-go/internal/ast"

type Environment struct {
	values map[string]ast.Value
}

func (e *Environment) define(name string, value ast.Value) {
	e.values[name] = value
}

func (e *Environment) get(name ast.Token) (ast.Value, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	return ast.NewNilValue(), newRuntimeError(name.Line, "Undefined variable '"+name.Lexeme+"'.")
}
