package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/parser"
	"codecrafters-interpreter-go/pkg/loxerrors"
)

type Interpreter struct {
	result       ast.Value
	runtimeError error
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) {
	var left, right ast.Value
	var err error

	operator := expr.Operator.Type

	i.result = nil
	i.runtimeError = nil

	left, err = i.evaluate(expr.Left)
	if err != nil {
		i.runtimeError = err
		return
	}

	right, err = i.evaluate(expr.Right)
	if err != nil {
		i.runtimeError = err
		return
	}

	leftType := left.GetType()
	rightType := right.GetType()
	bothNums := leftType == ast.NumberType && rightType == ast.NumberType

	switch operator {
	case ast.StarToken, ast.SlashToken, ast.MinusToken, ast.GreaterEqualToken, ast.GreaterToken, ast.LessEqualToken, ast.LessToken:
		if bothNums {
			leftNum := left.(*ast.NumberValue).Value
			rightNum := right.(*ast.NumberValue).Value

			switch operator {
			case ast.StarToken:
				i.result = ast.NewNumberValue(leftNum * rightNum)
			case ast.SlashToken:
				i.result = ast.NewNumberValue(leftNum / rightNum)
			case ast.MinusToken:
				i.result = ast.NewNumberValue(leftNum - rightNum)
			case ast.GreaterToken:
				i.result = ast.NewBooleanValue(leftNum > rightNum)
			case ast.GreaterEqualToken:
				i.result = ast.NewBooleanValue(leftNum >= rightNum)
			case ast.LessToken:
				i.result = ast.NewBooleanValue(leftNum < rightNum)
			case ast.LessEqualToken:
				i.result = ast.NewBooleanValue(leftNum <= rightNum)
			}
		} else {
			i.runtimeError = loxerrors.NewRuntimeError(expr.Operator.Line, "Operands must be numbers.")
			return
		}
	case ast.PlusToken:
		if bothNums {
			i.result = ast.NewNumberValue(left.(*ast.NumberValue).Value + right.(*ast.NumberValue).Value)
		} else if leftType == ast.StringType && rightType == ast.StringType {
			i.result = ast.NewStringValue(left.(*ast.StringValue).Value + right.(*ast.StringValue).Value)
		} else {
			i.runtimeError = loxerrors.NewRuntimeError(expr.Operator.Line, "Operands must be two numbers or two strings.")
		}
	case ast.EqualEqualToken:
		i.result = ast.NewBooleanValue(left.IsEqualTo(right))
	case ast.BangEqualToken:
		i.result = ast.NewBooleanValue(!left.IsEqualTo(right))
	default:
		i.runtimeError = loxerrors.NewRuntimeError(expr.Operator.Line, "Unknown operator.")
	}
}

func (i *Interpreter) VisitBooleanExpr(expr *ast.BooleanExpr) {
	i.result = ast.NewBooleanValue(expr.Value)
	i.runtimeError = nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) {
	i.result, i.runtimeError = i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitNilExpr() {
	i.result = ast.NewNilValue()
	i.runtimeError = nil
}

func (i *Interpreter) VisitNumberExpr(expr *ast.NumberExpr) {
	i.result = ast.NewNumberValue(expr.Value)
	i.runtimeError = nil
}

func (i *Interpreter) VisitStringExpr(expr *ast.StringExpr) {
	i.result = ast.NewStringValue(expr.Value)
	i.runtimeError = nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		i.runtimeError = err
		return
	}

	switch expr.Operator.Type {
	case ast.BangToken:
		switch right.GetType() {
		case ast.BooleanType:
			boolVal := right.(*ast.BooleanValue)
			i.result = ast.NewBooleanValue(!boolVal.Value)
			i.runtimeError = nil
		case ast.StringType:
			strVal := right.(*ast.StringValue)
			i.result = ast.NewBooleanValue(strVal.Value == "")
			i.runtimeError = nil
		case ast.NilType:
			i.result = ast.NewBooleanValue(true)
			i.runtimeError = nil
		case ast.NumberType:
			numVal := right.(*ast.NumberValue)
			i.result = ast.NewBooleanValue(numVal.Value == 0)
			i.runtimeError = nil
		default:
			i.result = nil
			i.runtimeError = loxerrors.NewRuntimeError(expr.Operator.Line, "Operand must be a boolean, string, nil, or number.")
		}
	case ast.MinusToken:
		if right.GetType() == ast.NumberType {
			num := right.(*ast.NumberValue)
			i.result = ast.NewNumberValue(-num.Value)
			i.runtimeError = nil
		} else {
			i.result = nil
			i.runtimeError = loxerrors.NewRuntimeError(expr.Operator.Line, "Operands must be numbers.")
		}
	}
}

func (i *Interpreter) Interpret(source string) (ast.Value, error) {
	parser := parser.NewParser(source)
	expr, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	value, err := i.evaluate(expr)
	if err != nil {
		return nil, err
	}
	return value, err
}

func (i *Interpreter) evaluate(ast ast.AST) (ast.Value, error) {
	ast.Accept(i)
	return i.result, i.runtimeError
}
