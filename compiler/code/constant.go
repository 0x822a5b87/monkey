package code

import (
	"0x822a5b87/monkey/interpreter/object"
)

type Index int

func (c Index) IntValue() int {
	return int(c)
}

type Constants struct {
	constantPool []object.Object
}

func NewConstants() *Constants {
	return &Constants{constantPool: make([]object.Object, 0)}
}

func (c *Constants) GetConstant(index Index) object.Object {
	return c.constantPool[index]
}

func (c *Constants) AddConstant(obj object.Object) Index {
	c.constantPool = append(c.constantPool, obj)
	return Index(len(c.constantPool) - 1)
}

func (c *Constants) Len() int {
	return len(c.constantPool)
}
