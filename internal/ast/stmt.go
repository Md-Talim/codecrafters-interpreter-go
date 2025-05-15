package ast

type Stmt interface {
	AST
}

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{Statements: statements}
}

func (e *BlockStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitBlockStmt(e)
}

type ExpressionStmt struct {
	Expression AST
}

func NewExpressionStmt(expression AST) *ExpressionStmt {
	return &ExpressionStmt{Expression: expression}
}

func (e *ExpressionStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitExpressionStmt(e)
}

type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func NewFunctionStmt(name Token, params []Token, body []Stmt) *FunctionStmt {
	return &FunctionStmt{Name: name, Params: params, Body: body}
}

func (stmt *FunctionStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitFunctionStmt(stmt)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *IfStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitIfStmt(p)
}

type PrintStmt struct {
	Expression AST
}

func NewPrintStmt(expression AST) *PrintStmt {
	return &PrintStmt{Expression: expression}
}

func (p *PrintStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func NewVarStmt(name Token, initializer Expr) *VarStmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

func (v *VarStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitVarStmt(v)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{Condition: condition, Body: body}
}

func (stmt *WhileStmt) Accept(visitor AstVisitor) (Value, error) {
	return visitor.VisitWhileStmt(stmt)
}
