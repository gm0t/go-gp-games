package functions

import (
	"math/rand"

	"lr1Go/pkg/tree"
)

type IfF struct {
	condition tree.BooleanNode
	success   tree.FloatNode
	fail      tree.FloatNode
}

func (i *IfF) Truncate(generator tree.Generator) {
	if _, isFunc := i.condition.(tree.FunctionNode); isFunc {
		i.condition = generator.BTerm()
	}
	if _, isFunc := i.success.(tree.FunctionNode); isFunc {
		i.success = generator.FTerm()
	}
	if _, isFunc := i.fail.(tree.FunctionNode); isFunc {
		i.fail = generator.FTerm()
	}
}

func (i *IfF) ReplaceA(cNode tree.ActionNode, nNode tree.ActionNode) bool {
	return false
}

func (i *IfF) Clone() tree.Node {
	return NewIfF(
		i.condition.Clone().(tree.BooleanNode),
		i.success.Clone().(tree.FloatNode),
		i.fail.Clone().(tree.FloatNode),
	)
}

func (i *IfF) ReplaceF(cNode tree.FloatNode, nNode tree.FloatNode) bool {
	if i.success == cNode {
		i.success = nNode
		return true
	} else if i.fail == cNode {
		i.fail = nNode
		return true
	}
	return false
}

func (i *IfF) ReplaceB(cNode tree.BooleanNode, nNode tree.BooleanNode) bool {
	if i.condition == cNode {
		i.condition = nNode
		return true
	}
	return false
}

func (i *IfF) String() string {
	return "if(" + i.condition.String() + ") { " + i.success.String() + "} else {" + i.fail.String() + "})"
}

func (i *IfF) Dfs(cb func(depth int, n tree.Node), extra ...int) {
	depth := extractDepth(extra)
	cb(depth, i)
	i.condition.Dfs(cb, depth+1)
	i.success.Dfs(cb, depth+1)
	i.fail.Dfs(cb, depth+1)
}

func (i *IfF) Grow(g tree.Generator) {
	i.condition = g.BFunc()
	i.success = g.FFunc()
	i.fail = g.FFunc()
}

func (i *IfF) Mutate(g tree.Generator) {
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

func (i *IfF) Resolve(args tree.ResolveArguments) float64 {
	if i.condition.Resolve(args) {
		return i.success.Resolve(args)
	}

	return i.fail.Resolve(args)
}

func NewIfF(condition tree.BooleanNode, success tree.FloatNode, fail tree.FloatNode) *IfF {
	if condition == nil {
		panic("condition can't be nil")
	}
	return &IfF{condition: condition, success: success, fail: fail}
}
