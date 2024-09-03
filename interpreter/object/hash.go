package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Hash struct {
	Pairs map[HashKey]*HashPair
}

func (h *Hash) Type() ObjType {
	return ObjHash
}

func (h *Hash) Inspect() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("{")

	elements := make([]string, 0)
	for _, v := range h.Pairs {
		elements = append(elements, fmt.Sprintf("%s:%s",
			v.Value.Inspect(), v.Value.Inspect()))
	}

	buffer.WriteString(strings.Join(elements, ", "))
	buffer.WriteString("}")
	return buffer.String()
}

func (h *Hash) Index(object Object) Object {
	hashable, ok := object.(Hashable)
	if !ok {
		return &Error{Message: fmt.Sprintf("unusable as hash key: %s", object.Type())}
	}
	o, ok := h.Pairs[hashable.HashKey()]
	if ok {
		return o.Value
	} else {
		return NativeNull
	}
}
