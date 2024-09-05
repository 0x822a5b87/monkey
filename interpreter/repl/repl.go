package repl

import (
	"0x822a5b87/monkey/interpreter/evaluator"
	"0x822a5b87/monkey/interpreter/lexer"
	"0x822a5b87/monkey/interpreter/object"
	"0x822a5b87/monkey/interpreter/parser"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment(nil)

	for {
		fmt.Print(PROMPT)

		sourceCode := readSourceCode(scanner)
		if len(sourceCode) == 0 {
			return
		}

		l := lexer.NewLexer(sourceCode)

		p := parser.NewParser(*l)
		program := p.ParseProgram()
		for _, stmt := range program.Statements {
			obj := evaluator.Eval(stmt, env)
			_, err := io.WriteString(out, obj.Inspect())
			fmt.Println()
			if err != nil {
				fmt.Printf(err.Error())
			}
		}
	}
}

func readSourceCode(scanner *bufio.Scanner) string {
	buffer := bytes.Buffer{}
	for {
		scan := scanner.Scan()
		if !scan {
			return ""
		}
		line := scanner.Text()
		if strings.HasSuffix(line, "\\") {
			buffer.WriteString(line[:len(line)-1])
			buffer.WriteString("\n")
			continue
		}
		buffer.WriteString(line)
		break
	}
	return buffer.String()
}
