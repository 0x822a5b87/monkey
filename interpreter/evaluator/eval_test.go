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

	for _, val := range tests {
		obj := testEval(val.input)
		testIntegerObject(t, obj, val.expected)
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

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	newLexer := lexer.NewLexer(input)
	newParser := parser.NewParser(*newLexer)
	program := newParser.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	if obj == nil {
		t.Fatalf("exepct not nil, got nil")
	}

	if obj.Type() != object.ObjInteger {
		t.Fatalf("expect ObjInteger, got [%s]", string(obj.Type()))
	}
	integerObj, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("expecte Integer, got [%s]", reflect.TypeOf(obj))
	}
	if integerObj.Value != expected {
		t.Fatalf("expect [%d], got [%d]", expected, integerObj.Value)
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
