package lexer

import (
	"0x822a5b87/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
let five = 5;
let ten= 10;
let add = fn(x, y) {
	return x + y;
};
let result = add(five, ten);

let sub = fn(x, y) {
	return x - y;
};

let cal = fn(x, y, z) {
	return x*y/z;	
};

five < ten;
ten > five;

if (five < ten) {
	return true;
} else {
	return false;
}

return five != 10;
return five == five;

"foobar";
"foo bar";

let y = "this is a sentence for testing! I want to say \"Hello World\"!";

["1", 2];

{"name": "0x822a5b87", "age": 30, true: "boolean", 99: "integer"};
`

	expectedTokens := []expectedToken{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "sub"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "x"},
		{token.SUB, "-"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "cal"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "z"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "x"},
		{token.ASTERISK, "*"},
		{token.IDENTIFIER, "y"},
		{token.SLASH, "/"},
		{token.IDENTIFIER, "z"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.IDENTIFIER, "five"},
		{token.LT, "<"},
		{token.IDENTIFIER, "ten"},
		{token.SEMICOLON, ";"},

		{token.IDENTIFIER, "ten"},
		{token.GT, ">"},
		{token.IDENTIFIER, "five"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.LT, "<"},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.IDENTIFIER, "five"},
		{token.NotEq, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.RETURN, "return"},
		{token.IDENTIFIER, "five"},
		{token.EQ, "=="},
		{token.IDENTIFIER, "five"},
		{token.SEMICOLON, ";"},

		{token.String, "foobar"},
		{token.SEMICOLON, ";"},

		{token.String, "foo bar"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "y"},
		{token.ASSIGN, "="},
		{token.String, `this is a sentence for testing! I want to say "Hello World"!`},
		{token.SEMICOLON, ";"},

		{token.LBRACKET, "["},
		{token.String, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.String, "name"},
		{token.COLON, ":"},
		{token.String, "0x822a5b87"},
		{token.COMMA, ","},
		{token.String, "age"},
		{token.COLON, ":"},
		{token.INT, "30"},
		{token.COMMA, ","},
		{token.TRUE, "true"},
		{token.COLON, ":"},
		{token.String, "boolean"},
		{token.COMMA, ","},
		{token.INT, "99"},
		{token.COLON, ":"},
		{token.String, "integer"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
	}

	expectedTokens = append(expectedTokens, expectedToken{
		expectedType:    token.EOF,
		expectedLiteral: string(LiteralEof),
	},
	)

	l := NewLexer(input)

	for i, expectedToken := range expectedTokens {
		nextToken, err := l.NextToken()
		if err != nil {
			t.Fatalf("tests[%d] - error get token, error = [%s], ch = [%s], token = [%q]", i, err.Error(), string(l.ch), nextToken)
		}

		if nextToken.Type != expectedToken.expectedType {
			t.Fatalf("tests[%d] - type wrong, expected = %q, got = %q", i, expectedToken.expectedType, nextToken.Type)
		}

		if nextToken.Literal != expectedToken.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected = %q, got = %q", i, expectedToken.expectedLiteral, nextToken.Literal)
		}
	}
}

type expectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}
