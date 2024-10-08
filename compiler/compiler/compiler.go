package compiler

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/common"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/token"
	"fmt"
	"sort"
)

type instructionIndex int

func (i instructionIndex) add(delta int) int {
	return int(i) + delta
}

type Compiler struct {
	constants   *code.Constants
	symbolTable *SymbolTable

	scopes     []*CompilationScope
	scopeIndex int
}

// ByteCode is what we'll pass to the VM and make assertions about in our compiler tests.
type ByteCode struct {
	Instructions code.Instructions
	Constants    *code.Constants
}

func NewCompiler() *Compiler {
	mainScope := &CompilationScope{
		instructions: make(code.Instructions, 0),
	}

	c := &Compiler{
		constants:   code.NewConstants(),
		symbolTable: NewGlobalSymbolTable(),
		scopes:      []*CompilationScope{mainScope},
		scopeIndex:  0,
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
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
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
	case *ast.BlockStatement:
		return c.compileBlockStatement(stmt)
	case *ast.LetStatement:
		return c.compileLetStatement(stmt)
	case *ast.ReturnStatement:
		return c.compileReturnStatement(stmt)
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
	c.emitSetScope(symbol)

	return nil
}

// we don't emit code.OpReturn or code.OpReturnValue, leave this responsibility to the function
func (c *Compiler) compileReturnStatement(statement *ast.ReturnStatement) error {
	err := c.Compile(statement.ReturnValue)
	if err != nil {
		return err
	}
	c.emit(code.OpReturnValue)
	return nil
}

func (c *Compiler) compileExpression(expr ast.Expression) error {
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return c.compileIntegerLiteral(expr)
	case *ast.BooleanExpression:
		return c.compileBooleanExpression(expr)
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
	case *ast.ArrayLiteral:
		return c.compileArrayLiteral(expr)
	case *ast.HashExpression:
		return c.compileHashExpression(expr)
	case *ast.IndexExpression:
		return c.compileIndexExpression(expr)
	case *ast.FnLiteral:
		return c.compileFnLiteral(expr)
	case *ast.CallExpression:
		return c.compileCallExpression(expr)

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
		c.replaceOperand(jumpIndex, c.currentInstructions().Len()-1)

		// jump to the end byte of NOT_MATTER_WHAT_JUMP
		c.replaceOperand(jumpNotTruthyIndex, jumpIndex.add(2))
	} else {
		var jumpIndex = instructionIndex(len(c.currentInstructions()))
		c.replaceOperand(jumpNotTruthyIndex, jumpIndex.add(-1))
	}

	return nil
}

func (c *Compiler) compileIdentifier(identifier *ast.Identifier) error {
	symbol, ok := c.symbolTable.Resolve(identifier.Value)
	if !ok {
		return common.NewUnresolvedVariable(identifier.Value)
	}
	c.emitGetScope(symbol)
	return nil
}

func (c *Compiler) compileStringLiteral(stringLiteral *ast.StringLiteral) error {
	stringObj := &object.StringObj{Value: stringLiteral.Literal}
	constantIndex := c.constants.AddConstant(stringObj)
	c.emit(code.OpConstant, constantIndex.IntValue())
	return nil
}

func (c *Compiler) compileArrayLiteral(arrayLiteral *ast.ArrayLiteral) error {
	for _, element := range arrayLiteral.Elements {
		err := c.Compile(element)
		if err != nil {
			return err
		}
	}
	c.emit(code.OpArray, len(arrayLiteral.Elements))
	return nil
}

func (c *Compiler) compileHashExpression(hashExpression *ast.HashExpression) error {
	// sort Paris for test
	pairs := make([]ExpressionPair, 0, len(hashExpression.Pairs))
	for key, value := range hashExpression.Pairs {
		pairs = append(pairs, ExpressionPair{key, value})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key.String() < pairs[j].Key.String()
	})

	for _, pair := range pairs {
		err := c.Compile(pair.Key)
		if err != nil {
			return err
		}
		err = c.Compile(pair.Value)
		if err != nil {
			return err
		}
	}
	c.emit(code.OpHash, len(hashExpression.Pairs)*2)
	return nil
}

func (c *Compiler) compileIndexExpression(indexExpr *ast.IndexExpression) error {
	err := c.Compile(indexExpr.Lhs)
	if err != nil {
		return err
	}
	err = c.Compile(indexExpr.Index)
	if err != nil {
		return err
	}
	c.emit(code.OpIndex)
	return nil
}

func (c *Compiler) compileFnLiteral(literal *ast.FnLiteral) error {
	c.enterScope()

	for _, param := range literal.Parameters {
		// just allocate memory for future arguments binding
		c.symbolTable.Define(param.Value)
	}

	err := c.Compile(literal.Body)
	if err != nil {
		return err
	}

	c.completeOpReturn(literal)

	return c.genClosure()
}

