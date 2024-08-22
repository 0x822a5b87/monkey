package object

import (
	"fmt"
	"strconv"
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

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjType {
	return ObjBoolean
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct {
}

func (n *Null) Type() ObjType {
	return ObjNull
}

func (n *Null) Inspect() string {
	return "null"
}
