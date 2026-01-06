package tsgoast

import (
	"testing"

	"github.com/ahmadramadhannn/tsgoast/ast"
)

func TestParseTree(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	tests := []struct {
		name           string
		source         string
		wantStatements int
		checkFunc      func(*testing.T, *Tree)
	}{
		{
			name: "Variable declarations",
			source: `
				const x = 42;
				let y = "hello";
				var z = true;
			`,
			wantStatements: 3,
			checkFunc: func(t *testing.T, tree *Tree) {
				varCount := 0
				for _, stmt := range tree.Statements {
					if _, ok := stmt.(*ast.VariableStatement); ok {
						varCount++
					}
				}
				if varCount < 1 {
					t.Errorf("Expected at least 1 variable statement, got %d", varCount)
				}
			},
		},
		{
			name: "Function declaration",
			source: `
				function test() {
					return 42;
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				funcCount := 0
				for _, stmt := range tree.Statements {
					if fn, ok := stmt.(*ast.FunctionDeclaration); ok {
						funcCount++
						if fn.Name != "test" {
							t.Errorf("Expected function name 'test', got '%s'", fn.Name)
						}
					}
				}
				if funcCount != 1 {
					t.Errorf("Expected 1 function declaration, got %d", funcCount)
				}
			},
		},
		{
			name: "Class declaration",
			source: `
				class MyClass {
					constructor() {}
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				classCount := 0
				for _, stmt := range tree.Statements {
					if cls, ok := stmt.(*ast.ClassDeclaration); ok {
						classCount++
						if cls.Name != "MyClass" {
							t.Errorf("Expected class name 'MyClass', got '%s'", cls.Name)
						}
					}
				}
				if classCount != 1 {
					t.Errorf("Expected 1 class declaration, got %d", classCount)
				}
			},
		},
		{
			name: "If statement",
			source: `
				if (true) {
					console.log("yes");
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				ifCount := 0
				for _, stmt := range tree.Statements {
					if _, ok := stmt.(*ast.IfStatement); ok {
						ifCount++
					}
				}
				if ifCount != 1 {
					t.Errorf("Expected 1 if statement, got %d", ifCount)
				}
			},
		},
		{
			name: "For loop",
			source: `
				for (let i = 0; i < 10; i++) {
					console.log(i);
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				forCount := 0
				for _, stmt := range tree.Statements {
					if _, ok := stmt.(*ast.ForStatement); ok {
						forCount++
					}
				}
				if forCount != 1 {
					t.Errorf("Expected 1 for statement, got %d", forCount)
				}
			},
		},
		{
			name: "For-of loop",
			source: `
				for (const item of items) {
					console.log(item);
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				forOfCount := 0
				for _, stmt := range tree.Statements {
					if _, ok := stmt.(*ast.ForOfStatement); ok {
						forOfCount++
					}
				}
				if forOfCount != 1 {
					t.Errorf("Expected 1 for-of statement, got %d", forOfCount)
				}
			},
		},
		{
			name: "Try-catch statement",
			source: `
				try {
					doSomething();
				} catch (error) {
					console.error(error);
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				tryCount := 0
				for _, stmt := range tree.Statements {
					if _, ok := stmt.(*ast.TryStatement); ok {
						tryCount++
					}
				}
				if tryCount != 1 {
					t.Errorf("Expected 1 try statement, got %d", tryCount)
				}
			},
		},
		{
			name: "Export declaration",
			source: `
				export function test() {
					return 42;
				}
			`,
			wantStatements: 1,
			checkFunc: func(t *testing.T, tree *Tree) {
				exportCount := 0
				for _, stmt := range tree.Statements {
					if exp, ok := stmt.(*ast.ExportDeclaration); ok {
						exportCount++
						if exp.IsDefault {
							t.Error("Expected non-default export")
						}
					}
				}
				if exportCount != 1 {
					t.Errorf("Expected 1 export declaration, got %d", exportCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := parser.ParseTree([]byte(tt.source))
			if err != nil {
				t.Fatalf("ParseTree() error = %v", err)
			}

			if tree == nil {
				t.Fatal("ParseTree() returned nil tree")
			}

			if tree.Root == nil {
				t.Error("Tree.Root is nil")
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, tree)
			}
		})
	}
}

func TestParseTreeFromFile(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	tree, err := parser.ParseTreeFromFile("testdata/functions.ts")
	if err != nil {
		t.Fatalf("ParseTreeFromFile() error = %v", err)
	}

	if tree == nil {
		t.Fatal("ParseTreeFromFile() returned nil tree")
	}

	if tree.Root == nil {
		t.Error("Tree.Root is nil")
	}

	// Count function declarations
	funcCount := 0
	for _, stmt := range tree.Statements {
		if _, ok := stmt.(*ast.FunctionDeclaration); ok {
			funcCount++
		}
	}

	if funcCount < 1 {
		t.Errorf("Expected at least 1 function declaration in functions.ts, got %d", funcCount)
	}
}

func TestStatementTypes(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	source := []byte(`
		// Variable declarations
		const x = 42;
		let y = "hello";
		
		// Function declaration
		function greet(name: string) {
			return "Hello, " + name;
		}
		
		// Async function
		async function fetchData() {
			return await fetch("/api");
		}
		
		// Class declaration
		class Person {
			name: string;
			constructor(name: string) {
				this.name = name;
			}
		}
		
		// If statement
		if (x > 0) {
			console.log("positive");
		}
		
		// While loop
		while (y < 10) {
			y++;
		}
		
		// For loop
		for (let i = 0; i < 5; i++) {
			console.log(i);
		}
		
		// Switch statement
		switch (x) {
			case 1:
				break;
			default:
				break;
		}
		
		// Try-catch
		try {
			doSomething();
		} catch (e) {
			console.error(e);
		}
		
		// Return statement
		function test() {
			return 42;
		}
		
		// Export
		export const PI = 3.14;
	`)

	tree, err := parser.ParseTree(source)
	if err != nil {
		t.Fatalf("ParseTree() error = %v", err)
	}

	// Count different statement types
	counts := map[string]int{
		"variable":   0,
		"function":   0,
		"class":      0,
		"if":         0,
		"while":      0,
		"for":        0,
		"switch":     0,
		"try":        0,
		"export":     0,
		"expression": 0,
	}

	for _, stmt := range tree.Statements {
		switch stmt.(type) {
		case *ast.VariableStatement:
			counts["variable"]++
		case *ast.FunctionDeclaration:
			counts["function"]++
		case *ast.ClassDeclaration:
			counts["class"]++
		case *ast.IfStatement:
			counts["if"]++
		case *ast.WhileStatement:
			counts["while"]++
		case *ast.ForStatement:
			counts["for"]++
		case *ast.SwitchStatement:
			counts["switch"]++
		case *ast.TryStatement:
			counts["try"]++
		case *ast.ExportDeclaration:
			counts["export"]++
		case *ast.ExpressionStatement:
			counts["expression"]++
		}
	}

	t.Logf("Statement counts: %+v", counts)

	// Verify we found at least some of each type
	if counts["variable"] < 1 {
		t.Error("Expected at least 1 variable statement")
	}
	if counts["function"] < 1 {
		t.Error("Expected at least 1 function declaration")
	}
}

func TestAsyncAndExportedFlags(t *testing.T) {
	parser, err := New()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	defer parser.Close()

	source := []byte(`
		async function asyncFunc() {
			return await Promise.resolve(42);
		}
		
		export function exportedFunc() {
			return 42;
		}
		
		export async function exportedAsyncFunc() {
			return await Promise.resolve(42);
		}
	`)

	tree, err := parser.ParseTree(source)
	if err != nil {
		t.Fatalf("ParseTree() error = %v", err)
	}

	asyncCount := 0
	exportedCount := 0
	exportedAsyncCount := 0

	for _, stmt := range tree.Statements {
		switch fn := stmt.(type) {
		case *ast.FunctionDeclaration:
			if fn.IsAsync {
				asyncCount++
			}
			if fn.IsExported {
				exportedCount++
			}
			if fn.IsAsync && fn.IsExported {
				exportedAsyncCount++
			}
		case *ast.ExportDeclaration:
			exportedCount++
		}
	}

	if asyncCount < 1 {
		t.Error("Expected at least 1 async function")
	}
	if exportedCount < 1 {
		t.Error("Expected at least 1 exported function")
	}
}
