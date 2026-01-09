package analyzer

import (
	"testing"

	"github.com/ahmadramadhannn/tsgoast"
	"github.com/ahmadramadhannn/tsgoast/ast"
)

func TestIsExported(t *testing.T) {
	parser, err := tsgoast.New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	source := []byte(`
		// 1. Regular exported function
		export function exportedRegular() {}

		// 2. Regular non-exported function
		function nonExportedRegular() {}

		// 3. Exported arrow function
		export const exportedArrow = () => {};

		// 4. Exported async arrow function
		export const exportedAsyncArrow = async () => {};

		// 5. Non-exported arrow function
		const nonExportedArrow = () => {};
		
		// 6. Default export function
		export default function defaultExported() {}

		// 7. Exported via declaration
		function separateExport() {}
		export { separateExport };
	`)

	root, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	a := New(root)
	functions := a.FindFunctions()

	tests := []struct {
		name       string
		isExported bool
	}{
		{"exportedRegular", true},
		{"nonExportedRegular", false},
		{"exportedArrow", true},
		{"exportedAsyncArrow", true},
		{"nonExportedArrow", false},
		{"defaultExported", true},
		// separateExport is tricky depending on how FindFunctions behaves with export declaration references
		// for now we focus on direct exports
	}

	// Helper to find function by name
	findFunc := func(name string) ast.Node {
		for _, fn := range functions {
			if GetFunctionName(fn) == name {
				return fn
			}
		}
		return nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := findFunc(tt.name)
			if fn == nil {
				t.Fatalf("Function %s not found", tt.name)
			}

			got := IsExported(fn)
			if got != tt.isExported {
				t.Errorf("IsExported(%s) = %v, want %v", tt.name, got, tt.isExported)
			}
		})
	}
}
