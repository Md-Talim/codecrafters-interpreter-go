package ast

type AST interface {
	Accept(visitor AstVisitor)
}

type AstVisitor interface {
	// Expression Nodes
	VisitBinaryExpr(expr *BinaryExpr)
	VisitBooleanExpr(expr *BooleanExpr)
	VisitGroupingExpr(expr *GroupingExpr)
	VisitNilExpr()
	VisitNumberExpr(expr *NumberExpr)
	VisitStringExpr(expr *StringExpr)
	VisitUnaryExpr(expr *UnaryExpr)
	VisitVariableExpr(expr *VariableExpr)

	// Statement Nodes
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitPrintStmt(stmt *PrintStmt)
	VisitVarStmt(stmt *VarStmt)
}
