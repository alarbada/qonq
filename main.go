package main

type TokenKind uint8

const (
	// Special tokens
	Illegal TokenKind = iota
	EOF
	TagNameStart
	TagNameEnd
	PropName
	PropValue
)

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

func (this *Lexer) readTimes(times int) string {
	s := this.input[this.pos:this.pos+times]
	this.pos += times
	return s
}


func (this *Lexer) skipWhitespace() {
	for this.peek() == ' ' || this.peek() == '\t' || this.peek() == '\n' || this.peek() == '\r' {
		this.read()
	}
}

func (this *Lexer) NextToken() Token {
	this.skipWhitespace()

	if ch := this.peek(); ch == '<' {
		if this.readTimes(2) == "q-" {

			// posSpace := strings.IndexRune(this.input[this.pos:], ' ')
			// posEndTag := strings.IndexRune(this.input[this.pos:], '>')


		}
	}

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
		println(tok.Kind, tok.Val)
	}
}
