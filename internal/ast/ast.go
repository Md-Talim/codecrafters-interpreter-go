package ast

type AST interface {
	Accept(visitor AstVisitor)
}

type AstVisitor interface {
	// Expression Nodes
	VisitAssignExpr(expr *AssignExpr)
	VisitBinaryExpr(expr *BinaryExpr)
	VisitBooleanExpr(expr *BooleanExpr)
	VisitCallExpr(expr *CallExpr)
	VisitGroupingExpr(expr *GroupingExpr)
	VisitLogicalExpr(expr *LogicalExpr)
	VisitNilExpr()
	VisitNumberExpr(expr *NumberExpr)
	VisitStringExpr(expr *StringExpr)
	VisitUnaryExpr(expr *UnaryExpr)
	VisitVariableExpr(expr *VariableExpr)

	// Statement Nodes
	VisitBlockStmt(stmt *BlockStmt)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitFunctionStmt(stmt *FunctionStmt)
	VisitIfStmt(stmt *IfStmt)
	VisitPrintStmt(stmt *PrintStmt)
	VisitVarStmt(stmt *VarStmt)
	VisitWhileStmt(stmt *WhileStmt)
}
