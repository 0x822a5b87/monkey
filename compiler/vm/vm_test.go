package vm

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/compiler/compiler"
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/parser"
	"testing"
)

func TestCompilerIntegerArithmetic(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: "1 + 2",
			expectedConstants: []any{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func TestIntegerArithmetic(t *testing.T) {
	testCases := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
	}

	runVmTests(t, testCases)
}

func runCompilerTests(t *testing.T, testCases []compilerTestCase) {
	t.Helper()

	for i, testCase := range testCases {
		vm := runVm(t, i, testCase.input)
		testVm(t, i, vm, testCase.expectedConstants, testCase.expectedInstructions)
	}
}

func runVmTests(t *testing.T, testCases []vmTestCase) {
	t.Helper()

	for caseIndex, testCase := range testCases {
		runVmTest(t, testCase, caseIndex)
	}
}

func runVmTest(t *testing.T, testCase vmTestCase, caseIndex int) {
	t.Helper()
	vm := runVm(t, caseIndex, testCase.input)
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

func runVm(t *testing.T, caseIndex int, input string) *Vm {
	t.Helper()

	program := parse(input)
	c := compiler.NewCompiler()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("test case [%d] compile error : [%s] for input = [%s]", caseIndex, err.Error(), input)
	}

	vm := NewVm(c.ByteCode())
	err = vm.Run()
	if err != nil {
		t.Fatalf("test case [%d] vm error : [%s] for input = [%s]", caseIndex, err.Error(), input)
	}

	return vm
}

func testVm(t *testing.T, caseIndex int, vm *Vm, expectedConstants []any, expectedInstructions []code.Instructions) {
	t.Helper()

	if vm.constants.Len() != len(expectedConstants) {
		t.Fatalf("test case [%d] contants length not match, expected = [%d], actual = [%d]",
			caseIndex, vm.constants.Len(), len(expectedConstants))
	}

	var expectedInstructionLen int
	for _, instruction := range expectedInstructions {
		expectedInstructionLen += instruction.Len()
	}
	if vm.instructions.Len() != expectedInstructionLen {
		t.Fatalf("test case [%d] instructions length not match, expected = [%d], actual = [%d]",
			caseIndex, expectedInstructionLen, vm.instructions.Len())
	}
}

type vmTestCase struct {
	input    string
	expected interface{}
}

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}
