package tsgoast

import (
	"strings"

	"github.com/ahmadramadhannn/tsgoast/ast"
)

// Tree represents the complete AST tree with typed statements.
type Tree struct {
	Root       *ast.BaseNode
	Statements []ast.Statement
}

// ParseTree parses TypeScript source code and returns a typed AST tree.
func (p *Parser) ParseTree(source []byte) (*Tree, error) {
	root, err := p.Parse(source)
	if err != nil {
		return nil, err
	}

	tree := &Tree{
		Root:       root,
		Statements: make([]ast.Statement, 0),
	}

	// Extract statements from the root
	tree.Statements = p.extractStatements(root)

	return tree, nil
}

// ParseTreeFromFile parses a TypeScript file and returns a typed AST tree.
func (p *Parser) ParseTreeFromFile(path string) (*Tree, error) {
	root, err := p.ParseFile(path)
	if err != nil {
		return nil, err
	}

	tree := &Tree{
		Root:       root,
		Statements: make([]ast.Statement, 0),
	}

	tree.Statements = p.extractStatements(root)

	return tree, nil
}

// extractStatements extracts typed statements from the AST.
func (p *Parser) extractStatements(node *ast.BaseNode) []ast.Statement {
	if node == nil {
		return nil
	}

	statements := make([]ast.Statement, 0)

	for _, child := range node.Children() {
		if stmt := p.buildStatement(child); stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements
}

// buildStatement builds a typed statement from an AST node.
func (p *Parser) buildStatement(node ast.Node) ast.Statement {
	if node == nil {
		return nil
	}

	baseNode, ok := node.(*ast.BaseNode)
	if !ok {
		return nil
	}

	text := baseNode.Text()

	// Use text-based detection since we're working with converted nodes
	// In a future version, we could store the original tree-sitter kind

	// Check for lexical_declaration (const, let)
	if strings.HasPrefix(strings.TrimSpace(text), "const ") ||
		strings.HasPrefix(strings.TrimSpace(text), "let ") ||
		strings.HasPrefix(strings.TrimSpace(text), "var ") {
		return p.buildVariableStatement(baseNode)
	}

	// Function declaration
	if strings.HasPrefix(strings.TrimSpace(text), "function ") ||
		strings.HasPrefix(strings.TrimSpace(text), "async function") {
		return p.buildFunctionDeclaration(baseNode)
	}

	// Class declaration
	if strings.HasPrefix(strings.TrimSpace(text), "class ") ||
		strings.HasPrefix(strings.TrimSpace(text), "abstract class") {
		return p.buildClassDeclaration(baseNode)
	}

	// If statement
	if strings.HasPrefix(strings.TrimSpace(text), "if ") ||
		strings.HasPrefix(strings.TrimSpace(text), "if(") {
		return p.buildIfStatement(baseNode)
	}

	// While statement
	if strings.HasPrefix(strings.TrimSpace(text), "while ") ||
		strings.HasPrefix(strings.TrimSpace(text), "while(") {
		return p.buildWhileStatement(baseNode)
	}

	// For statement (including for-of and for-in)
	if strings.HasPrefix(strings.TrimSpace(text), "for ") ||
		strings.HasPrefix(strings.TrimSpace(text), "for(") {
		return p.buildForStatement(baseNode)
	}

	// Switch statement
	if strings.HasPrefix(strings.TrimSpace(text), "switch ") ||
		strings.HasPrefix(strings.TrimSpace(text), "switch(") {
		return p.buildSwitchStatement(baseNode)
	}

	// Try statement
	if strings.HasPrefix(strings.TrimSpace(text), "try ") ||
		strings.HasPrefix(strings.TrimSpace(text), "try{") {
		return p.buildTryStatement(baseNode)
	}

	// Return statement
	if strings.HasPrefix(strings.TrimSpace(text), "return") {
		return p.buildReturnStatement(baseNode)
	}

	// Throw statement
	if strings.HasPrefix(strings.TrimSpace(text), "throw ") {
		return p.buildThrowStatement(baseNode)
	}

	// Break statement
	if strings.HasPrefix(strings.TrimSpace(text), "break") {
		return p.buildBreakStatement(baseNode)
	}

	// Continue statement
	if strings.HasPrefix(strings.TrimSpace(text), "continue") {
		return p.buildContinueStatement(baseNode)
	}

	// Import declaration
	if strings.HasPrefix(strings.TrimSpace(text), "import ") {
		return p.buildImportDeclaration(baseNode)
	}

	// Export declaration
	if strings.HasPrefix(strings.TrimSpace(text), "export ") {
		return p.buildExportDeclaration(baseNode)
	}

	// Enum declaration
	if strings.Contains(text, "enum ") {
		return p.buildEnumDeclaration(baseNode)
	}

	// Namespace declaration
	if strings.Contains(text, "namespace ") {
		return p.buildNamespaceDeclaration(baseNode)
	}

	// Expression statement (default for expressions)
	// Only create expression statements for actual expressions, not empty nodes
	if len(strings.TrimSpace(text)) > 0 && !strings.HasPrefix(text, "//") {
		return p.buildExpressionStatement(baseNode)
	}

	return nil
}

// buildVariableStatement builds a variable statement.
func (p *Parser) buildVariableStatement(node *ast.BaseNode) *ast.VariableStatement {
	text := node.Text()
	kind := "var"
	if strings.Contains(text, "const ") {
		kind = "const"
	} else if strings.Contains(text, "let ") {
		kind = "let"
	}

	return &ast.VariableStatement{
		BaseNode:     *node,
		Declarations: make([]*ast.VariableDeclarator, 0),
		Kind:         kind,
	}
}

// buildFunctionDeclaration builds a function declaration.
func (p *Parser) buildFunctionDeclaration(node *ast.BaseNode) *ast.FunctionDeclaration {
	text := node.Text()

	return &ast.FunctionDeclaration{
		BaseNode:    *node,
		Name:        p.extractFunctionName(node),
		Parameters:  make([]*ast.Parameter, 0),
		IsAsync:     strings.Contains(text, "async "),
		IsExported:  strings.HasPrefix(strings.TrimSpace(text), "export "),
		IsGenerator: strings.Contains(text, "function*"),
	}
}

// buildClassDeclaration builds a class declaration.
func (p *Parser) buildClassDeclaration(node *ast.BaseNode) *ast.ClassDeclaration {
	text := node.Text()

	return &ast.ClassDeclaration{
		BaseNode:   *node,
		Name:       p.extractClassName(node),
		IsAbstract: strings.Contains(text, "abstract "),
		IsExported: strings.HasPrefix(strings.TrimSpace(text), "export "),
	}
}

// buildIfStatement builds an if statement.
func (p *Parser) buildIfStatement(node *ast.BaseNode) *ast.IfStatement {
	return &ast.IfStatement{
		BaseNode: *node,
	}
}

// buildWhileStatement builds a while statement.
func (p *Parser) buildWhileStatement(node *ast.BaseNode) *ast.WhileStatement {
	return &ast.WhileStatement{
		BaseNode: *node,
	}
}

// buildForStatement builds a for statement.
func (p *Parser) buildForStatement(node *ast.BaseNode) ast.Statement {
	text := node.Text()

	if strings.Contains(text, " of ") {
		return &ast.ForOfStatement{
			BaseNode: *node,
			IsAwait:  strings.Contains(text, "await "),
		}
	} else if strings.Contains(text, " in ") {
		return &ast.ForInStatement{
			BaseNode: *node,
		}
	}

	return &ast.ForStatement{
		BaseNode: *node,
	}
}

// buildSwitchStatement builds a switch statement.
func (p *Parser) buildSwitchStatement(node *ast.BaseNode) *ast.SwitchStatement {
	return &ast.SwitchStatement{
		BaseNode: *node,
		Cases:    make([]*ast.SwitchCase, 0),
	}
}

// buildTryStatement builds a try statement.
func (p *Parser) buildTryStatement(node *ast.BaseNode) *ast.TryStatement {
	return &ast.TryStatement{
		BaseNode: *node,
	}
}

// buildReturnStatement builds a return statement.
func (p *Parser) buildReturnStatement(node *ast.BaseNode) *ast.ReturnStatement {
	return &ast.ReturnStatement{
		BaseNode: *node,
	}
}

// buildThrowStatement builds a throw statement.
func (p *Parser) buildThrowStatement(node *ast.BaseNode) *ast.ThrowStatement {
	return &ast.ThrowStatement{
		BaseNode: *node,
	}
}

// buildBreakStatement builds a break statement.
func (p *Parser) buildBreakStatement(node *ast.BaseNode) *ast.BreakStatement {
	return &ast.BreakStatement{
		BaseNode: *node,
	}
}

// buildContinueStatement builds a continue statement.
func (p *Parser) buildContinueStatement(node *ast.BaseNode) *ast.ContinueStatement {
	return &ast.ContinueStatement{
		BaseNode: *node,
	}
}

// buildExpressionStatement builds an expression statement.
func (p *Parser) buildExpressionStatement(node *ast.BaseNode) *ast.ExpressionStatement {
	return &ast.ExpressionStatement{
		BaseNode: *node,
	}
}

// buildImportDeclaration builds an import declaration.
func (p *Parser) buildImportDeclaration(node *ast.BaseNode) *ast.ImportDeclaration {
	return &ast.ImportDeclaration{
		BaseNode:   *node,
		Specifiers: make([]ast.Node, 0),
	}
}

// buildExportDeclaration builds an export declaration.
func (p *Parser) buildExportDeclaration(node *ast.BaseNode) *ast.ExportDeclaration {
	text := node.Text()

	return &ast.ExportDeclaration{
		BaseNode:   *node,
		Specifiers: make([]ast.Node, 0),
		IsDefault:  strings.Contains(text, "export default"),
	}
}

// buildEnumDeclaration builds an enum declaration.
func (p *Parser) buildEnumDeclaration(node *ast.BaseNode) *ast.EnumDeclaration {
	text := node.Text()

	return &ast.EnumDeclaration{
		BaseNode:   *node,
		Members:    make([]*ast.EnumMember, 0),
		IsConst:    strings.Contains(text, "const enum"),
		IsExported: strings.HasPrefix(strings.TrimSpace(text), "export "),
	}
}

// buildNamespaceDeclaration builds a namespace declaration.
func (p *Parser) buildNamespaceDeclaration(node *ast.BaseNode) *ast.NamespaceDeclaration {
	text := node.Text()

	return &ast.NamespaceDeclaration{
		BaseNode:   *node,
		Body:       make([]ast.Statement, 0),
		IsExported: strings.HasPrefix(strings.TrimSpace(text), "export "),
	}
}

// Helper functions

func (p *Parser) extractFunctionName(node *ast.BaseNode) string {
	// First try to find identifier in children
	for _, child := range node.Children() {
		if child.Type() == ast.NodeTypeIdentifier {
			return child.Text()
		}
	}

	// Fallback: extract from text
	text := strings.TrimSpace(node.Text())
	if strings.HasPrefix(text, "async ") {
		text = strings.TrimPrefix(text, "async ")
		text = strings.TrimSpace(text)
	}
	if strings.HasPrefix(text, "function ") {
		text = strings.TrimPrefix(text, "function ")
	} else if strings.HasPrefix(text, "function*") {
		text = strings.TrimPrefix(text, "function*")
	}
	text = strings.TrimSpace(text)

	// Extract name before (
	if idx := strings.Index(text, "("); idx > 0 {
		return strings.TrimSpace(text[:idx])
	}

	return ""
}

func (p *Parser) extractClassName(node *ast.BaseNode) string {
	// First try to find identifier in children
	for _, child := range node.Children() {
		if child.Type() == ast.NodeTypeIdentifier {
			return child.Text()
		}
	}

	// Fallback: extract from text
	text := strings.TrimSpace(node.Text())
	if strings.HasPrefix(text, "abstract ") {
		text = strings.TrimPrefix(text, "abstract ")
		text = strings.TrimSpace(text)
	}
	if strings.HasPrefix(text, "class ") {
		text = strings.TrimPrefix(text, "class ")
	}
	text = strings.TrimSpace(text)

	// Extract name before { or extends or implements
	for _, delim := range []string{"{", " extends", " implements", "<"} {
		if idx := strings.Index(text, delim); idx > 0 {
			return strings.TrimSpace(text[:idx])
		}
	}

	return strings.TrimSpace(text)
}
