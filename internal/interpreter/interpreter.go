package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/parser"
	"errors"
	"fmt"
	"os"
)

type Interpreter struct {
	environment *Environment
	globals     *Environment
}

func NewInterpreter() *Interpreter {
	globals := newEnvironment(nil)
	environment := globals
	globals.define("clock", NewClockFunction())
	return &Interpreter{environment: environment, globals: globals}
}

func newRuntimeError(line int, message string) error {
	text := fmt.Sprintf("%s\n[line %d]", message, line)
	return errors.New(text)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (ast.Value, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.environment.assign(expr.Name, value)
	return value, err
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (ast.Value, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	leftType := left.GetType()
	rightType := right.GetType()
	areOperandsNumeric := leftType == ast.NumberType && rightType == ast.NumberType
	areOperandsStrings := leftType == ast.StringType && rightType == ast.StringType
	operator := expr.Operator.Type

	switch operator {
	case ast.StarToken, ast.SlashToken, ast.MinusToken, ast.GreaterEqualToken, ast.GreaterToken, ast.LessEqualToken, ast.LessToken:
		if areOperandsNumeric {
			leftNum := left.(*ast.NumberValue).Value
			rightNum := right.(*ast.NumberValue).Value
			switch operator {
			case ast.StarToken:
				return ast.NewNumberValue(leftNum * rightNum), nil
			case ast.SlashToken:
				return ast.NewNumberValue(leftNum / rightNum), nil
			case ast.MinusToken:
				return ast.NewNumberValue(leftNum - rightNum), nil
			case ast.GreaterToken:
				return ast.NewBooleanValue(leftNum > rightNum), nil
			case ast.GreaterEqualToken:
				return ast.NewBooleanValue(leftNum >= rightNum), nil
			case ast.LessToken:
				return ast.NewBooleanValue(leftNum < rightNum), nil
			case ast.LessEqualToken:
				return ast.NewBooleanValue(leftNum <= rightNum), nil
			}
		} else {
			return nil, newRuntimeError(expr.Operator.Line, "Operands must be numbers.")
		}
	case ast.PlusToken:
		if areOperandsNumeric {
			return ast.NewNumberValue(left.(*ast.NumberValue).Value + right.(*ast.NumberValue).Value), nil
		} else if areOperandsStrings {
			return ast.NewStringValue(left.(*ast.StringValue).Value + right.(*ast.StringValue).Value), nil
		} else {
			return nil, newRuntimeError(expr.Operator.Line, "Operands must be two numbers or two strings.")
		}
	case ast.EqualEqualToken:
		return ast.NewBooleanValue(left.IsEqualTo(right)), nil
	case ast.BangEqualToken:
		return ast.NewBooleanValue(!left.IsEqualTo(right)), nil
	default:
		return nil, newRuntimeError(expr.Operator.Line, "Unknown operator.")
	}
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (ast.Value, error) {
	return i.executeBlock(stmt.Statements, &Environment{enclosing: i.environment, values: make(map[string]ast.Value)})
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (ast.Value, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	arguments := []ast.Value{}
	for _, argument := range expr.Arguments {
		value, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, value)
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, newRuntimeError(expr.Paren.Line, "Can only call function and classes")
	}
	if len(arguments) != function.arity() {
		message := fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments))
		return nil, newRuntimeError(expr.Paren.Line, message)
	}
	return function.call(i, arguments), nil
}

func (i *Interpreter) VisitBooleanExpr(expr *ast.BooleanExpr) (ast.Value, error) {
	return ast.NewBooleanValue(expr.Value), nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) (ast.Value, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) (ast.Value, error) {
	function := newLoxFunction(*stmt)
	i.environment.define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (ast.Value, error) {
	return i.evaluate(expr.Expression)
}

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

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (ast.Value, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Type == ast.OrKeyword {
		if left.IsTruthy() {
			return left, nil
		}
	} else {
		if !left.IsTruthy() {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitNilExpr() (ast.Value, error) {
	return ast.NewNilValue(), nil
}

func (i *Interpreter) VisitNumberExpr(expr *ast.NumberExpr) (ast.Value, error) {
	return ast.NewNumberValue(expr.Value), nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (ast.Value, error) {
	value, err := i.evaluate(stmt.Expression)
	if err == nil {
		fmt.Println(value)
	}
	return nil, err
}

func (i *Interpreter) VisitStringExpr(expr *ast.StringExpr) (ast.Value, error) {
	return ast.NewStringValue(expr.Value), nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (ast.Value, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case ast.BangToken:
		switch right.GetType() {
		case ast.BooleanType:
			boolVal := right.(*ast.BooleanValue)
			return ast.NewBooleanValue(!boolVal.Value), nil
		case ast.StringType:
			strVal := right.(*ast.StringValue)
			return ast.NewBooleanValue(strVal.Value == ""), nil
		case ast.NilType:
			return ast.NewBooleanValue(true), nil
		case ast.NumberType:
			numVal := right.(*ast.NumberValue)
			return ast.NewBooleanValue(numVal.Value == 0), nil
		default:
			return nil, newRuntimeError(expr.Operator.Line, "Operand must be a boolean, string, nil, or number.")
		}
	case ast.MinusToken:
		if right.GetType() == ast.NumberType {
			num := right.(*ast.NumberValue)
			return ast.NewNumberValue(-num.Value), nil
		} else {
			return nil, newRuntimeError(expr.Operator.Line, "Operands must be numbers.")
		}
	}
	return nil, nil
}

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

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) (ast.Value, error) {
	value, err := i.environment.get(expr.Name)
	return value, err
}

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

func (i *Interpreter) Interpret(source string) (ast.Value, error) {
	parser := parser.NewParser(source)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}
	return i.evaluate(expr)
}

func (i *Interpreter) Run(source string) {
	parser := parser.NewParser(source)
	statements, err := parser.GetStatements()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(70)
		}
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) (ast.Value, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) (ast.Value, error) {
	previous := i.environment
	i.environment = environment
	var lastValue ast.Value
	var err error
	for _, statement := range statements {
		lastValue, err = i.execute(statement)
		if err != nil {
			break
		}
	}
	i.environment = previous
	return lastValue, err
}

func (i *Interpreter) evaluate(ast ast.AST) (ast.Value, error) {
	return ast.Accept(i)
}
