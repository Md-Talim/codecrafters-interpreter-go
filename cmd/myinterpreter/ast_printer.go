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
