package common

import (
	"0x822a5b87/monkey/interpreter/object"
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

func NewErrUnsupportedBinaryExpr(nodeName string) error {
	return errUnsupportedBinaryOperator.format(nodeName)
}

func NewErrEmptyStack(nodeName string) error {
	return errEmptyStack.format(nodeName)
}

func NewErrTypeMismatch(expectType, actualType string) error {
	return errTypeMismatch.format(expectType, actualType)
}

func NewErrOperandsCount(expectedCount, actualCount int) error {
	return errOperandsCount.format(expectedCount, actualCount)
}

func NewUnresolvedVariable(name string) error {
	return errUnresolvedVariable.format(name)
}

func NewErrIndex(name object.ObjType) error {
	return errUnresolvedVariable.format(name)
}

func NewOpcodeUndefined(op byte) error {
	return errOpCodeUndefined.format(op)
}

func NewOperandWidthError(operandsCount int) error {
	return errOperandWidth.format(operandsCount)
}
