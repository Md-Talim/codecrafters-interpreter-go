package ast

type Expr interface {
	AST
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func NewAssignExpr(name Token, value Expr) *AssignExpr {
	return &AssignExpr{
		Name:  name,
		Value: value,
	}
}

func (expr *AssignExpr) Accept(visitor AstVisitor) {
	visitor.VisitAssignExpr(expr)
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

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator Token, right Expr) *LogicalExpr {
	return &LogicalExpr{Left: left, Operator: operator, Right: right}
}

func (expr *LogicalExpr) Accept(visitor AstVisitor) {
	visitor.VisitLogicalExpr(expr)
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

type VariableExpr struct {
	Name Token
}

func NewVariableExpr(name Token) *VariableExpr {
	return &VariableExpr{Name: name}
}

func (v *VariableExpr) Accept(visitor AstVisitor) {
	visitor.VisitVariableExpr(v)
}
