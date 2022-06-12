package functions

import (
	"math/rand"

	"lr1Go/pkg/tree"
)

type IfA struct {
	condition tree.BooleanNode
	success   tree.ActionNode
	fail      tree.ActionNode
}

func (i *IfA) Clone() tree.Node {
	return NewIfA(
		i.condition.Clone().(tree.BooleanNode),
		i.success.Clone().(tree.ActionNode),
		i.fail.Clone().(tree.ActionNode),
	)
}

func (i *IfA) ReplaceF(cNode tree.FloatNode, nNode tree.FloatNode) bool {
	return false
}

func (i *IfA) ReplaceB(cNode tree.BooleanNode, nNode tree.BooleanNode) bool {
	if i.condition == cNode {
		i.condition = nNode
		return true
	}
	return false
}

func (i *IfA) ReplaceA(cNode tree.ActionNode, nNode tree.ActionNode) bool {
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
	return "if(" + i.condition.String() + ") {" + i.success.String() + "} else {" + i.fail.String() + "})"
}

func (i *IfA) Dfs(cb func(depth int, n tree.Node), extra ...int) {
	depth := extractDepth(extra)
	cb(depth, i)
	i.condition.Dfs(cb, depth+1)
	i.success.Dfs(cb, depth+1)
	i.fail.Dfs(cb, depth+1)
}

func (i *IfA) Grow(g tree.Generator) {
	i.condition = g.BFunc()
	i.success = g.AFunc()
	i.fail = g.AFunc()
}

func (i *IfA) Mutate(g tree.Generator) {
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

func (i *IfA) Resolve(args tree.ResolveArguments) tree.Action {
	if i.condition.Resolve(args) {
		return i.success.Resolve(args)
	}

	return i.fail.Resolve(args)
}

func NewIfA(condition tree.BooleanNode, success tree.ActionNode, fail tree.ActionNode) *IfA {
	if condition == nil {
		panic("condition can't be nil")
	}
	return &IfA{condition: condition, success: success, fail: fail}
}
