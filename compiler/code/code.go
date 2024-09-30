package code

import (
	"bytes"
	"fmt"
)

const (
	// OpConstant it has one operand: the number we previously assigned to the constant.
	// 1. read the operand;
	// 2. retrieve constant value from constant pool through using operand as an index;
	// 3. push constant value onto the stack.
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
	// OpSetGlobal read in the operand, pop the topmost value off the stack and save it to the globals store
	// at the index encoded in the operand.
	OpSetGlobal
	// OpGetGlobal use the operand to retrieve the value from the globals store and push it onto the stack.
	OpGetGlobal
	// OpArray operate array with one operand:
	// Compile:
	// 1. compile all of its elements;
	// 2. compiling them results in instructions that leave N values on the VM's stack, where N is the number of element in the array literal;
	// 3. we're going to emit an OpArray instruction with the operand being N, the number of elements.
	// Run:
	// when the VM then executes the OpArray instruction it takes the N elements off the stack, builds an *object.Array out of them,
	// and pushes that on to the stack.
	OpArray
	// OpHash process object.Hash
	// Just like an array, its final value can't be determined during compile time.
	// Doubly so, actually because instead of having N elements, a hash in Monkey has N key and N value and all of them are created by expressions.
	// The operand specifies the number of keys and values sitting on the stack.
	// Compile :
	// 1. compile all of its keys and its values;
	// 2. compiling them results in instructions that leave 2 * N values on the VM's stack, where N is the number of keys in hash literals;
	// 3. emit an OpHash [2 * N] instruction.
	// Run :
	// 1. Read operand value;
	// 2. pop keys and values from stack and build an object.Hash and push that on to the stack.
	OpHash
	// OpIndex Instead of treating the index operator in combination with a specific data structure as a special case,
	// we'll build a generic index operator into the compiler and VM. So it's worked equally well with both array indices and hash indices.
	// for OpIndex to work, there need to be two values sitting on the top of the stack:
	// the object to be indexed and the object serving as the index.
	// When the VM executes OpIndex it should take both off the stack, perform the index operation, and put the result back on.
	OpIndex
	// OpCall calling convention, see https://en.wikipedia.org/wiki/Calling_convention
	// With OpCall defined, we are now able to get a function on to the stack of our VM and call it.
	// What we still need is a way to tell the VM to return from a called function.
	// More specifically, we need to differentiate between two cases where the Vm has to return from a function.
	// 1. the execution of a function actually returning something;
	// 2. the execution of a function ends without anything being returned;
	// --------------------------------------
	// Assuming we have a code snippet in monkey programming language:
	// fn(){} ()
	// The code snippet consists of two parts:
	// 1. the definition of the function;
	// 2. the calling of the function.
	// What really concerns me is the second part -- the calling of the function.
	// It produces two major instruction in compiling procedure :
	// OpConstant The first instruction retrieves CompiledFunction from constant pool and pushes result onto stack.
	// OpCall The second instruction pops the result from stack and execute it, that's why OpCall doesn't require an operand
	// --------------------------------------
	OpCall
	// OpReturnValue It doesn't have any arguments. The value to be returned has to sit on top of the stack.
	OpReturnValue
	// OpReturn representing a function that neither explicitly nor implicitly returns a value
	OpReturn
	// OpSetLocal binding local variable with one-byte operand.
	OpSetLocal
	// OpGetLocal retrieve local variable with one-byte operand.
	OpGetLocal
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
	OpJump:          {"OpJump", "", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", "", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", "", []int{2}},
	OpArray:         {"OpArray", "", []int{2}},
	OpHash:          {"OpHash", "", []int{2}},
	OpIndex:         {"OpIndex", "", []int{}},
	OpCall:          {"OpCall", "", []int{1}},
	OpReturnValue:   {"OpReturnValue", "", []int{}},
	OpReturn:        {"OpReturn", "", []int{}},
	OpSetLocal:      {"OpSetLocal", "", []int{1}},
	OpGetLocal:      {"OpGetLocal", "", []int{1}},
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

func (instruction Instructions) TestOnlyPrintln() {
	_, s := instruction.String()
	fmt.Println(s)
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
