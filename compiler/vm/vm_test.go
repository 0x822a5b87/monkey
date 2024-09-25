package vm

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/compiler/compiler"
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/parser"
	"reflect"
	"testing"
)

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

func TestGlobalLetStatement(t *testing.T) {
	testCases := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one", 1},
		{"let one = 1; let two = 2; two", 2},
		{"let one = 65536; let two = one; two", 65536},
		{"let one = 65536; let two = one; one + two", 65536 * 2},
	}

	runVmTests(t, testCases)
}

func TestStringExpressions(t *testing.T) {
	testCases := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + " banana"`, "monkey banana"},
	}

	runVmTests(t, testCases)
}

func TestArrayLiterals(t *testing.T) {
	testCases := []vmTestCase{
		{`[]`, []int{}},
		{`[1, 2, 3]`, []int{1, 2, 3}},
		{`[1 + 2, 3 - 4, 5 * 6, 8 / 2]`, []int{3, -1, 30, 4}},
	}

	runVmTests(t, testCases)
}

func TestHashLiterals(t *testing.T) {
	testCases := []vmTestCase{
		{
			input: `{}`,
			expected: &object.Hash{
				Pairs: map[object.HashKey]*object.HashPair{},
			},
		},
		{
			input: `{1:2, 2:3}`,
			expected: &object.Hash{
				Pairs: map[object.HashKey]*object.HashPair{
					(&object.Integer{Value: 1}).HashKey(): {
						Key:   &object.Integer{Value: 1},
						Value: &object.Integer{Value: 2},
					},
					(&object.Integer{Value: 2}).HashKey(): {
						Key:   &object.Integer{Value: 2},
						Value: &object.Integer{Value: 3},
					},
				},
			},
		},
		{
			input: `{1+1:2+2, 3+3:4*4}`,
			expected: &object.Hash{
				Pairs: map[object.HashKey]*object.HashPair{
					(&object.Integer{Value: 2}).HashKey(): {
						Key:   &object.Integer{Value: 2},
						Value: &object.Integer{Value: 4},
					},
					(&object.Integer{Value: 6}).HashKey(): {
						Key:   &object.Integer{Value: 6},
						Value: &object.Integer{Value: 16},
					},
				},
			},
		},
		{
			input: `{1+1:2+2, 3+3:4*4, 100 * 100: 200 - 500}`,
			expected: &object.Hash{
				Pairs: map[object.HashKey]*object.HashPair{
					(&object.Integer{Value: 2}).HashKey(): {
						Key:   &object.Integer{Value: 2},
						Value: &object.Integer{Value: 4},
					},
					(&object.Integer{Value: 6}).HashKey(): {
						Key:   &object.Integer{Value: 6},
						Value: &object.Integer{Value: 16},
					},
					(&object.Integer{Value: 10000}).HashKey(): {
						Key:   &object.Integer{Value: 10000},
						Value: &object.Integer{Value: -300},
					},
				},
			},
		},
	}

	runVmTests(t, testCases)
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
	case string:
		testStringObject(t, caseIndex, actual, expected)
	case []int:
		testArrayObject(t, caseIndex, expected, actual)
	case *object.Hash:
		testHashObject(t, caseIndex, expected, actual)
	case *object.Integer:
		testIntegerObject(t, caseIndex, expected.Value, actual)
	default:
		t.Fatalf("test case [%d] wrong type for test, expected [%+v] actual [%s]", caseIndex, expected, actual.Type())
	}
}

func testArrayObject(t *testing.T, caseIndex int, expected any, actual object.Object) {
	v := reflect.ValueOf(expected)
	actualArray, ok := actual.(*object.Array)
	if !ok {
		t.Errorf("object not match : expected [%T] actual [%+v]", expected, actual)
	}
	if len(actualArray.Elements) != v.Len() {
		t.Fatalf("test case [%d] length not match, expected = [%d], actual = [%d]", caseIndex, v.Len(), len(actualArray.Elements))
	}
	for i, element := range actualArray.Elements {
		testExpectedObject(t, caseIndex, v.Index(i).Interface(), element)
	}
}

func testHashObject(t *testing.T, caseIndex int, expected any, actual object.Object) {
	t.Helper()

	actualArray, ok := actual.(*object.Hash)
	if !ok {
		t.Errorf("object not match : expected [%T] actual [%+v]", expected, actual)
	}

	expectedHash, ok := expected.(*object.Hash)
	if !ok {
		t.Errorf("object not match : expected [%T] actual [%+v]", expected, actual)
	}
	if len(actualArray.Pairs) != len(expectedHash.Pairs) {
		t.Fatalf("test case [%d] length not match, expected = [%d], actual = [%d]", caseIndex, len(expectedHash.Pairs), len(actualArray.Pairs))
	}

	for k, element := range actualArray.Pairs {
		pair, ok := expectedHash.Pairs[k]
		if !ok {
			t.Fatalf("test case [%d] does't contains key = [%s]", caseIndex, k.Type)
		}
		testExpectedObject(t, caseIndex, pair.Value, element.Value)
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

func testStringObject(t *testing.T, caseIndex int, obj object.Object, expected string) {
	if obj == nil {
		t.Fatalf("case %d exepct string but got nil", caseIndex)
	}

	if obj.Type() != object.ObjString {
		t.Fatalf("case %d expect ObjString, got [%s], msg [%s]", caseIndex, string(obj.Type()), obj.Inspect())
	}

	stringObj, ok := obj.(*object.StringObj)
	if !ok {
		t.Fatalf("case %d expecte String, got [%s]", caseIndex, reflect.TypeOf(obj))
	}
	if stringObj.Value != expected {
		t.Fatalf("case %d expect [%s], got [%s]", caseIndex, expected, stringObj.Value)
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
