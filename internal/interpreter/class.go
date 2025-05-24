package interpreter

import "codecrafters-interpreter-go/internal/ast"

// LoxClass represents a class in the interpreter.
// It implements the ast.Value interface.
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

// LoxClassInstance is the runtime representation of an instance of a Lox class
type LoxClassInstance struct {
	class *LoxClass
}

// newLoxClassInstance creates a new LoxClassInstance for the given class.
func newLoxClassInstance(class *LoxClass) *LoxClassInstance {
	return &LoxClassInstance{class: class}
}

// String returns a string representation of the class instance.
func (i *LoxClassInstance) String() string {
	return i.class.String() + " instance"
}

func (i *LoxClassInstance) GetType() ast.ValueType {
	return ast.ClassInstanceType
}

func (i *LoxClassInstance) IsEqualTo(other ast.Value) bool {
	if other == nil || other.GetType() != i.GetType() {
		return false
	}
	otherInstance, ok := other.(*LoxClassInstance)
	if !ok {
		return false
	}
	return i.class.IsEqualTo(otherInstance.class)
}

func (i *LoxClassInstance) IsTruthy() bool {
	return true
}
