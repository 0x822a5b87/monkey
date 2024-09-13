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
				code.Make(code.OpPop),
			},
		},

		{
			input: "100; 200;",
			expectedConstants: []any{
				100,
				200,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
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
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-2147483648", -2147483648},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, testCases)
}

func TestBooleanArithmetic(t *testing.T) {
	testCases := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"true == false", false},
		{"true == true", true},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(1 < 2) != true", false},
		{"(1 > 2) != true", true},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!-2147483648", true},
	}

	runVmTests(t, testCases)
}

func TestConditionals(t *testing.T) {
	testCases := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (!true) { 10 } else { 20 }", 20},
		{"if (!true) { 10 } else { 20 } 30", 30},
		{"if (1) { 10 } else { 20 } 10", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (false) { 10; }", object.NativeFalse},
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

	topElement := vm.TestOnlyLastPoppedStackElement()

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
	case bool:
		testBooleanObject(t, caseIndex, expected, actual)
	case *object.Null:
		if actual != object.NativeNull {
			t.Errorf("object is not NativeNull: expected [%T] actual [%+v]", expected, actual)
		}
	case *object.Boolean:
		if actual != expected {
			t.Errorf("object not match : expected [%T] actual [%+v]", expected, actual)
		}
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

func testBooleanObject(t *testing.T, caseIndex int, expected bool, actual object.Object) {
	t.Helper()
	b, ok := actual.(*object.Boolean)
	if !ok {
		t.Fatalf("test case [%d] object is not Boolean. got=%T (%+v)", caseIndex, actual, actual)
	}

	if expected != b.Value {
		t.Fatalf("test case [%d] object has wrong value. expected = [%t], got = [%t]", caseIndex, expected, b.Value)
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
