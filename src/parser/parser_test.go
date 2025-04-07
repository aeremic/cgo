package parser

import (
	"fmt"
	"testing"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/tokenizer"
)

func setUpTest(t *testing.T, input string) *ast.ProgramRoot {
	inputTokenizer := tokenizer.New(input)
	parser := New(inputTokenizer)

	checkParseErrors(t, parser)

	return parser.ParseProgram()
}

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

	program := setUpTest(t, input)

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

	program := setUpTest(t, input)

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

	program := setUpTest(t, input)

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

	program := setUpTest(t, input)

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

func testIntegerLiteralExpression(t *testing.T, il ast.Expression, value int64) bool {
	integerLiteral, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Given integer literal is not type of IntegerLiteral. Got %T",
			il)
		return false
	}

	if integerLiteral.Value != value {
		t.Errorf("Wrong value given. Got %d, required is %d", integerLiteral.Value, value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integerLiteral.TokenLiteral is not %d. Got %s",
			value, integerLiteral.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, test := range prefixTests {
		program := setUpTest(t, test.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain required statements. Got %d",
				len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expression is not *ast.PrefixExpression. Got %T",
				statement.Expression)
		}

		if expression.Operator != test.operator {
			t.Fatalf("Invalid operator. Got %s and required is %s",
				expression.Operator, test.operator)
		}

		if !testIntegerLiteralExpression(t, expression.Right, test.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, test := range infixTests {
		program := setUpTest(t, test.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain required statements. Got %d",
				len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. Got %T",
				program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement.Expression is not *ast.InfixExpression. Got %T",
				statement.Expression)
		}

		if !testIntegerLiteralExpression(t, expression.Left, test.leftValue) {
			return
		}

		if expression.Operator != test.operator {
			t.Fatalf("Operator not valid. Got %s instead of %s",
				expression.Operator, test.operator)
		}

		if !testIntegerLiteralExpression(t, expression.Right, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, test := range tests {
		program := setUpTest(t, test.input)

		output := program.String()
		if output != test.expected {
			t.Errorf("Got wrong actual. Got %q; Expected %q",
				output, test.expected)
		}
	}
}
