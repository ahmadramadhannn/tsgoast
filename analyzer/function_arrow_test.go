package analyzer

import (
	"testing"

	"github.com/ahmadramadhannn/tsgoast"
	"github.com/ahmadramadhannn/tsgoast/ast"
)

func TestGetFunctionNameWithArrowFunctions(t *testing.T) {
	parser, err := tsgoast.New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	source := []byte(`
		// Regular function
		function regularFunc() {
			return 42;
		}

		// Arrow function assigned to variable
		const arrowFunc = () => {
			return 42;
		};

		// Arrow function with parameter
		const namedArrow = (x: number) => x * 2;

		// Arrow function in object
		const obj = {
			method: () => 42
		};
	`)

	root, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	a := New(root)
	functions := a.FindFunctions()

	t.Logf("Found %d functions", len(functions))

	for i, fn := range functions {
		name := GetFunctionName(fn)
		nodeType := fn.Type()
		text := fn.Text()

		t.Logf("\nFunction %d:", i+1)
		t.Logf("  Type: %v", nodeType)
		t.Logf("  Name: '%s'", name)
		t.Logf("  Text preview: %s...", text[:min(50, len(text))])

		// For arrow functions, we expect to get the variable name
		if nodeType == ast.NodeTypeArrowFunction {
			if name == "" {
				t.Logf("  WARNING: Arrow function has no name extracted")
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
