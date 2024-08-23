package repl

import (
	"0x822a5b87/monkey/evaluator"
	"0x822a5b87/monkey/lexer"
	"0x822a5b87/monkey/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		scan := scanner.Scan()
		if !scan {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)

		p := parser.NewParser(*l)
		program := p.ParseProgram()
		for _, stmt := range program.Statements {
			obj := evaluator.Eval(stmt)
			_, err := io.WriteString(out, obj.Inspect())
			fmt.Println()
			if err != nil {
				fmt.Printf(err.Error())
			}
		}
	}
}
