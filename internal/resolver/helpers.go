package resolver

import (
	"codecrafters-interpreter-go/internal/ast"
	"errors"
	"fmt"
)

// beginScope initializes a new scope for variable resolution.
func (r *Resolver) beginScope() {
	r.scopes.push(newScope())
}

// endScope removes the most recent scope from the scopes stack.
func (r *Resolver) endScope() {
	r.scopes.pop()
}

// declare marks a variable as declared in the current scope.
func (r *Resolver) declare(name ast.Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope := r.scopes.peek()
	scope.set(name.Lexeme, false)
}

// define marks a variable as defined in the current scope.
func (r *Resolver) define(name ast.Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope := r.scopes.peek()
	scope.set(name.Lexeme, true)
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
