package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/parser"
	"fmt"
	"os"
)

type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[ast.Expr]int
}

// NewInterpreter creates a new interpreter instance with a global environment.
func NewInterpreter() *Interpreter {
	globals := newEnvironment(nil)
	environment := globals
	globals.define("clock", NewClockFunction())
	return &Interpreter{environment: environment, globals: globals, locals: make(map[ast.Expr]int)}
}

// Interpret parses the source code and evaluates the expression.
func (i *Interpreter) Interpret(source string) (ast.Value, error) {
	parser := parser.NewParser(source)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}
	return i.evaluate(expr)
}

// Run parses the source code and executes the statements.
func (i *Interpreter) Run(statements []ast.Stmt) {
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(70)
		}
	}
}

// execute executes a statement and returns the result.
func (i *Interpreter) execute(stmt ast.Stmt) (ast.Value, error) {
	return stmt.Accept(i)
}

// executeBlock executes a block of statements in a new environment.
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

// evaluate evaluates an AST node and returns its value.
func (i *Interpreter) evaluate(ast ast.AST) (ast.Value, error) {
	return ast.Accept(i)
}

// Resolve resolves the variable's depth in the environment.
// This is used to keep track of the variable's scope.
func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}
