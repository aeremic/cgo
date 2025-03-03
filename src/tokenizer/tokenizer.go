package tokenizer

type Tokenizer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Position used for peeking after current position
	ch           byte // Current position character
}

func New(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.nextChar()

	return t
}

// Return next character and advance input position
func (t *Tokenizer) nextChar() {
	if t.readPosition >= len(t.input) {
		t.ch = 0 // Set to ASCII NUL if end is reached
	} else {
		t.ch = t.input[t.readPosition]
	}

	t.position = t.readPosition
	t.readPosition += 1
}

// Parse current character and move pointer to the next one
func (t *Tokenizer) Read() Token {
	var token Token

	t.skipWhitespaces()

	// Read and create current token which will be returned
	switch t.ch {
	case '=':
		token = Token{ASSIGN, string(t.ch)}
	case ';':
		token = Token{SEMICOLON, string(t.ch)}
	case '(':
		token = Token{LPAREN, string(t.ch)}
	case ')':
		token = Token{RPAREN, string(t.ch)}
	case '{':
		token = Token{LBRACE, string(t.ch)}
	case '}':
		token = Token{RBRACE, string(t.ch)}
	case ',':
		token = Token{COMMA, string(t.ch)}
	case '+':
		token = Token{PLUS, string(t.ch)}
	case 0:
		token = Token{EOF, ""}
	default:
		if isChLetter(t.ch) {
			token.Literal = t.readIdentifier()
			token.Type = GetKeywordByIdent(token.Literal)

			// Early return since moving char
			// since moving char pointer is not needed after readIdentifier call
			// (method itself moves a pointer)
			return token
		} else if isChDigit(t.ch) {
			token.Literal = t.readNumber()
			token.Type = INT

			// Early return since moving char
			// since moving char pointer is not needed after readNumber call
			// (method itself moves a pointer)
			return token
		}

		token = Token{ILLEGAL, string(t.ch)}
	}

	// Go to next char
	t.nextChar()

	return token
}

func (t *Tokenizer) skipWhitespaces() {
	for t.ch == ' ' || t.ch == '\t' || t.ch == '\n' || t.ch == '\r' {
		t.nextChar()
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

func isChLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isChDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
