package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

var line int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return fmt.Sprintf("%3d %s", line, strings.Repeat(traceIdentPlaceholder, traceLevel-1))
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() {
	traceLevel = traceLevel + 1
}

func decIdent() {
	traceLevel = traceLevel - 1
}

func trace(msg string) string {
	incIdent()
	tracePrint(fmt.Sprintf("BEGIN %s", msg))
	line++
	return msg
}

func untrace(msg string) {
	tracePrint(fmt.Sprintf("END %s", msg))
	line++
	decIdent()
}
