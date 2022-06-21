package tree

import (
	"math/rand"
)

func NewFunction(k string, t NodeType, children []*Node) *Node {
	return &Node{
		Key:      k,
		Type:     t,
		Children: children,
	}
}

func NewIf(condition *Node, success *Node, fail *Node) *Node {
	return NewFunction("if", IF, []*Node{condition, success, fail})
}
func NewTerminal(k string, t NodeType) *Node {
	return &Node{
		Key:      k,
		Type:     t,
		Children: nil,
	}
}

func NewFloat(key string) *Node {
	return NewTerminal(key, Float)
}

func NewBoolean(key string) *Node {
	return NewTerminal(key, Boolean)
}

func NewAction(key string) *Node {
	return NewTerminal(key, Action)
}

type Generator interface {
	FFunc() *Node
	BFunc() *Node
	AFunc() *Node
	FTerm() *Node
	BTerm() *Node
	ATerm() *Node
	Term(t NodeType) *Node
	FTree(depth int) *Node
	ATree(depth int) *Node
	BTree(depth int) *Node
	Tree(t NodeType, depth int) *Node
	Grow(node *Node)
	Truncate(node *Node)
}

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}

type BasicGenerator struct {
	floats   []string
	booleans []string
	actions  []string
}

func (g *BasicGenerator) Tree(nodeType NodeType, depth int) *Node {
	switch nodeType {
	case Boolean:
		return g.BTree(depth)
	case Float:
		return g.FTree(depth)
	case Action:
		return g.ATree(depth)
	}
	panic("Unknown nodeType: " + nodeType)
}

func (g *BasicGenerator) Grow(node *Node) {
	switch node.Type {
	case Float:
		update(node, g.FFunc())
		return
	case IF:
		update(node, NewIf(g.BFunc(), g.FFunc(), g.FFunc()))
	case Plus:
		fallthrough
	case Minus:
		fallthrough
	case Multiply:
		fallthrough
	case Divide:
		fallthrough
	case Eq:
		fallthrough
	case Gt:
		fallthrough
	case Lt:
		node.Children = []*Node{g.FFunc(), g.FFunc()}
		return
	}
}

func (g *BasicGenerator) Truncate(node *Node) {
	switch node.Type {
	case Float:
		fallthrough
	case Boolean:
		fallthrough
	case Action:
		// we can't truncate leaf nodes
		return
	case IF:
		update(node, g.Term(Float))
	case Plus:
		fallthrough
	case Minus:
		fallthrough
	case Multiply:
		fallthrough
	case Divide:
		update(node, g.FTerm())
		return
	case Eq:
		fallthrough
	case Gt:
		fallthrough
	case Lt:
		update(node, g.BTerm())
		return
	}
}

func update(dst *Node, source *Node) {
	dst.Type = source.Type
	dst.Children = source.Children
	dst.Key = source.Key
}

func (g *BasicGenerator) FFunc() *Node {
	return NewIf(
		g.BTerm(),
		g.FTerm(),
		g.FTerm(),
	)
}

func (g *BasicGenerator) BFunc() *Node {
	nType := getRandom(comparisonNodeTypes)
	return &Node{
		Key:      nType,
		Type:     nType,
		Children: []*Node{g.FTerm(), g.FTerm()},
	}
}

func (g *BasicGenerator) AFunc() *Node {
	return NewIf(
		g.BTerm(),
		g.ATerm(),
		g.ATerm(),
	)
}

func (g *BasicGenerator) FTerm() *Node {
	return NewFloat(getRandom(g.floats))
}

func (g *BasicGenerator) BTerm() *Node {
	if len(g.booleans) < 1 {
		// if there are no boolean keys in the state
		// then we have to simulate them using boolean functions
		return g.BFunc()
	}
	return NewBoolean(getRandom(g.booleans))
}

func (g *BasicGenerator) ATerm() *Node {
	return NewAction(getRandom(g.actions))
}

func (g *BasicGenerator) FTree(maxDepth int) *Node {
	root := g.FFunc()
	growTree(root, g, maxDepth)

	return root
}

func growTree(root *Node, g Generator, maxDepth int) {
	Dfs(root, func(node *Node, depth int) {
		if depth < maxDepth-1 {
			g.Grow(node)
		}
	})
}

func (g *BasicGenerator) ATree(depth int) *Node {
	//TODO implement me
	panic("implement me")
}

func (g *BasicGenerator) BTree(depth int) *Node {
	//TODO implement me
	panic("implement me")
}

func (g *BasicGenerator) Term(nodeType NodeType) *Node {
	switch nodeType {
	case Boolean:
		return g.BTerm()
	case Float:
		return g.FTerm()
	case Action:
		return g.ATerm()
	}
	panic("Unknown nodeType: " + nodeType)
}

func NewGenerator(booleans []string, floats []string, actions []string) Generator {
	return &BasicGenerator{
		floats:   floats,
		booleans: booleans,
		actions:  actions,
	}
}
