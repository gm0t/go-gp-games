package functions

import (
	"math/rand"

	"lr1Go/pkg/old-tree"
)

type IfA struct {
	condition old_tree.BooleanNode
	success   old_tree.ActionNode
	fail      old_tree.ActionNode
}

func (i *IfA) Truncate(generator old_tree.Generator) {
	if _, isFunc := i.condition.(old_tree.FunctionNode); isFunc {
		i.condition = generator.BTerm()
	}
	if _, isFunc := i.success.(old_tree.FunctionNode); isFunc {
		i.success = generator.ATerm()
	}
	if _, isFunc := i.fail.(old_tree.FunctionNode); isFunc {
		i.fail = generator.ATerm()
	}
}

func (i *IfA) Clone() old_tree.Node {
	return NewIfA(
		i.condition.Clone().(old_tree.BooleanNode),
		i.success.Clone().(old_tree.ActionNode),
		i.fail.Clone().(old_tree.ActionNode),
	)
}

func (i *IfA) ReplaceF(cNode old_tree.FloatNode, nNode old_tree.FloatNode) bool {
	return false
}

func (i *IfA) ReplaceB(cNode old_tree.BooleanNode, nNode old_tree.BooleanNode) bool {
	if i.condition == cNode {
		i.condition = nNode
		return true
	}
	return false
}

func (i *IfA) ReplaceA(cNode old_tree.ActionNode, nNode old_tree.ActionNode) bool {
	if i.success == cNode {
		i.success = nNode
		return true
	} else if i.fail == cNode {
		i.fail = nNode
		return true
	}
	return false
}

func (i *IfA) String() string {
	return "if(" + i.condition.String() + ") {" + i.success.String() + "} else {" + i.fail.String() + "}"
}

func (i *IfA) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	depth := extractDepth(extra)
	cb(depth, i)
	i.condition.Dfs(cb, depth+1)
	i.success.Dfs(cb, depth+1)
	i.fail.Dfs(cb, depth+1)
}

func (i *IfA) Grow(g old_tree.Generator) {
	i.condition = g.BFunc()
	i.success = g.AFunc()
	i.fail = g.AFunc()
}

func (i *IfA) Mutate(g old_tree.Generator) {
	switch rand.Intn(6) {
	case 0:
		i.condition = g.BTerm()
	case 1:
		i.condition = g.BTree(2)
	case 2:
		i.fail = g.AFunc()
	case 3:
		i.success = g.AFunc()
	case 4:
		i.fail = g.ATree(2)
	case 5:
		i.success = g.ATree(2)
	}
}

func (i *IfA) Resolve(args old_tree.ResolveArguments) old_tree.Action {
	if i.condition.Resolve(args) {
		return i.success.Resolve(args)
	}

	return i.fail.Resolve(args)
}

func NewIfA(condition old_tree.BooleanNode, success old_tree.ActionNode, fail old_tree.ActionNode) *IfA {
	if condition == nil {
		panic("condition can't be nil")
	}
	return &IfA{condition: condition, success: success, fail: fail}
}
