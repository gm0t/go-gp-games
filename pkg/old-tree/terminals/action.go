package terminals

import (
	"lr1Go/pkg/old-tree"
)

type ActionTerm struct {
	name     string
	valueKey string
	getter   func(args old_tree.ResolveArguments) old_tree.Action
}

func (n *ActionTerm) Serialize() old_tree.Serialized {
	return old_tree.Serialized{
		Key: old_tree.ActionT,
		Parameters: map[string]interface{}{
			"name": n.name,
		},
		Children: nil,
	}
}

func (n *ActionTerm) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	cb(extractDepth(extra), n)
}

func (n *ActionTerm) String() string {
	return n.name
}

func (n *ActionTerm) Clone() old_tree.Node {
	return NewActionTerm(n.name, n.getter)
}

func (n *ActionTerm) Resolve(args old_tree.ResolveArguments) old_tree.Action {
	return args.Action
}
func NewActionTerm(name string, getter func(args old_tree.ResolveArguments) old_tree.Action) *ActionTerm {
	return &ActionTerm{getter: getter, name: name}
}
