package vm

import "0x822a5b87/monkey/compiler/code"

// Frame short for call frame or stack frame
// A Frame has two fields: ip and fn.
// fn points to the compiled function referenced by the frame.
// ip is the instruction pointer in this frame, for the function.
type Frame struct {
	fn *code.Closure
	// ip means instruction pointer, we iterate through instructions by incrementing it.
	// Then fetch the current instruction by directly accessing instructions.
	ip int
	// basePointer the stack pointer of callee
	basePointer int
}

func NewFrame(f *code.Closure, stackPointer int) *Frame {
	return &Frame{
		fn:          f,
		ip:          0,
		basePointer: stackPointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Fn.Instructions
}
