package main

type Visitor interface {
	VisitLiteralExpr(expr *Literal) string
}

type Expr interface {
	Accept(visitor Visitor) string
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor Visitor) string {
	return visitor.VisitLiteralExpr(l)
}
