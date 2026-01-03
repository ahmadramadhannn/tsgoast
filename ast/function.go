package ast

// FunctionNode represents a function declaration.
type FunctionNode struct {
	BaseNode
	Name           string
	Parameters     []*Parameter
	ReturnType     string
	Body           string
	IsAsync        bool
	IsExported     bool
	IsGenerator    bool
	TypeParameters []string
}

// ArrowFunctionNode represents an arrow function expression.
type ArrowFunctionNode struct {
	BaseNode
	Parameters []*Parameter
	ReturnType string
	Body       string
	IsAsync    bool
}

// MethodNode represents a method in a class or interface.
type MethodNode struct {
	BaseNode
	Name       string
	Parameters []*Parameter
	ReturnType string
	Body       string
	IsAsync    bool
	IsStatic   bool
	IsAbstract bool
	Visibility string // "public", "private", "protected"
}

// Parameter represents a function or method parameter.
type Parameter struct {
	Name         string
	Type         string
	IsOptional   bool
	DefaultValue string
	IsRest       bool
}
