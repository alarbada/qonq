package main

import "fmt"

type Node interface{ isNode() }

type RootNode struct {
	Nodes []Node
}

type ComponentNode struct {
	Name     string
	Children []Node
}

func (ComponentNode) isNode() {}

type TextNode struct {
	Text string
}

func (TextNode) isNode() {}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (this *Parser) peek() Token {
	if this.pos >= len(this.tokens) {
		return Token{0, ""}
	}
	return this.tokens[this.pos]
}

func (this *Parser) backTimes(times int) {
	this.pos -= times
}

func (this *Parser) read() Token {
	token := this.peek()
	this.pos++
	return token
}

func (this *Parser) CheckTagOrder() error {
	stack := make([]string, 0, len(this.tokens))
	for _, token := range this.tokens {
		switch token.Kind {
		case TagNameStart:
			stack = append(stack, token.Val)
			continue
		case TagNameEnd:
			break
		default:
			continue
		}

		last := stack[len(stack)-1]
		if last == token.Val {
			stack = stack[:len(stack)-1]
		} else {
			return fmt.Errorf("Tag mismatch: %s != %s", last, token.Val)
		}
	}
	
	return nil
}

func (this *Parser) ParseComponent() (ComponentNode, error) {
	this.backTimes(1)
	t := this.peek()

	tagnameStart := t.Val
	node := ComponentNode{
		Name:     tagnameStart,
		Children: []Node{},
	}

	this.read() // consume the tagname

	children, err := this.ParseChildren()
	if err != nil {
		return node, err
	}

	node.Children = children

	tagnameEnd := this.read()

	if tagnameStart != tagnameEnd.Val {
		return node, fmt.Errorf("Tag mismatch: %s != %s", tagnameStart, tagnameEnd.Val)
	}

	return node, nil
}

func (this *Parser) ParseChildren() ([]Node, error) {
	children := []Node{}

	for {
		switch token := this.read(); token.Kind {
		case TagNameStart:
			component, err := this.ParseComponent()
			if err != nil {
				return nil, err
			}
			children = append(children, component)
		case Text:
			textNode := TextNode{
				Text: token.Val,
			}
			children = append(children, textNode)
		case TagNameEnd:
			this.backTimes(1)
			return children, nil
		default:
			return nil, fmt.Errorf("Unexpected token: %s", token)
		}

		if this.pos >= len(this.tokens) {
			return nil, fmt.Errorf("Unexpected EOF")
		}
	}
}

func (this *Parser) RootParse() (RootNode, error) {
	rootNode := RootNode{
		Nodes: []Node{},
	}

	if err := this.CheckTagOrder(); err != nil {
		return rootNode, err
	}

	for {
		switch token := this.read(); token.Kind {
		case TagNameStart:
			component, err := this.ParseComponent()
			if err != nil {
				return rootNode, err
			}
			rootNode.Nodes = append(rootNode.Nodes, component)
		case EOF:
			goto doneParsing
		default:
			return rootNode, fmt.Errorf("Unexpected token: %s", token)
		}

		if this.pos >= len(this.tokens) {
			break
		}
	}

doneParsing:
	return rootNode, nil
}

