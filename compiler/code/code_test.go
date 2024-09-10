package code

import (
	"0x822a5b87/monkey/compiler/util"
	"testing"
)

func TestMake(t *testing.T) {
	testCases := []struct {
		op       Opcode
		operand  []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for i, testCase := range testCases {
		instruction := Make(testCase.op, testCase.operand...)
		if len(instruction) != len(testCase.expected) {
			t.Fatalf("instruction has wrong length. expected=%d, got = %d", len(testCase.expected), len(instruction))
		}
		for j, expectedByte := range testCase.expected {
			b := instruction[j]
			if b != expectedByte {
				t.Fatalf("test case [%d] failed : byte [%d] not match, expected=%d, got=%d", i, j, expectedByte, b)
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}
	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
0010 OpConstant 2
0013 OpConstant 65535`

	concat := Instructions{}
	for _, ins := range instructions {
		concat = concat.Append(ins)
	}

	info := &util.TestCaseInfo{
		T:             t,
		TestFnName:    "TestInstructionsString",
		TestCaseIndex: 0,
	}
	err, s := concat.String()
	if err != nil {
		info.Fatalf(err.Error())
	}

	if s != expected {
		info.Fatalf("instructions wrongly formatted.\nexpect=[%q]\nactual=[%q]", expected, s)
	}
}

func TestReadOperands(t *testing.T) {
	testCases := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for i, tt := range testCases {
		info := &util.TestCaseInfo{
			T:             t,
			TestFnName:    "TestReadOperands",
			TestCaseIndex: i,
		}

		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(tt.op)
		if err != nil {
			info.Fatalf(err.Error())
		}

		operand, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			info.Fatalf("n wrong. expected=%d,got=%d", tt.bytesRead, n)
		}

		for j, expected := range tt.operands {
			if operand[j] != expected {
				info.Fatalf("operand wrong. expected=%d, got=%d", expected, operand[j])
			}
		}
	}
}
