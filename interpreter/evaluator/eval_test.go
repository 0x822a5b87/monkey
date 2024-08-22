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
