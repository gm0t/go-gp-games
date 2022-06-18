package terminals

import (
	"lr1Go/pkg/old-tree"
)

type BooleanTerm struct {
	name   string
	getter func(args old_tree.ResolveArguments) bool
}

func (n *BooleanTerm) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	cb(extractDepth(extra), n)
}

func (n *BooleanTerm) String() string {
	return n.name
}

func (n *BooleanTerm) Clone() old_tree.Node {
	return NewBooleanTerm(n.name, n.getter)
}

func (n *BooleanTerm) Resolve(args old_tree.ResolveArguments) bool {
	return n.getter(args)
}

func NewBooleanTerm(name string, getter func(args old_tree.ResolveArguments) bool) *BooleanTerm {
	return &BooleanTerm{getter: getter, name: name}
}
