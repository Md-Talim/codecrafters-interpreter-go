package interpreter

import "codecrafters-interpreter-go/internal/ast"

// LoxClass represents a class in the interpreter.
// It implements the ast.Value interface and the LoxCallable interface.
type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

// newLoxClass creates a new LoxClass instance.
func newLoxClass(name string, methods map[string]*LoxFunction, superclass *LoxClass) *LoxClass {
	return &LoxClass{name: name, methods: methods, superclass: superclass}
}

// String returns a string representation of the class.
func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

// call creates a new instance of the class.
func (c *LoxClass) call(interperter *Interpreter, arguments []ast.Value) (ast.Value, error) {
	instance := newLoxClassInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(interperter, arguments)
	}
	return instance, nil
}

// findMethod searches for a method with the given name in the class's methods map.
// It returns the method as a LoxFunction if found, or nil if not found.
func (c *LoxClass) findMethod(methodName string) *LoxFunction {
	if method, hasMethod := c.methods[methodName]; hasMethod {
		return method
	}
	if c.superclass != nil {
		return c.superclass.findMethod(methodName)
	}
	return nil
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
