package parser

import (
	"os"
	"slices"
	"strconv"
	"strings"

	"codecrafters-interpreter-go/internal/ast"
	"codecrafters-interpreter-go/internal/scanner"
	"codecrafters-interpreter-go/pkg/loxerrors"
)

type Parser struct {
	tokens  []*ast.Token
	current int
}

func NewParser(source string) *Parser {
	scanner := scanner.NewScanner(source)
	tokens, hadError := scanner.ScanTokens()
	if hadError {
		os.Exit(65)
	}
	return &Parser{tokens: tokens}
}

func (p *Parser) error(token ast.Token, message string) error {
	return loxerrors.NewParseError(token.Line, message)
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
