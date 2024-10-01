package parser

import (
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/common"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/token"
	"fmt"
	"strconv"
)

type Precedence int

const (
	LowestPrecedence      Precedence = 10
	EqualsPrecedence      Precedence = 20
	LessGreaterPrecedence Precedence = 30
	SumPrecedence         Precedence = 40
	ProductPrecedence     Precedence = 50
	PrefixPrecedence      Precedence = 60
	CallPrecedence        Precedence = 70
)

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression // infixParseFn the argument is "left side" of the infix operator which being parsed

type Parser struct {
	lex       lexer.Lexer
	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	precedences map[token.TokenType]Precedence

	tracing bool
}

func NewParser(l lexer.Lexer) *Parser {
	p := &Parser{
		lex:            l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
		precedences:    make(map[token.TokenType]Precedence),
	}

	p.precedences[token.INT] = LowestPrecedence
	p.precedences[token.IDENTIFIER] = LowestPrecedence
	p.precedences[token.COMMA] = LowestPrecedence
	p.precedences[token.BANG] = PrefixPrecedence
	p.precedences[token.SUB] = SumPrecedence
	p.precedences[token.PLUS] = SumPrecedence
	p.precedences[token.ASTERISK] = ProductPrecedence
	p.precedences[token.SLASH] = ProductPrecedence
	p.precedences[token.GT] = LessGreaterPrecedence
	p.precedences[token.LT] = LessGreaterPrecedence
	p.precedences[token.EQ] = EqualsPrecedence
	p.precedences[token.NotEq] = EqualsPrecedence
	p.precedences[token.TRUE] = LowestPrecedence
	p.precedences[token.FALSE] = LowestPrecedence
	p.precedences[token.LPAREN] = CallPrecedence
	p.precedences[token.RPAREN] = LowestPrecedence
	p.precedences[token.COLON] = LowestPrecedence

	p.precedences[token.LBRACE] = CallPrecedence
	p.precedences[token.RBRACE] = LowestPrecedence

	p.precedences[token.LBRACKET] = CallPrecedence
	p.precedences[token.RBRACKET] = LowestPrecedence

	p.precedences[token.RETURN] = LowestPrecedence

	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroup)
	p.registerPrefix(token.LBRACKET, p.parseArray)
	p.registerPrefix(token.LBRACE, p.parseMap)
	p.registerPrefix(token.IF, p.parseIfStmt)
	p.registerPrefix(token.FUNCTION, p.parseFn)
	p.registerPrefix(token.String, p.parseStringLiteral)

	p.registerInfix(token.PLUS, p.parseInfixOperator)
	p.registerInfix(token.SUB, p.parseInfixOperator)
	p.registerInfix(token.ASTERISK, p.parseInfixOperator)
	p.registerInfix(token.SLASH, p.parseInfixOperator)
	p.registerInfix(token.GT, p.parseInfixOperator)
	p.registerInfix(token.LT, p.parseInfixOperator)
	p.registerInfix(token.EQ, p.parseInfixOperator)
	p.registerInfix(token.NotEq, p.parseInfixOperator)
	p.registerInfix(token.LPAREN, p.parseCall)
	p.registerInfix(token.LBRACKET, p.parseIndex)

	// call next token twice so that current token and peek token are both set
	p.nextToken()
	p.nextToken()

	return p
}

func NewParserWithTracing(l lexer.Lexer) *Parser {
	parser := NewParser(l)
	parser.tracing = true

	return parser
}

