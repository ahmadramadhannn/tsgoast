// Package ast provides types and interfaces for representing TypeScript AST nodes.
package ast

// NodeType represents the type of an AST node.
type NodeType string

// Node type constants.
const (
	NodeTypeFunction      NodeType = "function"
	NodeTypeArrowFunction NodeType = "arrow_function"
	NodeTypeMethod        NodeType = "method"
	NodeTypeInterface     NodeType = "interface"
	NodeTypeTypeAlias     NodeType = "type_alias"
	NodeTypeExpression    NodeType = "expression"
	NodeTypeIdentifier    NodeType = "identifier"
	NodeTypeLiteral       NodeType = "literal"
	NodeTypeProperty      NodeType = "property"
	NodeTypeParameter     NodeType = "parameter"
	NodeTypeUnknown       NodeType = "unknown"
)

// Position represents a position in the source code.
type Position struct {
	Line   uint32
	Column uint32
	Offset uint32
}

// Range represents a range in the source code.
type Range struct {
	Start Position
	End   Position
}

// Node is the interface that all AST nodes implement.
type Node interface {
	// Type returns the type of the node.
	Type() NodeType

	// Text returns the text content of the node.
	Text() string

	// Children returns the child nodes.
	Children() []Node

	// Range returns the source range of the node.
	Range() Range

	// Parent returns the parent node, or nil if this is the root.
	Parent() Node
}

// BaseNode provides common functionality for all AST nodes.
type BaseNode struct {
	NodeType    NodeType
	Content     string
	ChildNodes  []Node
	SourceRange Range
	ParentNode  Node
}

// Type returns the type of the node.
func (n *BaseNode) Type() NodeType {
	return n.NodeType
}

// Text returns the text content of the node.
func (n *BaseNode) Text() string {
	return n.Content
}

// Children returns the child nodes.
func (n *BaseNode) Children() []Node {
	return n.ChildNodes
}

// Range returns the source range of the node.
func (n *BaseNode) Range() Range {
	return n.SourceRange
}

// Parent returns the parent node.
func (n *BaseNode) Parent() Node {
	return n.ParentNode
}
