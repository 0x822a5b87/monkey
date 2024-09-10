package common

import "fmt"

var (
	errUnsupportedCompilingNode  = errorPattern{100003, "unsupported compiling node for %s"}
	errEmptyStack                = errorPattern{100004, "the stack is empty, cannot do pop for %s"}
	errTypeMismatch              = errorPattern{100005, "type mismatch : expect [%s], actual [%s]"}
	errUnsupportedBinaryOperator = errorPattern{100006, "unsupported binary operator for %s"}
)

type errorPattern struct {
	code            ErrorCode
	errorMsgPattern string
}

func (e errorPattern) format(a ...any) ErrorInfo {
	return ErrorInfo{
		Code:     e.code,
		ErrorMsg: fmt.Sprintf(e.errorMsgPattern, a),
	}
}
