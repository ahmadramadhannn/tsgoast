package analyzer

import (
	"strings"

	"github.com/ahmadramadhannn/tsgoast/ast"
)

// FindFunctions finds all function declarations in the AST.
func (a *Analyzer) FindFunctions() []ast.Node {
	return a.FindNodes(func(node ast.Node) bool {
		t := node.Type()
		return t == ast.NodeTypeFunction || t == ast.NodeTypeArrowFunction
	})
}

// FindMethods finds all method definitions in the AST.
func (a *Analyzer) FindMethods() []ast.Node {
	return a.FindNodesByType(ast.NodeTypeMethod)
}

// IsAsync checks if a function node represents an async function.
// This is a simplified check based on the node's text content.
func IsAsync(node ast.Node) bool {
	if node == nil {
		return false
	}

	t := node.Type()
	if t != ast.NodeTypeFunction && t != ast.NodeTypeArrowFunction && t != ast.NodeTypeMethod {
		return false
	}

	text := node.Text()
	return strings.Contains(text, "async ")
}

// IsExported checks if a function node is exported.
// This is a simplified check based on the node's text content.
func IsExported(node ast.Node) bool {
	if node == nil {
		return false
	}

	text := node.Text()
	return strings.HasPrefix(strings.TrimSpace(text), "export ")
}

// IsGenerator checks if a function node is a generator function.
func IsGenerator(node ast.Node) bool {
	if node == nil {
		return false
	}

	t := node.Type()
	if t != ast.NodeTypeFunction {
		return false
	}

	text := node.Text()
	return strings.Contains(text, "function*")
}

// GetFunctionName attempts to extract the function name from a function node.
// Returns empty string if the name cannot be determined.
func GetFunctionName(node ast.Node) string {
	if node == nil {
		return ""
	}

	// Look for identifier child nodes
	for _, child := range node.Children() {
		if child.Type() == ast.NodeTypeIdentifier {
			return child.Text()
		}
	}

	return ""
}

// HasParameters checks if a function has parameters.
func HasParameters(node ast.Node) bool {
	if node == nil {
		return false
	}

	for _, child := range node.Children() {
		if child.Type() == ast.NodeTypeParameter {
			return true
		}
	}

	return false
}

// CountParameters counts the number of parameters in a function.
func CountParameters(node ast.Node) int {
	if node == nil {
		return 0
	}

	count := 0
	var countInNode func(ast.Node)
	countInNode = func(n ast.Node) {
		if n.Type() == ast.NodeTypeParameter {
			// Check if this is a formal_parameters container or actual parameter
			text := strings.TrimSpace(n.Text())
			if text != "" && !strings.HasPrefix(text, "(") {
				count++
			}
		}
		for _, child := range n.Children() {
			countInNode(child)
		}
	}

	countInNode(node)
	return count
}
