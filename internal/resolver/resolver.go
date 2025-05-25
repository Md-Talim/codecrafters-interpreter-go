package resolver

import (
	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/interpreter"
)

type FunctionType int

const (
	NoFunction FunctionType = iota
	Function
	Method
)

// Resolver is responsible for resolving variable references in the AST.
type Resolver struct {
	currentFunctionType FunctionType
	interpreter         *interpreter.Interpreter
	scopes              Scopes
}

// NewResolver creates a new Resolver instance.
func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	scopes := newScopes()
	return &Resolver{
		currentFunctionType: NoFunction,
		interpreter:         interpreter,
		scopes:              scopes,
	}
}

// VisitAssignExpr implements ast.AstVisitor.
func (r *Resolver) VisitAssignExpr(expr *ast.AssignExpr) (ast.Value, error) {
	if _, err := r.resolveExpression(expr.Value); err != nil {
		return ast.NewNilValue(), err
	}
	r.resolveLocal(expr, expr.Name)
	return ast.NewNilValue(), nil
}

// VisitBinaryExpr implements ast.AstVisitor.
func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) (ast.Value, error) {
	if _, err := r.resolveExpression(expr.Left); err != nil {
		return ast.NewNilValue(), err
	}
	if _, err := r.resolveExpression(expr.Right); err != nil {
		return ast.NewNilValue(), err
	}
	return ast.NewNilValue(), nil
}

// VisitBlockStmt implements ast.AstVisitor.
func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) (ast.Value, error) {
	r.beginScope()
	if _, err := r.resolveStatements(stmt.Statements); err != nil {
		r.endScope()
		return ast.NewNilValue(), err
	}
	r.endScope()
	return ast.NewNilValue(), nil
}

// VisitBooleanExpr implements ast.AstVisitor.
func (r *Resolver) VisitBooleanExpr(expr *ast.BooleanExpr) (ast.Value, error) {
	return ast.NewBooleanValue(expr.Value), nil
}

// VisitCallExpr implements ast.AstVisitor.
func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) (ast.Value, error) {
	if _, err := r.resolveExpression(expr.Callee); err != nil {
		return ast.NewNilValue(), err
	}
	for _, arg := range expr.Arguments {
		if _, err := r.resolveExpression(arg); err != nil {
			return ast.NewNilValue(), err
		}
	}
	return ast.NewNilValue(), nil
}

// VisitClassStmt implements ast.AstVisitor.
func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) (ast.Value, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	for _, method := range stmt.Methods {
		r.resolveFunction(&method, Method)
	}

	return ast.NewNilValue(), nil
}

// VisitExpressionStmt implements ast.AstVisitor.
func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) (ast.Value, error) {
	if _, err := r.resolveExpression(stmt.Expression); err != nil {
		return ast.NewNilValue(), err
	}
	return ast.NewNilValue(), nil
}

// VisitFunctionStmt implements ast.AstVisitor.
func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) (ast.Value, error) {
	if err := r.declare(stmt.Name); err != nil {
		return ast.NewNilValue(), err
	}
	r.define(stmt.Name)
	return r.resolveFunction(stmt, Function)
}

func (r *Resolver) VisitGetExpr(expr *ast.GetExpr) (ast.Value, error) {
	r.resolveExpression(expr.Object)
	return ast.NewNilValue(), nil
}

// VisitGroupingExpr implements ast.AstVisitor.
func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) (ast.Value, error) {
	if _, err := r.resolveExpression(expr.Expression); err != nil {
		return ast.NewNilValue(), err
	}
	return ast.NewNilValue(), nil
}

// VisitIfStmt implements ast.AstVisitor.
func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) (ast.Value, error) {
	if _, err := r.resolveExpression(stmt.Condition); err != nil {
		return ast.NewNilValue(), err
	}
	if _, err := r.resolveStatement(stmt.ThenBranch); err != nil {
		return ast.NewNilValue(), err
	}
	if stmt.ElseBranch != nil {
		if _, err := r.resolveStatement(stmt.ElseBranch); err != nil {
			return ast.NewNilValue(), err
		}
	}
	return ast.NewNilValue(), nil
}

// VisitLogicalExpr implements ast.AstVisitor.
func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) (ast.Value, error) {
	if _, err := r.resolveExpression(expr.Left); err != nil {
		return ast.NewNilValue(), err
	}
	if _, err := r.resolveExpression(expr.Right); err != nil {
		return ast.NewNilValue(), err
	}
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
	return r.resolveExpression(stmt.Expression)
}

// VisitReturnStmt implements ast.AstVisitor.
func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) (ast.Value, error) {
	if r.currentFunctionType == NoFunction {
		return nil, newSyntaxError(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		if _, err := r.resolveExpression(stmt.Value); err != nil {
			return ast.NewNilValue(), err
		}
	}
	return ast.NewNilValue(), nil
}

// VisitSetExpr implements ast.AstVisitor.
func (r *Resolver) VisitSetExpr(expr *ast.SetExpr) (ast.Value, error) {
	r.resolveExpression(expr.Object)
	r.resolveExpression(expr.Value)
	return ast.NewNilValue(), nil
}

// VisitStringExpr implements ast.AstVisitor.
func (r *Resolver) VisitStringExpr(expr *ast.StringExpr) (ast.Value, error) {
	return ast.NewStringValue(expr.Value), nil
}

// VisitUnaryExpr implements ast.AstVisitor.
func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) (ast.Value, error) {
	if _, err := r.resolveExpression(expr.Right); err != nil {
		return ast.NewNilValue(), err
	}
	return ast.NewNilValue(), nil
}

// VisitVarStmt implements ast.AstVisitor.
func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) (ast.Value, error) {
	if err := r.declare(stmt.Name); err != nil {
		return ast.NewNilValue(), err
	}
	var err error
	if stmt.Initializer != nil {
		if _, err = r.resolveExpression(stmt.Initializer); err != nil {
			return ast.NewNilValue(), err
		}
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
	if _, err := r.resolveExpression(stmt.Condition); err != nil {
		return ast.NewNilValue(), err
	}
	if _, err := r.resolveStatement(stmt.Body); err != nil {
		return ast.NewNilValue(), err
	}
	return ast.NewNilValue(), nil
}

// resolveExpression resolves a single expression.
func (r *Resolver) resolveExpression(expr ast.Expr) (ast.Value, error) {
	return expr.Accept(r)
}

// resolveFunction resolves a function statement.
// It declares the function parameters and resolves the function body.
func (r *Resolver) resolveFunction(function *ast.FunctionStmt, functionType FunctionType) (ast.Value, error) {
	enclosingFunctionType := r.currentFunctionType
	r.currentFunctionType = functionType
	r.beginScope()
	for _, param := range function.Params {
		if err := r.declare(param); err != nil {
			r.endScope()
			return ast.NewNilValue(), err
		}
		r.define(param)
	}
	if _, err := r.resolveStatements(function.Body); err != nil {
		r.endScope()
		return ast.NewNilValue(), err
	}
	r.endScope()
	r.currentFunctionType = enclosingFunctionType
	return ast.NewNilValue(), nil
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
			return nil, err
		}
	}
	return lastValue, err
}

// Resolve resolves a list of statements.
func (r *Resolver) Resolve(statements []ast.Stmt) (ast.Value, error) {
	return r.resolveStatements(statements)
}
