package main

import "fmt"

type Scanner struct {
	source  []rune
	tokens  []*Token
	start   int // points to the first character in the lexeme being scanned
	current int // points at the character currently being considered
	line    int // tracks what source line 'current' is on
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: []rune(source), line: 1}
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, &Token{Type: EOF, Lexeme: "", Literal: nil, Line: s.line})
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	// single-character tokens
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	// operators
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	// skip over other meaningless characters: newlines and whitespace
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	default:
		lox.error(s.line, fmt.Sprintf("Unexpected character: %c", c))
	}
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	token := s.source[s.current]
	s.current++
	return token
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, &Token{Type: t, Lexeme: text, Literal: literal, Line: s.line})
}