// nextToken return next token
func (p *Parser) nextToken() token.Token {
	tk, err := p.lex.NextToken()
	if err != nil {
		panic(err)
	}

	p.currToken = p.peekToken
	p.peekToken = tk

	return p.currToken
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// each time we parse a statement, we DO NOT skip the last token of the statement, normally is a semicolon
		// instead of, we skip the last token in current for loop due to avoid infinite loop
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENTIFIER:
		return p.parseAssignStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{
		Token:      p.currToken,
		Statements: make([]ast.Statement, 0),
	}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		stmt := p.parseStatement()
		blockStatement.Statements = append(blockStatement.Statements, stmt)
	}
	p.expectPeek(token.RBRACE)
	return blockStatement
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENTIFIER) {
		panic(common.ErrSyntax)
	}

	letStmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	p.expectPeek(token.ASSIGN)
	p.nextToken()

	letStmt.Value = p.parseExpression(LowestPrecedence)
	p.expectPeek(token.SEMICOLON)

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()
	returnStatement.ReturnValue = p.parseExpression(LowestPrecedence)
	p.expectPeek(token.SEMICOLON)
	return returnStatement
}

func (p *Parser) parseAssignStatement() ast.Statement {
	if p.peekTokenIs(token.ASSIGN) {
		letStmt := &ast.LetStatement{Token: p.currToken}

		letStmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

		p.expectPeek(token.ASSIGN)
		p.nextToken()

		letStmt.Value = p.parseExpression(LowestPrecedence)
		p.expectPeek(token.SEMICOLON)

		return letStmt
	} else {
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expr = p.parseExpression(LowestPrecedence)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	if p.tracing {
		defer untrace(trace(fmt.Sprintf("parseExpression : token [%s]", p.currToken.Literal)))
	}
	// start parse expression from prefix parse function
	prefixFn := p.getPrefixFn(p.currToken.Type)
	lhs := prefixFn()
	for !p.isEof() && precedence < p.peekPrecedence() {
		p.nextToken()
		infixFn := p.getInfixFn(p.currToken.Type)
		lhs = infixFn(lhs)
	}
	return lhs
}

func (p *Parser) getPrefixFn(tokenType token.TokenType) prefixParseFn {
	fn, ok := p.prefixParseFns[tokenType]
	if !ok {
		panic(fmt.Errorf("prefix parse fn not found for type [%s]", tokenType))
	}
	return fn
}

func (p *Parser) getInfixFn(tokenType token.TokenType) infixParseFn {
	fn, ok := p.infixParseFns[tokenType]
	if !ok {
		panic(fmt.Errorf("infix parse fn not found for type [%s]", tokenType))
	}
	return fn
}

func (p *Parser) currTokenIs(tokenType token.TokenType) bool {
	return p.currToken.Type == tokenType
}

func (p *Parser) isEof() bool {
	return p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.EOF)
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) getPrecedence(tokenType token.TokenType) Precedence {
	precedence, ok := p.precedences[tokenType]
	if !ok {
		info := p.lex.CurInfo()
		panic(fmt.Errorf("precedence not found for type [%s], row = [%d], col = [%d]", tokenType, info.RowNum, info.ColNum))
	}
	return precedence
}

func (p *Parser) peekPrecedence() Precedence {
	return p.getPrecedence(p.peekToken.Type)
}

func (p *Parser) expect(tokenType token.TokenType) bool {
	if p.currTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	panic(fmt.Errorf("expect current [%s], got [%s]", tokenType, p.currToken.Type))
}

// expectPeek step to next token if peek token type matches given token type
func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	panic(fmt.Errorf("expected [%s], got [%s]", tokenType, p.peekToken.Type))
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.tracing {
		defer untrace(trace(fmt.Sprintf("parsePrefixExpression : token [%s]", p.currToken.Literal)))
	}
	expr := ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PrefixPrecedence)
	return &expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	b, err := strconv.ParseBool(p.currToken.Literal)
	if err != nil {
		panic(err)
	}
	return &ast.BooleanExpression{Token: p.currToken, Value: b}
}

func (p *Parser) parseInteger() ast.Expression {
	integerLiteral := p.currToken.Literal
	integer, err := strconv.ParseInt(integerLiteral, 10, 64)
	if err != nil {
		panic(err)
	}
	return &ast.IntegerLiteral{Token: p.currToken, Value: integer}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token:   p.currToken,
		Literal: p.currToken.Literal,
	}
}

