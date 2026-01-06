package main

import (
	"fmt"
	"log"

	"github.com/ahmadramadhannn/tsgoast"
	"github.com/ahmadramadhannn/tsgoast/analyzer"
	"github.com/ahmadramadhannn/tsgoast/ast"
)

func main() {
	// Create parser
	parser, err := tsgoast.New()
	if err != nil {
		log.Fatal(err)
	}
	defer parser.Close()

	source := []byte(`
		// Variable declarations
		const PI = 3.14159;
		let count = 0;
		var name = "TypeScript";

		// Function declarations
		function greet(name: string): string {
			return "Hello, " + name;
		}

		async function fetchData(url: string): Promise<any> {
			const response = await fetch(url);
			return response.json();
		}

		// Class declaration
		class Person {
			constructor(public name: string, public age: number) {}
			
			greet(): string {
				return "Hi, I'm " + this.name;
			}
		}

		// Control flow
		if (count > 0) {
			console.log("Positive");
		}

		for (let i = 0; i < 10; i++) {
			console.log(i);
		}

		for (const item of items) {
			process(item);
		}

		// Switch statement
		switch (name) {
			case "TypeScript":
				console.log("Great choice!");
				break;
			default:
				console.log("Unknown");
		}

		// Try-catch
		try {
			doSomething();
		} catch (error) {
			console.error(error);
		}

		// Export
		export function exportedFunc() {
			return 42;
		}
	`)

	fmt.Println("=== Using New Tree API ===\n")

	// Parse into typed tree
	tree, err := parser.ParseTree(source)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over typed statements
	for i, stmt := range tree.Statements {
		fmt.Printf("%d. ", i+1)

		switch s := stmt.(type) {
		case *ast.VariableStatement:
			fmt.Printf("VariableStatement: %s\n", s.Kind)

		case *ast.FunctionDeclaration:
			fmt.Printf("FunctionDeclaration: %s", s.Name)
			if s.IsAsync {
				fmt.Print(" (async)")
			}
			if s.IsExported {
				fmt.Print(" (exported)")
			}
			fmt.Println()

		case *ast.ClassDeclaration:
			fmt.Printf("ClassDeclaration: %s", s.Name)
			if s.IsAbstract {
				fmt.Print(" (abstract)")
			}
			fmt.Println()

		case *ast.IfStatement:
			fmt.Println("IfStatement")

		case *ast.ForStatement:
			fmt.Println("ForStatement")

		case *ast.ForOfStatement:
			fmt.Print("ForOfStatement")
			if s.IsAwait {
				fmt.Print(" (await)")
			}
			fmt.Println()

		case *ast.SwitchStatement:
			fmt.Println("SwitchStatement")

		case *ast.TryStatement:
			fmt.Println("TryStatement")

		case *ast.ExportDeclaration:
			fmt.Print("ExportDeclaration")
			if s.IsDefault {
				fmt.Print(" (default)")
			}
			fmt.Println()

		case *ast.ExpressionStatement:
			fmt.Println("ExpressionStatement")

		default:
			fmt.Printf("Other: %T\n", s)
		}
	}

	fmt.Println("\n=== Using Existing Analyzer API ===\n")

	// You can still use the existing analyzer functions!
	a := analyzer.New(tree.Root)

	// Find all functions
	functions := a.FindFunctions()
	fmt.Printf("Found %d functions:\n", len(functions))
	for _, fn := range functions {
		name := analyzer.GetFunctionName(fn)
		isAsync := analyzer.IsAsync(fn)
		isExported := analyzer.IsExported(fn)
		fmt.Printf("  - %s (async: %v, exported: %v)\n", name, isAsync, isExported)
	}

	fmt.Println()

	// Count different node types
	fmt.Printf("Total identifiers: %d\n", a.CountNodesByType(ast.NodeTypeIdentifier))
	fmt.Printf("Total expressions: %d\n", a.CountNodesByType(ast.NodeTypeExpression))
}
