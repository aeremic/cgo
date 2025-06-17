package tokenizer

import (
	"testing"

	"github.com/aeremic/cgo/token"
)

func TestRead(t *testing.T) {
	// input := "=+(){},;"

	// expectedTokens := []struct {
	// 	expectedType    TokenType
	// 	expectedLiteral string
	// }{
	// 	{token.ASSIGN, "="},
	// 	{token.PLUS, "+"},
	// 	{token.LPAREN, "("},
	// 	{token.RPAREN, ")"},
	// 	{token.LBRACE, "{"},
	// 	{token.RBRACE, "}"},
	// 	{token.COMMA, ","},
	// 	{token.SEMICOLON, ";"},
	// 	{token.EOF, ""},
	// }

	input := `
		let a = 5;
		let b = 4;

		let add = fn(x, y) {
			x + y
		};

		let result = add(a, b);

		!-/*5;
		5 < 10 > 5

		if (5 < 10) {
			return true;
		} else {
			return false;
		}

		10 == 10;
		10 != 9;

		"foo \"barfoo\""

		[1, 2];

		{"foo": "bar"};
	`

	expectedTokens := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQUALS, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQUALS, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		//{token.STRING, "foobar"},
		//{token.STRING, "foo bar"},
		{token.STRING, `foo \"barfoo\"`},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
	}

	tokenizer := New(input)
	for i, expectedToken := range expectedTokens {
		parsedToken := tokenizer.NextToken()

		if parsedToken.Type != expectedToken.expectedType {
			t.Fatalf("expectedTokens[%d] - Token type is wrong. Expected %q, received %q",
				i, expectedToken.expectedType, parsedToken.Type)
		}

		if parsedToken.Literal != expectedToken.expectedLiteral {
			t.Fatalf("expectedTokens[%d] - Token literal is wrong. Expected %q, received %q",
				i, expectedToken.expectedLiteral, parsedToken.Literal)
		}
	}
}
