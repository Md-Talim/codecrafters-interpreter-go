package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

// VisitSuperExpr implements AstVisitor.
func (ap *AstPrinter) VisitSuperExpr(expr *SuperExpr) (Value, error) {
	fmt.Println(formatSExpr("super", expr.Method))
	return nil, nil
}

// VisitThisExpr implements AstVisitor.
func (ap *AstPrinter) VisitThisExpr(expr *ThisExpr) (Value, error) {
	fmt.Printf("this")
	return nil, nil
}

// VisitGetExpr implements AstVisitor.
func (ap *AstPrinter) VisitGetExpr(expr *GetExpr) (Value, error) {
	fmt.Println(formatSExpr("get", expr.Object, expr.Name))
	return nil, nil
}

// VisitSetExpr implements AstVisitor.
func (ap *AstPrinter) VisitSetExpr(expr *SetExpr) (Value, error) {
	fmt.Println(formatSExpr("set", expr.Object, expr.Name, expr.Value))
	return nil, nil
}

// VisitClassStmt implements AstVisitor.
func (ap *AstPrinter) VisitClassStmt(stmt *ClassStmt) (Value, error) {
	methods := make([]any, len(stmt.Methods))
	for i, method := range stmt.Methods {
		methods[i] = method
	}

	result := formatSExpr("class "+stmt.Name.Lexeme, methods...)
	fmt.Println(result)
	return nil, nil
}

// VisitReturnStmt implements AstVisitor.
func (ap *AstPrinter) VisitReturnStmt(stmt *ReturnStmt) (Value, error) {
	if stmt.Value != nil {
		fmt.Println(ap.parenthesize("return", stmt.Value))
	}
	fmt.Println("(return)")
	return nil, nil
}

// VisitCallExpr implements AstVisitor.
func (ap *AstPrinter) VisitCallExpr(expr *CallExpr) (Value, error) {
	fmt.Println(formatSExpr("call", expr.Callee, expr.Arguments))
	return nil, nil
}

// VisitFunctionStmt implements AstVisitor.
func (ap *AstPrinter) VisitFunctionStmt(stmt *FunctionStmt) (Value, error) {
	params := make([]string, len(stmt.Params))
	for i, token := range stmt.Params {
		params[i] = token.Lexeme
	}

	paramList := strings.Join(params, " ")
	result := formatSExpr("fun "+stmt.Name.Lexeme+"("+paramList+")", stmt.Body)
	fmt.Println(result)
	return nil, nil
}

// VisitWhileStmt implements AstVisitor.
func (ap *AstPrinter) VisitWhileStmt(stmt *WhileStmt) (Value, error) {
	fmt.Println(formatSExpr("while", stmt.Condition, stmt.Body))
	return nil, nil
}

// VisitLogicalExpr implements AstVisitor.
func (ap *AstPrinter) VisitLogicalExpr(expr *LogicalExpr) (Value, error) {
	fmt.Println(formatSExpr(expr.Operator.Lexeme, expr.Left, expr.Right))
	return nil, nil
}

// VisitIfStmt implements AstVisitor.
func (ap *AstPrinter) VisitIfStmt(stmt *IfStmt) (Value, error) {
	if stmt.ElseBranch != nil {
		fmt.Println(formatSExpr("if-else", stmt.Condition, stmt.ThenBranch, stmt.ElseBranch))
	} else {
		fmt.Println(formatSExpr("if", stmt.Condition, stmt.ThenBranch))
	}
	return nil, nil
}

// VisitBlockStmt implements AstVisitor.
func (ap *AstPrinter) VisitBlockStmt(stmt *BlockStmt) (Value, error) {
	var parts []any
	for _, s := range stmt.Statements {
		parts = append(parts, s)
	}
	fmt.Println(formatSExpr("block", parts...))
	return nil, nil
}

// VisitAssignExpr implements AstVisitor.
func (ap *AstPrinter) VisitAssignExpr(expr *AssignExpr) (Value, error) {
	fmt.Println(formatSExpr("=", expr.Name, expr.Value))
	return nil, nil
}

// VisitVarStmt implements AstVisitor.
func (ap *AstPrinter) VisitVarStmt(stmt *VarStmt) (Value, error) {
	if stmt.Initializer != nil {
		fmt.Println(formatSExpr("var", stmt.Name, "=", stmt.Initializer))
	} else {
		fmt.Println(formatSExpr("var", stmt.Name))
	}
	return nil, nil
}

// VisitVariableExpr implements AstVisitor.
func (ap *AstPrinter) VisitVariableExpr(expr *VariableExpr) (Value, error) {
	fmt.Printf("%s", expr.Name.Lexeme)
	return nil, nil
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) VisitNumberExpr(num *NumberExpr) (Value, error) {
	numStr := strings.TrimRight(fmt.Sprintf("%f", num.Value), "0")
	if numStr[len(numStr)-1] == uint8('.') {
		numStr = numStr + "0"
	}
	fmt.Print(numStr)
	return nil, nil
}

func (ap *AstPrinter) VisitBooleanExpr(be *BooleanExpr) (Value, error) {
	fmt.Printf("%t", be.Value)
	return nil, nil
}

func (ap *AstPrinter) VisitPrintStmt(ps *PrintStmt) (Value, error) {
	fmt.Print(ps.Expression)
	return nil, nil
}

func (ap *AstPrinter) VisitExpressionStmt(ex *ExpressionStmt) (Value, error) {
	fmt.Print(ex.Expression)
	return nil, nil
}

func (ap *AstPrinter) VisitNilExpr() (Value, error) {
	fmt.Printf("nil")
	return nil, nil
}

func (ap *AstPrinter) VisitStringExpr(str *StringExpr) (Value, error) {
	fmt.Printf("%s", str.Value)
	return nil, nil
}

func (ap *AstPrinter) VisitGroupingExpr(group *GroupingExpr) (Value, error) {
	fmt.Printf("(group ")
	group.Expression.Accept(ap)
	fmt.Printf(")")
	return nil, nil
}

func (ap *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) (Value, error) {
	opString := expr.Operator.Lexeme
	fmt.Printf("(%s ", opString)
	expr.Right.Accept(ap)
	fmt.Printf(")")
	return nil, nil
}

func (ap *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) (Value, error) {
	fmt.Printf("(%s ", expr.Operator.Lexeme)
	expr.Left.Accept(ap)
	fmt.Printf(" ")
	expr.Right.Accept(ap)
	fmt.Printf(")")
	return nil, nil
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var sb strings.Builder
	sb.WriteString("(" + name)
	for _, expr := range exprs {
		str, err := expr.Accept(ap)
		if err != nil {
			return "", err
		}
		sb.WriteString(" ")
		sb.WriteString(str.String())
	}
	sb.WriteString(")")
	return sb.String(), nil
}

func formatSExpr(name string, parts ...any) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	writeParts(&builder, parts...)
	builder.WriteString(")")
	return builder.String()
}

func writeParts(builder *strings.Builder, parts ...any) {
	for _, part := range parts {
		builder.WriteString(" ")
		switch v := part.(type) {
		case Stmt:
			str, _ := v.Accept(NewAstPrinter())
			builder.WriteString(fmt.Sprint(str))
		case Expr:
			str, _ := v.Accept(NewAstPrinter()) // or reuse existing visitor
			builder.WriteString(fmt.Sprint(str))
		case Token:
			builder.WriteString(v.Lexeme)
		case string:
			builder.WriteString(v)
		case []any:
			writeParts(builder, v...)
		default:
			builder.WriteString(fmt.Sprint(v))
		}
	}
}
