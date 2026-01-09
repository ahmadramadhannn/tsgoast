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
// It checks the node itself and its ancestors for "export" keywords.
func IsExported(node ast.Node) bool {
	if node == nil {
		return false
	}

	// Helper to check if a node text indicates export
	isExportNode := func(n ast.Node) bool {
		text := strings.TrimSpace(n.Text())
		return strings.HasPrefix(text, "export ")
	}

	// Check the node itself
	if isExportNode(node) {
		return true
	}

	// Traverse up the parent chain (up to 3 levels is usually enough)
	// Level 1: Parent (e.g. export_statement for function declaration)
	// Level 2: Grandparent (e.g. lexical_declaration for arrow function)
	// Level 3: Great-grandparent (e.g. export_statement for arrow function)
	current := node.Parent()
	for i := 0; i < 3 && current != nil; i++ {
		if isExportNode(current) {
			return true
		}
		current = current.Parent()
	}

	return false
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

// GetFunctionName extracts the function name from a function node.
// For regular functions and methods, it looks for an identifier child.
// For arrow functions, it traverses parent nodes to find the variable name.
// Returns an empty string if the name cannot be determined.
func GetFunctionName(node ast.Node) string {
	if node == nil {
		return ""
	}

	nodeType := node.Type()

	// For regular functions and methods, look for identifier child
	if nodeType == ast.NodeTypeFunction || nodeType == ast.NodeTypeMethod {
		return findIdentifierInChildren(node)
	}

	// For arrow functions, traverse parent chain to find the variable name
	if nodeType == ast.NodeTypeArrowFunction {
		return extractArrowFunctionName(node)
	}

	return ""
}

// findIdentifierInChildren searches for an identifier among direct children.
func findIdentifierInChildren(node ast.Node) string {
	for _, child := range node.Children() {
		if child.Type() == ast.NodeTypeIdentifier {
			return child.Text()
		}
	}
	return ""
}

// extractArrowFunctionName traverses parent nodes to find the variable name
// for an arrow function assignment pattern like: const name = () => {}
func extractArrowFunctionName(node ast.Node) string {
	// Traverse up the parent chain looking for a variable declaration
	current := node.Parent()
	for current != nil {
		// Look for identifier children that represent the variable name
		for _, child := range current.Children() {
			if child.Type() == ast.NodeTypeIdentifier {
				// Verify this is a variable name, not some other identifier
				// by checking it appears before the arrow function in the parent text
				if isVariableNameForArrowFunction(current, child, node) {
					return child.Text()
				}
			}
		}

		// Move up to the next parent
		current = current.Parent()

		// Safety limit to prevent infinite loops
		if current != nil && current.Parent() == current {
			break
		}
	}

	return ""
}

// isVariableNameForArrowFunction checks if the identifier is the variable name
// that the arrow function is assigned to.
func isVariableNameForArrowFunction(parent, identifier, arrowFunc ast.Node) bool {
	parentText := parent.Text()
	identifierText := identifier.Text()
	arrowText := arrowFunc.Text()

	// The identifier should come before the arrow function in the text
	identifierPos := strings.Index(parentText, identifierText)
	arrowPos := strings.Index(parentText, arrowText)

	// Identifier must be found before the arrow function
	if identifierPos == -1 || arrowPos == -1 {
		return false
	}

	return identifierPos < arrowPos
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
