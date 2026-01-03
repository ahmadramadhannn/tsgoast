package analyzer

import (
	"github.com/ahmadramadhannn/tsgoast/ast"
)

// FindExpressions finds all expression nodes in the AST.
func (a *Analyzer) FindExpressions() []ast.Node {
	return a.FindNodesByType(ast.NodeTypeExpression)
}

// FindIdentifiers finds all identifier nodes in the AST.
func (a *Analyzer) FindIdentifiers() []ast.Node {
	return a.FindNodesByType(ast.NodeTypeIdentifier)
}

// FindLiterals finds all literal nodes in the AST.
func (a *Analyzer) FindLiterals() []ast.Node {
	return a.FindNodesByType(ast.NodeTypeLiteral)
}

// GetIdentifierName returns the name of an identifier node.
func GetIdentifierName(node ast.Node) string {
	if node == nil || node.Type() != ast.NodeTypeIdentifier {
		return ""
	}
	return node.Text()
}

// GetLiteralValue returns the value of a literal node.
func GetLiteralValue(node ast.Node) string {
	if node == nil || node.Type() != ast.NodeTypeLiteral {
		return ""
	}
	return node.Text()
}
