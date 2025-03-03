package tokenizer

type TokenType string

const (
	// Special types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "ASSING"
	PLUS   = "PLUS"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNC = "FUNC"
	LET  = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNC,
	"let": LET,
}

func GetKeywordByIdent(ident string) TokenType {
	if tokenType, exists := keywords[ident]; exists {
		return tokenType
	}

	return IDENT
}
