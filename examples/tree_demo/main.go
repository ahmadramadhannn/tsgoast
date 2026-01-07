// Example tree_demo demonstrates the Tree API for parsing TypeScript.
//
// Run with: go run main.go
package main

import (
	"fmt"
	"log"

	"github.com/ahmadramadhannn/tsgoast"
	"github.com/ahmadramadhannn/tsgoast/analyzer"
	"github.com/ahmadramadhannn/tsgoast/ast"
)

func main() {
	parser, err := tsgoast.New()
	if err != nil {
		log.Fatal(err)
	}
	defer parser.Close()

	source := []byte(`
		const PI = 3.14159;
		let count = 0;

		function greet(name: string): string {
			return "Hello, " + name;
		}

		async function fetchData(url: string): Promise<any> {
			const response = await fetch(url);
			return response.json();
		}

		class Person {
			constructor(public name: string) {}
		}

		if (count > 0) {
			console.log("Positive");
		}

		for (let i = 0; i < 10; i++) {
			console.log(i);
		}

		export function exportedFunc() {
			return 42;
		}
	`)

	fmt.Println("=== Using Tree API ===")
	fmt.Println()

	tree, err := parser.ParseTree(source)
	if err != nil {
		log.Fatal(err)
	}

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
			fmt.Printf("ClassDeclaration: %s\n", s.Name)
		case *ast.IfStatement:
			fmt.Println("IfStatement")
		case *ast.ForStatement:
			fmt.Println("ForStatement")
		case *ast.ExportDeclaration:
			fmt.Println("ExportDeclaration")
		case *ast.ExpressionStatement:
			fmt.Println("ExpressionStatement")
		default:
			fmt.Printf("Other: %T\n", s)
		}
	}

	fmt.Println()
	fmt.Println("=== Using Analyzer API ===")
	fmt.Println()

	a := analyzer.New(tree.Root)
	functions := a.FindFunctions()
	fmt.Printf("Found %d functions:\n", len(functions))
	for _, fn := range functions {
		name := analyzer.GetFunctionName(fn)
		isAsync := analyzer.IsAsync(fn)
		fmt.Printf("  - %s (async: %v)\n", name, isAsync)
	}
}
