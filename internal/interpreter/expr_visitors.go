package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"fmt"
)

// lookUpVariable retrieves the value of a variable from the environment.
// If the variable is not found in the current environment, it checks the global environment.
func (i *Interpreter) lookUpVariable(name ast.Token, expr ast.Expr) (ast.Value, error) {
	distance, ok := i.locals[expr]
	if ok {
		// If the variable is found in the local environment, return its value.
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		// If the variable is not found in the local environment, check the global environment.
		return i.globals.get(name)
	}
}

// VisitAssignExpr implements ast.AstVisitor.
// It evaluates the value of the assignment expression and assigns it to the variable name.
func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (ast.Value, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	// Check if the variable is found in the local environment.
	// If found, assign the value to it.
	distance, ok := i.locals[expr]
	if ok {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		// If not found in the local environment, assign the value to the global environment.
		// This is a fallback mechanism to ensure that the variable is assigned correctly.
		err = i.globals.assign(expr.Name, value)
	}
	return value, err
}

// VisitBinaryExpr implements ast.AstVisitor.
// It evaluates the left and right expressions and performs the binary operation.
func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (ast.Value, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	operator := expr.Operator.Type
	switch operator {
	case ast.PlusToken:
		return performAddition(left, right, expr.Operator.Line)
	case ast.EqualEqualToken, ast.BangEqualToken:
		return performEqualityOperation(left, right, operator), nil
	case ast.StarToken, ast.SlashToken, ast.MinusToken, ast.GreaterEqualToken, ast.GreaterToken, ast.LessEqualToken, ast.LessToken:
		return performNumericOperation(left, right, operator, expr.Operator.Line)
	default:
		return nil, newRuntimeError(expr.Operator.Line, "Unknown operator.")
	}
}

// VisitCallExpr implements ast.AstVisitor.
// It evaluates the callee and arguments, and calls the function if valid.
func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (ast.Value, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []ast.Value{}
	for _, argument := range expr.Arguments {
		value, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, value)
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, newRuntimeError(expr.Paren.Line, "Can only call function and classes")
	}
	if len(arguments) != function.arity() {
		message := fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments))
		return nil, newRuntimeError(expr.Paren.Line, message)
	}
	return function.call(i, arguments)
}

// VisitGetExpr implements ast.AstVisitor.
// It retrieves the value of a property from the object instance.
func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) (ast.Value, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}
	if instance, ok := object.(*LoxClassInstance); ok {
		return instance.get(expr.Name)
	}
	return nil, newRuntimeError(expr.Name.Line, "Only instances have properties.")
}

// VisitBooleanExpr implements ast.AstVisitor.
// It returns a new Boolean value based on the expression.
func (i *Interpreter) VisitBooleanExpr(expr *ast.BooleanExpr) (ast.Value, error) {
	return ast.NewBooleanValue(expr.Value), nil
}

// VisitGroupingExpr implements ast.AstVisitor.
// It evaluates the expression inside the grouping and returns its value.
func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (ast.Value, error) {
	return i.evaluate(expr.Expression)
}

// VisitLogicalExpr implements ast.AstVisitor.
// It evaluates the left expression and, based on the operator, decides whether to evaluate the right expression.
// If the left expression is true and the operator is "or", it returns the left value.
// If the left expression is false and the operator is "and", it returns the left value.
func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (ast.Value, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Type == ast.OrKeyword {
		if left.IsTruthy() {
			return left, nil
		}
	} else {
		if !left.IsTruthy() {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

// VisitNilExpr implements ast.AstVisitor.
// It returns a new Nil value.
func (i *Interpreter) VisitNilExpr() (ast.Value, error) {
	return ast.NewNilValue(), nil
}

// VisitNumberExpr implements ast.AstVisitor.
// It returns a new Number value based on the expression.
func (i *Interpreter) VisitNumberExpr(expr *ast.NumberExpr) (ast.Value, error) {
	return ast.NewNumberValue(expr.Value), nil
}

// VisitSetExpr implements ast.AstVisitor.
// It evaluates the object and sets the value of the field in the object instance.
func (i *Interpreter) VisitSetExpr(expr *ast.SetExpr) (ast.Value, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	loxInstance, ok := object.(*LoxClassInstance)
	if !ok {
		return nil, newRuntimeError(expr.Name.Line, "Only instances have fields.")
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	loxInstance.set(expr.Name, value)
	return ast.NewNilValue(), nil
}

// VisitStringExpr implements ast.AstVisitor.
// It returns a new String value based on the expression.
func (i *Interpreter) VisitStringExpr(expr *ast.StringExpr) (ast.Value, error) {
	return ast.NewStringValue(expr.Value), nil
}

func (i *Interpreter) VisitSuperExpr(expr *ast.SuperExpr) (ast.Value, error) {
	distance := i.locals[expr]
	superclass, err := i.environment.getAt(distance, "super")
	if err != nil {
		return ast.NewNilValue(), err
	}
	instance, err := i.environment.getAt(distance-1, "this")
	if err != nil {
		return ast.NewNilValue(), err
	}

	superclassPtr, isLoxClass := superclass.(*LoxClass)
	if !isLoxClass {
		return ast.NewNilValue(), newRuntimeError(expr.Keyword.Line, "Superclass must be a class.")
	}
	instancePtr, isLoxClassInstance := instance.(*LoxClassInstance)
	if !isLoxClassInstance {
		return ast.NewNilValue(), newRuntimeError(expr.Keyword.Line, "")
	}

	method := superclassPtr.findMethod(expr.Method.Lexeme)
	if method == nil {
		return ast.NewNilValue(), newRuntimeError(expr.Method.Line, "Undefined property '"+expr.Method.Lexeme+"'.")
	}

	return method.bind(instancePtr), nil
}

// VisitThisExpr implements ast.AstVisitor
// It retrieves the value of the "this" keyword in the current environment.
func (i *Interpreter) VisitThisExpr(expr *ast.ThisExpr) (ast.Value, error) {
	return i.lookUpVariable(expr.Keyword, expr)
}

// VisitUnaryExpr implements ast.AstVisitor.
// It evaluates the right expression and performs the unary operation.
func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (ast.Value, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	return perfromUnaryOperation(right, expr.Operator.Type, expr.Operator.Line)
}

// VisitVariableExpr implements ast.AstVisitor.
// It retrieves the value of the variable from the environment.
func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) (ast.Value, error) {
	return i.lookUpVariable(expr.Name, expr)
}
