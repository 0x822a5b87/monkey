package code

import (
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
		{OpGetLocal, []int{254}, []byte{byte(OpGetLocal), 254}},
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
		Make(OpGetLocal, 0),
		Make(OpSetLocal, 255),
	}
	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
0010 OpConstant 2
0013 OpConstant 65535
0016 OpGetLocal 0
0018 OpSetLocal 255`

	concat := Instructions{}
	for _, ins := range instructions {
		concat = concat.Append(ins)
	}

	err, s := concat.String()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if s != expected {
		t.Fatalf("instructions wrongly formatted.\nexpect=[%q]\nactual=[%q]", expected, s)
	}
}

func TestReadOperands(t *testing.T) {
	testCases := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpGetLocal, []int{255}, 1},
	}

	for _, tt := range testCases {

		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(tt.op)
		if err != nil {
			t.Fatalf(err.Error())
		}

		operand, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. expected=%d,got=%d", tt.bytesRead, n)
		}

		for j, expected := range tt.operands {
			if operand[j] != expected {
				t.Fatalf("operand wrong. expected=%d, got=%d", expected, operand[j])
			}
		}
	}
}
