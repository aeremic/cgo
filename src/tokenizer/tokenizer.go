package tokenizer

import (
	"github.com/aeremic/cgo/token"
)

type Tokenizer struct {
	input        string
	position     int  // Current position in input
	nextPosition int  // Position used for peeking after current position
	ch           byte // Current position character
}

// New Constructor
func New(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.nextChar()

	return t
}

// Methods

// Return next character and advance input position
func (t *Tokenizer) nextChar() {
	if t.nextPosition >= len(t.input) {
		t.ch = 0 // Set to ASCII NUL if end is reached
	} else {
		t.ch = t.input[t.nextPosition]
	}

	t.position = t.nextPosition
	t.nextPosition += 1
}

// NextToken Parse current character and move pointer to the next one
func (t *Tokenizer) NextToken() token.Token {
	var parsedToken token.Token

	t.skipWhitespaces()

	// Read and create current token which will be returned
	switch t.ch {
	case '=':
		peekedChar := t.peekChar()
		if peekedChar == '=' {
			parsedToken = token.Token{Type: token.EQUALS,
				Literal: string(t.ch) + string(peekedChar)}
			t.nextChar()
		} else {
			parsedToken = token.Token{Type: token.ASSIGN, Literal: string(t.ch)}
		}
	case ';':
		parsedToken = token.Token{Type: token.SEMICOLON, Literal: string(t.ch)}
	case ':':
		parsedToken = token.Token{Type: token.COLON, Literal: string(t.ch)}
	case '(':
		parsedToken = token.Token{Type: token.LPAREN, Literal: string(t.ch)}
	case ')':
		parsedToken = token.Token{Type: token.RPAREN, Literal: string(t.ch)}
	case '{':
		parsedToken = token.Token{Type: token.LBRACE, Literal: string(t.ch)}
	case '}':
		parsedToken = token.Token{Type: token.RBRACE, Literal: string(t.ch)}
	case '[':
		parsedToken = token.Token{Type: token.LBRACKET, Literal: string(t.ch)}
	case ']':
		parsedToken = token.Token{Type: token.RBRACKET, Literal: string(t.ch)}
	case ',':
		parsedToken = token.Token{Type: token.COMMA, Literal: string(t.ch)}
	case '+':
		parsedToken = token.Token{Type: token.PLUS, Literal: string(t.ch)}
	case '-':
		parsedToken = token.Token{Type: token.MINUS, Literal: string(t.ch)}
	case '!':
		peekedChar := t.peekChar()
		if peekedChar == '=' {
			parsedToken = token.Token{Type: token.NOT_EQUALS,
				Literal: string(t.ch) + string(peekedChar)}
			t.nextChar()
		} else {
			parsedToken = token.Token{Type: token.BANG, Literal: string(t.ch)}
		}
	case '/':
		parsedToken = token.Token{Type: token.SLASH, Literal: string(t.ch)}
	case '*':
		parsedToken = token.Token{Type: token.ASTERISK, Literal: string(t.ch)}
	case '<':
		parsedToken = token.Token{Type: token.LT, Literal: string(t.ch)}
	case '>':
		parsedToken = token.Token{Type: token.GT, Literal: string(t.ch)}
	case 0:
		parsedToken = token.Token{Type: token.EOF, Literal: ""}
	case '"':
		parsedToken = token.Token{
			Type:    token.STRING,
			Literal: t.readString(),
		}
	default:
		if isChLetter(t.ch) {
			parsedToken.Literal = t.readIdentifier()
			parsedToken.Type = token.GetKeywordByIdent(parsedToken.Literal)

			// Early return since moving char
			// since moving char pointer is not needed after readIdentifier call
			// (method itself moves a pointer)
			return parsedToken
		} else if isChDigit(t.ch) {
			parsedToken.Literal = t.readNumber()
			parsedToken.Type = token.INT

			// Early return since moving char
			// since moving char pointer is not needed after readNumber call
			// (method itself moves a pointer)
			return parsedToken
		}

		parsedToken = token.Token{Type: token.ILLEGAL, Literal: string(t.ch)}
	}

	// Go to next char
	t.nextChar()

	return parsedToken
}

func (t *Tokenizer) skipWhitespaces() {
	for t.ch == ' ' || t.ch == '\t' || t.ch == '\n' || t.ch == '\r' {
		t.nextChar()
	}
}

func (t *Tokenizer) peekChar() byte {
	if t.nextPosition >= len(t.input) {
		return 0
	} else {
		return t.input[t.nextPosition]
	}
}

func (t *Tokenizer) readIdentifier() string {
	initialPosition := t.position
	for isChLetter(t.ch) {
		t.nextChar()
	}

	return t.input[initialPosition:t.position]
}

func (t *Tokenizer) readNumber() string {
	initialPosition := t.position
	for isChDigit(t.ch) {
		t.nextChar()
	}

	return t.input[initialPosition:t.position]
}

func (t *Tokenizer) readString() string {
	position := t.position + 1
	for {
		if t.ch == '\\' && t.input[t.position+1] == '"' {
			t.nextChar()
			t.nextChar()
		} else {
			t.nextChar()
		}

		if t.ch == '"' || t.ch == 0 {
			break
		}
	}

	return t.input[position:t.position]
}

func isChLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isChDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
