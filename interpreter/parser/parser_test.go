package parser

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/lexer"
	"testing"
)

func TestParseProgram(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo_bar = 510;
`
	l := lexer.NewLexer(input)
	p := NewParser(*l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatal("ParseProgram return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}

	expectedStatements := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo_bar"},
	}

	for i, expectedStatement := range expectedStatements {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, expectedStatement.expectedIdentifier) {
			t.Fatalf("line [%d], expected [%s], got [%s]", i, expectedStatement.expectedIdentifier, stmt.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("letStmt.TokenLiteral() not 'let', got [%s]", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement, got = %T", stmt)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value != name, expected [%s], got [%s]", name, letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() != name, expected [%s], got [%s]", name, letStmt.Name.TokenLiteral())
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.NewLexer(input)
	p := NewParser(*l)

	program := p.ParseProgram()

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got [%d]", len(program.Statements))
	}

	for i, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("program.Statements[%d] is not a return statement", i)
			continue
		}
		if returnStatement.Token.Literal != "return" {
			t.Errorf("program.Statements[%d]'s literal value is not 'return', got '%s'", i, returnStatement.Token.Literal)
		}
	}
}

func TestExpression_SingleIdentifier(t *testing.T) {
	input := `foo_bar`
	program := parseProgram(input)
	desc := "identifier"
	checkProgramSize(t, program, desc, 1, 0)
	expr := checkStatementTypeIsExpressionStatement(t, program, desc, 0)

	identifier, ok := expr.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is expected to be a Identfier yet not")
	}

	if identifier.Value != input {
		t.Fatalf("identifier's value is expected to be [%s] yet [%s]", input, identifier.Value)
	}

	if identifier.TokenLiteral() != input {
		t.Fatalf("identifier's token literal is expected to be [%s] yet [%s]", input, identifier.TokenLiteral())
	}
}

func TestExpression_Integer(t *testing.T) {
	input := `5;`
	program := parseProgram(input)

	desc := "integer"
	checkProgramSize(t, program, desc, 1, 0)
	expr := checkStatementTypeIsExpressionStatement(t, program, desc, 0)

	identifier, ok := expr.Expr.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression is expected to be a IntegerLigeral yet not")
	}

	if identifier.Value != 5 {
		t.Fatalf("integer's value is expected to be [%d] yet [%d]", 5, identifier.Value)
	}

	if identifier.TokenLiteral() != "5" {
		t.Fatalf("integer's token literal is expected to be [%s] yet [%s]", "5", identifier.TokenLiteral())
	}
}

func TestExpression_PrefixOperator(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"-5;", "-", 5},
		{"!10;", "!", 10},
	}

	for offset, prefixTestCase := range prefixTests {
		program := parseProgram(prefixTestCase.input)
		desc := "prefix operator"
		checkProgramSize(t, program, desc, 1, offset)
		exprStmt := checkStatementTypeIsExpressionStatement(t, program, desc, offset)
		expr, ok := exprStmt.Expr.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("%s[%d] program.Statments[0] is expected to be a PrefixExpression yet not", desc, offset)
		}
		if expr.Operator != prefixTestCase.operator {
			t.Fatalf("%s offset = [%d], operator expected [%s], got [%s]", desc, offset, prefixTestCase.operator, expr.Operator)
		}

		literal, ok := expr.Right.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("%s offset = [%d], right expression expected IntegerLiteral, got [%s]", desc, offset, expr.Right.TokenLiteral())
		}

		if literal.Value != prefixTestCase.integerValue {
			t.Fatalf("%s offset = [%d], right value expected [%d], got [%d]", desc, offset, prefixTestCase.integerValue, literal.Value)
		}
	}
}

func checkStatementTypeIsExpressionStatement(t *testing.T, program *ast.Program, desc string, offset int) *ast.ExpressionStatement {
	expr, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("%s[%d] program.Statments[0] is expected to be a ExpressionStatement yet not", desc, offset)
	}
	return expr
}

func checkProgramSize(t *testing.T, program *ast.Program, desc string, expectedSize, offset int) {
	if len(program.Statements) != expectedSize {
		t.Fatalf("%s[%d] does not have enough statements, expected [%d], got [%d]", desc, offset, expectedSize, len(program.Statements))
	}
}

func parseProgram(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := NewParser(*l)
	return p.ParseProgram()
}
