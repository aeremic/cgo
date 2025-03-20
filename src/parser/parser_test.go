package parser

import (
	"testing"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/tokenizer"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 12345;
	`

	tokenizer := tokenizer.New(input)
	parser := New(tokenizer)

	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements. Got %d",
			len(program.Statements))
	}

	checkParseErrors(t, parser)

	expectedLetIdentifiers := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range expectedLetIdentifiers {
		statement := program.Statements[i]
		if !checkLetStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

func checkParseErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %s", msg)
	}

	t.FailNow()
}

func checkLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Statement is not 'let'. Got %q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not *ast.LetStatement. Got %T",
			statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value is not %s. Got %s",
			name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral is not %s. Got %s",
			name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}
