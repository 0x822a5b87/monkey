package ast

import (
	"0x822a5b87/monkey/token"
	"bytes"
	"fmt"
)

// Node every node in our AST has to implement the Node interface
// meaning it has to provide a TokenLiteral() function that returns the literal value of the token it's associated with.
// TokenLiteral() will used only for debugging and testing.
// The AST we are going to construct consists solely of Nodes that are connected to each other.
type Node interface {
	TokenLiteral() string
	String() string // String convert Node to code as string
}

// Statement a statement is a complete unit of execution in a program.
// Statements typically perform an action, such as assigning a value to a variable,
// calling a function, or controlling the flow of the program.
// Statement DO NOT produce value.
type Statement interface {
	Node
	statementNode() // statementNode a Node implement this method to specify itself is a Statement
}

// Expression an expression is a combination of variables, constants, operators, and functions that are evaluated to produce a value.
// it can be simple, like a constant or variable, or complex, involving multiple operations.
// it can be nested with other expressions or statement.
// Expression produce value.
type Expression interface {
	Node
	expressionNode() // expressionNode a Node implement this method to specify itself is an Expression
}

// Program the program node is going to be the root node of every AST our parser produced.
// every valid monkey program is a serials of  statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	buffer := bytes.Buffer{}
	for _, stmt := range p.Statements {
		buffer.WriteString(stmt.String())
	}
	return buffer.String()
}

// Identifier note that identifier is an Expression
type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}
func (identifier *Identifier) String() string {
	return identifier.Value
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier // Name the name of variable
	Value Expression  // Value expression represent the right side of the let statement
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	return fmt.Sprintf("%s %s = %s;", ls.Token.Literal, ls.Name.String(), ls.Value.String())
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
func (r *ReturnStatement) String() string {
	var returnValue = ""
	if r.ReturnValue != nil {
		returnValue = r.ReturnValue.String()
	}

	return fmt.Sprintf("%s %s;", r.Token.Literal, returnValue)
}
func (r *ReturnStatement) statementNode() {}

// ExpressionStatement we need it because it's totally legal in monkey to write the following code:
// let x = 10;
// x + 10;
type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("%d;", i.Value)
}

func (i *IntegerLiteral) expressionNode() {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}

func (p *PrefixExpression) expressionNode() {}
