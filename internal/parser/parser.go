package parser

import (
	"slices"

	"codecrafters-interpreter-go/internal/ast"
)

type Parser[R any] struct {
	tokens  []*ast.Token
	current int
}

type ParseError struct {
	token   ast.Token
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

func (e *ParseError) Token() ast.Token {
	return e.token
}

func NewParser[R any](tokens []*ast.Token) *Parser[R] {
	return &Parser[R]{tokens: tokens}
}

func (p *Parser[R]) error(token ast.Token, message string) *ParseError {
	return &ParseError{token: token, message: message}
}

func (p *Parser[R]) primary() (ast.Expr[R], *ParseError) {
	if p.match(ast.TrueKeyword) {
		return &ast.Literal[R]{Value: true}, nil
	}
	if p.match(ast.FalseKeyword) {
		return &ast.Literal[R]{Value: false}, nil
	}
	if p.match(ast.NilKeyword) {
		return &ast.Literal[R]{Value: nil}, nil
	}
	if p.match(ast.StringToken, ast.NumberToken) {
		return &ast.Literal[R]{Value: p.previous().Literal}, nil
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
		return &ast.Grouping[R]{Expression: expr}, nil
	}
	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser[R]) unary() (ast.Expr[R], *ParseError) {
	if p.match(ast.BangToken, ast.MinusToken) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary[R]{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser[R]) factor() (ast.Expr[R], *ParseError) {
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
		expr = &ast.Binary[R]{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser[R]) term() (ast.Expr[R], *ParseError) {
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
		expr = &ast.Binary[R]{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser[R]) comparison() (ast.Expr[R], *ParseError) {
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
		expr = &ast.Binary[R]{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser[R]) equality() (ast.Expr[R], *ParseError) {
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
		expr = &ast.Binary[R]{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser[R]) expression() (ast.Expr[R], *ParseError) {
	return p.equality()
}

func (p *Parser[R]) Parse() (ast.Expr[R], *ParseError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser[R]) peek() ast.Token {
	return *p.tokens[p.current]
}

func (p *Parser[R]) isAtEnd() bool {
	return p.peek().Type == ast.EofToken
}

func (p *Parser[R]) previous() ast.Token {
	return *p.tokens[p.current-1]
}

func (p *Parser[R]) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser[R]) check(t ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser[R]) match(types ...ast.TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser[R]) consume(t ast.TokenType, message string) (ast.Token, *ParseError) {
	if p.check(t) {
		return p.advance(), nil
	}
	return ast.Token{}, p.error(p.peek(), message)
}
