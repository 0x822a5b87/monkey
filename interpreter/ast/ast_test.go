package ast

import (
	"0x822a5b87/monkey/interpreter/token"
	"testing"
)

func TestLetStatementString(t *testing.T) {
	program := Program{Statements: []Statement{
		&LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
				Value: "myVar",
			},
			Value: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
				Value: "anotherVar",
			},
		},
	}}

	expectedString := `let myVar = anotherVar;`

	if program.String() != expectedString {
		t.Fatalf("expected [%s], got [%s]", expectedString, program.String())
	}
}
