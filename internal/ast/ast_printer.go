package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

// VisitWhileStmt implements AstVisitor.
func (ap *AstPrinter) VisitWhileStmt(stmt *WhileStmt) {
	panic("unimplemented")
}

// VisitLogicalExpr implements AstVisitor.
func (ap *AstPrinter) VisitLogicalExpr(expr *LogicalExpr) {
	panic("unimplemented")
}

// VisitIfStmt implements AstVisitor.
func (ap *AstPrinter) VisitIfStmt(stmt *IfStmt) {
	panic("unimplemented")
}

// VisitBlockStmt implements AstVisitor.
func (ap *AstPrinter) VisitBlockStmt(stmt *BlockStmt) {
	panic("unimplemented")
}

// VisitAssignExpr implements AstVisitor.
func (ap *AstPrinter) VisitAssignExpr(expr *AssignExpr) {
	panic("unimplemented")
}

// VisitVarStmt implements AstVisitor.
func (ap *AstPrinter) VisitVarStmt(stmt *VarStmt) {
	panic("unimplemented")
}

// VisitVariableExpr implements AstVisitor.
func (ap *AstPrinter) VisitVariableExpr(expr *VariableExpr) {
	panic("unimplemented")
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) VisitNumberExpr(num *NumberExpr) {
	numStr := strings.TrimRight(fmt.Sprintf("%f", num.Value), "0")
	if numStr[len(numStr)-1] == uint8('.') {
		numStr = numStr + "0"
	}
	fmt.Print(numStr)
}

func (ap *AstPrinter) VisitBooleanExpr(be *BooleanExpr) {
	fmt.Printf("%t", be.Value)
}

func (ap *AstPrinter) VisitPrintStmt(ps *PrintStmt) {
	fmt.Print(ps.Expression)
}

func (ap *AstPrinter) VisitExpressionStmt(ex *ExpressionStmt) {
	fmt.Print(ex.Expression)
}

func (ap *AstPrinter) VisitNilExpr() {
	fmt.Printf("nil")
}

func (ap *AstPrinter) VisitStringExpr(str *StringExpr) {
	fmt.Printf("%s", str.Value)
}

func (ap *AstPrinter) VisitGroupingExpr(group *GroupingExpr) {
	fmt.Printf("(group ")
	group.Expression.Accept(ap)
	fmt.Printf(")")
}

func (ap *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) {
	opString := expr.Operator.Lexeme
	fmt.Printf("(%s ", opString)
	expr.Right.Accept(ap)
	fmt.Printf(")")
}

func (ap *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) {
	fmt.Printf("(%s ", expr.Operator.Lexeme)
	expr.Left.Accept(ap)
	fmt.Printf(" ")
	expr.Right.Accept(ap)
	fmt.Printf(")")
}
