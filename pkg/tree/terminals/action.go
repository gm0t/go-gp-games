package terminals

import (
	"lr1Go/pkg/tree"
)

type ActionTerm struct {
	name   string
	getter func(args tree.ResolveArguments) tree.Action
}

func (n *ActionTerm) Dfs(cb func(depth int, n tree.Node), extra ...int) {
	cb(extractDepth(extra), n)
}

func (n *ActionTerm) String() string {
	return n.name
}

func (n *ActionTerm) Clone() tree.Node {
	return NewActionTerm(n.name, n.getter)
}

func (n *ActionTerm) Resolve(args tree.ResolveArguments) tree.Action {
	return n.getter(args)
}
func NewActionTerm(name string, getter func(args tree.ResolveArguments) tree.Action) *ActionTerm {
	return &ActionTerm{getter: getter, name: name}
}
