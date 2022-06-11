package terminals

import (
	"lr1Go/pkg/tree"
)

type BooleanTerm struct {
	name   string
	getter func(args tree.ResolveArguments) bool
}

func (n *BooleanTerm) Dfs(cb func(depth int, n tree.Node), extra ...int) {
	cb(extractDepth(extra), n)
}

func (n *BooleanTerm) String() string {
	return n.name
}

func (n *BooleanTerm) Clone() tree.Node {
	return NewBooleanTerm(n.name, n.getter)
}

func (n *BooleanTerm) Resolve(args tree.ResolveArguments) bool {
	return n.getter(args)
}

func NewBooleanTerm(name string, getter func(args tree.ResolveArguments) bool) *BooleanTerm {
	return &BooleanTerm{getter: getter, name: name}
}
