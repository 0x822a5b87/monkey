package code

import (
	"0x822a5b87/monkey/interpreter/common"
	"encoding/binary"
)

func Lookup(opcode Opcode) (*Definition, error) {
	def, ok := definitions[opcode]
	if !ok {
		return nil, common.NewOpcodeUndefined(opcode.Byte())
	}
	return def, nil
}

func Make(op Opcode, operands ...int) Instructions {
	definition, err := Lookup(op)
	if err != nil {
		panic(err)
	}
	instruction := make([]byte, definition.TotalSize())
	instruction[0] = op.Byte()
	var offset = 1
	for i, operand := range operands {
		width := definition.NthOperandWidth(i)
		switch width {
		case 1:
			instruction[offset] = byte(operand)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		default:
			panic(common.NewOperandWidthError(width))
		}
		offset += width
	}
	return instruction
}

func ReadOperandsSelf(instruction []byte) (err error, operands []int, readBytes int) {
	op := Opcode(instruction[0])
	def, err := Lookup(op)
	if err != nil {
		return
	}
	operands, readBytes = ReadOperands(def, instruction[1:])
	return
}

// ReadOperands the Definition is very necessary, because the instruction is already skip Opcode
func ReadOperands(def *Definition, instruction []byte) (operands []int, readBytes int) {
	if def == nil {
		return
	}
	operands = make([]int, def.OperandLen())

	readBytes = 0
	var start = 0
	for i, width := range def.OperandWidths {
		start = readBytes
		readBytes += width
		switch width {
		case 1:
			operands[i] = ReadUint8(instruction[start:]).IntValue()
		case 2:
			operands[i] = ReadUint16(instruction[start:readBytes]).IntValue()
		default:
			panic(common.NewOperandWidthError(width))
		}
	}
	return
}

func BytesToInstruction(data []byte) Instructions {
	return data
}

func ReadUint16(instructions Instructions) Index {
	u := binary.BigEndian.Uint16(instructions)
	return Index(u)
}

func ReadUint8(instructions Instructions) Index {
	c := instructions[0]
	return Index(c)
}
