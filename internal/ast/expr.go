package ast

type Visitor[R any] interface {
	VisitLiteralExpr(expr *Literal[R]) R
	VisitGroupingExpr(expr *Grouping[R]) R
	VisitUnaryExpr(expr *Unary[R]) R
	VisitBinaryExpr(expr *Binary[R]) R
}

type Expr[R any] interface {
	Accept(visitor Visitor[R]) R
}

type Binary[R any] struct {
	Left     Expr[R]
	Operator Token
	Right    Expr[R]
}

func (b *Binary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitBinaryExpr(b)
}

type Literal[R any] struct {
	Value any
}

func (l *Literal[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitLiteralExpr(l)
}

type Grouping[R any] struct {
	Expression Expr[R]
}

func (g *Grouping[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitGroupingExpr(g)
}

type Unary[R any] struct {
	Operator Token
	Right    Expr[R]
}

func (u *Unary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitUnaryExpr(u)
}
