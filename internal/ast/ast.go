package ast

type AST interface {
	Accept(visitor AstVisitor)
}

type Stmt interface {
	Accept(visitor AstVisitor)
}

type ExpressionStmt struct {
	Expression AST
}

func NewExpressionStmt(expression AST) *ExpressionStmt {
	return &ExpressionStmt{Expression: expression}
}

func (e *ExpressionStmt) Accept(visitor AstVisitor) {
	visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression AST
}

func NewPrintStmt(expression AST) *PrintStmt {
	return &PrintStmt{Expression: expression}
}

func (p *PrintStmt) Accept(visitor AstVisitor) {
	visitor.VisitPrintStmt(p)
}

type Expr interface {
	AST
}

type AstVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr)
	VisitBooleanExpr(expr *BooleanExpr)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitGroupingExpr(expr *GroupingExpr)
	VisitNilExpr()
	VisitNumberExpr(expr *NumberExpr)
	VisitPrintStmt(stmt *PrintStmt)
	VisitStringExpr(expr *StringExpr)
	VisitUnaryExpr(expr *UnaryExpr)
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func NewBinaryExpr(operator Token, left, right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (b *BinaryExpr) Accept(visitor AstVisitor) {
	visitor.VisitBinaryExpr(b)
}

type BooleanExpr struct {
	Value bool
}

func NewBooleanExpr(value bool) *BooleanExpr {
	return &BooleanExpr{Value: value}
}

func (b *BooleanExpr) Accept(visitor AstVisitor) {
	visitor.VisitBooleanExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{Expression: expression}
}

func (g *GroupingExpr) Accept(visitor AstVisitor) {
	visitor.VisitGroupingExpr(g)
}

type NilExpr struct{}

func NewNilExpr() *NilExpr {
	return &NilExpr{}
}

func (nil *NilExpr) Accept(visitor AstVisitor) {
	visitor.VisitNilExpr()
}

type NumberExpr struct {
	Value float64
}

func NewNumberExpr(value float64) *NumberExpr {
	return &NumberExpr{value}
}

func (num *NumberExpr) Accept(visitor AstVisitor) {
	visitor.VisitNumberExpr(num)
}

type StringExpr struct {
	Value string
}

func NewStringExpr(value string) *StringExpr {
	return &StringExpr{value}
}

func (string *StringExpr) Accept(visitor AstVisitor) {
	visitor.VisitStringExpr(string)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func NewUnaryExpr(operator Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator, right}
}

func (un *UnaryExpr) Accept(visitor AstVisitor) {
	visitor.VisitUnaryExpr(un)
}
