package ast

// ExpressionType represents the type of an expression.
type ExpressionType string

// Expression type constants.
const (
	ExpressionTypeBinary      ExpressionType = "binary"
	ExpressionTypeUnary       ExpressionType = "unary"
	ExpressionTypeCall        ExpressionType = "call"
	ExpressionTypeMember      ExpressionType = "member"
	ExpressionTypeAssignment  ExpressionType = "assignment"
	ExpressionTypeConditional ExpressionType = "conditional"
	ExpressionTypeNew         ExpressionType = "new"
	ExpressionTypeAwait       ExpressionType = "await"
)

// ExpressionNode represents an expression in the code.
type ExpressionNode struct {
	BaseNode
	ExprType ExpressionType
	Operator string
	Left     Node
	Right    Node
}

// IdentifierNode represents an identifier.
type IdentifierNode struct {
	BaseNode
	Name  string
	Scope string // "local", "global", "parameter", etc.
}

// LiteralNode represents a literal value.
type LiteralNode struct {
	BaseNode
	Value       string
	LiteralType string // "string", "number", "boolean", "null", "undefined"
}
