package main

import (
	"fmt"
	"strings"
)

type TokenKind uint8

const (
	// Special tokens
	Illegal TokenKind = iota
	EOF
	TagNameStart
	TagNameEnd
	PropName
	PropValue
	Text
)

func (this TokenKind) String() string {
	switch this {
	case Illegal:
		return "Illegal"
	case EOF:
		return "EOF"
	case TagNameStart:
		return "TagNameStart"
	case TagNameEnd:
		return "TagNameEnd"
	case PropName:
		return "PropName"
	case PropValue:
		return "PropValue"
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

func (this *Lexer) skipWhitespace() {
	for this.peek() == ' ' || this.peek() == '\t' || this.peek() == '\n' || this.peek() == '\r' {
		this.read()
	}
}

func (this *Lexer) indexRune(r rune) int {
	return strings.IndexRune(this.input[this.pos:], r)
}

func (this *Lexer) NextToken() Token {
	this.skipWhitespace()

	if ch := this.peek(); ch == '<' {
		if this.readTimes(3) == "<q-" {
			posEndTag := this.indexRune('>')
			start := this.pos
			this.pos += posEndTag
			end := this.pos
			this.pos++ // skip >

			return Token{TagNameStart, this.input[start:end]}
		} else {
			this.backTimes(3)
		}

		if this.readTimes(4) == "</q-" {
			posEndTag := this.indexRune('>')
			start := this.pos
			this.pos += posEndTag
			end := this.pos
			this.pos++ // skip >

			return Token{TagNameEnd, this.input[start:end]}
		} else {
			this.backTimes(4)
		}

		this.read() // skip <
		nextTokenPos := this.indexRune('<')
		this.backTimes(1)
		if nextTokenPos != -1 {
			start := this.pos
			this.pos += nextTokenPos
			this.pos++
			end := this.pos

			return Token{Text, this.input[start:end]}
		}

		return Token{EOF, ""}
	}

	nextTokenPos := this.indexRune('<')
	if nextTokenPos != -1 {
		start := this.pos
		this.pos += nextTokenPos
		this.pos++
		end := this.pos

		return Token{Text, this.input[start:end]}
	}

	return Token{EOF, ""}
}

var contentsSimple = `
<q-item>
	<p>hello qonq</p>
</q-item>
`

var contentsAdvanced = `
<q-item>
	<q-prop name="name"></q-prop>
	<q-prop name="isUnavailable" type="boolean"></q-prop>
	<q-prop name="isAwesome" type="boolean" optional></q-prop>

	<li q-if="isUnavailable">
		{{ name }} - Unavailable
	</li>
	<li q-else>
		{{ name }}
	</li>

</q-item>
`

func main() {
	l := NewLexer(contentsSimple)

	for {
		tok := l.NextToken()
		if tok.Kind == EOF {
			break
		}
		fmt.Println(tok.Kind, tok.Val)
	}
}
