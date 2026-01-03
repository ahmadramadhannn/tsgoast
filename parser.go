// Package tsgoast provides a TypeScript AST parser and analyzer.
package tsgoast

import (
	"fmt"
	"os"

	"github.com/ahmadramadhannn/tsgoast/ast"
	sitter "github.com/tree-sitter/go-tree-sitter"
	typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

// Parser wraps the tree-sitter parser for TypeScript.
type Parser struct {
	parser   *sitter.Parser
	language *sitter.Language
}

// New creates a new TypeScript parser.
func New() (*Parser, error) {
	parser := sitter.NewParser()
	lang := sitter.NewLanguage(typescript.LanguageTypescript())

	if err := parser.SetLanguage(lang); err != nil {
		return nil, fmt.Errorf("failed to set language: %w", err)
	}

	return &Parser{
		parser:   parser,
		language: lang,
	}, nil
}

// Parse parses TypeScript source code and returns the root AST node.
func (p *Parser) Parse(source []byte) (*ast.BaseNode, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("source code is empty")
	}

	tree := p.parser.Parse(source, nil)
	if tree == nil {
		return nil, fmt.Errorf("failed to parse source code")
	}
	defer tree.Close()

	root := tree.RootNode()
	if root == nil {
		return nil, fmt.Errorf("failed to get root node")
	}

	return p.convertNode(root, source, nil), nil
}

// ParseFile parses a TypeScript file and returns the root AST node.
func (p *Parser) ParseFile(path string) (*ast.BaseNode, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.Parse(source)
}

// convertNode converts a tree-sitter node to our AST node.
func (p *Parser) convertNode(node *sitter.Node, source []byte, parent *ast.BaseNode) *ast.BaseNode {
	if node == nil {
		return nil
	}

	baseNode := &ast.BaseNode{
		NodeType: p.mapNodeType(node.Kind()),
		Content:  string(source[node.StartByte():node.EndByte()]),
		SourceRange: ast.Range{
			Start: ast.Position{
				Line:   uint32(node.StartPosition().Row),
				Column: uint32(node.StartPosition().Column),
				Offset: uint32(node.StartByte()),
			},
			End: ast.Position{
				Line:   uint32(node.EndPosition().Row),
				Column: uint32(node.EndPosition().Column),
				Offset: uint32(node.EndByte()),
			},
		},
		ParentNode: parent,
	}

	// Convert children
	childCount := node.ChildCount()
	if childCount > 0 {
		baseNode.ChildNodes = make([]ast.Node, 0, childCount)
		for i := uint(0); i < childCount; i++ {
			child := node.Child(i)
			if child != nil {
				childNode := p.convertNode(child, source, baseNode)
				if childNode != nil {
					baseNode.ChildNodes = append(baseNode.ChildNodes, childNode)
				}
			}
		}
	}

	return baseNode
}

// mapNodeType maps tree-sitter node types to our AST node types.
func (p *Parser) mapNodeType(tsType string) ast.NodeType {
	switch tsType {
	case "function_declaration":
		return ast.NodeTypeFunction
	case "arrow_function":
		return ast.NodeTypeArrowFunction
	case "method_definition":
		return ast.NodeTypeMethod
	case "interface_declaration":
		return ast.NodeTypeInterface
	case "type_alias_declaration":
		return ast.NodeTypeTypeAlias
	case "identifier":
		return ast.NodeTypeIdentifier
	case "property_signature":
		return ast.NodeTypeProperty
	case "formal_parameters", "required_parameter", "optional_parameter":
		return ast.NodeTypeParameter
	case "string", "number", "true", "false", "null", "undefined":
		return ast.NodeTypeLiteral
	default:
		// Check if it's an expression type
		if isExpressionType(tsType) {
			return ast.NodeTypeExpression
		}
		return ast.NodeTypeUnknown
	}
}

// isExpressionType checks if a tree-sitter type is an expression.
func isExpressionType(tsType string) bool {
	expressionTypes := []string{
		"binary_expression",
		"unary_expression",
		"call_expression",
		"member_expression",
		"assignment_expression",
		"ternary_expression",
		"new_expression",
		"await_expression",
	}

	for _, et := range expressionTypes {
		if tsType == et {
			return true
		}
	}
	return false
}

// Close releases resources held by the parser.
func (p *Parser) Close() {
	if p.parser != nil {
		p.parser.Close()
		p.parser = nil
	}
}
