package interpreter

import "codecrafters-interpreter-go/internal/ast"

type Environment struct {
	enclosing *Environment
	values    map[string]ast.Value
}

func newEnvironment(enclosing *Environment) *Environment {
	if enclosing == nil {
		return &Environment{enclosing: nil, values: make(map[string]ast.Value)}
	}
	return &Environment{enclosing: enclosing, values: make(map[string]ast.Value)}
}

func (e *Environment) assign(name ast.Token, value ast.Value) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.define(name.Lexeme, value)
		return nil
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return nil
	}
	return newRuntimeError(name.Line, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) define(name string, value ast.Value) {
	e.values[name] = value
}

func (e *Environment) get(name ast.Token) (ast.Value, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return ast.NewNilValue(), newRuntimeError(name.Line, "Undefined variable '"+name.Lexeme+"'.")
}
