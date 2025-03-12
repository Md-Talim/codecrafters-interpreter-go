package main

import "fmt"

type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.Accept(a)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) string {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	result := "(" + name
	for _, expr := range exprs {
		result += " " + expr.Accept(a)
	}
	result += ")"
	return result
}
