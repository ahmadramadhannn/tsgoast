# tsgoast

A Go library for parsing TypeScript AST using tree-sitter.

## Installation

```bash
go get github.com/ahmadramadhannn/tsgoast@v0.1.0
```

**Note:** Requires CGO and a C compiler.

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/ahmadramadhannn/tsgoast"
    "github.com/ahmadramadhannn/tsgoast/analyzer"
)

func main() {
    parser, _ := tsgoast.New()
    defer parser.Close()

    root, _ := parser.Parse([]byte(`
        function greet(name: string): string {
            return "Hello, " + name;
        }

        const add = (a: number, b: number) => a + b;
    `))

    a := analyzer.New(root)
    for _, fn := range a.FindFunctions() {
        fmt.Println(analyzer.GetFunctionName(fn))
    }
    // Output: greet, add
}
```

## Tree API

Parse into typed statements:

```go
tree, _ := parser.ParseTree(source)

for _, stmt := range tree.Statements {
    switch s := stmt.(type) {
    case *ast.FunctionDeclaration:
        fmt.Printf("Function: %s\n", s.Name)
    case *ast.VariableStatement:
        fmt.Printf("Variable: %s\n", s.Kind)
    case *ast.ClassDeclaration:
        fmt.Printf("Class: %s\n", s.Name)
    }
}
```

## Analyzer API

```go
a := analyzer.New(root)

// Find nodes
a.FindFunctions()      // All functions (including arrow functions)
a.FindInterfaces()     // All interfaces
a.FindTypeAliases()    // All type aliases

// Inspect functions
analyzer.GetFunctionName(fn)  // Works with arrow functions too
analyzer.IsAsync(fn)
analyzer.IsExported(fn)
```

## Examples

```bash
go run examples/tree_demo/main.go
go run examples/arrow_functions/main.go
```

## License

MIT
