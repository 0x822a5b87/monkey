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

func TestMinimalCompiler(t *testing.T) {
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
	}

	for i, testCase := range testCases {
		testCaseInfo := &util.TestCaseInfo{
			T:             t,
			TestFnName:    "TestMinimalCompiler",
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
				info.T.Fatalf("instructionOffset = [%d] not match\nexpect=%s\nactual=%s", i, expectedStr, curStr)
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
