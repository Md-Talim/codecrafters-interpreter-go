package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
)

// LoxFunction represents a function in Lox.
// It implements the LoxCallable & ast.Value interfaces.
type LoxFunction struct {
	closure       *Environment
	declaration   ast.FunctionStmt
	isInitializer bool
}

// newLoxFunction creates a new LoxFunction instance.
func newLoxFunction(declaration ast.FunctionStmt, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{closure: closure, declaration: declaration, isInitializer: isInitializer}
}

// arity returns the number of parameters the function takes.
func (f *LoxFunction) arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) bind(instance *LoxClassInstance) *LoxFunction {
	environment := newEnvironment(f.closure)
	environment.define("this", instance)
	return newLoxFunction(f.declaration, environment, f.isInitializer)
}

// call executes the function with the given arguments.
func (f *LoxFunction) call(interperter *Interpreter, arguments []ast.Value) (ast.Value, error) {
	// Create a new environment for the function call
	// This environment is a child of the closure environment.
	env := newEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		env.define(param.Lexeme, arguments[i])
	}

	// Execute the function body in the new environment
	_, err := interperter.executeBlock(f.declaration.Body, env)
	if err != nil {
		if returnErr, ok := err.(*ReturnError); ok {
			// This is a return statement, not a runtime error.
			// The actual return value of the function is returnErr.value
			return returnErr.value, nil
		}
		return ast.NewNilValue(), err
	}

	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}
	return ast.NewNilValue(), nil
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
