package main

import (
	"0x822a5b87/monkey/interpreter/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the compiler of Monkey Programming language!\n", u.Username)

	repl.Start("c", os.Stdin, os.Stdout)
}
