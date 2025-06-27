package parser

import "codecrafters-interpreter-go/internal/ast"

// expressionStatement parses an expression statement.
func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(ast.SemicolonToken, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ast.NewExpressionStmt(expr), nil
}

// printStatement parses a print statement.
// It handles the expression to be printed and the semicolon at the end.
func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(ast.SemicolonToken, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return ast.NewPrintStmt(value), nil
}

// ifStatement parses an if statement.
// It handles the condition and the then and else branches.
func (p *Parser) ifStatement() (ast.Stmt, error) {
	if _, err := p.consume(ast.LeftParenToken, "Expect '(' after 'if'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(ast.RightParenToken, "Expect ')' after if condition."); err != nil {
		return nil, err
	}

	thenBrach, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt = nil
	if p.match(ast.ElseKeyword) {
		if elseBranch, err = p.statement(); err != nil {
			return nil, err
		}
	}

	return ast.NewIfStmt(condition, thenBrach, elseBranch), nil
}

// whileStatement parses a while statement.
// It handles the condition and the body of the loop.
func (p *Parser) whileStatement() (ast.Stmt, error) {
	if _, err := p.consume(ast.LeftParenToken, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(ast.RightParenToken, "Expect ')' after while condition."); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return ast.NewWhileStmt(condition, body), nil
}

// forStatement parses a for statement.
// It handles the loop initializer, condition, increment, and body.
func (p *Parser) forStatement() (ast.Stmt, error) {
	var (
		initializer ast.Stmt
		condition   ast.Expr = nil
		increment   ast.Expr = nil
		err         error
	)

	if _, err = p.consume(ast.LeftParenToken, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}

	// Parse the loop initializer:
	// - If semicolon is found: no initializer (empty initialization)
	// - If 'var' keyword is found: parse variable declaration (e.g., var i = 0;)
	// - Otherwise: parse expression statement (e.g., i = 0;)
	// Note: Both varDeclaration and expressionStatement consume the trailing semicolon
	if p.match(ast.SemicolonToken) {
		initializer = nil
	} else if p.match(ast.VarKeyword) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// Parse the loop condition
	if !p.check(ast.SemicolonToken) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(ast.SemicolonToken, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	// Parse the increment statement
	if !p.check(ast.RightParenToken) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(ast.RightParenToken, "Expect ')' after for clauses."); err != nil {
		return nil, err
	}

	// Parse the loop body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		statements := []ast.Stmt{body, ast.NewExpressionStmt(increment)}
		body = ast.NewBlockStmt(statements)
	}
	if condition == nil {
		condition = ast.NewBooleanExpr(true)
	}
	body = ast.NewWhileStmt(condition, body)

	if initializer != nil {
		statements := []ast.Stmt{initializer, body}
		body = ast.NewBlockStmt(statements)
	}

	return body, nil
}

// block parses a block of statements.
// It handles the opening and closing braces and collects the statements in between.
func (p *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.check(ast.RightBraceToken) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}

	if _, err := p.consume(ast.RightBraceToken, "Expect '}' after block."); err != nil {
		return nil, err
	}
	return statements, nil
}

// returnStatement parses a return statement.
// It handles the return keyword, optional return value, and the semicolon at the end.
func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()
	var value ast.Expr = nil
	var err error
	if !p.check(ast.SemicolonToken) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(ast.SemicolonToken, "Expect ';' after return value."); err != nil {
		return nil, err
	}
	return ast.NewReturnStmt(keyword, value), nil
}

// statement parses a statement.
// It checks for various statement types (for, while, if, print, block) and delegates to the appropriate parsing function.
func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(ast.ForKeyword) {
		return p.forStatement()
	}
	if p.match(ast.WhileKeyword) {
		return p.whileStatement()
	}
	if p.match(ast.IfKeyword) {
		return p.ifStatement()
	}
	if p.match(ast.PrintKeyword) {
		return p.printStatement()
	}
	if p.match(ast.ReturnKeyword) {
		return p.returnStatement()
	}
	if p.match(ast.LeftBraceToken) {
		blocks, err := p.block()
		if err != nil {
			return nil, err
		}
		return ast.NewBlockStmt(blocks), nil
	}

	return p.expressionStatement()
}

// varDeclaration parses a variable declaration.
// It handles the variable name, optional initializer, and the semicolon at the end.
func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(ast.IdentifierToken, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(ast.EqualToken) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(ast.SemicolonToken, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return ast.NewVarStmt(name, initializer), nil
}

// function parses a function declaration.
// It handles the function name, parameter list, and body.
func (p *Parser) function(kind string) (ast.Stmt, error) {
	// Consume the identifier token for the function’s name
	name, err := p.consume(ast.IdentifierToken, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	// Parse the parameter list and the pair of parentheses wrapped around it
	if _, err := p.consume(ast.LeftParenToken, "Expect '(' after "+kind+" name."); err != nil {
		return nil, err
	}
	parameters := []ast.Token{}
	if !p.check(ast.RightParenToken) {
		for {
			if len(parameters) >= 255 {
				return nil, newSyntaxError(p.peek(), "Can't have more than 255 parameters.")
			}

			param, err := p.consume(ast.IdentifierToken, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, param)

			if !p.match(ast.CommaToken) {
				break
			}
		}
	}
	if _, err := p.consume(ast.RightParenToken, "Expect ')' after parameters."); err != nil {
		return nil, err
	}

	// parse the body and wrap it all up in a function node
	if _, err := p.consume(ast.LeftBraceToken, "Expect '{' before "+kind+" body."); err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return ast.NewFunctionStmt(name, parameters, body), nil
}

// classDeclaration parses a class declaration.
// It handles the class name, optional superclass, and method
func (p *Parser) classDeclaration() (ast.Stmt, error) {
	name, err := p.consume(ast.IdentifierToken, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *ast.VariableExpr = nil
	if p.match(ast.LessToken) {
		if _, err := p.consume(ast.IdentifierToken, "Expect superclass name."); err != nil {
			return nil, err
		}
		superclass = ast.NewVariableExpr(p.previous())
	}

	if _, err := p.consume(ast.LeftBraceToken, "Expect '{' before class body."); err != nil {
		return nil, err
	}

	methods := []ast.FunctionStmt{}
	for !p.check(ast.RightBraceToken) && !p.isAtEnd() {
		method, err := p.function("method")
		if err != nil {
			return nil, err
		}
		functionStmt, ok := method.(*ast.FunctionStmt)
		if !ok {
			return nil, newSyntaxError(p.peek(), "Expect method declaration.")
		}
		methods = append(methods, *functionStmt)
	}
	if _, err := p.consume(ast.RightBraceToken, "Expect '}' after class body."); err != nil {
		return nil, err
	}

	return ast.NewClassStmt(name, methods, superclass), nil
}

// declaration parses a declaration statement.
// It checks for function declarations, variable declarations, and other statements.
func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(ast.ClassKeyword) {
		return p.classDeclaration()
	}
	if p.match(ast.FunKeyword) {
		return p.function("function")
	}
	if p.match(ast.VarKeyword) {
		return p.varDeclaration()
	}
	return p.statement()
}
