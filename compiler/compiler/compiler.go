package compiler

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/common"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/token"
	"fmt"
)

type instructionIndex int

type Compiler struct {
	instructions code.Instructions
	constants    *code.Constants
}

func NewCompiler() *Compiler {
	c := &Compiler{
		instructions: make(code.Instructions, 0),
		constants:    code.NewConstants(),
	}

	return c
}

// Compile
// Assume that we are going to build a minimal compiler for adding 1 + 2
// 1. traverse the AST we pass in, find the *ast.IntegerLiteral nodes.
// 2. evaluate them by turning them into *object.Integer objects, add the object to constant pool.
// 3. emit code.OpConstant instructions that reference the Constants in said pool.
func (c *Compiler) Compile(root ast.Node) error {
	switch node := root.(type) {
	case *ast.Program:
		return c.compileProgram(node)
	case ast.Expression:
		return c.compileExpression(node)
	case ast.Statement:
		return c.compileStatement(node)
	}
	return nil
}

// ByteCode transform a compiled AST into bytecode
func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// ByteCode is what we'll pass to the VM and make assertions about in our compiler tests.
type ByteCode struct {
	Instructions code.Instructions
	Constants    *code.Constants
}

func (c *Compiler) compileProgram(program *ast.Program) error {
	for _, stmt := range program.Statements {
		err := c.Compile(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileStatement(statement ast.Statement) error {
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		return c.compileExpressionStatement(stmt)
		// TODO support more statement type
	}

	return common.NewErrUnsupportedCompilingNode(statement.String())
}

func (c *Compiler) compileExpressionStatement(statement *ast.ExpressionStatement) error {
	err := c.Compile(statement.Expr)
	if err != nil {
		return err
	}
	// expression will produce an object on stack which can produce a stack overflow.
	// we assert that the compiled expression statement should be followed by an OpPop instruction.
	c.emit(code.OpPop)
	return nil
}

func (c *Compiler) compileExpression(expr ast.Expression) error {
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return c.compileIntegerLiteral(expr)
	case *ast.BooleanExpression:
		return c.compileBooleanExpression(expr)
		// TODO support more expression type
	case *ast.InfixExpression:
		return c.compileInfixExpression(expr)
	case *ast.PrefixExpression:
		return c.compilePrefixExpression(expr)
	}
	return common.NewErrUnsupportedCompilingNode(expr.String())
}

func (c *Compiler) compileIntegerLiteral(literal *ast.IntegerLiteral) error {
	integer := &object.Integer{Value: literal.Value}
	index := c.constants.AddConstant(integer)
	c.emit(code.OpConstant, index.IntValue())
	return nil
}

func (c *Compiler) compileBooleanExpression(literal *ast.BooleanExpression) error {
	if literal.Value {
		c.emit(code.OpTrue)
	} else {
		c.emit(code.OpFalse)
	}
	return nil
}

func (c *Compiler) compileInfixExpression(infixExpr *ast.InfixExpression) error {
	err := c.compileExpression(infixExpr.Lhs)
	if err != nil {
		return err
	}
	err = c.compileExpression(infixExpr.Rhs)
	if err != nil {
		return err
	}
	return c.compileInfixOperator(infixExpr.Operator)
}

func (c *Compiler) compilePrefixExpression(prefixExpr *ast.PrefixExpression) error {
	err := c.compileExpression(prefixExpr.Right)
	if err != nil {
		return err
	}

	return c.compilePrefixOperator(prefixExpr.Operator)
}

func (c *Compiler) compileInfixOperator(operator string) error {
	switch operator {
	case string(token.PLUS):
		c.emit(code.OpAdd)
		return nil
	case string(token.SUB):
		c.emit(code.OpSub)
		return nil
	case string(token.ASTERISK):
		c.emit(code.OpMul)
		return nil
	case string(token.SLASH):
		c.emit(code.OpDiv)
		return nil
	case string(token.GT):
		c.emit(code.OpGreaterThan)
		return nil
	case string(token.LT):
		c.emit(code.OpLessThan)
		return nil
	case string(token.EQ):
		c.emit(code.OpEqual)
		return nil
	case string(token.NotEq):
		c.emit(code.OpNotEqual)
		return nil
		// TODO support more operator
	}

	return common.NewErrUnsupportedCompilingNode(fmt.Sprintf(" infix [%s]", operator))
}

func (c *Compiler) compilePrefixOperator(operator string) error {
	switch operator {
	case string(token.SUB):
		c.emit(code.OpMinus)
		return nil
	case string(token.BANG):
		c.emit(code.OpBang)
		return nil
		// TODO support more operator
	}

	return common.NewErrUnsupportedCompilingNode(fmt.Sprintf(" prefix [%s]", operator))
}

func (c *Compiler) emit(op code.Opcode, operands ...int) instructionIndex {
	instruction := code.Make(op, operands...)
	var index = instructionIndex(len(c.instructions))
	c.instructions = append(c.instructions, instruction...)
	return index
}

func (c *Compiler) compileOperatorPlus() error {
	c.emit(code.OpAdd)
	return nil
}
