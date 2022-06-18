package generator

import (
	"lr1Go/pkg/tree"
)

func NewFunction(k string, t tree.NodeType, children []*tree.Node) *tree.Node {
	return &tree.Node{
		Key:      k,
		Type:     t,
		Children: children,
	}
}

func NewIf(condition *tree.Node, success *tree.Node, fail *tree.Node) *tree.Node {
	return NewFunction("if", tree.IF, []*tree.Node{condition, success, fail})
}

func NewComparison(key string) *tree.Node {
	return NewTerminal(key, tree.)
}

func NewAction(key string) *tree.Node {
	return NewTerminal(key, tree.Action)
}
