package compiler

import "0x822a5b87/monkey/compiler/code"

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position instructionIndex
}

type CompilationScope struct {
	instructions code.Instructions
	last         *EmittedInstruction
	previous     *EmittedInstruction
}

func NewCompilationScope() *CompilationScope {
	return &CompilationScope{
		instructions: make(code.Instructions, 0),
		last:         nil,
		previous:     nil,
	}
}

func NewEmittedInstruction(opcode code.Opcode, pos instructionIndex) *EmittedInstruction {
	return &EmittedInstruction{
		Opcode:   opcode,
		Position: pos,
	}
}
