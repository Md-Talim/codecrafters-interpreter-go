package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
)

// LoxCallable defines the interface for callable lox objects.
type LoxCallable interface {
	arity() int
	call(interperter *Interpreter, arguments []ast.Value) ast.Value
}

// LoxFunction represents a function in Lox.
// It implements the LoxCallable & ast.Value interfaces.
type LoxFunction struct {
	declaration ast.FunctionStmt
}

// newLoxFunction creates a new LoxFunction instance.
func newLoxFunction(declaration ast.FunctionStmt) *LoxFunction {
	return &LoxFunction{declaration: declaration}
}

// arity returns the number of parameters the function takes.
func (f *LoxFunction) arity() int {
	return len(f.declaration.Params)
}

// call executes the function with the given arguments.
func (f *LoxFunction) call(interperter *Interpreter, arguments []ast.Value) ast.Value {
	// Create a new environment for the function call
	env := newEnvironment(interperter.globals)
	for i, param := range f.declaration.Params {
		env.define(param.Lexeme, arguments[i])
	}

	// Execute the function body in the new environment
	interperter.executeBlock(f.declaration.Body, env)

	return nil
}

// String returns a string representation of the function.
func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}

// GetType implements ast.Value.
func (f *LoxFunction) GetType() ast.ValueType {
	return ast.FunctionType
}

// IsEqualTo implements ast.Value.
func (f *LoxFunction) IsEqualTo(other ast.Value) bool {
	_, ok := other.(*LoxFunction)
	return ok
}

// IsTruthy implements ast.Value.
func (f *LoxFunction) IsTruthy() bool {
	return true
}
