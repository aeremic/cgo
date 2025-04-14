package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/aeremic/cgo/token"
	"github.com/aeremic/cgo/tokenizer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, err := fmt.Fprintf(out, PROMPT)
		if err != nil {
			return
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		t := tokenizer.New(line)
		for parsedToken := t.NextToken(); parsedToken.Type != token.EOF; parsedToken = t.NextToken() {
			_, err := fmt.Fprintf(out, "%+v\n", parsedToken)
			if err != nil {
				return
			}
		}
	}
}
