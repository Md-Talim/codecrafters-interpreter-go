package interpreter

import "codecrafters-interpreter-go/internal/ast"

// LoxCallable defines the interface for callable lox objects.
type LoxCallable interface {
	arity() int
	call(interperter *Interpreter, arguments []ast.Value) (ast.Value, error)
}
