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

func (e *BlockStmt) Accept(visitor AstVisitor) {
	visitor.VisitBlockStmt(e)
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

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *IfStmt) Accept(visitor AstVisitor) {
	visitor.VisitIfStmt(p)
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

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func NewVarStmt(name Token, initializer Expr) *VarStmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

func (v *VarStmt) Accept(visitor AstVisitor) {
	visitor.VisitVarStmt(v)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{Condition: condition, Body: body}
}

func (stmt *WhileStmt) Accept(visitor AstVisitor) {
	visitor.VisitWhileStmt(stmt)
}
