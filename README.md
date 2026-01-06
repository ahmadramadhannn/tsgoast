# tsgoast

A lightweight Go library for analyzing TypeScript Abstract Syntax Trees (AST) with minimal dependencies.

## Features

- 🚀 **Fast parsing** using tree-sitter
- 🎯 **Simple API** for querying TypeScript code structure
- 📦 **Minimal dependencies** (only tree-sitter)
- ✅ **Well-tested** with comprehensive unit tests
- 🔍 **Type information extraction** for functions, interfaces, type aliases, and more
- 📝 **Idiomatic Go** following official style guide

## Installation

```bash
go get github.com/ahmadramadhannn/tsgoast
```

**Note:** This library uses tree-sitter which requires CGO. Make sure you have a C compiler installed.

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
    // Create a new parser
    parser, err := tsgoast.New()
    if err != nil {
        log.Fatal(err)
    }
    defer parser.Close()

    // Parse TypeScript code
    source := []byte(`
        function greet(name: string): string {
            return "Hello, " + name;
        }

        interface User {
            id: number;
            name: string;
        }
    `)

    root, err := parser.Parse(source)
    if err != nil {
        log.Fatal(err)
    }

    // Analyze the AST
    a := analyzer.New(root)

    // Find all functions
    functions := a.FindFunctions()
    fmt.Printf("Found %d functions\n", len(functions))

    // Find all interfaces
    interfaces := a.FindInterfaces()
    fmt.Printf("Found %d interfaces\n", len(interfaces))
}
```

## Usage Examples

### Parsing a TypeScript File

```go
parser, err := tsgoast.New()
if err != nil {
    log.Fatal(err)
}
defer parser.Close()

root, err := parser.ParseFile("example.ts")
if err != nil {
    log.Fatal(err)
}
```

### Parsing with Typed Statements (NEW!)

The library now supports parsing TypeScript into a typed statement tree:

```go
parser, err := tsgoast.New()
if err != nil {
    log.Fatal(err)
}
defer parser.Close()

// Parse into a typed tree
tree, err := parser.ParseTree([]byte(`
    const x = 42;
    
    function greet(name: string) {
        return "Hello, " + name;
    }
    
    class Person {
        constructor(public name: string) {}
    }
    
    if (x > 0) {
        console.log("positive");
    }
    
    for (let i = 0; i < 10; i++) {
        console.log(i);
    }
`))

if err != nil {
    log.Fatal(err)
}

// Iterate over typed statements
for _, stmt := range tree.Statements {
    switch s := stmt.(type) {
    case *ast.VariableStatement:
        fmt.Printf("Variable declaration: %s\n", s.Kind)
    case *ast.FunctionDeclaration:
        fmt.Printf("Function: %s (async: %v, exported: %v)\n", 
            s.Name, s.IsAsync, s.IsExported)
    case *ast.ClassDeclaration:
        fmt.Printf("Class: %s (abstract: %v)\n", s.Name, s.IsAbstract)
    case *ast.IfStatement:
        fmt.Println("If statement")
    case *ast.ForStatement:
        fmt.Println("For loop")
    case *ast.ForOfStatement:
        fmt.Printf("For-of loop (await: %v)\n", s.IsAwait)
    case *ast.WhileStatement:
        fmt.Println("While loop")
    case *ast.SwitchStatement:
        fmt.Println("Switch statement")
    case *ast.TryStatement:
        fmt.Println("Try-catch statement")
    case *ast.ExportDeclaration:
        fmt.Printf("Export (default: %v)\n", s.IsDefault)
    }
}
```



### Finding Functions

```go
a := analyzer.New(root)

// Find all functions
functions := a.FindFunctions()

for _, fn := range functions {
    fmt.Println("Function:", fn.Text())
    
    // Check if async
    if analyzer.IsAsync(fn) {
        fmt.Println("  - This is an async function")
    }
    
    // Check if exported
    if analyzer.IsExported(fn) {
        fmt.Println("  - This function is exported")
    }
    
    // Get function name
    name := analyzer.GetFunctionName(fn)
    fmt.Println("  - Name:", name)
}
```

### Finding Interfaces and Type Aliases

```go
a := analyzer.New(root)

// Find all interfaces
interfaces := a.FindInterfaces()
for _, iface := range interfaces {
    name := analyzer.GetInterfaceName(iface)
    fmt.Println("Interface:", name)
    
    // Check if it extends another interface
    if analyzer.HasExtends(iface) {
        fmt.Println("  - Extends another interface")
    }
    
    // Count properties
    propCount := analyzer.CountProperties(iface)
    fmt.Printf("  - Has %d properties\n", propCount)
}

// Find all type aliases
typeAliases := a.FindTypeAliases()
for _, ta := range typeAliases {
    name := analyzer.GetTypeAliasName(ta)
    fmt.Println("Type alias:", name)
}
```

### Custom AST Traversal

```go
a := analyzer.New(root)

// Visit all nodes
a.Visit(func(node ast.Node) bool {
    fmt.Printf("Node type: %s, Text: %s\n", node.Type(), node.Text())
    return true // continue traversal
})

// Find nodes with custom predicate
identifiers := a.FindNodes(func(node ast.Node) bool {
    return node.Type() == ast.NodeTypeIdentifier
})
```

### Counting Nodes

```go
a := analyzer.New(root)

// Count all functions
functionCount := a.CountNodesByType(ast.NodeTypeFunction)
fmt.Printf("Total functions: %d\n", functionCount)

