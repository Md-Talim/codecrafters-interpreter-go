package interpreter

import "codecrafters-interpreter-go/internal/ast"

// LoxClassInstance is the runtime representation of an instance of a Lox class.
// It implements the ast.Value interface.
type LoxClassInstance struct {
	class  *LoxClass
	fields map[string]ast.Value
}

// newLoxClassInstance creates a new LoxClassInstance for the given class.
func newLoxClassInstance(class *LoxClass) *LoxClassInstance {
	return &LoxClassInstance{class: class, fields: make(map[string]ast.Value)}
}

// get retrieves a property from the class instance.
func (i *LoxClassInstance) get(name ast.Token) (ast.Value, error) {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value, nil
	}
	return nil, newRuntimeError(name.Line, "Undefined property '"+name.Lexeme+"'.")
}

// set sets a property on the class instance.
func (i *LoxClassInstance) set(name ast.Token, value ast.Value) {
	i.fields[name.Lexeme] = value
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
