package main

import (
	"slices"
)

type Parser struct {
	tokens  []*Token
	current int
}

type ParseError struct {
	token   Token
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) error(token Token, message string) *ParseError {
	return &ParseError{token: token, message: message}
}

func (p *Parser) primary() (Expr, *ParseError) {
	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping{Expression: expr}, nil
	}
	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) unary() (Expr, *ParseError) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) factor() (Expr, *ParseError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, *ParseError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, *ParseError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, *ParseError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, *ParseError) {
	return p.equality()
}

func (p *Parser) Parse() (Expr, *ParseError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) peek() Token {
	return *p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) previous() Token {
	return *p.tokens[p.current-1]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) consume(t TokenType, message string) (Token, *ParseError) {
	if p.check(t) {
		return p.advance(), nil
	}
	return Token{}, p.error(p.peek(), message)
}
