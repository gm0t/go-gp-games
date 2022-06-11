package terminals

import (
	"lr1Go/pkg/tree"
)

type FloatTerm struct {
	name   string
	getter func(args tree.ResolveArguments) float64
}

func (n *FloatTerm) Dfs(cb func(depth int, n tree.Node), extra ...int) {
	cb(extractDepth(extra), n)
}

func (n *FloatTerm) String() string {
	return n.name
}

func (n *FloatTerm) Clone() tree.Node {
	return NewFloatTerm(n.name, n.getter)
}

func (n *FloatTerm) Resolve(args tree.ResolveArguments) float64 {
	return n.getter(args)
}
func NewFloatTerm(name string, getter func(args tree.ResolveArguments) float64) *FloatTerm {
	return &FloatTerm{getter: getter, name: name}
}
