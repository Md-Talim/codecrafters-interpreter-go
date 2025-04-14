package scanner

import (
	"fmt"
	"os"

	"codecrafters-interpreter-go/internal/ast"
)

type ErrorHandler func(line int, message string)

type Scanner struct {
	source   []rune
	tokens   []*ast.Token
	start    int // points to the first character in the lexeme being scanned
	current  int // points at the character currently being considered
	line     int // tracks what source line 'current' is on
	hadError bool
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: []rune(source), line: 1, hadError: false}
}

func (s *Scanner) ScanTokens() ([]*ast.Token, bool) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, &ast.Token{Type: ast.EofToken, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens, s.hadError
}

func (s *Scanner) lexicalError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
	s.hadError = true
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	// single-character tokens
	case '(':
		s.addToken(ast.LeftParenToken)
	case ')':
		s.addToken(ast.RightParenToken)
	case '{':
		s.addToken(ast.LeftBraceToken)
	case '}':
		s.addToken(ast.RightBraceToken)
	case ',':
		s.addToken(ast.CommaToken)
	case '.':
		s.addToken(ast.DotToken)
	case '-':
		s.addToken(ast.MinusToken)
	case '+':
		s.addToken(ast.PlusToken)
	case ';':
		s.addToken(ast.SemicolonToken)
	case '*':
		s.addToken(ast.StarToken)
	// operators
	case '!':
		if s.match('=') {
			s.addToken(ast.BangEqualToken)
		} else {
			s.addToken(ast.BangToken)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EqualEqualToken)
		} else {
			s.addToken(ast.EqualToken)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LessEqualToken)
		} else {
			s.addToken(ast.LessToken)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GreaterEqualToken)
		} else {
			s.addToken(ast.GreaterToken)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(ast.SlashToken)
		}
	// skip over other meaningless characters: newlines and whitespace
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			s.lexicalError(s.line, fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lexicalError(s.line, "Unterminated string.")
		return
	}

	// the closing ".
	s.advance()

	// trim the surrounding quotes
	value := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(ast.StringToken, value)
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// look for fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// consume the .
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	// grab the number as text & parse into float
	text := string(s.source[s.start:s.current])
	num, err := formatFloat(text)
	if err != nil {
		s.lexicalError(s.line, "Invalid number literal.")
		return
	}

	s.addTokenWithLiteral(ast.NumberToken, num)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, ok := ast.Keywords[text]
	if !ok {
		tokenType = ast.IdentifierToken
	}
	s.addToken(tokenType)
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

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	token := s.source[s.current]
	s.current++
	return token
}

func (s *Scanner) addToken(t ast.TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t ast.TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, &ast.Token{Type: t, Lexeme: text, Literal: literal, Line: s.line})
}
