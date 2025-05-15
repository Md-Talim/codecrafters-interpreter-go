package parser

import (
	"codecrafters-interpreter-go/internal/ast"
	"strconv"
	"strings"
)

// primary parses the primary expressions.
// It handles literals (true, false, nil), identifiers, string literals, numbers, and grouping.
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
		value, err := strconv.ParseFloat(p.previous().Lexeme, 64)
		if err != nil {
			return nil, newSyntaxError(p.previous(), "Invalid number.")
		}
		return ast.NewNumberExpr(value), nil
	}
	if p.match(ast.IdentifierToken) {
		name := p.previous()
		return ast.NewVariableExpr(name), nil
	}

	if p.match(ast.LeftParenToken) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err = p.consume(ast.RightParenToken, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return ast.NewGroupingExpr(expr), nil
	}
	return nil, newSyntaxError(p.peek(), "Expect expression.")
}

// call parses function calls.
func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(ast.LeftParenToken) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

// finishCall finishes parsing a function call.
// It handles the arguments and the closing parenthesis.
func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}
	if !p.check(ast.RightParenToken) {
		for {
			if len(arguments) >= 255 {
				return nil, newSyntaxError(p.peek(), "Can't have more than 255 arguments.")
			}
			argument, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argument)
			if !p.match(ast.CommaToken) {
				break
			}
		}
	}

	paren, err := p.consume(ast.RightParenToken, "Expect ')' after argument.")
	if err != nil {
		return nil, err
	}
	return ast.NewCallExpr(callee, paren, arguments), nil
}

// unary parses unary expressions.
// It handles negation and logical NOT.
func (p *Parser) unary() (ast.Expr, error) {
	if p.match(ast.BangToken, ast.MinusToken) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnaryExpr(operator, right), nil
	}
	return p.call()
}

// factor parses factor expressions.
// It handles multiplication and division.
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

// term parses term expressions.
// It handles addition and subtraction.
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

// comparison parses comparison expressions.
// It handles greater than, less than, greater than or equal to, and less than or equal to.
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

// equality parses equality expressions.
// It handles equality and inequality.
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

// and parses logical AND expressions.
func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(ast.AndKeyword) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = ast.NewLogicalExpr(expr, operator, right)
	}

	return expr, nil
}

// or parses logical OR expressions.
func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(ast.OrKeyword) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = ast.NewLogicalExpr(expr, operator, right)
	}

	return expr, nil
}

// assignment parses assignment expressions.
// It handles variable assignments.
func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(ast.EqualToken) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*ast.VariableExpr); ok {
			return ast.NewAssignExpr(v.Name, value), nil
		}

		return nil, newSyntaxError(equals, "Invalid assignment target.")
	}

	return expr, nil
}

// expression parses a complete expression.
func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}
