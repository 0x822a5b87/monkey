package parser

import (
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/token"
	"fmt"
	"reflect"
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
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for i, complexExpr := range complexExpressions {
		program := parseProgram(complexExpr.input)
		if complexExpr.expectedOutput != program.String() {
			t.Fatalf("offset = [%d], expected [%s], got [%s]", i, complexExpr.expectedOutput, program.String())
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
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
	}

	for _, tt := range tests {
		program := parseProgram(tt.input)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

// TestTracing tracing the execution of parseExpression to understand the function's principles
func TestTracing(t *testing.T) {
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := parseProgram(input)
	if len(program.Statements) != 1 {
		t.Errorf("if expression expected [%d] statements, got [%d] statement", 1, len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("if expression expected ExpressionStatement, actually = [%s]", reflect.TypeOf(ifStmt))
	}

	if ifStmt.Token.Type != token.IF {
		t.Errorf("if expression's token expected IF, actually = [%s]", ifStmt.Token.Type)
	}

	ifExpr, ok := ifStmt.Expr.(*ast.IfExpression)
	if !ok {
		t.Errorf("if expression's Expr expected compileIfExpression, actually = [%s]", reflect.TypeOf(ifStmt.Expr))
	}

	if !testInfixExpression(t, "x < y", ifExpr.Condition, "x", "<", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Errorf("Consequence expected 1 Statements, actually [%d]", len(ifExpr.Consequence.Statements))
	}

	consequenceStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Consequence.Statements[0] expected ExpressionStatement, actually [%s]", reflect.TypeOf(ifExpr.Consequence.Statements[0]))
	}

	if testIdentifier(t, consequenceStmt.Expr, "x") {
		return
	}

	if ifExpr.Alternative == nil {
		t.Errorf("Alternative expected not nil")
	}

	if len(ifExpr.Alternative.Statements) != 1 {
		t.Errorf("Alternative expected 1 Statements, actually [%d]", len(ifExpr.Consequence.Statements))
	}

	alternativeStmt, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Alternative.Statements[0] expected ExpressionStatement, actually [%s]", reflect.TypeOf(ifExpr.Consequence.Statements[0]))
	}

	if testIdentifier(t, alternativeStmt.Expr, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := parseProgram(input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expr.(*ast.FnLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expr)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, "expr", bodyStmt.Expr, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := parseProgram(tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expr.(*ast.FnLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		program := parseProgram(tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expr.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expr)
		}

		if !testIdentifier(t, exp.Fn, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d", len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i, arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	program := parseProgram(input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expr.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expr)
	}

	if !testIdentifier(t, exp.Fn, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, "TestCallExpressionParsing", exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, "TestCallExpressionParsing", exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := parseProgram(input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expr)
	}

	if literal.Literal != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Literal)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	program := parseProgram(input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expr)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, input, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, input, array.Elements[2], 3, "+", 3)
}

func TestEmptyArrayLiterals(t *testing.T) {
	input := "[]"
	program := parseProgram(input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expr)
	}

	if len(array.Elements) != 0 {
		t.Fatalf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := `myArray[1 + 1]`
	program := parseProgram(input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expr.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expr)
	}

	if !testIdentifier(t, indexExp.Lhs, "myArray") {
		return
	}

	if !testInfixExpression(t, input, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	program := parseProgram(input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expr.(*ast.HashExpression)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	program := parseProgram(input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expr.(*ast.HashExpression)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	program := parseProgram(input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expr.(*ast.HashExpression)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		boolean, ok := key.(*ast.BooleanExpression)
		if !ok {
			t.Errorf("key is not ast.BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`
	program := parseProgram(input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expr.(*ast.HashExpression)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		integer, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[integer.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	program := parseProgram(input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expr.(*ast.HashExpression)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expr)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, input, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, input, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, input, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
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

func tracingParseProgram(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := NewParserWithTracing(*l)
	return p.ParseProgram()
}

func parseProgram(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := NewParser(*l)
	return p.ParseProgram()
}
