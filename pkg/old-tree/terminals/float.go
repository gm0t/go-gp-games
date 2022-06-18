package terminals

import (
	"lr1Go/pkg/old-tree"
)

type FloatTerm struct {
	name   string
	getter func(args old_tree.ResolveArguments) float64
}

func (n *FloatTerm) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	cb(extractDepth(extra), n)
}

func (n *FloatTerm) String() string {
	return n.name
}

func (n *FloatTerm) Clone() old_tree.Node {
	return NewFloatTerm(n.name, n.getter)
}

func (n *FloatTerm) Resolve(args old_tree.ResolveArguments) float64 {
	return n.getter(args)
}
func NewFloatTerm(name string, getter func(args old_tree.ResolveArguments) float64) *FloatTerm {
	return &FloatTerm{getter: getter, name: name}
}
