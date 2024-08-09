package parser

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/common"
	"0x822a5b87/monkey/lexer"
	"0x822a5b87/monkey/token"
)

type Parser struct {
	lex       lexer.Lexer
	currToken token.Token
	peekToken token.Token
}

func NewParser(l lexer.Lexer) *Parser {
	p := &Parser{
		lex: l,
	}

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
		return nil
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

func (p *Parser) currTokenIs(tokenType token.TokenType) bool {
	return p.currToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
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
