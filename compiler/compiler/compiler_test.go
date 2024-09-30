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

func TestLetStatement(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: `
		let num = 65535;
		fn() { num; }
				`,
			expectedConstants: []any{
				65535,
				[]code.Instructions{
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
fn() {
	let num = 65535;
	num;
}
		`,
			expectedConstants: []any{
				65535,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
fn() {
	let a = 55;
	let b = 77;
	a + b;
}
		`,
			expectedConstants: []any{
				55,
				77,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestArrayLiterals(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input:             `[]`,
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `[1, 2, 3]`,
			expectedConstants: []any{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `[1 + 2, 3 - 4, 5 * 6]`,
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestHashLiterals(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input:             `{}`,
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `{1: 2, 3: 4, 5: 6}`,
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `{1: 2 + 3, 4: 5 * 6}`,
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestIndexExpressions(t *testing.T) {
	testCases := []compilerTestCase{
		// test array index
		{
			input:             `[1, 2, 3][1 + 1]`,
			expectedConstants: []any{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		// test hash index
		{
			input:             `{1 : 2}[2 -1]`,
			expectedConstants: []any{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestFunctions(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: `fn() { return 5 + 10; }`,
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() {}`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() {"hello world"}`,
			expectedConstants: []any{
				"hello world",
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { return fn() { "hello world" }; }`,
			expectedConstants: []any{
				"hello world",
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { 1; 2 }`,
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestFunctionCall(t *testing.T) {
	testCases := []compilerTestCase{
		{
			input: `fn() {}()`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { return 5 + 10; }()`,
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
				let noArgFn = fn() { 24; };
				noArgFn();
				`,
			expectedConstants: []any{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 1),  // the compiled function
				code.Make(code.OpSetGlobal, 0), // insert compiled function to global store
				code.Make(code.OpGetGlobal, 0), // retrieve compiled function from global store
				code.Make(code.OpCall, 0),      // call function
				code.Make(code.OpPop),          // pop value
			},
		},
		{
			input: `
		let earlyReturn = fn() { return 0; 1 };
		`,
			expectedConstants: []any{
				0,
				1,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpSetGlobal, 0),
			},
		},
		{
			input: `
let oneArg = fn(a) { };
oneArg(65535)
`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
				65535,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),

				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
				let oneArg = fn(a) { a };
				oneArg(24);
				`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				24,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
let manyArg = fn(a, b, c) { };
manyArg(24, 25, 26);
`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []code.Instructions{
				// pushing the function retrieved from constant pool onto the stack
				code.Make(code.OpConstant, 0),
				// binding the topmost object we just retrieved from constant pool to global variable
				code.Make(code.OpSetGlobal, 0),

				// calling function
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
		},
	}

	for i, testCase := range testCases {
		runCompilerTest(t, i, &testCase)
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{}

	for i, testCase := range tests {
		runCompilerTest(t, i, &testCase)
	}
}

func TestCompilerScopes(t *testing.T) {

	compiler := NewCompiler()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}

	globalSymbolTable := compiler.symbolTable
	compiler.emit(code.OpMul)
	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSub)
	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].last
	if last.Opcode != code.OpSub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d", last.Opcode, code.OpSub)
	}
	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.exitScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}
	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}
	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(code.OpAdd)
	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].last
	if last.Opcode != code.OpAdd {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d", last.Opcode, code.OpAdd)
	}

	previous := compiler.scopes[compiler.scopeIndex].previous
	if previous.Opcode != code.OpMul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d", previous.Opcode, code.OpMul)
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
		t.Errorf("case [%d] error [%s] compile program for input : %s", caseIndex, err.Error(), testCase.input)
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
		t.Errorf("case index [%d] wrong instructions length.\nexpect=%d\nactual=%d", caseIndex, expectedLen, len(actualInstruction))
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

func testCompiledFunction(t *testing.T, caseIndex int, expected []code.Instructions, actual object.Object) {
	t.Helper()
	compiledFn, ok := actual.(*code.CompiledFunction)
	if !ok {
		t.Errorf("case %d type error\nexpect=%s\nactual=%s", caseIndex, code.ObjCompiledFunction, actual.Type())
		return
	}
	testInstructions(t, caseIndex, expected, compiledFn.Instructions)
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
		case []code.Instructions:
			testCompiledFunction(t, caseIndex, expected, constants.GetConstant(code.Index(i)))
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
	t.Helper()
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
