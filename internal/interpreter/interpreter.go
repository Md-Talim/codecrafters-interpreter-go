package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/parser"
	"errors"
	"fmt"
	"os"
)

type Interpreter struct {
	environment  *Environment
	result       ast.Value
	runtimeError error
}

func NewInterpreter() *Interpreter {
	environment := &Environment{enclosing: nil, values: make(map[string]ast.Value)}
	return &Interpreter{environment: environment}
}

func newRuntimeError(line int, message string) error {
	text := fmt.Sprintf("%s\n[line %d]", message, line)
	return errors.New(text)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) {
	value, _ := i.evaluate(expr.Value)
	i.environment.assign(expr.Name, value)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) {
	var left, right ast.Value
	var err error

	operator := expr.Operator.Type

	i.result = nil
	i.runtimeError = nil

	left, err = i.evaluate(expr.Left)
	if err != nil {
		i.runtimeError = err
		return
	}

	right, err = i.evaluate(expr.Right)
	if err != nil {
		i.runtimeError = err
		return
	}

	leftType := left.GetType()
	rightType := right.GetType()
	bothNums := leftType == ast.NumberType && rightType == ast.NumberType

	switch operator {
	case ast.StarToken, ast.SlashToken, ast.MinusToken, ast.GreaterEqualToken, ast.GreaterToken, ast.LessEqualToken, ast.LessToken:
		if bothNums {
			leftNum := left.(*ast.NumberValue).Value
			rightNum := right.(*ast.NumberValue).Value

			switch operator {
			case ast.StarToken:
				i.result = ast.NewNumberValue(leftNum * rightNum)
			case ast.SlashToken:
				i.result = ast.NewNumberValue(leftNum / rightNum)
			case ast.MinusToken:
				i.result = ast.NewNumberValue(leftNum - rightNum)
			case ast.GreaterToken:
				i.result = ast.NewBooleanValue(leftNum > rightNum)
			case ast.GreaterEqualToken:
				i.result = ast.NewBooleanValue(leftNum >= rightNum)
			case ast.LessToken:
				i.result = ast.NewBooleanValue(leftNum < rightNum)
			case ast.LessEqualToken:
				i.result = ast.NewBooleanValue(leftNum <= rightNum)
			}
		} else {
			i.runtimeError = newRuntimeError(expr.Operator.Line, "Operands must be numbers.")
		}
	case ast.PlusToken:
		if bothNums {
			i.result = ast.NewNumberValue(left.(*ast.NumberValue).Value + right.(*ast.NumberValue).Value)
		} else if leftType == ast.StringType && rightType == ast.StringType {
			i.result = ast.NewStringValue(left.(*ast.StringValue).Value + right.(*ast.StringValue).Value)
		} else {
			i.runtimeError = newRuntimeError(expr.Operator.Line, "Operands must be two numbers or two strings.")
		}
	case ast.EqualEqualToken:
		i.result = ast.NewBooleanValue(left.IsEqualTo(right))
	case ast.BangEqualToken:
		i.result = ast.NewBooleanValue(!left.IsEqualTo(right))
	default:
		i.runtimeError = newRuntimeError(expr.Operator.Line, "Unknown operator.")
	}
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) {
	i.executeBlock(stmt.Statements, &Environment{enclosing: i.environment, values: make(map[string]ast.Value)})
}

func (i *Interpreter) VisitBooleanExpr(expr *ast.BooleanExpr) {
	i.result = ast.NewBooleanValue(expr.Value)
	i.runtimeError = nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) {
	i.result, i.runtimeError = i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) {
	value, _ := i.evaluate(stmt.Condition)
	if value.IsTruthy() {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) {
	left, _ := i.evaluate(expr.Left)
	if expr.Operator.Type == ast.OrKeyword {
		if left.IsTruthy() {
			i.result = left
			return
		}
	}
	right, _ := i.evaluate(expr.Right)
	i.result = right
}

func (i *Interpreter) VisitNilExpr() {
	i.result = ast.NewNilValue()
	i.runtimeError = nil
}

func (i *Interpreter) VisitNumberExpr(expr *ast.NumberExpr) {
	i.result = ast.NewNumberValue(expr.Value)
	i.runtimeError = nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) {
	if value, err := i.evaluate(stmt.Expression); err == nil {
		fmt.Println(value)
	}
}

func (i *Interpreter) VisitStringExpr(expr *ast.StringExpr) {
	i.result = ast.NewStringValue(expr.Value)
	i.runtimeError = nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		i.runtimeError = err
		return
	}

	switch expr.Operator.Type {
	case ast.BangToken:
		switch right.GetType() {
		case ast.BooleanType:
			boolVal := right.(*ast.BooleanValue)
			i.result = ast.NewBooleanValue(!boolVal.Value)
			i.runtimeError = nil
		case ast.StringType:
			strVal := right.(*ast.StringValue)
			i.result = ast.NewBooleanValue(strVal.Value == "")
			i.runtimeError = nil
		case ast.NilType:
			i.result = ast.NewBooleanValue(true)
			i.runtimeError = nil
		case ast.NumberType:
			numVal := right.(*ast.NumberValue)
			i.result = ast.NewBooleanValue(numVal.Value == 0)
			i.runtimeError = nil
		default:
			i.result = nil
			i.runtimeError = newRuntimeError(expr.Operator.Line, "Operand must be a boolean, string, nil, or number.")
		}
	case ast.MinusToken:
		if right.GetType() == ast.NumberType {
			num := right.(*ast.NumberValue)
			i.result = ast.NewNumberValue(-num.Value)
			i.runtimeError = nil
		} else {
			i.result = nil
			i.runtimeError = newRuntimeError(expr.Operator.Line, "Operands must be numbers.")
		}
	}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) {
	var value ast.Value = ast.NewNilValue()
	if stmt.Initializer != nil {
		value, _ = i.evaluate(stmt.Initializer)
	}
	i.environment.define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) {
	value, err := i.environment.get(expr.Name)
	i.result = value
	i.runtimeError = err
}

func (i *Interpreter) Interpret(source string) (ast.Value, error) {
	parser := parser.NewParser(source)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	value, err := i.evaluate(expr)
	if err != nil {
		return nil, err
	}
	return value, err
}

func (i *Interpreter) Run(source string) {
	parser := parser.NewParser(source)
	statements, err := parser.GetStatements()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	for _, stmt := range statements {
		if _, err := i.execute(stmt); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(70)
		}
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) (ast.Value, error) {
	stmt.Accept(i)
	return i.result, i.runtimeError
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
	previous := i.environment
	i.environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}
	i.environment = previous
}

func (i *Interpreter) evaluate(ast ast.AST) (ast.Value, error) {
	ast.Accept(i)
	return i.result, i.runtimeError
}
