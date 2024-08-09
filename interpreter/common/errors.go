package common

import (
	"fmt"
)

var (
	ErrUnknownToken = ErrorInfo{100000, "unknown token"} // encounter unknown token
	ErrSyntax       = ErrorInfo{100001, "syntax error"}
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
