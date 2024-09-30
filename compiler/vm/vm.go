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
const MaxFrameSize = 1024

type Vm struct {
	constants   *code.Constants
	stack       []object.Object
	globalStore []object.Object

	// sp means stack pointer, we'll increment or decrement to grow or shrink the stack, instead of modifying the stack slice itself.
	// for the convention, it will always point to the next free slot in the stack
	// If there's one element on the stack, located at index 0, the value of stackPointer would be 1
	sp int

	frames      []*Frame
	framesIndex int
}

func NewVm(c *compiler.ByteCode) *Vm {
	main := &code.CompiledFunction{Instructions: c.Instructions}

	v := &Vm{
		constants: c.Constants,

		stack:       make([]object.Object, StackSize),
		globalStore: make([]object.Object, GlobalStoreSize),
		sp:          0,

		frames:      make([]*Frame, MaxFrameSize),
		framesIndex: 0,
	}

	v.pushFrame(NewFrame(main, 0))

	return v
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
	for v.hasNext() {
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
		case code.OpArray:
			err = v.executeOpArray(op)
		case code.OpHash:
			err = v.executeHash(op)
		case code.OpIndex:
			err = v.executeIndex(op)
		case code.OpCall:
			err = v.executeCall(op)
		case code.OpReturnValue:
			err = v.executeReturnValue(op)
		case code.OpReturn:
			err = v.executeReturn(op)
		case code.OpSetLocal:
			err = v.executeSetLocal(op)
		case code.OpGetLocal:
			err = v.executeGetLocal(op)
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
	return v.currentIp() < len(v.currentInstructions())
}

func (v *Vm) currentOpcode() code.Opcode {
	instructions := v.currentInstructions()
	return code.Opcode(instructions[v.currentIp()])
}

func (v *Vm) incrementIp(delta int) {
	v.currentFrame().ip += delta
}

func (v *Vm) decrementIp(delta int) {
	v.currentFrame().ip -= delta
}

func (v *Vm) pop() object.Object {
	topElement := v.StackTop()
	if topElement != nil {
		v.sp--
	}
	return topElement
}

func (v *Vm) opConstant() error {
	defer v.incrementIp(1)
	// ATTENTION: let's recap what we did with the object.Integer in the compiler
	// we construct an Integer object and added it to the constant pool.
	// Then, we emitted an instruction with the index of the Integer object in the constant pool to bytecode
	constantIndex := v.readUint16AndIncIp()
	return v.push(v.constants.GetConstant(constantIndex))
}

func (v *Vm) opBoolean(op code.Opcode) error {
	defer v.incrementIp(1)
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
	defer v.incrementIp(1)
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
	defer v.incrementIp(1)
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
	defer v.incrementIp(1)

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
	defer v.incrementIp(1)

	definition, _ := code.Lookup(op)
	return v.doJump(definition)
}

func (v *Vm) executeSetGlobal(op code.Opcode) error {
	defer v.incrementIp(1)

	globalIndex := v.readUint16AndIncIp()
	v.globalStore[globalIndex] = v.pop()
	return nil
}

func (v *Vm) executeGetGlobal(op code.Opcode) error {
	defer v.incrementIp(1)

	globalIndex := v.readUint16AndIncIp()
	value := v.globalStore[globalIndex]
	return v.push(value)
}

func (v *Vm) executeOpArray(op code.Opcode) error {
	defer v.incrementIp(1)

	n := v.readUint16AndIncIp()
	array := &object.Array{Elements: make([]object.Object, n.IntValue())}
	for i := n.IntValue() - 1; i >= 0; i-- {
		array.Elements[i] = v.pop()
	}
	return v.push(array)
}

func (v *Vm) executeHash(op code.Opcode) error {
	defer v.incrementIp(1)

	doubleN := v.readUint16AndIncIp()
	hash := &object.Hash{Pairs: make(map[object.HashKey]*object.HashPair)}
	for i := 0; i < doubleN.IntValue(); i += 2 {
		value := v.pop()
		key := v.pop()
		hashable := key.(object.Hashable)
		hash.Pairs[hashable.HashKey()] = &object.HashPair{
			Key:   key,
			Value: value,
		}
	}
	return v.push(hash)
}

func (v *Vm) executeIndex(op code.Opcode) error {
	defer v.incrementIp(1)

	index := v.pop()
	obj := v.pop()

	indexed, ok := obj.(object.Index)
	if !ok {
		return common.NewErrIndex(indexed.Type())
	}

	return v.push(indexed.Index(index))
}

func (v *Vm) executeCall(op code.Opcode) error {
	// NumOfLocalVars = NumOfArguments + NumOfVariablesDefinedInFunction
	numOfArgs := v.readUint8AndIncIp()

	obj := v.stack[v.sp-numOfArgs.IntValue()-1]
	fn, ok := obj.(*code.CompiledFunction)
	if !ok {
		return common.NewErrTypeMismatch(object.ObjFunction.String(), obj.Type().String())
	}

	// base pointer points to the start position of local variable
	basePointer := v.sp - numOfArgs.IntValue()
	// stack pointer points to the start position of the new frame's stack
	stackPointer := v.sp + fn.NumOfLocalVars

	frame := NewFrame(fn, basePointer)
	v.sp = stackPointer
	v.pushFrame(frame)

	return nil
}

func (v *Vm) executeReturnValue(op code.Opcode) error {
	defer v.incrementIp(1)
	returnValue := v.pop()
	oldFrame := v.popFrame()
	v.sp = oldFrame.basePointer - 1
	_ = v.push(returnValue)
	return nil
}

func (v *Vm) executeReturn(op code.Opcode) error {
	defer v.incrementIp(1)
	oldFrame := v.popFrame()
	v.sp = oldFrame.basePointer - 1
	return v.push(object.NativeNull)
}

func (v *Vm) executeSetLocal(op code.Opcode) error {
	defer v.incrementIp(1)
	relativeIndex := v.readUint8AndIncIp()
	realIndex := v.currentFrame().basePointer + relativeIndex.IntValue()
	v.stack[realIndex] = v.pop()
	return nil
}

func (v *Vm) executeGetLocal(op code.Opcode) error {
	defer v.incrementIp(1)
	relativeIndex := v.readUint8AndIncIp()
	realIndex := v.currentFrame().basePointer + relativeIndex.IntValue()
	o := v.stack[realIndex]
	return v.push(o)
}

func (v *Vm) readUint16() code.Index {
	return code.ReadUint16(v.currentInstructions()[v.currentIp()+1:])
}

func (v *Vm) readUint16AndIncIp() code.Index {
	val := v.readUint16()
	v.incrementIp(2)
	return val
}

func (v *Vm) readUint8() code.Index {
	nextIp := v.currentIp() + 1
	val := v.currentInstructions()[nextIp]
	return code.Index(val)
}

func (v *Vm) readUint8AndIncIp() code.Index {
	val := v.readUint8()
	v.incrementIp(1)
	return val
}

func (v *Vm) doJump(definition *code.Definition) error {
	operands, _ := code.ReadOperands(definition, v.currentInstructions()[v.currentIp()+1:])
	if len(operands) != 1 {
		return common.NewErrOperandsCount(1, len(operands))
	}
	v.currentFrame().ip = operands[0]
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
	defer v.incrementIp(1)
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

func (v *Vm) currentInstructions() code.Instructions {
	return v.currentFrame().Instructions()
}

func (v *Vm) currentIp() int {
	return v.currentFrame().ip
}

func (v *Vm) currentFrame() *Frame {
	return v.frames[v.framesIndex-1]
}

func (v *Vm) pushFrame(newFrame *Frame) {
	v.frames[v.framesIndex] = newFrame
	v.framesIndex++
}

func (v *Vm) popFrame() *Frame {
	oldFrame := v.currentFrame()
	v.framesIndex--
	return oldFrame
}
