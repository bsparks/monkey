package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"bytes"

	"github.com/bsparks/monkey/evaluator"
	"github.com/bsparks/monkey/lexer"
	"github.com/bsparks/monkey/object"
	"github.com/bsparks/monkey/parser"
	"github.com/bsparks/monkey/repl"
)

func readSource(filename string) ([]byte, error) {
	if filename == "" || filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

type MonkeyError struct {
	Messages []string
}

func (me *MonkeyError) Error() string {
	var out bytes.Buffer

	for _, msg := range me.Messages {
		out.WriteString("\t" + msg + "\n")
	}

	return out.String()
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		// REPL mode
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Hello %s welcome to the Monkey programming language.\n", user.Username)
		fmt.Println("Type code:")
		repl.Start(os.Stdin, os.Stdout)
	} else {
		err := func() error {
			src, err := readSource(flag.Arg(0))
			if err != nil {
				return err
			}

			context := object.NewEnvironment()
			lex := lexer.New(string(src))
			par := parser.New(lex)
			program := par.ParseProgram()

			if len(par.Errors()) != 0 {
				err := &MonkeyError{Messages: par.Errors()}
				return err
			}

			evaluated := evaluator.Eval(program, context)
			if evaluated != nil {
				os.Stdout.WriteString(evaluated.Inspect())
			}

			return err
		}()
		if err != nil {
			switch err := err.(type) {
			case (*MonkeyError):
				fmt.Print(err)
			default:
				fmt.Println(err)
			}
			os.Exit(64)
		}
	}
}
