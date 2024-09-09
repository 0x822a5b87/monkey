package code

import (
	"fmt"
	"strconv"
	"strings"
)

func joinInt(separator string, numbers []int) string {
	s := make([]string, 0)
	for _, number := range numbers {
		s = append(s, strconv.Itoa(number))
	}
	return strings.Join(s, separator)
}

func formatInstruction(offsetBytes int, name string, operands []int) string {
	offset := fmt.Sprintf("%04d", offsetBytes)
	operandsStr := joinInt(" ", operands)
	return fmt.Sprintf("%s %s %s", offset, name, operandsStr)
}
