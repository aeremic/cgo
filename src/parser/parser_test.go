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

	// input := `
	// let x 5;
	// let = 10;
	// let 12345;
	// `

	tokenizer := tokenizer.New(input)
	parser := New(tokenizer)

	checkParseErrors(t, parser)

	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements. Got %d",
			len(program.Statements))
	}

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

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 101;`

	tokenizer := tokenizer.New(input)
	parser := New(tokenizer)

	checkParseErrors(t, parser)

	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements. Got %d",
			len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Statement is not ReturnStatement type. Got %T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("Return statemen token literal is not 'return'. Got %q",
				returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	tokenizer := tokenizer.New(input)
	parser := New(tokenizer)

	checkParseErrors(t, parser)

	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Incorrect number of statements. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not expression statement. Got %T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression is not identifier. Got %T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Fatalf("Identifier invalid. Got %s", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Fatalf("Token literal invalid. Got %s", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	tokenizer := tokenizer.New(input)
	parser := New(tokenizer)

	checkParseErrors(t, parser)

	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Incorrect number of statements. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not expression statement. Got %T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expression is not IntegerLiteral type. Got %T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value incorrect. Got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral incorrect. Got %s", literal.TokenLiteral())
	}
}
