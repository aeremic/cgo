package parser

import (
	"fmt"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/token"
	"github.com/aeremic/cgo/tokenizer"
)

type Parser struct {
	tokenizer    *tokenizer.Tokenizer // Pointer to the lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
}

// Constructor
func New(t *tokenizer.Tokenizer) *Parser {
	p := &Parser{tokenizer: t, errors: []string{}}

	// Call nextToken two times to initialize
	// both current token and next token
	p.nextToken()
	p.nextToken()

	return p
}

// Methods

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) LogPeekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token %s. Got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokenizer.NextToken()
}

func (p *Parser) ParseProgram() *ast.ProgramRoot {
	program := &ast.ProgramRoot{}
	program.Statements = []ast.Statement{}

	for !p.checkCurrentTokenType(token.EOF) {
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.peekAndMove(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.peekAndMove(token.ASSIGN) {
		return nil
	}

	// Skipping expression part
	for !p.checkCurrentTokenType(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// Skipping expression part
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

func (p *Parser) peekAndMove(t token.TokenType) bool {
	if p.checkPeekTokenType(t) {
		p.nextToken()
		return true
	}

	p.LogPeekError(t)

	return false
}
