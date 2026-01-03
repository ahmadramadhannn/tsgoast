package ast

// InterfaceNode represents an interface declaration.
type InterfaceNode struct {
	BaseNode
	Name           string
	Properties     []*PropertySignature
	Methods        []*MethodSignature
	Extends        []string
	TypeParameters []string
	IsExported     bool
}

// TypeAliasNode represents a type alias declaration.
type TypeAliasNode struct {
	BaseNode
	Name           string
	TypeDefinition string
	TypeParameters []string
	IsExported     bool
}

// PropertySignature represents a property in an interface or type.
type PropertySignature struct {
	Name       string
	Type       string
	IsOptional bool
	IsReadonly bool
}

// MethodSignature represents a method signature in an interface.
type MethodSignature struct {
	Name       string
	Parameters []*Parameter
	ReturnType string
	IsOptional bool
}

// TypeReference represents a reference to a type.
type TypeReference struct {
	Name          string
	TypeArguments []string
}
