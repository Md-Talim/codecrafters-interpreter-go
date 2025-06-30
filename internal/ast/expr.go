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

func (expr *AssignExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitAssignExpr(expr)
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

func (b *BinaryExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitBinaryExpr(b)
}

type BooleanExpr struct {
	Value bool
}

func NewBooleanExpr(value bool) *BooleanExpr {
	return &BooleanExpr{Value: value}
}

func (b *BooleanExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitBooleanExpr(b)
}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func NewCallExpr(callee Expr, paren Token, arguments []Expr) *CallExpr {
	return &CallExpr{Callee: callee, Paren: paren, Arguments: arguments}
}

func (expr *CallExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitCallExpr(expr)
}

type GetExpr struct {
	Object Expr
	Name   Token
}

func NewGetExpr(object Expr, name Token) *GetExpr {
	return &GetExpr{Object: object, Name: name}
}

func (g *GetExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitGetExpr(g)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{Expression: expression}
}

func (g *GroupingExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitGroupingExpr(g)
}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator Token, right Expr) *LogicalExpr {
	return &LogicalExpr{Left: left, Operator: operator, Right: right}
}

func (expr *LogicalExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitLogicalExpr(expr)
}

type NilExpr struct{}

func NewNilExpr() *NilExpr {
	return &NilExpr{}
}

func (nil *NilExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitNilExpr()
}

type NumberExpr struct {
	Value float64
}

func NewNumberExpr(value float64) *NumberExpr {
	return &NumberExpr{value}
}

func (num *NumberExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitNumberExpr(num)
}

type SetExpr struct {
	Object Expr
	Name   Token
	Value  Expr
}

func NewSetExpr(object Expr, name Token, value Expr) *SetExpr {
	return &SetExpr{Object: object, Name: name, Value: value}
}

func (s *SetExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitSetExpr(s)
}

type StringExpr struct {
	Value string
}

func NewStringExpr(value string) *StringExpr {
	return &StringExpr{value}
}

func (string *StringExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitStringExpr(string)
}

type SuperExpr struct {
	Keyword Token
	Method  Token
}

func NewSuperExpr(keyword, method Token) *SuperExpr {
	return &SuperExpr{keyword, method}
}

func (expr *SuperExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitSuperExpr(expr)
}

type ThisExpr struct {
	Keyword Token
}

func NewThisExpr(keyword Token) *ThisExpr {
	return &ThisExpr{Keyword: keyword}
}

func (expr *ThisExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitThisExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func NewUnaryExpr(operator Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator, right}
}

func (un *UnaryExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitUnaryExpr(un)
}

type VariableExpr struct {
	Name Token
}

func NewVariableExpr(name Token) *VariableExpr {
	return &VariableExpr{Name: name}
}

func (v *VariableExpr) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitVariableExpr(v)
}
