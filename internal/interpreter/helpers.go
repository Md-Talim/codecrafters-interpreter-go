package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"fmt"
)

// performNumericOperation performs numeric operations on two values.
// It returns the result of the operation or an error if the operands are not numbers.
// Supported operations are addition, subtraction, multiplication, division,
// greater than, greater than or equal to, less than, and less than or equal to.
func performNumericOperation(left ast.Value, right ast.Value, operator ast.TokenType, line int) (ast.Value, error) {
	leftNum, leftOk := left.(*ast.NumberValue)
	rightNum, rightOk := right.(*ast.NumberValue)
	if !leftOk || !rightOk {
		return nil, newRuntimeError(line, "Operands must be numbers.")
	}
	switch operator {
	case ast.StarToken:
		return ast.NewNumberValue(leftNum.Value * rightNum.Value), nil
	case ast.SlashToken:
		return ast.NewNumberValue(leftNum.Value / rightNum.Value), nil
	case ast.MinusToken:
		return ast.NewNumberValue(leftNum.Value - rightNum.Value), nil
	case ast.GreaterToken:
		return ast.NewBooleanValue(leftNum.Value > rightNum.Value), nil
	case ast.GreaterEqualToken:
		return ast.NewBooleanValue(leftNum.Value >= rightNum.Value), nil
	case ast.LessToken:
		return ast.NewBooleanValue(leftNum.Value < rightNum.Value), nil
	case ast.LessEqualToken:
		return ast.NewBooleanValue(leftNum.Value <= rightNum.Value), nil
	default:
		return nil, newRuntimeError(line, "Unknown numeric operator.")
	}
}

// performAddition performs addition on two values.
// It returns the result of the addition or an error if the operands are not numbers or strings.
func performAddition(left ast.Value, right ast.Value, line int) (ast.Value, error) {
	if left.GetType() == ast.NumberType && right.GetType() == ast.NumberType {
		return ast.NewNumberValue(left.(*ast.NumberValue).Value + right.(*ast.NumberValue).Value), nil
	} else if left.GetType() == ast.StringType && right.GetType() == ast.StringType {
		return ast.NewStringValue(left.(*ast.StringValue).Value + right.(*ast.StringValue).Value), nil
	}
	return nil, newRuntimeError(line, "Operands must be two numbers or two strings.")
}

// performEqualityOperation performs equality operations on two values.
// It returns true if the values are equal or false if they are not.
func performEqualityOperation(left ast.Value, right ast.Value, operator ast.TokenType) ast.Value {
	if operator == ast.EqualEqualToken {
		return ast.NewBooleanValue(left.IsEqualTo(right))
	}
	return ast.NewBooleanValue(!left.IsEqualTo(right))
}

// performUnaryOperation performs unary operations on a value.
// It returns the result of the operation or an error if the operand is not valid.
func perfromUnaryOperation(right ast.Value, operator ast.TokenType, line int) (ast.Value, error) {
	switch operator {
	case ast.BangToken:
		switch right.GetType() {
		case ast.BooleanType:
			return ast.NewBooleanValue(!right.(*ast.BooleanValue).Value), nil
		case ast.StringType:
			return ast.NewBooleanValue(right.(*ast.StringValue).Value == ""), nil
		case ast.NilType:
			return ast.NewBooleanValue(true), nil
		case ast.NumberType:
			return ast.NewBooleanValue(right.(*ast.NumberValue).Value == 0), nil
		default:
			return nil, newRuntimeError(line, "Operand must be a boolean, string, nil, or number.")
		}
	case ast.MinusToken:
		if right.GetType() == ast.NumberType {
			return ast.NewNumberValue(-right.(*ast.NumberValue).Value), nil
		}
		return nil, newRuntimeError(line, "Operands must be numbers.")
	}
	return nil, nil
}

// newRuntimeError creates a new runtime error with the given line number and message.
// It formats the error message to include the line number for better debugging.
func newRuntimeError(line int, message string) error {
	text := fmt.Sprintf("%s\n[line %d]", message, line)
	return fmt.Errorf("%s", text)
}
