package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"fmt"
)

// returnError is a special error used to unwind the call stack for return statements.
// It's not a "real" error but a control flow mechanism.
type ReturnError struct {
	value ast.Value
}

// Error implements the error interface for ReturnError.
func (r *ReturnError) Error() string {
	if r.value != nil && r.value.GetType() != ast.NilType {
		return fmt.Sprintf("Return value: %s", r.value)
	}
	return "return"
}

// VisitBlockStmt implements the ast.AstVisitor.
// It executes the block of statements in a new environment.
func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (ast.Value, error) {
	return i.executeBlock(stmt.Statements, newEnvironment(i.environment))
}

// VisitClassStmt implements the ast.AstVisitor.
// It defines a new class in the current environment.
func (i *Interpreter) VisitClassStmt(stmt *ast.ClassStmt) (ast.Value, error) {
	i.environment.define(stmt.Name.Lexeme, ast.NewNilValue())
	class := newLoxClass(stmt.Name.Lexeme)
	i.environment.define(stmt.Name.Lexeme, class)
	return ast.NewNilValue(), nil
}

// VisitExpressionStmt implements the ast.AstVisitor.
// It evaluates the expression statement.
func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) (ast.Value, error) {
	return i.evaluate(stmt.Expression)
}

// VisitFunctionStmt implements the ast.AstVisitor.
// It defines a new function in the current environment.
func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) (ast.Value, error) {
	// The function is a closure, capturing the current environment.
	// The function is not executed immediately; it is stored in the environment.
	function := newLoxFunction(*stmt, i.environment)
	i.environment.define(stmt.Name.Lexeme, function)
	return nil, nil
}

// VisitIfStmt implements the ast.AstVisitor.
// It evaluates the condition and executes the appropriate branch.
func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (ast.Value, error) {
	value, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if value.IsTruthy() {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

// VisitPrintStmt implements the ast.AstVisitor.
// It evaluates the expression and prints its value to the console.
func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (ast.Value, error) {
	value, err := i.evaluate(stmt.Expression)
	if err == nil {
		fmt.Println(value)
	}
	return nil, err
}

// VisitReturnStmt implements the ast.AstVisitor.
// It raises a ReturnError with the value to be returned.
// This is used to exit from a function and return a value.
func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) (ast.Value, error) {
	var value ast.Value = ast.NewNilValue()
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, &ReturnError{value: value}
}

// VisitVarStmt implements the ast.AstVisitor.
// It defines a new variable in the current environment with an optional initializer.
// If no initializer is provided, the variable is initialized to nil.
func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (ast.Value, error) {
	var value ast.Value = ast.NewNilValue()
	var err error
	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.define(stmt.Name.Lexeme, value)
	return nil, nil
}

// VisitWhileStmt implements the ast.AstVisitor.
// It evaluates the condition and executes the body of the loop until the condition is false.
func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) (ast.Value, error) {
	var lastValue ast.Value
	for {
		condition, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !condition.IsTruthy() {
			break
		}
		lastValue, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return lastValue, nil
}
