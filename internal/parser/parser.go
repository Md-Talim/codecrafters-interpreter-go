package parser

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/scanner"
)

type Parser struct {
	tokens      []*ast.Token
	current     int
	syntaxError bool
}

func NewParser(source string) *Parser {
	scanner := scanner.NewScanner(source)
	tokens, hadError := scanner.ScanTokens()
	if hadError {
		os.Exit(65)
	}
	return &Parser{tokens: tokens, syntaxError: false}
}

func (p *Parser) error(token ast.Token, message string) error {
	var where string
	if token.Type == ast.EofToken {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%s' ", token.Lexeme)
	}

	text := fmt.Sprintf("[line %d] Error %s: %s", token.Line, where, message)
	return errors.New(text)
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(ast.TrueKeyword) {
		return ast.NewBooleanExpr(true), nil
	}
	if p.match(ast.FalseKeyword) {
		return ast.NewBooleanExpr(false), nil
	}
	if p.match(ast.NilKeyword) {
		return ast.NewNilExpr(), nil
	}
	if p.match(ast.StringToken) {
		value := strings.Trim(p.previous().Lexeme, "\"")
		return ast.NewStringExpr(value), nil
	}
	if p.match(ast.NumberToken) {
		value, _ := strconv.ParseFloat(p.previous().Lexeme, 64)
		return ast.NewNumberExpr(value), nil
	}
	if p.match(ast.IdentifierToken) {
		name := p.previous()
		return ast.NewVariableExpr(name), nil
	}

	if p.match(ast.LeftParenToken) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(ast.RightParenToken, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return ast.NewGroupingExpr(expr), nil
	}
	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(ast.BangToken, ast.MinusToken) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnaryExpr(operator, right), nil
	}
	return p.primary()
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(ast.SlashToken, ast.StarToken) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinaryExpr(operator, expr, right)
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(ast.PlusToken, ast.MinusToken) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinaryExpr(operator, expr, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(ast.GreaterToken, ast.GreaterEqualToken, ast.LessToken, ast.LessEqualToken) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinaryExpr(operator, expr, right)
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(ast.EqualEqualToken, ast.BangEqualToken) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinaryExpr(operator, expr, right)
	}

	return expr, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(ast.SemicolonToken, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ast.NewExpressionStmt(expr), nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(ast.SemicolonToken, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ast.NewPrintStmt(value), nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(ast.PrintKeyword) {
		statement, err := p.printStatement()
		if err != nil {
			return nil, err
		}
		return statement, nil
	}

	return p.expressionStatement()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(ast.IdentifierToken, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(ast.EqualToken) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(ast.SemicolonToken, "Expect ';' after variable declaration.")
	return ast.NewVarStmt(name, initializer), nil
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(ast.VarKeyword) {
		return p.varDeclaration()
	}
	return p.statement()
}

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

func (p *Parser) Parse() (ast.Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) peek() ast.Token {
	return *p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.EofToken
}

func (p *Parser) previous() ast.Token {
	return *p.tokens[p.current-1]
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(types ...ast.TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) consume(t ast.TokenType, message string) (ast.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return ast.Token{}, p.error(p.peek(), message)
}
