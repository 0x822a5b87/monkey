package common

import (
	"fmt"
)

var (
	ErrUnknownToken            = ErrorInfo{100000, "unknown token"} // encounter unknown token
	ErrSyntax                  = ErrorInfo{100001, "syntax error"}
	ErrUnknownTypeOfExpression = ErrorInfo{100002, "unknown type of expression"}
)

type ErrorCode int // ErrorCode 错误码

type ErrorInfo struct {
	Code     ErrorCode
	ErrorMsg string
}

func (e ErrorInfo) Error() string {
	return e.ErrorMsg
}

func (e ErrorCode) Sting() string {
	return fmt.Sprintf("%d", e)
}

func NewErrUnsupportedCompilingNode(nodeName string) error {
	return errUnsupportedCompilingNode.format(nodeName)
}

func NewErrEmptyStack(nodeName string) error {
	return errEmptyStack.format(nodeName)
}

func NewErrTypeMismatch(expectType, actualType string) error {
	return errTypeMismatch.format(expectType, actualType)
}
