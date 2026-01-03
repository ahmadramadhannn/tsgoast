// Package analyzer provides high-level analysis functions for TypeScript AST.
package analyzer

import (
	"github.com/ahmadro/tsgoast/ast"
)

// Analyzer provides high-level AST analysis capabilities.
type Analyzer struct {
	root *ast.BaseNode
}

// New creates a new analyzer for the given AST root node.
func New(root *ast.BaseNode) *Analyzer {
	return &Analyzer{
		root: root,
	}
}

// Root returns the root node of the AST.
func (a *Analyzer) Root() *ast.BaseNode {
	return a.root
}

// Visit traverses the AST and calls the visitor function for each node.
// If the visitor returns false, traversal of that subtree is stopped.
func (a *Analyzer) Visit(visitor func(node ast.Node) bool) {
	if a.root == nil {
		return
	}
	a.visitNode(a.root, visitor)
}

func (a *Analyzer) visitNode(node ast.Node, visitor func(ast.Node) bool) {
	if node == nil {
		return
	}

	// Call visitor, if it returns false, stop traversing this subtree
	if !visitor(node) {
		return
	}

	// Visit children
	for _, child := range node.Children() {
		a.visitNode(child, visitor)
	}
}

// FindNodes finds all nodes matching the given predicate.
func (a *Analyzer) FindNodes(predicate func(node ast.Node) bool) []ast.Node {
	var results []ast.Node
	a.Visit(func(node ast.Node) bool {
		if predicate(node) {
			results = append(results, node)
		}
		return true
	})
	return results
}

// FindNodesByType finds all nodes of the given type.
func (a *Analyzer) FindNodesByType(nodeType ast.NodeType) []ast.Node {
	return a.FindNodes(func(node ast.Node) bool {
		return node.Type() == nodeType
	})
}

// CountNodes counts all nodes matching the given predicate.
func (a *Analyzer) CountNodes(predicate func(node ast.Node) bool) int {
	count := 0
	a.Visit(func(node ast.Node) bool {
		if predicate(node) {
			count++
		}
		return true
	})
	return count
}

// CountNodesByType counts all nodes of the given type.
func (a *Analyzer) CountNodesByType(nodeType ast.NodeType) int {
	return a.CountNodes(func(node ast.Node) bool {
		return node.Type() == nodeType
	})
}
