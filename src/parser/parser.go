package parser

import (
	"fmt"
	"strconv"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/token"
	"github.com/aeremic/cgo/tokenizer"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // Input represents left side of op.
)

type Parser struct {
	tokenizer      *tokenizer.Tokenizer // Pointer to the lexer
	currentToken   token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
}

// Constructor
func New(t *tokenizer.Tokenizer) *Parser {
	p := &Parser{tokenizer: t, errors: []string{}}

	// Call nextToken two times to initialize
	// both current token and next token
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// p.infixParseFns = make(map[token.TokenType]infixParseFn)

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
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	statement.Expression = p.parseExpression(LOWEST)

	if p.checkPeekTokenType(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExpression := prefix()

	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("Could not parse %q as integer.", p.currentToken.Literal)
		p.errors = append(p.errors, message)

		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}
