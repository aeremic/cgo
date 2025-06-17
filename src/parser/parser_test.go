package parser

import (
	"fmt"
	"testing"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/tokenizer"
)

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

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Invalid token literal. Got %s instead of let",
			statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Invalid let statement type. Got %T",
			statement.(*ast.LetStatement))
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value invalid. Got %s instead of %s",
			letStatement.Name.Value, name)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLitral() invalid. Got %s instead of %s",
			letStatement.Name.TokenLiteral(), name)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not ast.Identifier. Got %T",
			expression)

		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %s. Got %s", value, identifier.Value)

		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() is not %s. Got %s",
			value, identifier.TokenLiteral())

		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression is not ast.Identifier. Got %T",
			expression)

		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t. Got %t", value, boolean.Value)

		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("identifier.TokenLiteral() is not %t. Got %s",
			value, boolean.TokenLiteral())

		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int:
		return testIntegerLiteralExpression(t, expression, int64(value))
	case int64:
		return testIntegerLiteralExpression(t, expression, value)
	case string:
		return testIdentifier(t, expression, value)
	case bool:
		return testBooleanLiteral(t, expression, value)
	}

	t.Errorf("type of expression not handled. Got %T", expression)

	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{},
	operator string, right interface{}) bool {
	infixExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.InfixExpression. Got %T(%s)",
			expression, expression)

		return false
	}

	if !testLiteralExpression(t, infixExpression.Left, left) {
		return false
	}

	if infixExpression.Operator != operator {
		t.Errorf("Operator invalid. Got %s", operator)

		return false
	}

	if !testLiteralExpression(t, infixExpression.Right, right) {
		return false
	}

	return true
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

func setUpTest(t *testing.T, input string) *ast.ProgramRoot {
	inputTokenizer := tokenizer.New(input)
	parser := New(inputTokenizer)

	checkParseErrors(t, parser)

	return parser.ParseProgram()
}

