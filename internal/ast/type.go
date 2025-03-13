package ast

type TokenType int

const (
	LeftParenToken TokenType = iota
	RightParenToken
	LeftBraceToken
	RightBraceToken
	CommaToken
	DotToken
	MinusToken
	PlusToken
	SemicolonToken
	StarToken
	EqualToken
	EqualEqualToken
	BangToken
	BangEqualToken
	LessToken
	LessEqualToken
	GreaterToken
	GreaterEqualToken
	SlashToken
	StringToken
	NumberToken
	IdentifierToken

	AndKeyword
	ClassKeyword
	ElseKeyword
	FalseKeyword
	ForKeyword
	FunKeyword
	IfKeyword
	NilKeyword
	OrKeyword
	PrintKeyword
	ReturnKeyword
	SuperKeyword
	ThisKeyword
	TrueKeyword
	VarKeyword
	WhileKeyword

	EofToken
)

var tokenTypeNames = map[TokenType]string{
	LeftParenToken:    "LEFT_PAREN",
	RightParenToken:   "RIGHT_PAREN",
	LeftBraceToken:    "LEFT_BRACE",
	RightBraceToken:   "RIGHT_BRACE",
	CommaToken:        "COMMA",
	DotToken:          "DOT",
	MinusToken:        "MINUS",
	PlusToken:         "PLUS",
	SemicolonToken:    "SEMICOLON",
	StarToken:         "SLASH",
	EqualToken:        "STAR",
	EqualEqualToken:   "BANG",
	BangToken:         "BANG_EQUAL",
	BangEqualToken:    "EQUAL",
	LessToken:         "EQUAL_EQUAL",
	LessEqualToken:    "GREATER",
	GreaterToken:      "GREATER_EQUAL",
	GreaterEqualToken: "LESS",
	SlashToken:        "LESS_EQUAL",
	StringToken:       "STRING",
	NumberToken:       "NUMBER",
	IdentifierToken:   "IDENTIFIER",
	AndKeyword:        "AND",
	ClassKeyword:      "CLASS",
	ElseKeyword:       "ELSE",
	FalseKeyword:      "FALSE",
	ForKeyword:        "FOR",
	FunKeyword:        "FUN",
	IfKeyword:         "IF",
	NilKeyword:        "NIL",
	OrKeyword:         "OR",
	PrintKeyword:      "PRINT",
	ReturnKeyword:     "RETURN",
	SuperKeyword:      "SUPER",
	ThisKeyword:       "THIS",
	TrueKeyword:       "TRUE",
	VarKeyword:        "VAR",
	WhileKeyword:      "WHILE",
	EofToken:          "EOF",
}

var Keywords = map[string]TokenType{
	"and":    AndKeyword,
	"class":  ClassKeyword,
	"else":   ElseKeyword,
	"false":  FalseKeyword,
	"for":    ForKeyword,
	"fun":    FunKeyword,
	"if":     IfKeyword,
	"nil":    NilKeyword,
	"or":     OrKeyword,
	"print":  PrintKeyword,
	"return": ReturnKeyword,
	"super":  SuperKeyword,
	"this":   ThisKeyword,
	"true":   TrueKeyword,
	"var":    VarKeyword,
	"while":  WhileKeyword,
}
