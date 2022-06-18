package generator

import (
	"lr1Go/pkg/tree"
)

func NewTerminal(k string, t tree.NodeType) *tree.Node {
	return &tree.Node{
		Key:      k,
		Type:     t,
		Children: nil,
	}
}

func NewFloat(key string) *tree.Node {
	return NewTerminal(key, tree.Float)
}

func NewBoolean(key string) *tree.Node {
	return NewTerminal(key, tree.Boolean)
}

func NewAction(key string) *tree.Node {
	return NewTerminal(key, tree.Action)
}
