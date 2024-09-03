package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjType {
	return ObjArray
}

func (a *Array) Inspect() string {
	buffer := bytes.Buffer{}
	elements := make([]string, 0)
	for _, element := range a.Elements {
		elements = append(elements, element.Inspect())
	}
	buffer.WriteString("[")
	buffer.WriteString(strings.Join(elements, ", "))
	buffer.WriteString("]")
	return buffer.String()
}

func (a *Array) Index(o Object) Object {
	other, ok := o.(*Integer)
	if !ok {
		return NativeNull
	}
	if other.Value >= int64(len(a.Elements)) || other.Value < 0 {
		return NativeNull
	}
	return a.Elements[other.Value]
}

func (a *Array) First() Object {
	if a.Len().Value == 0 {
		return NativeNull
	}
	return a.Index(&Integer{Value: 0})
}

func (a *Array) Last() Object {
	if a.Len().Value == 0 {
		return NativeNull
	}
	return a.Index(&Integer{Value: a.Len().Value - 1})
}

func (a *Array) Len() Integer {
	return Integer{Value: int64(len(a.Elements))}
}

func (a *Array) Rest() Object {
	length := len(a.Elements)
	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, a.Elements[1:length])
		restArray := &Array{Elements: newElements}
		return restArray
	}
	return NativeNull
}

func (a *Array) Push(obj Object) Object {
	length := len(a.Elements)
	newElements := make([]Object, length)
	copy(newElements, a.Elements)
	newElements = append(newElements, obj)
	newArray := &Array{Elements: newElements}
	return newArray
}
