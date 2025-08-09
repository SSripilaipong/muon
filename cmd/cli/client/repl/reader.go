package repl

import (
	"fmt"

	"github.com/chzyer/readline"
)

type consoleReader struct {
	reader *readline.Instance
}

func newConsoleReader() consoleReader {
	reader, err := readline.New("µ> ")
	if err != nil {
		panic(fmt.Errorf("unexpected error while creating consoleReader: %w", err))
	}
	return consoleReader{reader: reader}
}

func (r consoleReader) ReadLine() (string, error) {
	return r.reader.Readline()
}
