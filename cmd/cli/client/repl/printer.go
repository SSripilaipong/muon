package repl

import "fmt"

type consolePrinter struct{}

func newConsolePrinter() consolePrinter {
	return consolePrinter{}
}

func (consolePrinter) Print(s string) { fmt.Println(s) }
