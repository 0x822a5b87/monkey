package parser

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/common"
	"0x822a5b87/monkey/lexer"
	"0x822a5b87/monkey/token"
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
}

func NewParser(l lexer.Lexer) *Parser {
	p := &Parser{
		lex:            l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
		precedences:    make(map[token.TokenType]Precedence),
	}

	p.precedences[token.INT] = LowestPrecedence
	p.precedences[token.BANG] = PrefixPrecedence
	p.precedences[token.SUB] = PrefixPrecedence

	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseInteger)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)

	// call next token twice so that current token and peek token are both set
	p.nextToken()
	p.nextToken()

	return p
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
	// TODO support more statement parser
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENTIFIER) {
		panic(common.ErrSyntax)
	}

	letStmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		panic(common.ErrSyntax)
	}

	// TODO parse expression
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()
	// TODO parse expression as return statement's value
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return returnStatement
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
	// start parse expression from prefix parse function
	prefixFn := p.getPrefixFn(p.currToken.Type)
	lhs := prefixFn()
	for !p.currTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
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

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) getPrecedence(tokenType token.TokenType) Precedence {
	return p.precedences[tokenType]
}

func (p *Parser) peekPrecedence() Precedence {
	return p.getPrecedence(p.nextToken().Type)
}

// expectPeek step to next token if peek token type matches given token type
func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parsePrefixExpression() ast.Expression {
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

func (p *Parser) parseInteger() ast.Expression {
	integerLiteral := p.currToken.Literal
	integer, err := strconv.ParseInt(integerLiteral, 10, 64)
	if err != nil {
		panic(err)
	}
	return &ast.IntegerLiteral{Token: p.currToken, Value: integer}
}
