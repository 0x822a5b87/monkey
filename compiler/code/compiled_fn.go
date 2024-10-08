package code

import (
	"0x822a5b87/monkey/interpreter/object"
	"fmt"
	"strings"
)

const (
	ObjCompiledFunction object.ObjType = "COMPILED_FUNCTION"
	ObjClosure          object.ObjType = "Closure"
)

// CompiledFunction a function object that holds bytecode instead of AST nodes.
// It can hold the Instructions we get from the compilation of a function literal, and it's an object.Object, which means
// we can add it as a constant to our compiler.ByteCode and load it in the VM
type CompiledFunction struct {
	Instructions   Instructions
	NumOfLocalVars int
}

func (c *CompiledFunction) Type() object.ObjType {
	return ObjCompiledFunction
}

func (c *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction [%p]", c.Instructions)
}

type Closure struct {
	Fn   *CompiledFunction
	Free []object.Object
}

func (c *Closure) Type() object.ObjType {
	return ObjClosure
}

func (c *Closure) Inspect() string {
	objs := make([]string, 0)
	for _, o := range c.Free {
		objs = append(objs, o.Inspect())
	}
	return fmt.Sprintf("Closure [%s], Free [%s]", c.Inspect(), strings.Join(objs, ","))
}
