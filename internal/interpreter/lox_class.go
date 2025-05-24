package interpreter

import "codecrafters-interpreter-go/internal/ast"

// LoxClass represents a class in the interpreter.
// It implements the ast.Value interface and the LoxCallable interface.
type LoxClass struct {
	name string
}

// newLoxClass creates a new LoxClass instance.
func newLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

// String returns a string representation of the class.
func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) arity() int {
	return 0
}

// call creates a new instance of the class.
func (c *LoxClass) call(interperter *Interpreter, arguments []ast.Value) ast.Value {
	return newLoxClassInstance(c)
}

// GetType returns the type of the class.
func (c *LoxClass) GetType() ast.ValueType {
	return ast.ClassType
}

// IsEqualTo checks if the class is equal to another value.
func (c *LoxClass) IsEqualTo(other ast.Value) bool {
	if other == nil || other.GetType() != c.GetType() {
		return false
	}
	return c.name == other.(*LoxClass).name
}

// IsTruthy returns true for all class instances.
func (c *LoxClass) IsTruthy() bool {
	return true
}
