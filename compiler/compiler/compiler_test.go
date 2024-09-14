package compiler

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/interpreter/ast"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/parser"
	"reflect"
	"testing"
)

func TestIntegerCompiler(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input:             `1 + 2`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `65535 - 65534`,
			expectedConstants: []interface{}{65535, 65534},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `256 * 0`,
			expectedConstants: []interface{}{256, 0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `512 / 256`,
			expectedConstants: []interface{}{512, 256},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			// the prefix operator minus.
			// The expected behavior is that when we encounter a negative integer, we push the absolute value to the stack.
			// we compile the index of the integer and OpMinus into an instruction, During runtime, the vm will decompile it.
			input:             `-2147483648`,
			expectedConstants: []interface{}{2147483648},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `-1`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestCompilerIntegerArithmetic(t *testing.T) {
	testCases := []*compilerTestCase{
		{
			input: "1 + 2",
			expectedConstants: []any{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
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

	for i, testCase := range testCases {
		runCompilerTest(t, i, testCase)
	}
}

func TestBooleanExpressions(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},

		{
			input: "1 > 2",
			expectedConstants: []any{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},

		{
			input: "100 < 200",
			expectedConstants: []any{
				100,
				200,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpLessThan),
				code.Make(code.OpPop),
			},
		},

		{
			input: "123 == 321",
			expectedConstants: []any{
				123,
				321,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},

		{
			input: "2147483647 != 2147483647",
			expectedConstants: []any{
				2147483647,
				2147483647,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestCondition(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: `
if (true) {
	10;
};
3333;
`,
			expectedConstants: []any{10, 3333},
			// note that OpJumpNotTruthy tell VM jump to 0007, that's because in Compiler#Run()'s for loop,
			// it will increment the instruction pointer by one, so we only need to jump to the byte preceding the target instruction.
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),             // 0000
				code.Make(code.OpJumpNotTruthy, 7), // 0001
				code.Make(code.OpConstant, 0),      // 0004
				code.Make(code.OpPop),              // 0007
				code.Make(code.OpConstant, 1),      // 0008
				code.Make(code.OpPop),              // 0009
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestExpressionStatement(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: `
if (true) {
	10;
};
`,
			expectedConstants: []any{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),             // 0000
				code.Make(code.OpJumpNotTruthy, 7), // 0001
				code.Make(code.OpConstant, 0),      // 0004
				code.Make(code.OpPop),              // 0007
			},
		},
		{
			input:             `if (true) { 10 } else { 20 }; 3333`,
			expectedConstants: []any{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),              // 0000
				code.Make(code.OpJumpNotTruthy, 10), // 0001
				code.Make(code.OpConstant, 0),       // 0004
				code.Make(code.OpPop),               // 0007
				code.Make(code.OpJump, 14),          // 008
				code.Make(code.OpConstant, 1),       // 011
				code.Make(code.OpPop),               // 014
				code.Make(code.OpConstant, 2),       // 015
				code.Make(code.OpPop),               // 018
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestGlobalLetStatements(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: `
		let one = 1;
		let two = 2;
		`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				// The compiler executes according to the following steps for compile `let one = 1;`:
				// 1. compile the expression `1` and push result onto topmost constant pool;
				// 2. retrieve the constant value from constant pool and push it onto the topmost stack - `OpConstant 0`
				// 3. pop the topmost value off the stack and save it to the global store at the index encoded in the operand - `OpSetGlobal 0`
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
let one = 1;
one;
`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
let one = 1;
let two = one;
two;
`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				// let's think about this, why do these statements have three expressions
				// yet construct only one pop instruction?
				// this is because OpConstant produces a variable on the stack, but the OpSetGlobal pops it already
				// therefore, there is no need to generate an additional explicit pop instruction.
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestStringExpression(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []any{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `"mon" + "key"`,
			expectedConstants: []any{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func runCompilerTest(t *testing.T, caseIndex int, testCase *compilerTestCase) {
	t.Helper()
	c := NewCompiler()
	program := testParseProgram(testCase.input)
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("case [%d] error [%s] compile program for input : %s", caseIndex, err.Error(), testCase.input)
	}
	byteCode := c.ByteCode()

	testInstructions(t, caseIndex, testCase.expectedInstructions, byteCode.Instructions)
	testConstants(t, caseIndex, testCase.expectedConstants, byteCode.Constants)
}

func testInstructions(t *testing.T, caseIndex int, expectedInstructions []code.Instructions, actualInstruction code.Instructions) {
	t.Helper()
	var expectedLen = 0
	for _, instruction := range expectedInstructions {
		expectedLen += instruction.Len()
	}

	if expectedLen != len(actualInstruction) {
		t.Fatalf("case index [%d] wrong instructions length.\nexpect=%d\nactual=%d", caseIndex, expectedLen, len(actualInstruction))
	}

	var byteOffset = 0
	for i, expected := range expectedInstructions {
		nextInstructionByteOffset := byteOffset + len(expected)
		current := actualInstruction[byteOffset:nextInstructionByteOffset]
		byteOffset += len(expected)
		for j, ins := range expected {
			if current[j] != ins {
				err, curStr := code.BytesToInstruction(current).String()
				if err != nil {
					t.Fatalf("case %d error lookup op : %s", caseIndex, err.Error())
				}
				err, expectedStr := expected.String()
				if err != nil {
					t.Fatalf("case %d error lookup op : %s", caseIndex, err.Error())
				}
				t.Errorf("case %d instructionOffset = [%d] not match\nexpect=%s\nactual=%s", caseIndex, i, expectedStr, curStr)
			}
		}
	}
}

func testConstants(t *testing.T, caseIndex int, expectedConstants []interface{}, constants *code.Constants) {
	t.Helper()
	if len(expectedConstants) != constants.Len() {
		t.Fatalf("case %d wrong number of constants. expect=%d,actual=%d", caseIndex, len(expectedConstants), constants.Len())
	}
	for i, constant := range expectedConstants {
		switch expected := constant.(type) {
		case int:
			testIntegerObject(t, caseIndex, constants.GetConstant(code.Index(i)), int64(expected))
		case string:
			testStringObject(t, caseIndex, constants.GetConstant(code.Index(i)), expected)
		}
	}
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	bytecode := make(code.Instructions, 0)
	for _, instruction := range instructions {
		bytecode = bytecode.Append(instruction)
	}
	return bytecode
}

func testParseProgram(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(*l)
	return p.ParseProgram()
}

func testIntegerObject(t *testing.T, caseIndex int, obj object.Object, expected int64) {
	if obj == nil {
		t.Fatalf("case %d exepct int but got nil", caseIndex)
	}

	if obj.Type() != object.ObjInteger {
		t.Fatalf("case %d expect ObjInteger, got [%s], msg [%s]", caseIndex, string(obj.Type()), obj.Inspect())
	}
	integerObj, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("case %d expecte Integer, got [%s]", caseIndex, reflect.TypeOf(obj))
	}
	if integerObj.Value != expected {
		t.Fatalf("case %d expect [%d], got [%d]", caseIndex, expected, integerObj.Value)
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

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}
