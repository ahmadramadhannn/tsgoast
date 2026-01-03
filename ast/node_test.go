package ast

import (
	"testing"
)

func TestNodeType(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		expected NodeType
	}{
		{"Function", NodeTypeFunction, "function"},
		{"ArrowFunction", NodeTypeArrowFunction, "arrow_function"},
		{"Interface", NodeTypeInterface, "interface"},
		{"TypeAlias", NodeTypeTypeAlias, "type_alias"},
		{"Identifier", NodeTypeIdentifier, "identifier"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nodeType != tt.expected {
				t.Errorf("NodeType = %v, want %v", tt.nodeType, tt.expected)
			}
		})
	}
}

func TestBaseNode(t *testing.T) {
	t.Run("Type", func(t *testing.T) {
		node := &BaseNode{
			NodeType: NodeTypeFunction,
		}
		if got := node.Type(); got != NodeTypeFunction {
			t.Errorf("Type() = %v, want %v", got, NodeTypeFunction)
		}
	})

	t.Run("Text", func(t *testing.T) {
		content := "function test() {}"
		node := &BaseNode{
			Content: content,
		}
		if got := node.Text(); got != content {
			t.Errorf("Text() = %v, want %v", got, content)
		}
	})

	t.Run("Children", func(t *testing.T) {
		child1 := &BaseNode{NodeType: NodeTypeIdentifier}
		child2 := &BaseNode{NodeType: NodeTypeParameter}
		node := &BaseNode{
			ChildNodes: []Node{child1, child2},
		}
		children := node.Children()
		if len(children) != 2 {
			t.Errorf("Children() length = %v, want 2", len(children))
		}
	})

	t.Run("Range", func(t *testing.T) {
		r := Range{
			Start: Position{Line: 1, Column: 0, Offset: 0},
			End:   Position{Line: 1, Column: 10, Offset: 10},
		}
		node := &BaseNode{
			SourceRange: r,
		}
		if got := node.Range(); got != r {
			t.Errorf("Range() = %v, want %v", got, r)
		}
	})

	t.Run("Parent", func(t *testing.T) {
		parent := &BaseNode{NodeType: NodeTypeFunction}
		child := &BaseNode{
			NodeType:   NodeTypeIdentifier,
			ParentNode: parent,
		}
		if got := child.Parent(); got != parent {
			t.Errorf("Parent() = %v, want %v", got, parent)
		}
	})
}

func TestPosition(t *testing.T) {
	pos := Position{
		Line:   10,
		Column: 5,
		Offset: 105,
	}

	if pos.Line != 10 {
		t.Errorf("Line = %v, want 10", pos.Line)
	}
	if pos.Column != 5 {
		t.Errorf("Column = %v, want 5", pos.Column)
	}
	if pos.Offset != 105 {
		t.Errorf("Offset = %v, want 105", pos.Offset)
	}
}

func TestRange(t *testing.T) {
	r := Range{
		Start: Position{Line: 1, Column: 0, Offset: 0},
		End:   Position{Line: 5, Column: 10, Offset: 50},
	}

	if r.Start.Line != 1 {
		t.Errorf("Start.Line = %v, want 1", r.Start.Line)
	}
	if r.End.Line != 5 {
		t.Errorf("End.Line = %v, want 5", r.End.Line)
	}
}

func TestNodeHierarchy(t *testing.T) {
	// Create a simple hierarchy: parent -> child1, child2
	parent := &BaseNode{
		NodeType: NodeTypeFunction,
		Content:  "function test() {}",
	}

	child1 := &BaseNode{
		NodeType:   NodeTypeIdentifier,
		Content:    "test",
		ParentNode: parent,
	}

	child2 := &BaseNode{
		NodeType:   NodeTypeParameter,
		Content:    "()",
		ParentNode: parent,
	}

	parent.ChildNodes = []Node{child1, child2}

	// Test parent has correct children
	children := parent.Children()
	if len(children) != 2 {
		t.Fatalf("Parent should have 2 children, got %d", len(children))
	}

	// Test children have correct parent
	if child1.Parent() != parent {
		t.Error("Child1 parent is incorrect")
	}
	if child2.Parent() != parent {
		t.Error("Child2 parent is incorrect")
	}
}
