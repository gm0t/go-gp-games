package functions

import (
	"math/rand"

	"lr1Go/pkg/old-tree"
)

type IfF struct {
	condition old_tree.BooleanNode
	success   old_tree.FloatNode
	fail      old_tree.FloatNode
}

func (i *IfF) Truncate(generator old_tree.Generator) {
	if _, isFunc := i.condition.(old_tree.FunctionNode); isFunc {
		i.condition = generator.BTerm()
	}
	if _, isFunc := i.success.(old_tree.FunctionNode); isFunc {
		i.success = generator.FTerm()
	}
	if _, isFunc := i.fail.(old_tree.FunctionNode); isFunc {
		i.fail = generator.FTerm()
	}
}

func (i *IfF) ReplaceA(cNode old_tree.ActionNode, nNode old_tree.ActionNode) bool {
	return false
}

func (i *IfF) Clone() old_tree.Node {
	return NewIfF(
		i.condition.Clone().(old_tree.BooleanNode),
		i.success.Clone().(old_tree.FloatNode),
		i.fail.Clone().(old_tree.FloatNode),
	)
}

func (i *IfF) ReplaceF(cNode old_tree.FloatNode, nNode old_tree.FloatNode) bool {
	if i.success == cNode {
		i.success = nNode
		return true
	} else if i.fail == cNode {
		i.fail = nNode
		return true
	}
	return false
}

func (i *IfF) ReplaceB(cNode old_tree.BooleanNode, nNode old_tree.BooleanNode) bool {
	if i.condition == cNode {
		i.condition = nNode
		return true
	}
	return false
}

func (i *IfF) String() string {
	return "if(" + i.condition.String() + ") {" + i.success.String() + "} else {" + i.fail.String() + "}"
}

func (i *IfF) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	depth := extractDepth(extra)
	cb(depth, i)
	i.condition.Dfs(cb, depth+1)
	i.success.Dfs(cb, depth+1)
	i.fail.Dfs(cb, depth+1)
}

func (i *IfF) Grow(g old_tree.Generator) {
	i.condition = g.BFunc()
	i.success = g.FFunc()
	i.fail = g.FFunc()
}

func (i *IfF) Mutate(g old_tree.Generator) {
	switch rand.Intn(6) {
	case 0:
		i.condition = g.BTerm()
	case 1:
		i.condition = g.BTree(2)
	case 2:
		i.fail = g.FFunc()
	case 3:
		i.success = g.FFunc()
	case 4:
		i.fail = g.FTree(2)
	case 5:
		i.success = g.FTree(2)
	}
}

func (i *IfF) Resolve(args old_tree.ResolveArguments) float64 {
	if i.condition.Resolve(args) {
		return i.success.Resolve(args)
	}

	return i.fail.Resolve(args)
}

func NewIfF(condition old_tree.BooleanNode, success old_tree.FloatNode, fail old_tree.FloatNode) *IfF {
	if condition == nil {
		panic("condition can't be nil")
	}
	return &IfF{condition: condition, success: success, fail: fail}
}
