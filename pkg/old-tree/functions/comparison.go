package functions

import (
	"math/rand"

	"lr1Go/pkg/old-tree"
)

type CompOp string

const (
	Eq CompOp = "=="
	Lt CompOp = "<"
	Gt CompOp = ">"
)

var AllCompOps = []CompOp{Eq, Lt, Gt}

type Comparison struct {
	op    CompOp
	left  old_tree.FloatNode
	right old_tree.FloatNode
}

func (n *Comparison) Truncate(generator old_tree.Generator) {
	if _, isFunc := n.left.(old_tree.FunctionNode); isFunc {
		n.left = generator.FTerm()
	}
	if _, isFunc := n.right.(old_tree.FunctionNode); isFunc {
		n.right = generator.FTerm()
	}
}

func (n *Comparison) ReplaceA(cNode old_tree.ActionNode, nNode old_tree.ActionNode) bool {
	return false
}

func (n *Comparison) ReplaceF(cNode old_tree.FloatNode, nNode old_tree.FloatNode) bool {
	if n.left == cNode {
		n.left = nNode
		return true
	} else if n.right == cNode {
		n.right = nNode
		return true
	}
	return false
}

func (n *Comparison) ReplaceB(old_tree.BooleanNode, old_tree.BooleanNode) bool {
	return false
}

func (n *Comparison) Clone() old_tree.Node {
	return NewComparison(
		n.op,
		n.left.Clone().(old_tree.FloatNode),
		n.right.Clone().(old_tree.FloatNode),
	)
}

func (n *Comparison) String() string {
	return n.left.String() + string(n.op) + n.right.String()
}

func (n *Comparison) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	depth := extractDepth(extra)
	cb(depth, n)
	n.left.Dfs(cb, depth+1)
	n.right.Dfs(cb, depth+1)
}

func (n *Comparison) Resolve(args old_tree.ResolveArguments) bool {
	r := n.right.Resolve(args)
	l := n.left.Resolve(args)
	switch n.op {
	case Lt:
		return r < l
	case Eq:
		return r == l
	case Gt:
		return r > l
	}

	panic("Unknown comparison op: " + n.op)
}

func (n *Comparison) Mutate(g old_tree.Generator) {
	switch rand.Intn(5) {
	case 0:
		n.op = getRandom(AllCompOps)
	case 1:
		n.left = g.FFunc()
	case 2:
		n.right = g.FFunc()
	case 3:
		n.left = g.FTree(2)
	case 4:
		n.right = g.FTree(2)
	}
}

func (n *Comparison) Grow(g old_tree.Generator) {
	n.left = g.FFunc()
	n.right = g.FFunc()
}

func NewComparison(op CompOp, left old_tree.FloatNode, right old_tree.FloatNode) *Comparison {
	return &Comparison{op: op, left: left, right: right}
}
