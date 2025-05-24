package interpreter

import "codecrafters-interpreter-go/internal/ast"

// LoxClassInstance is the runtime representation of an instance of a Lox class.
// It implements the ast.Value interface.
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

// GetType returns the type of the class instance.
func (i *LoxClassInstance) GetType() ast.ValueType {
	return ast.ClassInstanceType
}

// IsEqualTo checks if the class instance is equal to another instance.
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

// IsTruthy returns true for all class instances.
func (i *LoxClassInstance) IsTruthy() bool {
	return true
}
