package main

import (
	"strings"
)

type TokenKind uint8

const (
	EOF TokenKind = iota
	TagNameStart
	TagNameEnd
	// will be implemented later
	// PropName
	// PropValue
	Text
)

func (this TokenKind) String() string {
	switch this {
	case EOF:
		return "EOF"
	case TagNameStart:
		return "TagNameStart"
	case TagNameEnd:
		return "TagNameEnd"
	// case PropName:
	// 	return "PropName"
	// case PropValue:
	// 	return "PropValue"
	case Text:
		return "Text"
	default:
		return "Unknown"
	}
}

type Token struct {
	Kind TokenKind
	Val  string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (this *Lexer) peek() byte {
	if this.pos >= len(this.input) {
		return 0
	}
	return this.input[this.pos]
}

func (this *Lexer) read() string {
	ch := this.input[this.pos]
	this.pos++
	return string(ch)
}

func (this *Lexer) backTimes(times int) {
	this.pos -= times
}

func (this *Lexer) readTimes(times int) string {
	s := this.input[this.pos : this.pos+times]
	this.pos += times
	return s
}

func (this *Lexer) indexRune(r rune) int {
	return strings.IndexRune(this.input[this.pos:], r)
}

func (this *Lexer) indexExact(pattern string) int {
	return strings.Index(this.input[this.pos:], pattern)
}

func (this *Lexer) NextToken() Token {
	if nextPos := this.indexExact("<q-"); nextPos != -1 {
		if nextPos != 0 {
			start := this.pos
			this.pos += nextPos
			end := this.pos

			return Token{Text, this.input[start:end]}
		}

		this.readTimes(3)

		start := this.pos
		posEndTag := this.indexRune('>')
		this.pos += posEndTag
		end := this.pos
		this.pos++ // skip >

		return Token{TagNameStart, this.input[start:end]}
	}

	if nextPos := this.indexExact("</q-"); nextPos != -1 {
		if nextPos != 0 {
			start := this.pos
			this.pos += nextPos
			end := this.pos

			return Token{Text, this.input[start:end]}
		}

		this.readTimes(4)

		start := this.pos
		posEndTag := this.indexRune('>')
		this.pos += posEndTag
		end := this.pos
		this.pos++ // skip >

		return Token{TagNameEnd, this.input[start:end]}
	}

	if this.pos < len(this.input) {
		start := this.pos
		this.pos = len(this.input)
		end := this.pos

		return Token{Text, this.input[start:end]}
	}

	return Token{EOF, ""}
}

func (this *Lexer) Lex() []Token {
	var tokens []Token
	for {
		tok := this.NextToken()
		if tok.Kind == EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens
}
