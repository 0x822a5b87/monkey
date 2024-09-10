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
)

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	// tells the VM to pop the two topmost elements off the stack, add them and push the result back on the stack
	OpAdd: {"OpAdd", []int{}},
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
