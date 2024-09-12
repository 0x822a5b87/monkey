package compiler

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/compiler/util"
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
		testCaseInfo := &util.TestCaseInfo{
			T:             t,
			TestFnName:    "TestIntegerCompiler",
			TestCaseIndex: i,
		}
		runCompilerTest(testCaseInfo, &testCase)
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
		testCaseInfo := &util.TestCaseInfo{
			T:             t,
			TestFnName:    "TestBooleanExpressions",
			TestCaseIndex: i,
		}
		runCompilerTest(testCaseInfo, &testCase)
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
		testCaseInfo := &util.TestCaseInfo{
			T:             t,
			TestFnName:    "TestCondition",
			TestCaseIndex: i,
		}
		runCompilerTest(testCaseInfo, &testCase)
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
	}

	for i, testCase := range testCases {
		testCaseInfo := &util.TestCaseInfo{
			T:             t,
			TestFnName:    "TestExpressionStatement",
			TestCaseIndex: i,
		}
		runCompilerTest(testCaseInfo, &testCase)
	}
}

func runCompilerTest(testCaseInfo *util.TestCaseInfo, testCase *compilerTestCase) {
	testCaseInfo.T.Helper()
	c := NewCompiler()
	program := testParseProgram(testCase.input)
	err := c.Compile(program)
	if err != nil {
		testCaseInfo.Fatalf("error [%s] compile program for input : %s", err.Error(), testCase.input)
	}
	byteCode := c.ByteCode()

	testInstructions(testCaseInfo, testCase.expectedInstructions, byteCode.Instructions)
	testConstants(testCaseInfo, testCase.expectedConstants, byteCode.Constants)
}

func testInstructions(info *util.TestCaseInfo, expectedInstructions []code.Instructions, actualInstruction code.Instructions) {
	info.Helper()
	info.T.Helper()
	var expectedLen = 0
	for _, instruction := range expectedInstructions {
		expectedLen += instruction.Len()
	}

	if expectedLen != len(actualInstruction) {
		info.T.Fatalf("wrong instructions length.\nexpect=%d\nactual=%d", expectedLen, len(actualInstruction))
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
					info.T.Fatalf("error lookup op : %s", err.Error())
				}
				err, expectedStr := expected.String()
				if err != nil {
					info.T.Fatalf("error lookup op : %s", err.Error())
				}
				info.T.Errorf("instructionOffset = [%d] not match\nexpect=%s\nactual=%s", i, expectedStr, curStr)
			}
		}
	}
}

func testConstants(info *util.TestCaseInfo, expectedConstants []interface{}, constants *code.Constants) {
	if len(expectedConstants) != constants.Len() {
		info.T.Fatalf("wrong number of constants. expect=%d,actual=%d", len(expectedConstants), constants.Len())
	}
	for i, constant := range expectedConstants {
		switch expected := constant.(type) {
		case int:
			testIntegerObject(info, constants.GetConstant(code.Index(i)), int64(expected))
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

func testIntegerObject(info *util.TestCaseInfo, obj object.Object, expected int64) {
	if obj == nil {
		info.Fatalf("exepct not nil, got nil")
	}

	if obj.Type() != object.ObjInteger {
		info.Fatalf("expect ObjInteger, got [%s], msg [%s]", string(obj.Type()), obj.Inspect())
	}
	integerObj, ok := obj.(*object.Integer)
	if !ok {
		info.Fatalf("expecte Integer, got [%s]", reflect.TypeOf(obj))
	}
	if integerObj.Value != expected {
		info.Fatalf("expect [%d], got [%d]", expected, integerObj.Value)
	}
}

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}
