package tokenizer

import (
	"testing"
)

func TestRead(t *testing.T) {
	// input := "=+(){},;"

	// expectedTokens := []struct {
	// 	expectedType    TokenType
	// 	expectedLiteral string
	// }{
	// 	{ASSIGN, "="},
	// 	{PLUS, "+"},
	// 	{LPAREN, "("},
	// 	{RPAREN, ")"},
	// 	{LBRACE, "{"},
	// 	{RBRACE, "}"},
	// 	{COMMA, ","},
	// 	{SEMICOLON, ";"},
	// 	{EOF, ""},
	// }

	input := `
		let a = 5;
		let b = 4;

		let add = fn(x, y) {
			x + y
		};

		let result = add(a, b);
	`

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "a"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "b"},
		{ASSIGN, "="},
		{INT, "4"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNC, "fn"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COMMA, ","},
		{IDENT, "y"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
	}

	tokenizer := New(input)
	for i, expectedToken := range expectedTokens {
		token := tokenizer.Read()

		if token.Type != expectedToken.expectedType {
			t.Fatalf("expectedTokens[%d] - Token type is wrong. Expected %q, received %q",
				i, expectedToken.expectedType, token.Type)
		}

		if token.Literal != expectedToken.expectedLiteral {
			t.Fatalf("expectedTokens[%d] - Token literal is wrong. Expected %q, received %q",
				i, expectedToken.expectedLiteral, token.Literal)
		}
	}
}
