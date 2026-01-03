package tsgoast

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ahmadro/tsgoast/ast"
)

func TestNew(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}
	if parser == nil {
		t.Fatal("New() returned nil parser")
	}
	defer parser.Close()

	if parser.parser == nil {
		t.Error("Parser.parser is nil")
	}
	if parser.language == nil {
		t.Error("Parser.language is nil")
	}
}

func TestParse(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name:    "Simple function",
			source:  "function test() {}",
			wantErr: false,
		},
		{
			name:    "Arrow function",
			source:  "const add = (a: number, b: number) => a + b;",
			wantErr: false,
		},
		{
			name:    "Interface",
			source:  "interface User { name: string; }",
			wantErr: false,
		},
		{
			name:    "Type alias",
			source:  "type Point = { x: number; y: number; };",
			wantErr: false,
		},
		{
			name:    "Empty source",
			source:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parser.Parse([]byte(tt.source))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && node == nil {
				t.Error("Parse() returned nil node")
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "Simple TypeScript file",
			filename: "simple.ts",
			wantErr:  false,
		},
		{
			name:     "Functions file",
			filename: "functions.ts",
			wantErr:  false,
		},
		{
			name:     "Types file",
			filename: "types.ts",
			wantErr:  false,
		},
		{
			name:     "Edge cases file",
			filename: "edge_cases.ts",
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			filename: "nonexistent.ts",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("testdata", tt.filename)
			node, err := parser.ParseFile(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if node == nil {
					t.Error("ParseFile() returned nil node")
				}
				if node.Text() == "" {
					t.Error("ParseFile() returned node with empty text")
				}
			}
		})
	}
}

func TestMapNodeType(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	tests := []struct {
		tsType   string
		expected ast.NodeType
	}{
		{"function_declaration", ast.NodeTypeFunction},
		{"arrow_function", ast.NodeTypeArrowFunction},
		{"method_definition", ast.NodeTypeMethod},
		{"interface_declaration", ast.NodeTypeInterface},
		{"type_alias_declaration", ast.NodeTypeTypeAlias},
		{"identifier", ast.NodeTypeIdentifier},
		{"property_signature", ast.NodeTypeProperty},
		{"required_parameter", ast.NodeTypeParameter},
		{"string", ast.NodeTypeLiteral},
		{"binary_expression", ast.NodeTypeExpression},
		{"unknown_type", ast.NodeTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.tsType, func(t *testing.T) {
			got := parser.mapNodeType(tt.tsType)
			if got != tt.expected {
				t.Errorf("mapNodeType(%s) = %v, want %v", tt.tsType, got, tt.expected)
			}
		})
	}
}

func TestIsExpressionType(t *testing.T) {
	tests := []struct {
		tsType   string
		expected bool
	}{
		{"binary_expression", true},
		{"unary_expression", true},
		{"call_expression", true},
		{"member_expression", true},
		{"assignment_expression", true},
		{"ternary_expression", true},
		{"new_expression", true},
		{"await_expression", true},
		{"function_declaration", false},
		{"identifier", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.tsType, func(t *testing.T) {
			got := isExpressionType(tt.tsType)
			if got != tt.expected {
				t.Errorf("isExpressionType(%s) = %v, want %v", tt.tsType, got, tt.expected)
			}
		})
	}
}

func TestParseRealFiles(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	// Test parsing the simple.ts file and verify structure
	t.Run("Parse simple.ts and verify structure", func(t *testing.T) {
		node, err := parser.ParseFile("testdata/simple.ts")
		if err != nil {
			t.Fatalf("Failed to parse simple.ts: %v", err)
		}

		if node == nil {
			t.Fatal("Parsed node is nil")
		}

		// Verify we have children (the file should have declarations)
		if len(node.Children()) == 0 {
			t.Error("Expected parsed file to have child nodes")
		}

		// Verify the content is not empty
		if node.Text() == "" {
			t.Error("Expected parsed node to have text content")
		}
	})
}

func TestParserClose(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	// Should not panic
	parser.Close()
	parser.Close() // Calling twice should be safe
}

func TestParseWithSyntaxError(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	// Tree-sitter is error-tolerant, so it won't fail on syntax errors
	// but we can still parse and get a tree with error nodes
	source := []byte("function {{{{{ invalid syntax")
	node, err := parser.Parse(source)

	// Should not return an error (tree-sitter is error-tolerant)
	if err != nil {
		t.Errorf("Parse() with syntax error returned error: %v", err)
	}

	if node == nil {
		t.Error("Parse() with syntax error returned nil node")
	}
}

// Benchmark tests
func BenchmarkParse(b *testing.B) {
	parser, err := New()
	if err != nil {
		b.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	source := []byte(`
		function test(x: number, y: number): number {
			return x + y;
		}
		
		interface User {
			id: number;
			name: string;
		}
	`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(source)
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}

func BenchmarkParseFile(b *testing.B) {
	parser, err := New()
	if err != nil {
		b.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseFile("testdata/simple.ts")
		if err != nil {
			b.Fatalf("ParseFile error: %v", err)
		}
	}
}

// Test helper to create a temporary file
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test-*.ts")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func TestParseFileWithTempFile(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	content := "const x: number = 42;"
	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	node, err := parser.ParseFile(tmpfile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if node == nil {
		t.Fatal("ParseFile() returned nil node")
	}
}