func TestLetStatement(t *testing.T) {
	// input := `
	// let x = 5;
	// let y = 10;
	// let foobar = 12345;
	// `

	// input := `
	// let x 5;
	// let = 10;
	// let 12345;
	// `

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		program := setUpTest(t, test.input)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
				len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}

		value := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, test.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, test := range tests {
		program := setUpTest(t, test.input)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
				len(program.Statements))
		}

		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement is not ast.ReturnStatement. Got %T",
				statement.(*ast.ReturnStatement))
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Fatalf("returnStatement.TokenLiteral() is not return. Got %s",
				returnStatement.TokenLiteral())
		}

		if testLiteralExpression(t, returnStatement.ReturnValue, test.expectedValue) {
			return
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

func TestBooleanExpression(t *testing.T) {
	input := "true;"

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

	literal, ok := statement.Expression.(*ast.Boolean)
	if !ok {
		t.Errorf("Expression is not Boolean type. Got %T", statement.Expression)
	}

	if literal.Value != true {
		t.Errorf("literal.Value incorrect.")
	}

	if literal.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral incorrect. Got %s", literal.TokenLiteral())
	}
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

		if !testInfixExpression(t, expression, test.leftValue, test.operator, test.rightValue) {
			t.Fatalf("testInfixExpression failed for %T.", expression)
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
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"5 + 6 * 7", "(5 + (6 * 7))"},
		{"5 + 6 - 7", "((5 + 6) - 7)"},
		{"(5 + 6) * 7", "((5 + 6) * 7)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!true == true", "((!true) == true)"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := setUpTest(t, input)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement type. Got %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression is not IfExpression type. Got %T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements. Got %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement type. Got %T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("Expression alternative not expected. Got %v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	program := setUpTest(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("Invalid number of statements. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement type. Got %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expression is not IfExpression type. Got %T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements. Got %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement type. Got %T",
			expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative == nil {
		t.Errorf("Expression alternative expected but missing. Got %v", expression.Alternative)
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statements. Got %d", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not ExpressionStatement type. Got %T",
			expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := setUpTest(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. Got %d",
			1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement type. Got %T",
			program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not ast.FunctionLiteral. Got %T",
			statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal params wrong. Got %d instead of %d",
			len(function.Parameters), 2)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain %d statements. Got %d",
			1, len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] invalid type. Got %T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		program := setUpTest(t, test.input)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expectedParams) {
			t.Errorf("function.Parameters len invalid. Got %d instead of %d",
				len(function.Parameters), len(test.expectedParams))
		}

		for i, expectedParam := range test.expectedParams {
			testLiteralExpression(t, function.Parameters[i], expectedParam)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	program := setUpTest(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("Invalid number of statements. Got %d instead of %d",
			len(program.Statements), 1)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] invalid type. Got %T",
			program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement.Expression invalid type. Got %T",
			statement.Expression)
	}

	if !testIdentifier(t, expression.Function, "add") {
		return
	}

	if len(expression.Arguments) != 3 {
		t.Errorf("Invalid number of args in the call. Got %d instead of %d",
			len(expression.Arguments), 3)
	}

	// for _, arg := range expression.Arguments {
	// 	t.Errorf("%T %s", arg, arg.String())
	// }

	testLiteralExpression(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParametersParsing(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedArguments  []string
	}{
		{
			input:              "add();",
			expectedIdentifier: "add",
			expectedArguments:  []string{},
		},
		{
			input:              "add(1);",
			expectedIdentifier: "add",
			expectedArguments:  []string{"1"},
		},
		{
			input:              "add(1, 2 * 3, 4 + 5);",
			expectedIdentifier: "add",
			expectedArguments:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, test := range tests {
		program := setUpTest(t, test.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, test.expectedIdentifier) {
			return
		}

		if len(exp.Arguments) != len(test.expectedArguments) {
			t.Fatalf("Wrong number of arguments. want=%d, got=%d",
				len(test.expectedArguments), len(exp.Arguments))
		}

		for i, arg := range test.expectedArguments {
			if exp.Arguments[i].String() != arg {
				t.Errorf("Argument %d wrong. want=%q, got=%q",
					i, arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := setUpTest(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
			len(program.Statements))
	}

	statement := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expression is not StringLiteral type. Got %T",
			statement.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value is not %s. Got %s",
			"hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	program := setUpTest(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := statement.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("statement is not ArrayLiteral type. Got %T", statement.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array.Elements len invalid. Got %d instead of %d",
			len(array.Elements), 3)
	}

	testIntegerLiteralExpression(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "arr[1 + 1]"

	program := setUpTest(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
			len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExpression, ok := statement.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("statement.Expression is not IndexExpression type. Got %T",
			statement.Expression)
	}

	if !testIdentifier(t, indexExpression.Left, "arr") {
		t.Fatalf("invalid index expression left identifier")
		return
	}

	if !testInfixExpression(t, indexExpression.Index, 1, "+", 1) {
		t.Fatalf("invalid index expression index")
		return
	}
}

func TestParsingDictLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	program := setUpTest(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
			len(program.Statements))
	}

	statement := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := statement.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not DictLiteral. Got %T",
			statement.Expression)
	}

	if len(dict.Elements) != 3 {
		t.Errorf("invalid number of elements in dict. got %d instead of %d",
			len(dict.Elements), 3)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range dict.Elements {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("invalid element key. got %T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteralExpression(t, value, expectedValue)
	}
}

func TestParsingEmptyDictLiteral(t *testing.T) {
	input := `{}`

	program := setUpTest(t, input)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements doesn't contain 1 statements. Got %d",
			len(program.Statements))
	}

	statement := program.Statements[0].(*ast.ExpressionStatement)
	dict, ok := statement.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not DictLiteral. Got %T",
			statement.Expression)
	}

	if len(dict.Elements) != 0 {
		t.Errorf("invalid number of elements in dict. got %d instead of %d",
			len(dict.Elements), 0)
	}

}
