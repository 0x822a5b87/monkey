package util

import (
	"fmt"
	"strings"
)

func AnyJoin(sep string, objs ...any) string {
	elements := make([]string, 0)
	for _, obj := range objs {
		elements = append(elements, fmt.Sprintf("%s", obj))
	}
	return strings.Join(elements, sep)
}
