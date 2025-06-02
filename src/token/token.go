package token

type Type string

const (
	// Special types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN     = "="
	PLUS       = "+"
	MINUS      = "-"
	BANG       = "!"
	ASTERISK   = "*"
	SLASH      = "/"
	LT         = "<"
	GT         = ">"
	EQUALS     = "=="
	NOT_EQUALS = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNC   = "FUNC"
	LET    = "LET"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
)

type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fn":     FUNC,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func GetKeywordByIdent(ident string) Type {
	if tokenType, exists := keywords[ident]; exists {
		return tokenType
	}

	return IDENT
}
