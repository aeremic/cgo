package parser

import (
	"fmt"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/token"
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

var precedences = map[token.Type]int{
	token.EQUALS:     EQUALS,
	token.NOT_EQUALS: EQUALS,
	token.LT:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.SLASH:      PRODUCT,
	token.ASTERISK:   PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // Input represents left side of op.
)

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokenizer.NextToken()
}

func (p *Parser) checkCurrentTokenType(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) checkPeekTokenType(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekAndMove(t token.Type) bool {
	if p.checkPeekTokenType(t) {
		p.nextToken()
		return true
	}

	p.LogPeekError(t)

	return false
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("No prefix parse function found for type %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) currentTokenPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}
