package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/aeremic/cgo/evaluator"
	"github.com/aeremic/cgo/parser"
	"github.com/aeremic/cgo/tokenizer"
	"github.com/aeremic/cgo/value"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := value.NewEnvironment()

	for {
		_, err := fmt.Fprintf(out, PROMPT)
		if err != nil {
			return
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		l := scanner.Text()
		t := tokenizer.New(l)
		p := parser.New(t)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			io.WriteString(out, "Parse error:\n")
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}

			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			// io.WriteString(out, program.String())
			io.WriteString(out, evaluated.Sprintf())
			io.WriteString(out, "\n")
		}
	}
}
