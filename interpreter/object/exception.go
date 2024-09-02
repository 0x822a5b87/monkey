package object

import "fmt"

func newWrongArgumentSizeError(actualArgumentSize, expectedArgumentSize int) Object {
	return &Error{
		Message: fmt.Sprintf("wrong number of arguments. got=%d, want=%d", actualArgumentSize, expectedArgumentSize),
	}
}

func newWrongArgumentTypeError(fnName string, actualTypeName ObjType) Object {
	return &Error{
		Message: fmt.Sprintf("argument to `%s` not supported, got %s", fnName, actualTypeName),
	}
}
