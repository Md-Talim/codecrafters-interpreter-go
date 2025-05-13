package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"time"
)

// ClockFunction is a native function that returns the current time.
type ClockFunction struct{}

// NewClockFunction creates a new instance of ClockFunction.
func NewClockFunction() *ClockFunction {
	return &ClockFunction{}
}

// arity returns the number of arguments the function takes.
func (c *ClockFunction) arity() int {
	return 0
}

// call executes the function with the given arguments.
func (c *ClockFunction) call(interpreter *Interpreter, arguments []ast.Value) ast.Value {
	currentTime := time.Now().Unix()
	return ast.NewNumberValue(float64(currentTime))
}

// String for ClockFunction
func (c *ClockFunction) String() string {
	return "<native fn>"
}

// GetType for ClockFunction (returns NativeFunctionType)
func (c *ClockFunction) GetType() ast.ValueType {
	return ast.NativeFunctionType
}

// IsEqualTo for ClockFunction. Native functions can be considered equal if they are the same type/instance.
func (n *ClockFunction) IsEqualTo(other ast.Value) bool {
	_, ok := other.(*ClockFunction)
	return ok
}

// IsTruthy for ClockFunction. Native functions are always truthy.
func (n *ClockFunction) IsTruthy() bool {
	return true
}
