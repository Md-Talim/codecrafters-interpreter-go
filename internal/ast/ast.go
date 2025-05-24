package ast

type AST interface {
	Accept(visitor AstVisitor) (Value, error)
}

type AstVisitor interface {
	// Expression Nodes
	VisitAssignExpr(expr *AssignExpr) (Value, error)
	VisitBinaryExpr(expr *BinaryExpr) (Value, error)
	VisitBooleanExpr(expr *BooleanExpr) (Value, error)
	VisitCallExpr(expr *CallExpr) (Value, error)
	VisitGetExpr(expr *GetExpr) (Value, error)
	VisitGroupingExpr(expr *GroupingExpr) (Value, error)
	VisitLogicalExpr(expr *LogicalExpr) (Value, error)
	VisitNilExpr() (Value, error)
	VisitNumberExpr(expr *NumberExpr) (Value, error)
	VisitSetExpr(expr *SetExpr) (Value, error)
	VisitStringExpr(expr *StringExpr) (Value, error)
	VisitUnaryExpr(expr *UnaryExpr) (Value, error)
	VisitVariableExpr(expr *VariableExpr) (Value, error)

	// Statement Nodes
	VisitBlockStmt(stmt *BlockStmt) (Value, error)
	VisitClassStmt(stmt *ClassStmt) (Value, error)
	VisitExpressionStmt(stmt *ExpressionStmt) (Value, error)
	VisitFunctionStmt(stmt *FunctionStmt) (Value, error)
	VisitIfStmt(stmt *IfStmt) (Value, error)
	VisitPrintStmt(stmt *PrintStmt) (Value, error)
	VisitReturnStmt(stmt *ReturnStmt) (Value, error)
	VisitVarStmt(stmt *VarStmt) (Value, error)
	VisitWhileStmt(stmt *WhileStmt) (Value, error)
}
