package vm

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/compiler/compiler"
	"0x822a5b87/monkey/interpreter/common"
	"0x822a5b87/monkey/interpreter/object"
	"errors"
	"fmt"
)

const StackSize = 2048

type Vm struct {
	constants    *code.Constants
	instructions code.Instructions
	stack        []object.Object

	// sp means stack pointer, we'll increment or decrement to grow or shrink the stack, instead of modifying the stack slice itself.
	// for the convention, it will always point to the next free slot in the stack
	// If there's one element on the stack, located at index 0, the value of stackPointer would be 1
	sp int
	// ip means instruction pointer, we iterate through instructions by incrementing it.
	// Then fetch the current instruction by directly accessing instructions.
	ip int
}

func NewVm(code *compiler.ByteCode) *Vm {
	return &Vm{
		constants:    code.Constants,
		instructions: code.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,
		ip:    0,
	}
}

// Run what turns the Vm into a real virtual machine
// it read single instruction from Vm's instructions then execute it until the end of instructions.
//
// In this method, we use code.ReadUint16 instead of code.ReadOperands for the same reason we don't use code.Lookup
// when fetching the instruction: performance.
func (v *Vm) Run() error {
	var err error
	// In every loop, we reach the end of a single instruction and increment by 1 byte to move to the next instruction
	for ; v.hasNext(); v.incrementIp(1) {
		op := v.currentOpcode()
		switch op {
		case code.OpConstant:
			err = v.opConstant()
		case code.OpAdd:
			err = v.opAdd()
		case code.OpPop:
			err = v.opPop()
		default:
			err = fmt.Errorf("wrong type of Opcode : [%d]", op)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Vm) StackTop() object.Object {
	if v.sp == 0 {
		return nil
	}

	return v.stack[v.sp-1]
}

// TestOnlyLastPoppedStackElement this method is for test only
// Sometimes, we want to assert that "this should have been on the stack, right before you popped it off"
func (v *Vm) TestOnlyLastPoppedStackElement() object.Object {
	return v.stack[v.sp]
}

func (v *Vm) push(o object.Object) error {
	if v.sp >= StackSize {
		return errors.New("stack overflow")
	}

	v.stack[v.sp] = o
	v.sp++

	return nil
}

func (v *Vm) hasNext() bool {
	return v.ip < len(v.instructions)
}

func (v *Vm) currentOpcode() code.Opcode {
	return code.Opcode(v.instructions[v.ip])
}

func (v *Vm) incrementIp(delta int) {
	v.ip += delta
}

func (v *Vm) decrementIp(delta int) {
	v.ip -= delta
}

func (v *Vm) pop() object.Object {
	topElement := v.StackTop()
	if topElement != nil {
		v.sp--
	}
	return topElement
}

func (v *Vm) opConstant() error {
	// ATTENTION: let's recap what we did with the object.Integer in the compiler
	// we construct an Integer object and added it to the constant pool.
	// Then, we emitted an instruction with the index of the Integer object in the constant pool to bytecode
	constantIndex := code.ReadUint16(v.instructions[v.ip+1:])
	v.incrementIp(2)
	return v.push(v.constants.GetConstant(constantIndex))
}

// opAdd add two number from the stack and push the result back onto the stack
func (v *Vm) opAdd() error {
	lhs := v.pop()
	rhs := v.pop()
	if lhs == nil || rhs == nil {
		return common.NewErrEmptyStack("OpAdd")
	}

	left, ok := lhs.(object.Add)
	if !ok {
		return common.NewErrTypeMismatch("integer", string(left.Type()))
	}

	right, ok := rhs.(object.Add)
	if !ok {
		return common.NewErrTypeMismatch("integer", string(right.Type()))
	}

	return v.push(left.Add(right))
}

func (v *Vm) opPop() error {
	v.pop()
	return nil
}
