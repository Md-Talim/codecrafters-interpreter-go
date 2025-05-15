package parser

import (
	"codecrafters-interpreter-go/internal/ast"
	"slices"
)

// peek returns the current token without consuming it.
func (p *Parser) peek() ast.Token {
	return *p.tokens[p.current]
}

// isAtEnd checks if the parser has reached the end of the tokens.
// It returns true if the current token is of type EofToken.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.EofToken
}

// previous returns the previous token.
// It is used to get the last token that was consumed.
func (p *Parser) previous() ast.Token {
	return *p.tokens[p.current-1]
}

// advance consumes the current token and returns it.
func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// check returns true if the current token is of the given type
func (p *Parser) check(expectedTokenType ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == expectedTokenType
}

// match checks if the current token has any of the given types
func (p *Parser) match(types ...ast.TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

// consume checks if the next token is of the expected type.
// If it is, it consumes the token and returns it.
func (p *Parser) consume(expectedTokenType ast.TokenType, message string) (ast.Token, error) {
	if p.check(expectedTokenType) {
		return p.advance(), nil
	}
	return ast.Token{}, newSyntaxError(p.peek(), message)
}
