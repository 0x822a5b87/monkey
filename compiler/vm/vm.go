package vm

import (
	"0x822a5b87/monkey/compiler/code"
	"0x822a5b87/monkey/compiler/compiler"
	"0x822a5b87/monkey/interpreter/common"
	"0x822a5b87/monkey/interpreter/evaluator"
	"0x822a5b87/monkey/interpreter/object"
	"errors"
	"fmt"
)

const StackSize = 2048
const GlobalStoreSize = 65535

type Vm struct {
	constants    *code.Constants
	instructions code.Instructions
	stack        []object.Object
	globalStore  []object.Object

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

		stack:       make([]object.Object, StackSize),
		globalStore: make([]object.Object, GlobalStoreSize),
		sp:          0,
		ip:          0,
	}
}

func NewVmWithState(code *compiler.ByteCode, prev *Vm) *Vm {
	v := NewVm(code)
	v.globalStore = prev.globalStore
	return v
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
		case code.OpPop:
			err = v.opPop()
		case code.OpTrue, code.OpFalse:
			err = v.opBoolean(op)
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpLessThan:
			err = v.executeBinaryOperation(op)
		case code.OpBang, code.OpMinus:
			err = v.executePrefixOpcode(op)
		case code.OpJumpNotTruthy:
			err = v.executeNotTruthyJump(op)
		case code.OpJump:
			err = v.executeJump(op)
		case code.OpSetGlobal:
			err = v.executeSetGlobal(op)
		case code.OpGetGlobal:
			err = v.executeGetGlobal(op)
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

func (v *Vm) opBoolean(op code.Opcode) error {
	switch op {
	case code.OpTrue:
		return v.push(object.NativeTrue)
	case code.OpFalse:
		return v.push(object.NativeFalse)
	default:
		return common.NewErrTypeMismatch("Boolean", string(op))
	}
}

// opAdd add two number from the stack and push the result back onto the stack
func (v *Vm) opAdd(lhs, rhs object.Object) error {
	left := lhs.(object.Add)
	right := rhs.(object.Add)
	return v.push(left.Add(right))
}

func (v *Vm) opSub(lhs, rhs object.Object) error {
	left := lhs.(object.Subtract)
	right := rhs.(object.Subtract)
	return v.push(left.Sub(right))
}

func (v *Vm) opMul(lhs, rhs object.Object) error {
	left := lhs.(object.Multiply)
	right := rhs.(object.Multiply)
	return v.push(left.Mul(right))
}

func (v *Vm) opDiv(lhs, rhs object.Object) error {
	left := lhs.(object.Divide)
	right := rhs.(object.Divide)
	return v.push(left.Divide(right))
}

func (v *Vm) opEqual(lhs, rhs object.Object) error {
	left := lhs.(object.Equatable)
	right := rhs.(object.Equatable)
	return v.push(left.Equal(right))
}
func (v *Vm) opNotEqual(lhs, rhs object.Object) error {
	left := lhs.(object.Equatable)
	right := rhs.(object.Equatable)
	return v.push(left.NotEqual(right))
}
func (v *Vm) opGreaterThan(lhs, rhs object.Object) error {
	left := lhs.(object.Comparable)
	right := rhs.(object.Comparable)
	return v.push(left.GreaterThan(right))
}
func (v *Vm) opLessThan(lhs, rhs object.Object) error {
	left := lhs.(object.Comparable)
	right := rhs.(object.Comparable)
	return v.push(left.LessThan(right))
}

func (v *Vm) executeBinaryOperation(op code.Opcode) error {
	rhs := v.pop()
	lhs := v.pop()
	definition, _ := code.Lookup(op)
	if lhs == nil || rhs == nil {
		return common.NewErrEmptyStack(definition.Name)
	}

	err := evaluator.InfixExpressionTypeCheck(definition.Operator, lhs, rhs)
	if err != nil {
		return errors.New(err.Message)
	}

	switch op {
	case code.OpAdd:
		return v.opAdd(lhs, rhs)
	case code.OpSub:
		return v.opSub(lhs, rhs)
	case code.OpMul:
		return v.opMul(lhs, rhs)
	case code.OpDiv:
		return v.opDiv(lhs, rhs)
	case code.OpEqual:
		return v.opEqual(lhs, rhs)
	case code.OpNotEqual:
		return v.opNotEqual(lhs, rhs)
	case code.OpGreaterThan:
		return v.opGreaterThan(lhs, rhs)
	case code.OpLessThan:
		return v.opLessThan(lhs, rhs)
	default:
		return common.NewErrUnsupportedBinaryExpr(definition.Name)
	}
}

func (v *Vm) executePrefixOpcode(op code.Opcode) error {
	lhs := v.pop()
	definition, _ := code.Lookup(op)
	if lhs == nil {
		return common.NewErrEmptyStack(definition.Name)
	}

	switch op {
	case code.OpBang:
		return v.opBang(lhs)
	case code.OpMinus:
		return v.opMinus(lhs)
	default:
		return common.NewErrUnsupportedBinaryExpr(definition.Name)
	}
}

func (v *Vm) executeNotTruthyJump(op code.Opcode) error {
	definition, _ := code.Lookup(op)

	lhs := v.pop()
	if lhs == nil {
		return common.NewErrEmptyStack(definition.Name)
	}

	if !v.isTruthy(lhs) {
		return v.doJump(definition)
	} else {
		// skip operands
		v.incrementIp(2)
		return nil
	}
}

func (v *Vm) executeJump(op code.Opcode) error {
	definition, _ := code.Lookup(op)
	return v.doJump(definition)
}

func (v *Vm) executeSetGlobal(op code.Opcode) error {
	globalIndex := v.readUint16()
	v.incrementIp(2)
	v.globalStore[globalIndex] = v.pop()
	return nil
}

func (v *Vm) executeGetGlobal(op code.Opcode) error {
	globalIndex := v.readUint16()
	v.incrementIp(2)
	value := v.globalStore[globalIndex]
	return v.push(value)
}

func (v *Vm) readUint16() code.Index {
	return code.ReadUint16(v.instructions[v.ip+1:])
}

func (v *Vm) doJump(definition *code.Definition) error {
	operands, _ := code.ReadOperands(definition, v.instructions[v.ip+1:])
	if len(operands) != 1 {
		return common.NewErrOperandsCount(1, len(operands))
	}
	v.ip = operands[0]
	return nil
}

func (v *Vm) isTruthy(obj object.Object) bool {
	switch obj {
	case object.NativeNull, object.NativeFalse:
		return false
	default:
		return true
	}
}

func (v *Vm) opPop() error {
	v.pop()
	return nil
}

func (v *Vm) opBang(lhs object.Object) error {
	switch lhs {
	case object.NativeFalse, object.NativeNull:
		return v.push(object.NativeTrue)
	case object.NativeTrue:
		return v.push(object.NativeFalse)
	default:
		return v.push(object.NativeFalse)
	}
}

func (v *Vm) opMinus(lhs object.Object) error {
	left := lhs.(object.Negative)
	return v.push(left.Negative())
}

func (v *Vm) operands() {

}
