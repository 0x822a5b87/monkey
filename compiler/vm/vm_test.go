package vm

import (
	"0x822a5b87/monkey/compiler/compiler"
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/parser"
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	testCases := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 2},
	}

	runVmTests(t, testCases)
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, testCases []vmTestCase) {
	t.Helper()

	for caseIndex, testCase := range testCases {
		runVmTest(t, testCase, caseIndex)
	}
}

func runVmTest(t *testing.T, testCase vmTestCase, caseIndex int) {
	t.Helper()
	program := parse(testCase.input)
	c := compiler.NewCompiler()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("compile error : [%s] for input = [%s]", err.Error(), testCase.input)
	}

	vm := NewVm(c.ByteCode())
	err = vm.Run()
	if err != nil {
		t.Fatalf("vm error : [%s] for input = [%s]", err.Error(), testCase.input)
	}

	topElement := vm.StackTop()

	testExpectedObject(t, caseIndex, testCase.expected, topElement)
}

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(*l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, caseIndex int, expected interface{}, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, caseIndex, int64(expected), actual)
	default:
		t.Fatalf("test case [%d] wrong type [%s] for test", caseIndex, expected)
	}
}

func testIntegerObject(t *testing.T, caseIndex int, expected int64, actual object.Object) {
	t.Helper()
	integerObj, ok := actual.(*object.Integer)
	if !ok {
		t.Fatalf("test case [%d] object is not Integer. got=%T (%+v)", caseIndex, actual, actual)
	}

	if expected != integerObj.Value {
		t.Fatalf("test case [%d] object has wrong value. expected = [%d], got = [%d]", caseIndex, expected, integerObj.Value)
	}
}
