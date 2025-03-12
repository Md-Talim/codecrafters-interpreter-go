package main

import (
	"fmt"
	"slices"
)

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) primary() (Expr, error) {
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
		expr, _ := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{Expression: expr}, nil
	}
	return nil, fmt.Errorf("unexpected token")
}

func (p *Parser) expression() (Expr, error) {
	return p.primary()
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
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

func (p *Parser) consume(t TokenType, msg string) Token {
	if p.check(t) {
		return p.advance()
	}
	panic(msg)
}
