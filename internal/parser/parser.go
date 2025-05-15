package parser

import (
	"errors"
	"fmt"
	"os"

	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/scanner"
)

// Parser implements a recursive descent parser for the Lox language.
// The grammar follows this precedence hierarchy:
// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment | logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | call ;
// call           → primary ( "(" arguments? ")" )* ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
type Parser struct {
	tokens  []*ast.Token
	current int
}

// NewParser creates a new Parser instance with the given source code.
func NewParser(source string) *Parser {
	scanner := scanner.NewScanner(source)
	tokens, hadError := scanner.ScanTokens()
	if hadError {
		os.Exit(65)
	}
	return &Parser{tokens: tokens}
}

// GetStatements parses the source code and returns a slice of statements.
func (p *Parser) GetStatements() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

// Parse parses the source code and returns the root expression of the AST.
func (p *Parser) Parse() (ast.Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

// newSyntaxError creates a new syntax error with the given token and message.
// It formats the error message to include the line number and the token's lexeme.
func newSyntaxError(token ast.Token, message string) error {
	var where string
	if token.Type == ast.EofToken {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%s' ", token.Lexeme)
	}

	text := fmt.Sprintf("[line %d] Error %s: %s", token.Line, where, message)
	return errors.New(text)
}
