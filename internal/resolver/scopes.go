package resolver

import "errors"

// Scope represents a single scope in the resolver.
type Scope map[string]bool

func newScope() Scope {
	return Scope{}
}

// set adds a value to the scope.
func (s *Scope) set(name string, value bool) {
	(*s)[name] = value
}

// get retrieves a value from the scope.
func (s *Scope) get(name string) bool {
	if value, exists := (*s)[name]; exists {
		return value
	}
	return false
}

// hasKey checks if the scope contains a key.
func (s *Scope) hasKey(name string) bool {
	_, exists := (*s)[name]
	return exists
}

// Scopes is a stack of scopes used for variable resolution.
type Scopes []Scope

func newScopes() Scopes {
	return Scopes{}
}

// get return the scope at the given index.
func (s *Scopes) get(index int) Scope {
	if index < 0 || index >= len(*s) {
		return nil
	}
	return (*s)[index]
}

// push adds a new scope to the stack.
func (s *Scopes) push(scope Scope) {
	*s = append(*s, scope)
}

// pop removes the top scope from the stack and returns it.
func (s *Scopes) pop() (Scope, error) {
	if len(*s) == 0 {
		return nil, errors.New("no scopes to pop")
	}
	scope := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return scope, nil
}

// peek returns the top scope without removing it.
func (s *Scopes) peek() Scope {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

// isEmpty checks if the stack of scopes is empty.
func (s *Scopes) isEmpty() bool {
	return len(*s) == 0
}

// size returns the number of scopes in the stack.
func (s *Scopes) size() int {
	return len(*s)
}
