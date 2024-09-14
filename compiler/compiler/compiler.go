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

func (i instructionIndex) add(delta int) int {
	return int(i) + delta
}

type Compiler struct {
	instructions code.Instructions
	constants    *code.Constants
	symbolTable  *SymbolTable
}

func NewCompiler() *Compiler {
	c := &Compiler{
		instructions: make(code.Instructions, 0),
		constants:    code.NewConstants(),
		symbolTable:  NewSymbolTable(),
	}

	return c
}

func NewCompilerWithState(prev *Compiler) *Compiler {
	c := NewCompiler()
	c.symbolTable = prev.symbolTable
	c.constants = prev.constants
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
	case *ast.BlockStatement:
		return c.compileBlockStatement(stmt)
	case *ast.LetStatement:
		return c.compileLetStatement(stmt)
	}

	return common.NewErrUnsupportedCompilingNode(statement.String())
}

func (c *Compiler) compileExpressionStatement(statement *ast.ExpressionStatement) error {
	err := c.Compile(statement.Expr)
	if err != nil {
		return err
	}

	switch statement.Expr.(type) {
	case *ast.IfExpression:
		return nil
	default:
		// expression will produce an object on stack which can produce a stack overflow.
		// we assert that the compiled expression statement should be followed by an OpPop instruction.
		c.emit(code.OpPop)
		return nil
	}
}

func (c *Compiler) compileBlockStatement(statement *ast.BlockStatement) error {
	for _, stmt := range statement.Statements {
		err := c.Compile(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileLetStatement(statement *ast.LetStatement) error {
	err := c.Compile(statement.Value)
	if err != nil {
		return err
	}

	symbol := c.symbolTable.Define(statement.Name.Value)
	c.emit(code.OpSetGlobal, symbol.Index)

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
	case *ast.IfExpression:
		return c.compileIfExpression(expr)
	case *ast.Identifier:
		return c.compileIdentifier(expr)
	case *ast.StringLiteral:
		return c.compileStringLiteral(expr)
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

func (c *Compiler) compileIfExpression(ifExpr *ast.IfExpression) error {
	err := c.compileExpression(ifExpr.Condition)
	if err != nil {
		return err
	}

	jumpNotTruthyIndex := c.emit(code.OpJumpNotTruthy, 0)
	err = c.Compile(ifExpr.Consequence)
	if err != nil {
		return err
	}

	if ifExpr.Alternative != nil {
		jumpIndex := c.emit(code.OpJump, 0)
		err = c.Compile(ifExpr.Alternative)
		if err != nil {
			return err
		}

		// jump to length - 1, as the for loop will increment 1 automatically
		c.replaceOperand(jumpIndex, c.instructions.Len()-1)

		// jump to the end byte of NOT_MATTER_WHAT_JUMP
		c.replaceOperand(jumpNotTruthyIndex, jumpIndex.add(2))
	} else {
		var jumpIndex = instructionIndex(len(c.instructions))
		c.replaceOperand(jumpNotTruthyIndex, jumpIndex.add(-1))
	}

	return nil
}

func (c *Compiler) compileIdentifier(identifier *ast.Identifier) error {
	symbol, ok := c.symbolTable.Resolve(identifier.Value)
	if !ok {
		return common.NewUnresolvedVariable(identifier.Value)
	}
	c.emit(code.OpGetGlobal, symbol.Index)
	return nil
}

func (c *Compiler) compileStringLiteral(stringLiteral *ast.StringLiteral) error {
	stringObj := &object.StringObj{Value: stringLiteral.Literal}
	constantIndex := c.constants.AddConstant(stringObj)
	c.emit(code.OpConstant, constantIndex.IntValue())
	return nil
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

func (c *Compiler) replaceOperand(instructionBeginIndex instructionIndex, operands ...int) {
	op := c.instructions[instructionBeginIndex]
	instruction := code.Make(code.Opcode(op), operands...)
	c.replaceInstruction(instructionBeginIndex, instruction)
}

func (c *Compiler) replaceInstruction(instructionBeginIndex instructionIndex, another code.Instructions) {
	for i := 0; i < another.Len(); i++ {
		c.instructions[instructionBeginIndex.add(i)] = another[i]
	}
}

func (c *Compiler) compileOperatorPlus() error {
	c.emit(code.OpAdd)
	return nil
}
