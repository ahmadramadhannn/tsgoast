package analyzer

import (
	"strings"

	"github.com/ahmadramadhannn/tsgoast/ast"
)

// FindInterfaces finds all interface declarations in the AST.
func (a *Analyzer) FindInterfaces() []ast.Node {
	return a.FindNodesByType(ast.NodeTypeInterface)
}

// FindTypeAliases finds all type alias declarations in the AST.
func (a *Analyzer) FindTypeAliases() []ast.Node {
	return a.FindNodesByType(ast.NodeTypeTypeAlias)
}

// GetInterfaceName attempts to extract the interface name.
func GetInterfaceName(node ast.Node) string {
	if node == nil || node.Type() != ast.NodeTypeInterface {
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

// GetTypeAliasName attempts to extract the type alias name.
func GetTypeAliasName(node ast.Node) string {
	if node == nil || node.Type() != ast.NodeTypeTypeAlias {
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

// HasExtends checks if an interface extends another interface.
func HasExtends(node ast.Node) bool {
	if node == nil || node.Type() != ast.NodeTypeInterface {
		return false
	}

	text := node.Text()
	return strings.Contains(text, " extends ")
}

// IsReadonly checks if a property is marked as readonly.
func IsReadonly(node ast.Node) bool {
	if node == nil {
		return false
	}

	text := node.Text()
	return strings.Contains(text, "readonly ")
}

// IsOptionalProperty checks if a property is optional.
func IsOptionalProperty(node ast.Node) bool {
	if node == nil {
		return false
	}

	text := node.Text()
	return strings.Contains(text, "?:")
}

// CountProperties counts the number of properties in an interface or type.
func CountProperties(node ast.Node) int {
	if node == nil {
		return 0
	}

	count := 0
	for _, child := range node.Children() {
		if child.Type() == ast.NodeTypeProperty {
			count++
		}
	}

	return count
}

// IsGenericType checks if a type has type parameters.
func IsGenericType(node ast.Node) bool {
	if node == nil {
		return false
	}

	t := node.Type()
	if t != ast.NodeTypeInterface && t != ast.NodeTypeTypeAlias {
		return false
	}

	text := node.Text()
	return strings.Contains(text, "<") && strings.Contains(text, ">")
}
