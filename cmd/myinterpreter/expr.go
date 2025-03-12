package main

type Visitor interface {
	VisitLiteralExpr(expr *Literal) string
	VisitGroupingExpr(expr *Grouping) string
	VisitUnaryExpr(expr *Unary) string
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

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor Visitor) string {
	return visitor.VisitGroupingExpr(g)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(visitor Visitor) string {
	return visitor.VisitUnaryExpr(u)
}
