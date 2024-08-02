package repl

import (
	"0x822a5b87/monkey/lexer"
	"0x822a5b87/monkey/token"
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

		for nextToken, err := l.NextToken(); err == nil && nextToken.Type != token.EOF; nextToken, err = l.NextToken() {
			fmt.Printf("%+v\n", nextToken)
		}
	}
}
