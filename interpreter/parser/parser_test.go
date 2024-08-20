package parser

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/lexer"
	"fmt"
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
	input := `foo_bar;`
	expected := "foo_bar"
	program := parseProgram(input)
	desc := "identifier"
	checkProgramSize(t, program, desc, 1, 0)
	expr := checkStatementTypeIsExpressionStatement(t, program, desc, 0)

	identifier, ok := expr.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is expected to be a Identfier yet not")
	}

	if identifier.Value != expected {
		t.Fatalf("identifier's value is expected to be [%s] yet [%s]", expected, identifier.Value)
	}

	if identifier.TokenLiteral() != expected {
		t.Fatalf("identifier's token literal is expected to be [%s] yet [%s]", expected, identifier.TokenLiteral())
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

func TestExpression_InfixOperator(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
	}

	for _, infixTest := range infixTests {
		program := parseProgram(infixTest.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements [%s] does not contain %d statements. got=%d\n",
				infixTest.input, 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, infixTest.input, stmt.Expr, infixTest.leftValue, infixTest.operator, infixTest.rightValue) {
			return
		}
	}

}

func TestExpression_ComplexExpression(t *testing.T) {
	complexExpressions := []struct {
		input          string
		expectedOutput string
	}{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		//{
		//	"a * [1, 2, 3, 4][b * c] * d",
		//	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		//},
		//{
		//	"add(a * b[2], b[1], 2 * [1, 2][1])",
		//	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		//},
	}

	for i, complexExpr := range complexExpressions {
		program := parseProgram(complexExpr.input)
		if complexExpr.expectedOutput != program.String() {
			t.Fatalf("offset = [%d], expected [%s], got [%s]", i, complexExpr.expectedOutput, program.String())
		}
	}
}

func testInfixExpression(t *testing.T, input string, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp for input [%s] is not ast.OperatorExpression. got=%T(%s)", input, exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Lhs, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Rhs, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value,
			integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
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
