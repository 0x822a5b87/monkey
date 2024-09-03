package evaluator

import (
	"0x822a5b87/monkey/lexer"
	"0x822a5b87/monkey/object"
	"0x822a5b87/monkey/parser"
	"reflect"
	"testing"
)

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			input:    `5`,
			expected: 5,
		},
		{
			input:    `10`,
			expected: 10,
		},
		{
			input:    `123456789`,
			expected: 123456789,
		},
	}

	for i, val := range tests {
		obj := testEval(val.input)
		testIntegerObject(t, i, obj, val.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1", 1},
		{"1 + 1", 2},
		{"123 - 100", 23},
		{"328 * 21", 6888},
		{"20 / 2", 10},
		{"1 + 2 * 3 / 2 + 4", 8},
		{"(1 + 2) * 6 / (2 * (4 + 5))", 1},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, i, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, i, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
		//		{
		//			`
		//let f = fn(x) {
		//  return x;
		//  x + 10;
		//};
		//f(10);`,
		//			10,
		//		},
		//		{
		//			`
		//let f = fn(x) {
		//   let result = x + 10;
		//   return result;
		//   return 10;
		//};
		//f(10);`,
		//			20,
		//		},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, i, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
		if (10 > 1) {
		 if (10 > 1) {
		   return true + false;
		 }
		
		 return 1;
		}
		`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		//{
		//	`{"name": "Monkey"}[fn(x) { x }];`,
		//	"unusable as hash key: FUNCTION",
		//},
		//{
		//	`999[1]`,
		//	"index operator not supported: INTEGER",
		//},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v), input = [%s]",
				evaluated, evaluated, tt.input)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q, input = [%s]",
				tt.expectedMessage, errObj.Message, tt.input)
		}
	}
}

func TestLetStatement(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c= a + b + 10; c;", 20},
	}

	for i, tc := range testCases {
		eval := testEval(tc.input)
		testIntegerObject(t, i, eval, tc.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Fn)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Params)
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Params[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for i, tt := range tests {
		testIntegerObject(t, i, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
let first = 10;
let second = 10;
let third = 10;

let ourFunction = fn(first) {
  let second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, 0, testEval(input), 70)
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let x = 10;
let addTwo = newAdder(2);
addTwo(2);
`
	testIntegerObject(t, 0, testEval(input), 4)
}

func TestClosures2(t *testing.T) {
	input := `
let self = fn(x) {
	return x;
};
let x = 10;
self(x);
`
	testIntegerObject(t, 0, testEval(input), 10)
}

func TestClosures3(t *testing.T) {
	input := `
let self = fn(x) {
	return x;
}(10);
self;
`
	testIntegerObject(t, 0, testEval(input), 10)
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.StringObj)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltInFunctions(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{`len("")`, &object.Integer{Value: 0}},
		{`len("four")`, &object.Integer{Value: 4}},
		{`len("hello world")`, &object.Integer{Value: 11}},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for i, testCase := range testCases {
		evaluated := testEval(testCase.input)
		switch expected := testCase.expected.(type) {
		case *object.Integer:
			v := testCase.expected.(*object.Integer)
			testIntegerObject(t, i, evaluated, v.Value)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("index = [%d], object is not Error. got=%T (%+v)", i, evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Errorf("index = [%d], wrong error message. expected=%q, got=%q", i, expected, errObj.Message)
			}
		default:
			t.Errorf("index = [%d], unexpected return type. expected=%q, got=%q", i, expected, evaluated.Type())
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, 0, result.Elements[0], 1)
	testIntegerObject(t, 0, result.Elements[1], 4)
	testIntegerObject(t, 0, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
		{
			`"hello world!"[0]`,
			"h",
		},
		{
			`"hello world!"[11]`,
			"!",
		},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, i, evaluated, int64(expected))
		case string:
			testStringObj(t, i, evaluated, expected)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"first([1, 2, 3])",
			1,
		},
		{
			"last([1, 2, 3])",
			3,
		},
		{
			`first("hello world!")`,
			"h",
		},
		{
			`last("hello world!")`,
			"!",
		},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, i, evaluated, int64(expected))
		case string:
			testStringObj(t, i, evaluated, expected)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func testEval(input string) object.Object {
	newLexer := lexer.NewLexer(input)
	newParser := parser.NewParser(*newLexer)
	program := newParser.ParseProgram()
	env := object.NewEnvironment(nil)
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, testCaseIndex int, obj object.Object, expected int64) {
	if obj == nil {
		t.Fatalf("test case [%d], exepct not nil, got nil", testCaseIndex)
	}

	if obj.Type() != object.ObjInteger {
		t.Fatalf("test case [%d], expect ObjInteger, got [%s], msg [%s]", testCaseIndex, string(obj.Type()), obj.Inspect())
	}
	integerObj, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("test case [%d], expecte Integer, got [%s]", testCaseIndex, reflect.TypeOf(obj))
	}
	if integerObj.Value != expected {
		t.Fatalf("test case [%d], expect [%d], got [%d]", testCaseIndex, expected, integerObj.Value)
	}
}

func testStringObj(t *testing.T, testCaseIndex int, obj object.Object, expected string) {
	if obj == nil {
		t.Fatalf("test case [%d], exepct not nil, got nil", testCaseIndex)
	}

	if obj.Type() != object.ObjString {
		t.Fatalf("test case [%d], expect ObjInteger, got [%s]", testCaseIndex, string(obj.Type()))
	}

	str, ok := obj.(*object.StringObj)
	if !ok {
		t.Fatalf("test case [%d], expecte Integer, got [%s]", testCaseIndex, reflect.TypeOf(obj))
	}

	if str.Value != expected {
		t.Fatalf("test case [%d], expect [%s], got [%s]", testCaseIndex, expected, str.Value)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NativeNull {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
