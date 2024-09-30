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

func TestIndexExpressions(t *testing.T) {
	testCases := []vmTestCase{
		{`[1][0]`, &object.Integer{Value: 1}},
		{`[1, 2, 3][1]`, &object.Integer{Value: 2}},
		{`[1, 2, 3][0 + 2]`, &object.Integer{Value: 3}},
		{`[[1, 2, 3],[4,5,6]][0][0]`, &object.Integer{Value: 1}},
		{`[[1, 2, 3], [4, 5,6]][0][1]`, &object.Integer{Value: 2}},
		{`[[1, 2, 3], [4, 5,6]][0][2]`, &object.Integer{Value: 3}},
		{`[[1, 2, 3], [4, 5,6]][1][0]`, &object.Integer{Value: 4}},
		{`[[1, 2, 3], [4, 5,6]][1][1]`, &object.Integer{Value: 5}},
		{`[[1, 2, 3], [4, 5,6]][1][2]`, &object.Integer{Value: 6}},
		// TODO It should throw an exception when the index is greater than the length.
		{`[][0]`, &object.Null{}},
		{`[][1]`, &object.Null{}},
		{`{1:2, 3:4}[1]`, &object.Integer{Value: 2}},
		{`{1:2, 3:4}[3]`, &object.Integer{Value: 4}},
		{`{}[3]`, &object.Null{}},
	}

	runVmTests(t, testCases)
}

func TestCallingFunctionsWithoutArgument(t *testing.T) {
	testCases := []vmTestCase{
		{
			input: `
		let fivePlusTen = fn() { 5 + 10; };
		fivePlusTen();
		`,
			expected: &object.Integer{Value: 15},
		},
		{
			input: `
				let one = fn() { 1; };
				let two = fn() { 2; };
				two();
				`,
			expected: &object.Integer{Value: 2},
		},
		{
			input: `
				let one = fn() { 1; };
				let two = fn() { 2; };
				two();
				one();
				`,
			expected: &object.Integer{Value: 1},
		},
		{
			input: `
		let one = fn() { 1; };
		let two = fn() { 2; };
		one() + two();
		`,
			expected: &object.Integer{Value: 3},
		},
		{
			input: `
		let a = fn() { 1; };
		let b = fn() { a() + 2; };
		let c = fn() { b() + 3; };
		c();
		`,
			expected: &object.Integer{Value: 6},
		},
	}

	runVmTests(t, testCases)
}

func TestCallingFunctionWithBindings(t *testing.T) {
	testCases := []vmTestCase{
		{
			input: `
		let one = fn() { let x = 1; return x; };
		one();
		`,
			expected: &object.Integer{Value: 1},
		},
		{
			input: `
		let one = fn() { let one = 1; one; };
		one();
		`,
			expected: &object.Integer{Value: 1},
		},
		{
			input: `
		let oneAndTwo = fn() { let one = 1; let two = 2; one + two; }; oneAndTwo();
		`,
			expected: &object.Integer{Value: 3},
		},
		{
			input: `
		let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
		let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
		oneAndTwo() + threeAndFour();
		`,
			expected: &object.Integer{Value: 10},
		},
		{
			input: `
		let firstFoobar = fn() { let foobar = 50; foobar; };
		let secondFoobar = fn() { let foobar = 100; foobar; };
		firstFoobar() + secondFoobar();
		`,
			expected: &object.Integer{Value: 150},
		},
		{
			input: `
		let globalSeed = 50;
		
		let minusOne = fn() {
		 let num = 1;
		 globalSeed - num;
		};
		
		let minusTwo = fn() {
		 let num = 2;
		 globalSeed - num;
		};
		
		minusOne() + minusTwo();
		`,
			expected: &object.Integer{Value: 97},
		},
		{
			input: `
			let noReturnFn = fn() { };
			noReturnFn();
		`,
			expected: object.NativeNull,
		},
	}

	runVmTests(t, testCases)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	testCases := []vmTestCase{
		{
			input: `
let earlyExit = fn() { return 99; 100; };
earlyExit();
`,
			expected: &object.Integer{Value: 99},
		},
		{
			input: `
let earlyExit = fn() { return 99; return 100; };
earlyExit();
`,
			expected: &object.Integer{Value: 99},
		},
		{
			input: `
let earlyExit = fn() { 99; return 100; };
earlyExit();
`,
			expected: &object.Integer{Value: 100},
		},
	}

	runVmTests(t, testCases)
}

func TestFunctionWithArguments(t *testing.T) {
	testCases := []vmTestCase{
		{
			input: `
				let oneArg = fn(a) { a };
				oneArg(24);
`,
			expected: &object.Integer{Value: 24},
		},
		{
			input: `
let sum = fn(a, b) {
	let c = a + b;
	c;
};
sum(1, 2);
`,
			expected: &object.Integer{Value: 3},
		},
		{
			input: `
let sum = fn(a, b) {
	let c = a + b;
	c;
};
sum(1, 2) + sum(3, 4);
`,
			expected: &object.Integer{Value: 10},
		},
		{
			input: `
let globalNum = 10;

let sum = fn(a, b) {
	let c = a + b;
	c + globalNum;
};

sum(1, 2);
`,
			expected: &object.Integer{Value: 13},
		},
		{
			input: `
let sum = fn(a, b) {
	let c = a + b;
	c;
};

let outer = fn() {
	sum(-100, 15) + sum(95, -10);
};

outer();
`,
			expected: &object.Integer{Value: 0},
		},
	}

	runVmTests(t, testCases)
}

func TestFunctionWithoutReturnValue(t *testing.T) {
	testCases := []vmTestCase{
		{
			input: `
let noReturn = fn() { };
noReturn();
`,
			expected: object.NativeNull,
		},
		{
			input: `
let noReturn = fn() { };
let noReturnTwo = fn() { noReturn(); };
noReturnTwo();
noReturn();
noReturnTwo();
`,
			expected: object.NativeNull,
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
	t.Helper()

	v := reflect.ValueOf(expected)
	actualArray, ok := actual.(*object.Array)
	if !ok {
		t.Errorf("object not match : expected [%T] actual [%+v]", expected, actual)
		return
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
	if vm.currentInstructions().Len() != expectedInstructionLen {
		t.Fatalf("test case [%d] instructions length not match, expected = [%d], actual = [%d]",
			caseIndex, expectedInstructionLen, vm.currentInstructions().Len())
	}
}

type vmTestCase struct {
	input    string
	expected interface{}
}
