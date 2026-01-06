package ast

// Statement represents any TypeScript statement.
type Statement interface {
	Node
	statementNode()
}

// Declaration represents any TypeScript declaration.
type Declaration interface {
	Statement
	declarationNode()
}

// VariableStatement represents a variable declaration statement (var, let, const).
type VariableStatement struct {
	BaseNode
	Declarations []*VariableDeclarator
	Kind         string // "var", "let", or "const"
}

func (v *VariableStatement) statementNode() {}

// VariableDeclarator represents a single variable declarator.
type VariableDeclarator struct {
	BaseNode
	Name        string
	Type        string
	Initializer Node
}

// FunctionDeclaration represents a function declaration statement.
type FunctionDeclaration struct {
	BaseNode
	Name           string
	Parameters     []*Parameter
	ReturnType     string
	Body           *BlockStatement
	IsAsync        bool
	IsExported     bool
	IsGenerator    bool
	TypeParameters []string
}

func (f *FunctionDeclaration) statementNode()   {}
func (f *FunctionDeclaration) declarationNode() {}

// ClassDeclaration represents a class declaration.
type ClassDeclaration struct {
	BaseNode
	Name           string
	SuperClass     string
	Body           *ClassBody
	TypeParameters []string
	IsAbstract     bool
	IsExported     bool
	Decorators     []string
}

func (c *ClassDeclaration) statementNode()   {}
func (c *ClassDeclaration) declarationNode() {}

// ClassBody represents the body of a class.
type ClassBody struct {
	BaseNode
	Members []Node
}

// ExpressionStatement represents an expression statement.
type ExpressionStatement struct {
	BaseNode
	Expression Node
}

func (e *ExpressionStatement) statementNode() {}

// IfStatement represents an if statement.
type IfStatement struct {
	BaseNode
	Condition   Node
	Consequence *BlockStatement
	Alternative Node // can be *IfStatement or *BlockStatement
}

func (i *IfStatement) statementNode() {}

// WhileStatement represents a while loop.
type WhileStatement struct {
	BaseNode
	Condition Node
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode() {}

// ForStatement represents a for loop.
type ForStatement struct {
	BaseNode
	Initializer Node
	Condition   Node
	Increment   Node
	Body        *BlockStatement
}

func (f *ForStatement) statementNode() {}

// ForInStatement represents a for...in loop.
type ForInStatement struct {
	BaseNode
	Left  Node
	Right Node
	Body  *BlockStatement
}

func (f *ForInStatement) statementNode() {}

// ForOfStatement represents a for...of loop.
type ForOfStatement struct {
	BaseNode
	Left    Node
	Right   Node
	Body    *BlockStatement
	IsAwait bool
}

func (f *ForOfStatement) statementNode() {}

// SwitchStatement represents a switch statement.
type SwitchStatement struct {
	BaseNode
	Discriminant Node
	Cases        []*SwitchCase
}

func (s *SwitchStatement) statementNode() {}

// SwitchCase represents a case in a switch statement.
type SwitchCase struct {
	BaseNode
	Test       Node // nil for default case
	Consequent []Statement
}

// TryStatement represents a try-catch-finally statement.
type TryStatement struct {
	BaseNode
	Body      *BlockStatement
	Handler   *CatchClause
	Finalizer *BlockStatement
}

func (t *TryStatement) statementNode() {}

// CatchClause represents a catch clause.
type CatchClause struct {
	BaseNode
	Parameter string
	ParamType string
	Body      *BlockStatement
}

// ThrowStatement represents a throw statement.
type ThrowStatement struct {
	BaseNode
	Argument Node
}

func (t *ThrowStatement) statementNode() {}

// ReturnStatement represents a return statement.
type ReturnStatement struct {
	BaseNode
	Argument Node
}

func (r *ReturnStatement) statementNode() {}

// BreakStatement represents a break statement.
type BreakStatement struct {
	BaseNode
	Label string
}

func (b *BreakStatement) statementNode() {}

// ContinueStatement represents a continue statement.
type ContinueStatement struct {
	BaseNode
	Label string
}

func (c *ContinueStatement) statementNode() {}

// BlockStatement represents a block statement.
type BlockStatement struct {
	BaseNode
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

// EmptyStatement represents an empty statement (;).
type EmptyStatement struct {
	BaseNode
}

func (e *EmptyStatement) statementNode() {}

// LabeledStatement represents a labeled statement.
type LabeledStatement struct {
	BaseNode
	Label     string
	Statement Statement
}

func (l *LabeledStatement) statementNode() {}

// WithStatement represents a with statement.
type WithStatement struct {
	BaseNode
	Object Node
	Body   *BlockStatement
}

func (w *WithStatement) statementNode() {}

// DebuggerStatement represents a debugger statement.
type DebuggerStatement struct {
	BaseNode
}

func (d *DebuggerStatement) statementNode() {}

// ImportDeclaration represents an import statement.
type ImportDeclaration struct {
	BaseNode
	Specifiers []Node
	Source     string
}

func (i *ImportDeclaration) statementNode()   {}
func (i *ImportDeclaration) declarationNode() {}

// ExportDeclaration represents an export statement.
type ExportDeclaration struct {
	BaseNode
	Declaration Node
	Specifiers  []Node
	Source      string
	IsDefault   bool
}

func (e *ExportDeclaration) statementNode()   {}
func (e *ExportDeclaration) declarationNode() {}

// EnumDeclaration represents an enum declaration.
type EnumDeclaration struct {
	BaseNode
	Name       string
	Members    []*EnumMember
	IsConst    bool
	IsExported bool
}

func (e *EnumDeclaration) statementNode()   {}
func (e *EnumDeclaration) declarationNode() {}

// EnumMember represents a member of an enum.
type EnumMember struct {
	BaseNode
	Name        string
	Initializer Node
}

// NamespaceDeclaration represents a namespace declaration.
type NamespaceDeclaration struct {
	BaseNode
	Name       string
	Body       []Statement
	IsExported bool
}

func (n *NamespaceDeclaration) statementNode()   {}
func (n *NamespaceDeclaration) declarationNode() {}
