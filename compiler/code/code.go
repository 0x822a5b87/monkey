package code

import (
	"bytes"
)

const (
	// OpConstant it has one operand: the number we previously assigned to the constant.
	// When the VM executes OpConstant it retrieves the constant using the operand as an index and pushes
	// it on to the stack.
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLessThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
)

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", "", []int{2}},
	// tells the VM to pop the two topmost elements off the stack, add them and push the result back on the stack
	OpAdd: {"OpAdd", "+", []int{}},
	// OpPop doesn't need to use any operands. It's only job is to tell the VM to pop
	// the topmost element off the stack and for that it doesn't need an operand.
	// For further information, let's have a quick refresher, there are three types of statement in monkey:
	// 1. let statement
	// 2. return statement
	// 3. expression statement
	// Whereas the first two explicitly reuse the value their child-expression nodes produce,
	// but expression statement merely wrap expressions so the can occur on their own. The value they produce is not
	// reuse, by definition. So, we need to emit a OpPop for every expression statement to clear it up.
	OpPop:         {"OpPop", "", []int{}},
	OpSub:         {"OpSub", "-", []int{}},
	OpMul:         {"OpMul", "*", []int{}},
	OpDiv:         {"OpDiv", "/", []int{}},
	OpTrue:        {"OpTrue", "", []int{}},
	OpFalse:       {"OpFalse", "", []int{}},
	OpEqual:       {"OpEqual", "==", []int{}},
	OpNotEqual:    {"OpNotEqual", "!=", []int{}},
	OpGreaterThan: {"OpGreaterThan", ">", []int{}},
	OpLessThan:    {"OpLessThan", "<", []int{}},
	OpMinus:       {"OpMinus", "-", []int{}},
	OpBang:        {"OpBang", "!", []int{}},
	// tell the VM to only jump if the value on top of stack is not monkey truthy
	// The operand of OpJumpNotTruthy and OpJump is 16-bit wide.
	OpJumpNotTruthy: {"OpJumpNotTruthy", "", []int{2}},
	OpJump:          {"OpJump", "!", []int{2}},
}

// Instructions the instructions are a series of bytes and a single instruction
// consist of an opcode and an optional number of operands.
//
// We didn't specify two respective types: "Instruction" for a single instruction and "instructions" for a sequence of instructions, like this:
//
// type Instruction []byte
// type Instructions []Instruction
//
// This is for a more convenient way to access instructions, where we define two types that transform a one-dimensional array into a two-dimensional array.
type Instructions []byte

func (instruction Instructions) Append(appended Instructions) Instructions {
	return append(instruction, appended...)
}

func (instruction Instructions) Len() int {
	return len(instruction)
}

func (instruction Instructions) String() (error, string) {
	buffer := bytes.Buffer{}
	for offsetBytes := 0; offsetBytes < instruction.Len(); {
		op := instruction.Opcode(offsetBytes)
		def, err := Lookup(op)
		if err != nil {
			return err, ""
		}
		operands, readBytes := ReadOperands(def, instruction[offsetBytes+1:])
		buffer.WriteString(formatInstruction(offsetBytes, def.Name, operands))
		offsetBytes += readBytes + 1
		if offsetBytes != instruction.Len() {
			buffer.WriteString("\n")
		}
	}
	return nil, buffer.String()
}

func (instruction Instructions) Opcode(index int) Opcode {
	return Opcode(instruction[index])
}

// Opcode opcode
type Opcode byte

func (c Opcode) Byte() byte {
	return byte(c)
}

// Definition the definition for an Opcode has two fields: Name and OperandWidths
// Name helps to make an Opcode readable
// OperandWidths contains the number of bytes each operand takes up.
type Definition struct {
	Name          string
	Operator      string
	OperandWidths []int
}

func (d *Definition) TotalSize() int {
	return d.OperandSize() + 1
}

func (d *Definition) OperandSize() int {
	var instructionLen int
	for _, w := range d.OperandWidths {
		instructionLen += w
	}
	return instructionLen
}

func (d *Definition) OperandLen() int {
	return len(d.OperandWidths)
}

func (d *Definition) NthOperandWidth(nth int) int {
	return d.OperandWidths[nth]
}
