package object

import (
	"0x822a5b87/monkey/ast"
	"0x822a5b87/monkey/util"
	"bytes"
	"fmt"
	"strconv"
)

var (
	NativeNull  = &Null{}
	NativeFalse = &Boolean{Value: false} // NativeFalse native false
	NativeTrue  = &Boolean{Value: true}  // NativeTrue native true
)

type ObjType string

// Object represent an object
// looks a lot like we did in the token package with the Token and TokenType types.
// Except that instead of being a structure like Token the Object type is an interface.
// The reason is that every value needs a different internal representation, and it's easier to define two
// different struct types than trying to fit booleans and integers into the same struct field.
type Object interface {
	Type() ObjType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjType {
	return ObjInteger
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) Add(o Object) Object {
	if other, ok := o.(*Integer); ok {
		i.Value = i.Value + other.Value
		return i
	}
	return NativeNull
}

func (i *Integer) Sub(o Object) Object {
	if other, ok := o.(*Integer); ok {
		i.Value = i.Value - other.Value
		return i
	}
	return NativeNull
}

func (i *Integer) Mul(o Object) Object {
	if other, ok := o.(*Integer); ok {
		i.Value = i.Value * other.Value
		return i
	}
	return NativeNull
}

func (i *Integer) Divide(o Object) Object {
	if other, ok := o.(*Integer); ok {
		i.Value = i.Value / other.Value
		return i
	}
	return NativeNull
}

func (i *Integer) Equal(o Object) *Boolean {
	var other *Integer
	var ok bool
	if other, ok = o.(*Integer); !ok {
		return NativeFalse
	}

	if i.Value == other.Value {
		return NativeTrue
	} else {
		return NativeFalse
	}
}

func (i *Integer) NotEqual(o Object) *Boolean {
	if i.Equal(o).Value {
		return NativeFalse
	} else {
		return NativeTrue
	}
}

func (i *Integer) GreaterThan(o Object) *Boolean {
	var other *Integer
	var ok bool
	if other, ok = o.(*Integer); !ok {
		return NativeFalse
	}

	if i.Value > other.Value {
		return NativeTrue
	}

	return NativeFalse
}

func (i *Integer) LessThan(o Object) *Boolean {
	if !i.Equal(o).Value && !i.GreaterThan(o).Value {
		return NativeTrue
	}
	return NativeFalse
}

func (i *Integer) Negative() Object {
	return &Integer{Value: -i.Value}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjType {
	return ObjBoolean
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Equal(o Object) *Boolean {
	var other *Boolean
	var ok bool
	if other, ok = o.(*Boolean); !ok {
		return NativeFalse
	}

	if b.Value == other.Value {
		return NativeTrue
	} else {
		return NativeFalse
	}
}

func (b *Boolean) NotEqual(o Object) *Boolean {
	if b.Equal(o).Value {
		return NativeFalse
	} else {
		return NativeTrue
	}
}

type Null struct {
}

func (n *Null) Type() ObjType {
	return ObjNull
}

func (n *Null) Inspect() string {
	return "null"
}

type Return struct {
	Object
}

func (n *Return) Type() ObjType {
	return ObjReturn
}

func (n *Return) Inspect() string {
	return n.Object.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjType {
	return ObjError
}

func (e *Error) Inspect() string {
	return e.Message
}

type Fn struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (f *Fn) Type() ObjType {
	return ObjFunction
}

func (f *Fn) Inspect() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("fn(")
	buffer.WriteString(util.AnyJoin(" ,", f.Params))
	buffer.WriteString(")")
	buffer.WriteString("{\n")
	buffer.WriteString(f.Body.String())
	buffer.WriteString("\n}")
	return buffer.String()
}

type StringObj struct {
	Value string
	Env   *Environment
}

func (s *StringObj) Type() ObjType {
	return ObjString
}

func (s *StringObj) Inspect() string {
	return s.Value
}

func (s *StringObj) Add(object Object) Object {
	s.Value = s.Value + object.Inspect()
	return s
}
