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
// For arrow functions, it looks at the parent variable declaration or property assignment.
// Returns empty string if the name cannot be determined.
func GetFunctionName(node ast.Node) string {
	if node == nil {
		return ""
	}

	nodeType := node.Type()

	// For regular functions and methods, look for identifier child
	if nodeType == ast.NodeTypeFunction || nodeType == ast.NodeTypeMethod {
		for _, child := range node.Children() {
			if child.Type() == ast.NodeTypeIdentifier {
				return child.Text()
			}
		}
		return ""
	}

	// For arrow functions, we need to look at the parent context
	if nodeType == ast.NodeTypeArrowFunction {
		// Arrow functions are typically assigned to variables or properties
		// We need to traverse up to find the variable name
		parent := node.Parent()
		if parent == nil {
			return "" // Anonymous arrow function
		}

		// Check if parent is a variable declarator
		// Pattern: const name = () => {}
		parentText := parent.Text()

		// Try to extract variable name from parent
		if strings.Contains(parentText, "=") {
			// Look for identifier in parent's children
			for _, child := range parent.Children() {
				if child.Type() == ast.NodeTypeIdentifier {
					// Make sure this identifier comes before the arrow function
					childText := child.Text()
					arrowText := node.Text()
					if strings.Index(parentText, childText) < strings.Index(parentText, arrowText) {
						return childText
					}
				}
			}

			// Fallback: parse from text
			parts := strings.Split(parentText, "=")
			if len(parts) >= 2 {
				leftSide := strings.TrimSpace(parts[0])
				// Remove const/let/var
				leftSide = strings.TrimPrefix(leftSide, "const ")
				leftSide = strings.TrimPrefix(leftSide, "let ")
				leftSide = strings.TrimPrefix(leftSide, "var ")
				leftSide = strings.TrimSpace(leftSide)

				// Extract just the identifier (before any : type annotation)
				if idx := strings.Index(leftSide, ":"); idx > 0 {
					leftSide = leftSide[:idx]
				}
				leftSide = strings.TrimSpace(leftSide)

				// Return if it looks like a valid identifier
				if leftSide != "" && !strings.Contains(leftSide, " ") {
					return leftSide
				}
			}
		}

		// Check if it's a property in an object
		// Pattern: { methodName: () => {} }
		if strings.Contains(parentText, ":") {
			parts := strings.Split(parentText, ":")
			if len(parts) >= 2 {
				potentialName := strings.TrimSpace(parts[0])
				// Check if this looks like a property name (not a type annotation)
				if !strings.Contains(potentialName, " ") &&
					!strings.Contains(potentialName, "(") &&
					potentialName != "" {
					return potentialName
				}
			}
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
