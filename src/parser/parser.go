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
	program := &ast.ProgramRoot{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.checkPeekTokenAndMove(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.checkPeekTokenAndMove(token.ASSIGN) {
		return nil
	}

	for !p.checkCurrentTokenType(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) checkCurrentTokenType(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) checkPeekTokenType(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) checkPeekTokenAndMove(t token.TokenType) bool {
	if p.checkPeekTokenType(t) {
		p.nextToken()
		return true
	}

	return false
}
