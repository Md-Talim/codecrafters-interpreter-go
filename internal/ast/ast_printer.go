package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

// VisitCallExpr implements AstVisitor.
func (ap *AstPrinter) VisitCallExpr(expr *CallExpr) (Value, error) {
	panic("unimplemented")
}

// VisitFunctionStmt implements AstVisitor.
func (ap *AstPrinter) VisitFunctionStmt(stmt *FunctionStmt) (Value, error) {
	panic("unimplemented")
}

// VisitWhileStmt implements AstVisitor.
func (ap *AstPrinter) VisitWhileStmt(stmt *WhileStmt) (Value, error) {
	panic("unimplemented")
}

// VisitLogicalExpr implements AstVisitor.
func (ap *AstPrinter) VisitLogicalExpr(expr *LogicalExpr) (Value, error) {
	panic("unimplemented")
}

// VisitIfStmt implements AstVisitor.
func (ap *AstPrinter) VisitIfStmt(stmt *IfStmt) (Value, error) {
	panic("unimplemented")
}

// VisitBlockStmt implements AstVisitor.
func (ap *AstPrinter) VisitBlockStmt(stmt *BlockStmt) (Value, error) {
	panic("unimplemented")
}

// VisitAssignExpr implements AstVisitor.
func (ap *AstPrinter) VisitAssignExpr(expr *AssignExpr) (Value, error) {
	panic("unimplemented")
}

// VisitVarStmt implements AstVisitor.
func (ap *AstPrinter) VisitVarStmt(stmt *VarStmt) (Value, error) {
	panic("unimplemented")
}

// VisitVariableExpr implements AstVisitor.
func (ap *AstPrinter) VisitVariableExpr(expr *VariableExpr) (Value, error) {
	panic("unimplemented")
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
