package lexer

import (
	"0x822a5b87/monkey/common"
	"0x822a5b87/monkey/token"
)

const LiteralEof byte = 0

type Lexer struct {
	sourceCode   string
	position     int  // current position in input(points to current char)
	readPosition int  // current reading position in input(after current char)
	ch           byte // current char under examination
}

func NewLexer(source string) *Lexer {
	l := &Lexer{sourceCode: source}
	// init lexer
	l.readChar()
	return l
}

// NextToken get next token from source code
// token.EOF is a const untyped string, and it can be used as token.TokenType due to implicit type conversion
// for further info please read my blog:
// https://0x822a5b87.github.io/2024/07/26/%E5%85%B3%E4%BA%8Egolang%E7%9A%84%E7%B1%BB%E5%9E%8B%E6%8E%A8%E5%AF%BC%E3%80%81%E9%9A%90%E5%BC%8F%E7%B1%BB%E5%9E%8B%E8%BD%AC%E6%8D%A2%E7%9A%84%E4%B8%80%E4%BA%9B%E6%80%9D%E8%80%83/
func (l *Lexer) NextToken() (token.Token, error) {
	// before we parse token, we should skip the whitespace
	l.skipWhitespace()

	var tok token.Token
	var err error
	switch l.ch {
	case LiteralEof:
		tok, err = newToken(token.EOF, l.ch)
	case '=':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			tok, err = newTokenForBinary(token.EQ, ch, l.ch)
		} else {
			tok, err = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			tok, err = newTokenForBinary(token.NotEq, ch, l.ch)
		} else {
			tok, err = newToken(token.BANG, l.ch)
		}
	case '+':
		tok, err = newToken(token.PLUS, l.ch)
	case ',':
		tok, err = newToken(token.COMMA, l.ch)
	case ';':
		tok, err = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok, err = newToken(token.LPAREN, l.ch)
	case ')':
		tok, err = newToken(token.RPAREN, l.ch)
	case '{':
		tok, err = newToken(token.LBRACE, l.ch)
	case '}':
		tok, err = newToken(token.RBRACE, l.ch)
	case '-':
		tok, err = newToken(token.SUB, l.ch)
	case '*':
		tok, err = newToken(token.ASTERISK, l.ch)
	case '/':
		tok, err = newToken(token.SLASH, l.ch)
	case '>':
		tok, err = newToken(token.GT, l.ch)
	case '<':
		tok, err = newToken(token.LT, l.ch)
	default:
		if l.isLetter() {
			tok, err = l.readIdentifier(), nil
		} else if l.isDigit() {
			tok, err = l.readNumber(), nil
		} else {
			tok, err = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}, common.ErrUnknownToken
		}

		// return immediately after parse identifier/number or other specific object
		// because readChar() already called in the function
		return tok, err
	}

	// call readChar() after parse token, because we need to move l.position to next char
	l.readChar()

	return tok, err
}

func (l *Lexer) CurrentPos() int {
	return l.position
}

func (l *Lexer) readIdentifier() token.Token {
	cur := l.position
	for l.isLetter() {
		l.readChar()
	}
	literal := l.sourceCode[cur:l.position]
	identifier := token.LookupIdentifier(literal)
	return token.Token{
		Type:    identifier,
		Literal: literal,
	}
}

func (l *Lexer) readNumber() token.Token {
	cur := l.position
	for l.isDigit() {
		l.readChar()
	}
	return token.Token{
		Type:    token.INT,
		Literal: l.sourceCode[cur:l.position],
	}
}

func (l *Lexer) isLetter() bool {
	return ('a' <= l.ch && l.ch <= 'z') || ('A' <= l.ch && l.ch <= 'Z') || l.ch == '_'
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9'
}

func newTokenForBinary(tokenType token.TokenType, first, second byte) (token.Token, error) {
	return token.Token{
		Type:    tokenType,
		Literal: string(first) + string(second),
	}, nil
}

func newToken(tokenType token.TokenType, ch byte) (token.Token, error) {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}, nil
}

func (l *Lexer) readChar() {
	if !l.hasNextChar() {
		l.ch = LiteralEof
	} else {
		l.ch = l.sourceCode[l.readPosition]
	}

	// should we use l.position++ instead of the code? absolutely not!
	// if we use l.position++, then position and readPosition will be permanently the same
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peakChar() byte {
	if !l.hasNextChar() {
		return LiteralEof
	} else {
		return l.sourceCode[l.readPosition]
	}
}

func (l *Lexer) hasNextChar() bool {
	return l.readPosition < len(l.sourceCode)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// curStr convert current ch to string, mainly for debugging
func (l *Lexer) curStr() string {
	return string(l.ch)
}
