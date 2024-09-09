package common

import "fmt"

var (
	errUnsupportedCompilingNode = errorPattern{100003, "unsupported compiling node for %s"}
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
