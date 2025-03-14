package ast

import "fmt"

type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr[string]) string {
	return expr.Accept(a)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary[string]) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal[string]) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping[string]) string {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary[string]) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	result := "(" + name
	for _, expr := range exprs {
		result += " " + expr.Accept(a)
	}
	result += ")"
	return result
}
