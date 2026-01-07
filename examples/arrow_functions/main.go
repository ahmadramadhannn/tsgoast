// Example arrow_functions demonstrates GetFunctionName with arrow functions.
//
// Run with: go run main.go
package main

import (
	"fmt"
	"log"

	"github.com/ahmadramadhannn/tsgoast"
	"github.com/ahmadramadhannn/tsgoast/analyzer"
)

func main() {
	parser, err := tsgoast.New()
	if err != nil {
		log.Fatal(err)
	}
	defer parser.Close()

	source := []byte(`
		// Regular arrow function
		const regularArrow = () => {
			return 42;
		};

		// Exported arrow function
		export const exportedArrow = () => {
			return 100;
		};

		// Exported async arrow function
		export const exportedAsync = async () => {
			return await Promise.resolve(200);
		};

		// Arrow with parameters
		const withParams = (x: number, y: number) => x + y;

		// Regular function for comparison
		function regularFunc() {
			return 42;
		}
	`)

	root, err := parser.Parse(source)
	if err != nil {
		log.Fatal(err)
	}

	a := analyzer.New(root)
	functions := a.FindFunctions()

	fmt.Println("=== Arrow Function Name Extraction ===")
	fmt.Println()
	for i, fn := range functions {
		name := analyzer.GetFunctionName(fn)
		isAsync := analyzer.IsAsync(fn)
		nodeType := fn.Type()

		fmt.Printf("%d. %s\n", i+1, name)
		fmt.Printf("   Type: %v\n", nodeType)
		fmt.Printf("   Async: %v\n", isAsync)
		fmt.Println()
	}
}
