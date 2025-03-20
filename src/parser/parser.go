package parser

import (
	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/token"
	"github.com/aeremic/cgo/tokenizer"
)

type Parser struct {
	tokenizer    *tokenizer.Tokenizer // Pointer to the lexer
	currentToken token.Token
	peekToken    token.Token
}

// Constructor
func New(t *tokenizer.Tokenizer) *Parser {
	p := &Parser{tokenizer: t}

	// Call nextToken two times to initialize
	// both current token and next token
	p.nextToken()
	p.nextToken()

	return p
}

// Methods

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokenizer.NextToken()
}

func (p *Parser) ParseProgram() *ast.ProgramRoot {
	return nil
}