func (c *Compiler) compileCallExpression(call *ast.CallExpression) error {
	err := c.compileExpression(call.Fn)
	if err != nil {
		return err
	}

	for _, argument := range call.Arguments {
		err = c.compileExpression(argument)
		if err != nil {
			return err
		}
	}

	c.emit(code.OpCall, len(call.Arguments))
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
	}

	return common.NewErrUnsupportedCompilingNode(fmt.Sprintf(" prefix [%s]", operator))
}

func (c *Compiler) emit(op code.Opcode, operands ...int) instructionIndex {
	instruction := code.Make(op, operands...)
	pos := c.addInstruction(instruction)
	c.setLastEmitInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) instructionIndex {
	var index = instructionIndex(len(c.currentInstructions()))
	c.updateCurrentInstructions(ins)
	return index
}

func (c *Compiler) setLastEmitInstruction(op code.Opcode, pos instructionIndex) {
	scope := c.currentScope()
	scope.previous = scope.last
	scope.last = NewEmittedInstruction(op, pos)
}

func (c *Compiler) isLastInstructionMatch(op code.Opcode) bool {
	return c.currentScope().last != nil && c.currentScope().last.Opcode == op
}

func (c *Compiler) updateLastPopInstruction(op code.Opcode, operands ...int) {
	last := c.currentScope().last
	newInstruction := code.Make(op, operands...)
	c.replaceInstruction(last.Position, newInstruction)
	c.setLastEmitInstruction(op, last.Position)
}

func (c *Compiler) currentScope() *CompilationScope {
	return c.scopes[c.scopeIndex]
}

func (c *Compiler) replaceOperand(instructionBeginIndex instructionIndex, operands ...int) {
	op := c.currentInstructions()[instructionBeginIndex]
	instruction := code.Make(code.Opcode(op), operands...)
	c.replaceInstruction(instructionBeginIndex, instruction)
}

func (c *Compiler) replaceInstruction(instructionBeginIndex instructionIndex, another code.Instructions) {
	for i := 0; i < another.Len(); i++ {
		c.currentInstructions()[instructionBeginIndex.add(i)] = another[i]
	}
}

func (c *Compiler) compileOperatorPlus() error {
	c.emit(code.OpAdd)
	return nil
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.currentScope().instructions
}

func (c *Compiler) updateCurrentInstructions(instruction code.Instructions) {
	instructions := c.currentInstructions()
	instructions = append(instructions, instruction...)
	c.currentScope().instructions = instructions
}

func (c *Compiler) enterScope() {
	scope := NewCompilationScope()
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) exitScope() code.Instructions {
	instructions := c.currentInstructions()
	c.scopes = c.scopes[:c.scopeIndex]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return instructions
}

func (c *Compiler) completeOpReturn(literal *ast.FnLiteral) {
	if c.isLastInstructionMatch(code.OpReturnValue) {
		return
	}

	if len(literal.Body.Statements) == 0 {
		c.emit(code.OpReturn)
		return
	}

	if c.isLastInstructionMatch(code.OpPop) {
		c.updateLastPopInstruction(code.OpReturnValue)
		return
	}

	c.emit(code.OpReturnValue)
}

func (c *Compiler) genClosure() error {
	subSymbolTable := c.symbolTable
	fnInstructions := c.exitScope()

	fnCompiled := &code.CompiledFunction{
		Instructions:   fnInstructions,
		NumOfLocalVars: subSymbolTable.numDefinitions,
	}

	// the closure is inside another function
	closure := &code.Closure{
		Fn: fnCompiled,
		// TODO add num of free variables
		Free: make([]object.Object, subSymbolTable.numFree),
	}

	for _, s := range subSymbolTable.Free {
		c.emitGetScope(s)
	}

	index := c.constants.AddConstant(closure).IntValue()
	// TODO add number of free variables
	c.emit(code.OpClosure, index, subSymbolTable.numFree)

	return nil
}

func (c *Compiler) emitGetScope(symbol Symbol) {
	switch symbol.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, symbol.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, symbol.Index)
	case BuiltInScope:
		c.emit(code.OpGetBuiltIn, symbol.Index)
	case FreeScope:
		c.emit(code.OpGetFree, symbol.Index)
	default:
		panic(common.NewUnknownScope(symbol.Name))
	}
}

func (c *Compiler) emitSetScope(symbol Symbol) {
	switch symbol.Scope {
	case GlobalScope:
		c.emit(code.OpSetGlobal, symbol.Index)
	case LocalScope:
		c.emit(code.OpSetLocal, symbol.Index)
	case BuiltInScope:
		c.emit(code.OpSetBuiltIn, symbol.Index)
	default:
		panic(common.NewUnknownScope(symbol.Name))
	}
}
