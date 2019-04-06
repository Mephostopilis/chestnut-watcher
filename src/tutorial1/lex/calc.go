package lex

import (
	"strconv"
	_ "strings"
)

const (
	ILLEGAL  = "ILLEGAL"
	EOF      = "EOF"
	INT      = "INT"
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LPAREN   = "("
	RPAREN   = ")"
)

type Token struct {
	Type    string
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(t string, l int) Token {
	literal := strconv.Itoa(l)
	return Token{t, literal}
}

func NewLex(input string) *Lexer {
	l := new(Lexer)
	l.input = input
	l.position = 0
	l.readPosition = 0
	l.ch = 0
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case '0':
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken("", "")
		}
	}
	l.readChar()
	return tok
}
