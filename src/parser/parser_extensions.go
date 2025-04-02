package parser

import (
	"fmt"

	"github.com/aeremic/cgo/token"
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokenizer.NextToken()
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function found for type %s", t)
	p.errors = append(p.errors, msg)
}
