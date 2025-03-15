package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"fmt"
	"strconv"
	"strings"
)

type Interpreter struct{}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping[any]) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary[any]) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	// Check if both operands are strings
	leftStr, leftIsString := left.(string)
	rightStr, rightIsString := right.(string)

	if (expr.Operator.Type == ast.PlusToken) && leftIsString && rightIsString {
		leftNum, leftNumErr := strconv.ParseFloat(leftStr, 64)
		rightNum, rightNumErr := strconv.ParseFloat(rightStr, 64)

		if leftNumErr == nil && rightNumErr == nil {
			return leftNum + rightNum
		}
		return leftStr + rightStr
	}

	leftNum, leftOk := toFloat64(left)
	rightNum, rightOk := toFloat64(right)

	if !leftOk || !rightOk {
		return nil
	}

	switch expr.Operator.Type {
	case ast.LessToken:
		return leftNum < rightNum
	case ast.LessEqualToken:
		return leftNum <= rightNum
	case ast.GreaterToken:
		return leftNum > rightNum
	case ast.GreaterEqualToken:
		return leftNum >= rightNum
	case ast.MinusToken:
		return leftNum - rightNum
	case ast.PlusToken:
		return leftNum + rightNum
	case ast.StarToken:
		return leftNum * rightNum
	case ast.SlashToken:
		return leftNum / rightNum
	}

	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary[any]) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case ast.BangToken:
		return !isTruthy(right)
	case ast.MinusToken:
		i, ok := toFloat64(right)
		if !ok {
			return nil
		}
		return -i
	}

	// Unreachable code
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal[any]) any {
	return expr.Value
}

func (i *Interpreter) Interpret(expr ast.Expr[any]) {
	result := i.evaluate(expr)
	fmt.Println(i.stringify(result))
}

func (i *Interpreter) evaluate(expr ast.Expr[any]) any {
	return expr.Accept(i)
}

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case float64:
		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%g", v)
	case string:
		if strings.HasSuffix(v, ".0") {
			return v[:len(v)-2]
		}
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}
