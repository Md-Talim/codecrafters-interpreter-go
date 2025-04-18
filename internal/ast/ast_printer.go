package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

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
