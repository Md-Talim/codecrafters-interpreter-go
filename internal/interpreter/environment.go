package interpreter

import "codecrafters-interpreter-go/internal/ast"

// Environment represents a single environment in the interpreter.
type Environment struct {
	enclosing *Environment
	values    map[string]ast.Value
}

// newEnvironment creates a new environment with the given enclosing environment.
func newEnvironment(enclosing *Environment) *Environment {
	if enclosing == nil {
		return &Environment{enclosing: nil, values: make(map[string]ast.Value)}
	}
	return &Environment{enclosing: enclosing, values: make(map[string]ast.Value)}
}

// ancestor traverses up the environment chain to find the ancestor at the specified distance.
// The distance is the number of environments to traverse up.
func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for range distance {
		environment = environment.enclosing
	}
	return environment
}

// assign assigns a value to a variable in the environment.
// If the variable is not found in the current environment, it checks the enclosing environment.
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

// assignAt assigns a value to a variable at a specific distance in the environment chain.
func (e *Environment) assignAt(distance int, name ast.Token, value ast.Value) {
	e.ancestor(distance).values[name.Lexeme] = value
}

// define defines a new variable in the environment.
func (e *Environment) define(name string, value ast.Value) {
	e.values[name] = value
}

// get retrieves a value from the environment.
// If the variable is not found in the current environment, it checks the enclosing environment.
func (e *Environment) get(name ast.Token) (ast.Value, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return ast.NewNilValue(), newRuntimeError(name.Line, "Undefined variable '"+name.Lexeme+"'.")
}

// getAt retrieves a value from a specific distance in the environment chain.
func (e *Environment) getAt(distance int, name string) (ast.Value, error) {
	return e.ancestor(distance).values[name], nil
}
