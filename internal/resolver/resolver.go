package resolver

import (
	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/interpreter"
)

// Resolver is responsible for resolving variable references in the AST.
type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      Scopes
}

// NewResolver creates a new Resolver instance.
func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	scopes := newScopes()
	return &Resolver{
		interpreter: interpreter,
		scopes:      scopes,
	}
}

// VisitAssignExpr implements ast.AstVisitor.
func (r *Resolver) VisitAssignExpr(expr *ast.AssignExpr) (ast.Value, error) {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return ast.NewNilValue(), nil
}

// VisitBinaryExpr implements ast.AstVisitor.
func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) (ast.Value, error) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return ast.NewNilValue(), nil
}

// VisitBlockStmt implements ast.AstVisitor.
func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (ast.Value, error) {
	r.beginScope()
	r.resolveStatements(stmt.Statements)
	r.endScope()
	return ast.NewNilValue(), nil
}

// VisitBooleanExpr implements ast.AstVisitor.
func (r *Resolver) VisitBooleanExpr(expr *ast.BooleanExpr) (ast.Value, error) {
	return ast.NewBooleanValue(expr.Value), nil
}

// VisitCallExpr implements ast.AstVisitor.
func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) (ast.Value, error) {
	r.resolveExpression(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpression(arg)
	}
	return ast.NewNilValue(), nil
}

// VisitExpressionStmt implements ast.AstVisitor.
func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (ast.Value, error) {
	r.resolveExpression(stmt.Expression)
	return ast.NewNilValue(), nil
}

// VisitFunctionStmt implements ast.AstVisitor.
func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) (ast.Value, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt)
	return ast.NewNilValue(), nil
}

// VisitGroupingExpr implements ast.AstVisitor.
func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) (ast.Value, error) {
	r.resolveExpression(expr.Expression)
	return ast.NewNilValue(), nil
}

// VisitIfStmt implements ast.AstVisitor.
func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (ast.Value, error) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return ast.NewNilValue(), nil
}

// VisitLogicalExpr implements ast.AstVisitor.
func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) (ast.Value, error) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return ast.NewNilValue(), nil
}

// VisitNilExpr implements ast.AstVisitor.
func (r *Resolver) VisitNilExpr() (ast.Value, error) {
	return ast.NewNilValue(), nil
}

// VisitNumberExpr implements ast.AstVisitor.
func (r *Resolver) VisitNumberExpr(expr *ast.NumberExpr) (ast.Value, error) {
	return ast.NewNumberValue(expr.Value), nil
}

// VisitPrintStmt implements ast.AstVisitor.
func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) (ast.Value, error) {
	r.resolveExpression(stmt.Expression)
	return ast.NewNilValue(), nil
}

// VisitReturnStmt implements ast.AstVisitor.
func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (ast.Value, error) {
	if stmt.Value != nil {
		r.resolveExpression(stmt.Value)
	}
	return ast.NewNilValue(), nil
}

// VisitStringExpr implements ast.AstVisitor.
func (r *Resolver) VisitStringExpr(expr *ast.StringExpr) (ast.Value, error) {
	return ast.NewStringValue(expr.Value), nil
}

// VisitUnaryExpr implements ast.AstVisitor.
func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) (ast.Value, error) {
	r.resolveExpression(expr.Right)
	return ast.NewNilValue(), nil
}

// VisitVarStmt implements ast.AstVisitor.
func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) (ast.Value, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
	return ast.NewNilValue(), nil
}

// VisitVariableExpr implements ast.AstVisitor.
func (r *Resolver) VisitVariableExpr(expr *ast.VariableExpr) (ast.Value, error) {
	if !r.scopes.isEmpty() {
		scope := r.scopes.peek()
		if scope.hasKey(expr.Name.Lexeme) && !scope.get(expr.Name.Lexeme) {
			return nil, newSyntaxError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.Name)
	return ast.NewNilValue(), nil
}

// VisitWhileStmt implements ast.AstVisitor.
func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) (ast.Value, error) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return ast.NewNilValue(), nil
}

// resolveExpression resolves a single expression.
func (r *Resolver) resolveExpression(expr ast.Expr) {
	expr.Accept(r)
}

// resolveFunction resolves a function statement.
// It declares the function parameters and resolves the function body.
func (r *Resolver) resolveFunction(function *ast.FunctionStmt) {
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.Body)
	r.endScope()
}

// resolveLocal resolves a local variable reference.
// It checks the current scope and its enclosing scopes to find the variable.
func (r *Resolver) resolveLocal(expr ast.Expr, name ast.Token) {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		scope := r.scopes.get(i)
		if scope.hasKey(name.Lexeme) {
			depth := r.scopes.size() - 1 - i
			r.interpreter.Resolve(expr, depth)
			return
		}
	}
}

// resolveStatement resolves a single statement.
func (r *Resolver) resolveStatement(stmt ast.Stmt) (ast.Value, error) {
	return stmt.Accept(r)
}

// resolveStatements resolves a list of statements.
func (r *Resolver) resolveStatements(statements []ast.Stmt) (ast.Value, error) {
	var lastValue ast.Value
	var err error
	for _, stmt := range statements {
		lastValue, err = r.resolveStatement(stmt)
		if err != nil {
			break
		}
	}
	return lastValue, err
}

// Resolve resolves a list of statements.
func (r *Resolver) Resolve(statements []ast.Stmt) (ast.Value, error) {
	return r.resolveStatements(statements)
}