func (p *Parser) parseGroup() ast.Expression {
	// skip left parentheses
	p.expect(token.LPAREN)
	groupExpr := p.parseExpression(LowestPrecedence)
	p.expectPeek(token.RPAREN)
	return groupExpr
}

func (p *Parser) parseArray() ast.Expression {
	arrayLiteral := &ast.ArrayLiteral{Token: p.currToken}
	arrayLiteral.Elements = p.parseExpressionList(token.RBRACKET)
	p.expectPeek(token.RBRACKET)
	return arrayLiteral
}

//func (p *Parser) parseMap() ast.Expression {
//	m := &ast.HashExpression{
//		Token: p.currToken,
//		Pairs:   make(map[ast.Expression]ast.Expression),
//	}
//
//	p.expect(token.LBRACE)
//	for !p.currTokenIs(token.RBRACE) {
//		k := p.parseExpression(LowestPrecedence)
//		p.expect(token.COLON)
//		v := p.parseExpression(LowestPrecedence)
//		m.Pairs[k] = v
//	}
//
//	return m
//}

func (p *Parser) parseMap() ast.Expression {
	m := &ast.HashExpression{
		Token: p.currToken,
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		k := p.parseExpression(LowestPrecedence)
		p.expectPeek(token.COLON)
		p.nextToken()
		v := p.parseExpression(LowestPrecedence)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
		m.Pairs[k] = v
	}
	p.expectPeek(token.RBRACE)

	return m
}

func (p *Parser) parseIndex(lhs ast.Expression) ast.Expression {
	indexExpression := &ast.IndexExpression{
		Token: p.currToken,
		Lhs:   lhs,
	}
	p.expect(token.LBRACKET)
	indexExpression.Index = p.parseExpression(LowestPrecedence)
	p.expectPeek(token.RBRACKET)
	return indexExpression
}

func (p *Parser) parseExpressionList(terminalTokenType token.TokenType) []ast.Expression {
	expressions := make([]ast.Expression, 0)
	for !p.peekTokenIs(terminalTokenType) {
		// skip the left bracket if this is first element; otherwise, skip the comma
		p.nextToken()
		expr := p.parseExpression(LowestPrecedence)
		expressions = append(expressions, expr)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	return expressions
}

func (p *Parser) parseIfStmt() ast.Expression {
	ifStmt := &ast.IfExpression{
		Token: p.currToken,
	}

	p.expect(token.IF)
	ifStmt.Condition = p.parseGroup()
	p.expect(token.RPAREN)

	ifStmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		// parse else statement
		p.expectPeek(token.ELSE)
		p.nextToken()
		ifStmt.Alternative = p.parseBlockStatement()
	}

	return ifStmt
}

func (p *Parser) parseInfixOperator(lhs ast.Expression) ast.Expression {
	if p.tracing {
		defer untrace(trace(fmt.Sprintf("parseInfixOperator : token [%s]", p.currToken.Literal)))
	}
	expr := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Lhs:      lhs,
	}

	precedence := p.getPrecedence(p.currToken.Type)
	p.nextToken()

	expr.Rhs = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseCall(lhs ast.Expression) ast.Expression {
	call := &ast.CallExpression{
		Token: p.currToken,
		Fn:    lhs,
	}
	call.Arguments = p.parseExpressionList(token.RPAREN)
	p.expectPeek(token.RPAREN)

	return call
}

func (p *Parser) parseFn() ast.Expression {
	fn := &ast.FnLiteral{
		Token:      p.currToken,
		Parameters: make([]*ast.Identifier, 0),
	}

	// parse parameters
	p.expect(token.FUNCTION)
	for !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		// Identifier inherits from Expression, so we can't convert an Expression to an Identifier.
		// Therefor, we can't simply use p.parseIdentifier
		identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		fn.Parameters = append(fn.Parameters, identifier)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.expectPeek(token.RPAREN)
	p.expectPeek(token.LBRACE)
	fn.Body = p.parseBlockStatement()

	return fn
}
