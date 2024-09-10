package repl

import (
	"0x822a5b87/monkey/compiler/compiler"
	"0x822a5b87/monkey/compiler/vm"
	"0x822a5b87/monkey/interpreter/ast"
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
const Interpreter = "i"
const Compiler = "c"

func Start(typed string, in io.Reader, out io.Writer) {
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
			switch typed {
			case Interpreter:
				interpret(out, stmt, env)
			case Compiler:
				compile(out, stmt)
			}
		}
	}
}

func compile(out io.Writer, stmt ast.Statement) {
	c := compiler.NewCompiler()
	err := c.Compile(stmt)
	if err != nil {
		silentWrite(out, err.Error())
		silentWrite(out, "\n")
		return
	}
	v := vm.NewVm(c.ByteCode())
	err = v.Run()
	if err != nil {
		silentWrite(out, err.Error())
		silentWrite(out, "\n")
		return
	}

	stackTop := v.StackTop()

	silentWrite(out, stackTop.Inspect())
	silentWrite(out, "\n")
}

func interpret(out io.Writer, stmt ast.Statement, env *object.Environment) {
	obj := evaluator.Eval(stmt, env)
	_, err := io.WriteString(out, obj.Inspect())
	fmt.Println()
	if err != nil {
		fmt.Println(err.Error())
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

func silentWrite(out io.Writer, msg string) {
	_, _ = io.WriteString(out, msg)
}