// Count with custom predicate
asyncCount := a.CountNodes(func(node ast.Node) bool {
    return analyzer.IsAsync(node)
})
fmt.Printf("Async functions: %d\n", asyncCount)
```

## API Overview

### Parser

- `New() (*Parser, error)` - Create a new TypeScript parser
- `Parse(source []byte) (*ast.BaseNode, error)` - Parse TypeScript source code
- `ParseFile(path string) (*ast.BaseNode, error)` - Parse a TypeScript file
- `ParseTree(source []byte) (*Tree, error)` - Parse into a typed statement tree (NEW!)
- `ParseTreeFromFile(path string) (*Tree, error)` - Parse file into a typed statement tree (NEW!)
- `Close()` - Release parser resources

### Tree (NEW!)

- `Root *ast.BaseNode` - The root AST node
- `Statements []ast.Statement` - Typed statement list

### Analyzer

- `New(root *ast.BaseNode) *Analyzer` - Create a new analyzer
- `Visit(visitor func(node ast.Node) bool)` - Traverse the AST
- `FindNodes(predicate func(node ast.Node) bool) []ast.Node` - Find nodes matching predicate
- `FindNodesByType(nodeType ast.NodeType) []ast.Node` - Find nodes by type
- `CountNodes(predicate func(node ast.Node) bool) int` - Count matching nodes
- `CountNodesByType(nodeType ast.NodeType) int` - Count nodes by type

### Function Analysis

- `FindFunctions() []ast.Node` - Find all function declarations
- `FindMethods() []ast.Node` - Find all method definitions
- `IsAsync(node ast.Node) bool` - Check if function is async
- `IsExported(node ast.Node) bool` - Check if function is exported
- `IsGenerator(node ast.Node) bool` - Check if function is a generator
- `GetFunctionName(node ast.Node) string` - Get function name
- `HasParameters(node ast.Node) bool` - Check if function has parameters
- `CountParameters(node ast.Node) int` - Count function parameters

### Type Analysis

- `FindInterfaces() []ast.Node` - Find all interface declarations
- `FindTypeAliases() []ast.Node` - Find all type alias declarations
- `GetInterfaceName(node ast.Node) string` - Get interface name
- `GetTypeAliasName(node ast.Node) string` - Get type alias name
- `HasExtends(node ast.Node) bool` - Check if interface extends another
- `IsReadonly(node ast.Node) bool` - Check if property is readonly
- `IsOptionalProperty(node ast.Node) bool` - Check if property is optional
- `CountProperties(node ast.Node) int` - Count properties in interface/type
- `IsGenericType(node ast.Node) bool` - Check if type has type parameters

### Expression Analysis

- `FindExpressions() []ast.Node` - Find all expressions
- `FindIdentifiers() []ast.Node` - Find all identifiers
- `FindLiterals() []ast.Node` - Find all literals
- `GetIdentifierName(node ast.Node) string` - Get identifier name
- `GetLiteralValue(node ast.Node) string` - Get literal value

## AST Node Types

The library recognizes the following TypeScript constructs:

- `NodeTypeFunction` - Function declarations
- `NodeTypeArrowFunction` - Arrow functions
- `NodeTypeMethod` - Method definitions
- `NodeTypeInterface` - Interface declarations
- `NodeTypeTypeAlias` - Type alias declarations
- `NodeTypeExpression` - Expressions
- `NodeTypeIdentifier` - Identifiers
- `NodeTypeLiteral` - Literal values
- `NodeTypeProperty` - Properties
- `NodeTypeParameter` - Parameters

## Statement Types (NEW!)

The library now provides typed statement nodes for more structured AST analysis:

### Declarations

- `VariableStatement` - Variable declarations (const, let, var)
- `FunctionDeclaration` - Function declarations
- `ClassDeclaration` - Class declarations
- `EnumDeclaration` - Enum declarations
- `ImportDeclaration` - Import statements
- `ExportDeclaration` - Export statements
- `NamespaceDeclaration` - Namespace declarations

### Control Flow

- `IfStatement` - If statements
- `WhileStatement` - While loops
- `ForStatement` - For loops
- `ForInStatement` - For...in loops
- `ForOfStatement` - For...of loops
- `SwitchStatement` - Switch statements
- `TryStatement` - Try-catch-finally statements

### Other Statements

- `ExpressionStatement` - Expression statements
- `ReturnStatement` - Return statements
- `ThrowStatement` - Throw statements
- `BreakStatement` - Break statements
- `ContinueStatement` - Continue statements
- `BlockStatement` - Block statements
- `EmptyStatement` - Empty statements
- `LabeledStatement` - Labeled statements
- `WithStatement` - With statements
- `DebuggerStatement` - Debugger statements

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Requirements

- Go 1.21 or later
- CGO enabled (for tree-sitter)
- C compiler (gcc, clang, etc.)

## Dependencies

- [tree-sitter/go-tree-sitter](https://github.com/tree-sitter/go-tree-sitter) - Tree-sitter Go bindings
- [tree-sitter/tree-sitter-typescript](https://github.com/tree-sitter/tree-sitter-typescript) - TypeScript grammar

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see LICENSE file for details

## Roadmap

Future enhancements may include:

- [ ] More detailed type information extraction
- [ ] Support for decorators
- [ ] Class analysis utilities
- [ ] Import/export analysis
- [ ] JSDoc comment extraction
- [ ] Source map generation
- [ ] Pretty printing utilities

## Acknowledgments

This library is built on top of the excellent [tree-sitter](https://tree-sitter.github.io/) parsing library.
