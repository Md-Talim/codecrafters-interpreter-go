package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"fmt"
	"strings"
)

type Interpreter struct{}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping[any]) any {
	return ""
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary[any]) any {
	return ""
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary[any]) any {
	return ""
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
