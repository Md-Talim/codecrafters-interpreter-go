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

	var leftd float64
	var rightd float64
	var err error
	var ok bool

	if leftd, ok = left.(float64); !ok {
		leftd, err = strconv.ParseFloat(left.(string), 64)
		if err != nil {
			return nil
		}
	}

	if rightd, ok = right.(float64); !ok {
		rightd, err = strconv.ParseFloat(right.(string), 64)
		if err != nil {
			return nil
		}
	}

	switch expr.Operator.Type {
	case ast.MinusToken:
		return fmt.Sprintf("%g", leftd-rightd)
	case ast.PlusToken:
		return fmt.Sprintf("%g", leftd+rightd)
	case ast.StarToken:
		return fmt.Sprintf("%g", leftd*rightd)
	case ast.SlashToken:
		return fmt.Sprintf("%g", leftd/rightd)
	}

	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary[any]) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case ast.BangToken:
		return !isTruthy(right)
	case ast.MinusToken:
		i, err := strconv.ParseFloat(right.(string), 64)
		if err != nil {
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

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	switch v := object.(type) {
	case bool:
		return v
	default:
		return true
	}
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
