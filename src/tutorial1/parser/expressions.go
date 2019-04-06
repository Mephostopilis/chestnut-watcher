package parser

import (
	"bytes"
	"tutorial1/lex"
)

type Expression interface {
	String() string
}

type IntegerLiteralExpression struct {
	Token Token
	Value int64
}

func (il *IntegerLiteralExpression) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type prefixParseFn func()
type infixParseFn func()

type Parser struct {
	l              *Lexer
	errors         []string
	curToken       Token
	peekToken      Token
	prefixParseFns map[string]prefixParseFn
	infixParseFns  map[string]infixParseFn
}

func (p *Parser) registerPrefix(tokenType string, fn prefixParseFn) {

}

func (p *Parser) registerInfix(tokenType string, fn infixParseFn) {

}

func (p *Parser) ParseExpression(precedence int) Expression {

}
