package main

import (
	"github.com/davecgh/go-spew/spew"
)

var contentsSimple = `<q-item>
	<p>hello qonq</p>
</q-item>`

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

func init() {
	spew.Config.DisableCapacities = true
	spew.Config.Indent = "    "
}

func main() {
	l := NewLexer(contentsSimple)
	tokens := l.Lex()

	spew.Dump(tokens)

	p := NewParser(tokens)

	ast, err := p.RootParse()
	if err != nil {
		panic(err)
	}

	spew.Dump(ast)
}
