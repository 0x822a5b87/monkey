package token

// TokenType define tye type of the token
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

// system info
const (
	ILLEGAL TokenType = "ILLEGAL" // signifies a token/character we don't know
	EOF     TokenType = "EOF"     // end of file
)

// identifiers and literals
const (
	IDENTIFIER TokenType = "IDENTIFIER" // identifier
	INT        TokenType = "INT"        // int
	String     TokenType = "STRING"     // string
)

// operators
const (
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	BANG     TokenType = "!"
	SUB      TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	GT       TokenType = ">"
	LT       TokenType = "<"
	EQ       TokenType = "=="
	NotEq    TokenType = "!="
)

// delimiters
const (
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
)

// preserved keywords
const (
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
	RETURN   TokenType = "RETURN"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
)

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
