package analyzer

import (
	"path/filepath"
	"testing"

	"github.com/ahmadro/tsgoast"
	"github.com/ahmadro/tsgoast/ast"
)

func TestNew(t *testing.T) {
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}

	analyzer := New(root)
	if analyzer == nil {
		t.Fatal("New() returned nil")
	}

	if analyzer.Root() != root {
		t.Error("Root() did not return the expected node")
	}
}

func TestVisit(t *testing.T) {
	// Create a simple tree
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}
	child1 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child2 := &ast.BaseNode{
		NodeType: ast.NodeTypeParameter,
	}
	root.ChildNodes = []ast.Node{child1, child2}

	analyzer := New(root)

	// Count visited nodes
	count := 0
	analyzer.Visit(func(node ast.Node) bool {
		count++
		return true
	})

	// Should visit root + 2 children = 3 nodes
	if count != 3 {
		t.Errorf("Visit() visited %d nodes, want 3", count)
	}
}

func TestVisitWithEarlyStop(t *testing.T) {
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}
	child1 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child2 := &ast.BaseNode{
		NodeType: ast.NodeTypeParameter,
	}
	root.ChildNodes = []ast.Node{child1, child2}

	analyzer := New(root)

	// Stop after first node
	count := 0
	analyzer.Visit(func(node ast.Node) bool {
		count++
		return false // Stop traversal
	})

	// Should only visit root
	if count != 1 {
		t.Errorf("Visit() with early stop visited %d nodes, want 1", count)
	}
}

func TestFindNodes(t *testing.T) {
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}
	child1 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child2 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child3 := &ast.BaseNode{
		NodeType: ast.NodeTypeParameter,
	}
	root.ChildNodes = []ast.Node{child1, child2, child3}

	analyzer := New(root)

	// Find all identifiers
	identifiers := analyzer.FindNodes(func(node ast.Node) bool {
		return node.Type() == ast.NodeTypeIdentifier
	})

	if len(identifiers) != 2 {
		t.Errorf("FindNodes() found %d identifiers, want 2", len(identifiers))
	}
}

func TestFindNodesByType(t *testing.T) {
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}
	child1 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child2 := &ast.BaseNode{
		NodeType: ast.NodeTypeParameter,
	}
	root.ChildNodes = []ast.Node{child1, child2}

	analyzer := New(root)

	tests := []struct {
		name     string
		nodeType ast.NodeType
		want     int
	}{
		{"Find identifiers", ast.NodeTypeIdentifier, 1},
		{"Find parameters", ast.NodeTypeParameter, 1},
		{"Find functions", ast.NodeTypeFunction, 1},
		{"Find interfaces", ast.NodeTypeInterface, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := analyzer.FindNodesByType(tt.nodeType)
			if len(nodes) != tt.want {
				t.Errorf("FindNodesByType(%v) found %d nodes, want %d", tt.nodeType, len(nodes), tt.want)
			}
		})
	}
}

func TestCountNodes(t *testing.T) {
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}
	child1 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child2 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	root.ChildNodes = []ast.Node{child1, child2}

	analyzer := New(root)

	count := analyzer.CountNodes(func(node ast.Node) bool {
		return node.Type() == ast.NodeTypeIdentifier
	})

	if count != 2 {
		t.Errorf("CountNodes() = %d, want 2", count)
	}
}

func TestCountNodesByType(t *testing.T) {
	root := &ast.BaseNode{
		NodeType: ast.NodeTypeFunction,
	}
	child1 := &ast.BaseNode{
		NodeType: ast.NodeTypeIdentifier,
	}
	child2 := &ast.BaseNode{
		NodeType: ast.NodeTypeParameter,
	}
	root.ChildNodes = []ast.Node{child1, child2}

	analyzer := New(root)

	tests := []struct {
		name     string
		nodeType ast.NodeType
		want     int
	}{
		{"Count identifiers", ast.NodeTypeIdentifier, 1},
		{"Count parameters", ast.NodeTypeParameter, 1},
		{"Count all nodes", ast.NodeTypeFunction, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := analyzer.CountNodesByType(tt.nodeType)
			if count != tt.want {
				t.Errorf("CountNodesByType(%v) = %d, want %d", tt.nodeType, count, tt.want)
			}
		})
	}
}

// Integration tests with real TypeScript files
func TestAnalyzerWithRealFiles(t *testing.T) {
	parser, err := tsgoast.New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	tests := []struct {
		name     string
		filename string
		checks   func(*testing.T, *Analyzer)
	}{
		{
			name:     "Simple file",
			filename: "simple.ts",
			checks: func(t *testing.T, a *Analyzer) {
				functions := a.FindFunctions()
				if len(functions) < 1 {
					t.Error("Expected to find at least 1 function")
				}

				interfaces := a.FindInterfaces()
				if len(interfaces) < 1 {
					t.Error("Expected to find at least 1 interface")
				}
			},
		},
		{
			name:     "Functions file",
			filename: "functions.ts",
			checks: func(t *testing.T, a *Analyzer) {
				functions := a.FindFunctions()
				if len(functions) < 5 {
					t.Errorf("Expected to find at least 5 functions, got %d", len(functions))
				}
			},
		},
		{
			name:     "Types file",
			filename: "types.ts",
			checks: func(t *testing.T, a *Analyzer) {
				interfaces := a.FindInterfaces()
				if len(interfaces) < 3 {
					t.Errorf("Expected to find at least 3 interfaces, got %d", len(interfaces))
				}

				typeAliases := a.FindTypeAliases()
				if len(typeAliases) < 3 {
					t.Errorf("Expected to find at least 3 type aliases, got %d", len(typeAliases))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "testdata", tt.filename)
			root, err := parser.ParseFile(path)
			if err != nil {
				t.Fatalf("Failed to parse file: %v", err)
			}

			analyzer := New(root)
			tt.checks(t, analyzer)
		})
	}
}

func TestVisitNilRoot(t *testing.T) {
	analyzer := New(nil)

	// Should not panic
	count := 0
	analyzer.Visit(func(node ast.Node) bool {
		count++
		return true
	})

	if count != 0 {
		t.Errorf("Visit() on nil root visited %d nodes, want 0", count)
	}
}

func BenchmarkVisit(b *testing.B) {
	// Create a tree with some depth
	root := &ast.BaseNode{NodeType: ast.NodeTypeFunction}
	for i := 0; i < 10; i++ {
		child := &ast.BaseNode{NodeType: ast.NodeTypeIdentifier}
		root.ChildNodes = append(root.ChildNodes, child)
	}

	analyzer := New(root)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Visit(func(node ast.Node) bool {
			return true
		})
	}
}

func BenchmarkFindNodesByType(b *testing.B) {
	root := &ast.BaseNode{NodeType: ast.NodeTypeFunction}
	for i := 0; i < 100; i++ {
		child := &ast.BaseNode{NodeType: ast.NodeTypeIdentifier}
		root.ChildNodes = append(root.ChildNodes, child)
	}

	analyzer := New(root)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analyzer.FindNodesByType(ast.NodeTypeIdentifier)
	}
}
