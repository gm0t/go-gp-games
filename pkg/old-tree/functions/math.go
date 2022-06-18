package functions

import (
	"math/rand"

	"lr1Go/pkg/old-tree"
)

type MathOp string

const (
	Plus           MathOp = "+"
	Minus          MathOp = "-"
	Multiplication MathOp = "*"
	Division       MathOp = "/"
)

var AllMathOps = []MathOp{Plus, Minus, Multiplication, Division}

type Math struct {
	op    MathOp
	left  old_tree.FloatNode
	right old_tree.FloatNode
}

func (n *Math) Truncate(generator old_tree.Generator) {
	if _, isFunc := n.left.(old_tree.FunctionNode); isFunc {
		n.left = generator.FTerm()
	}
	if _, isFunc := n.right.(old_tree.FunctionNode); isFunc {
		n.right = generator.FTerm()
	}
}

func (n *Math) ReplaceA(cNode old_tree.ActionNode, nNode old_tree.ActionNode) bool {
	return false
}

func (n *Math) Clone() old_tree.Node {
	return NewMath(
		n.op,
		n.left.Clone().(old_tree.FloatNode),
		n.right.Clone().(old_tree.FloatNode),
	)
}

func (n *Math) ReplaceF(cNode old_tree.FloatNode, nNode old_tree.FloatNode) bool {
	if n.left == cNode {
		n.left = nNode
		return true
	} else if n.right == cNode {
		n.right = nNode
		return true
	}
	return false
}

func (n *Math) ReplaceB(old_tree.BooleanNode, old_tree.BooleanNode) bool {
	return false
}

func (n *Math) String() string {
	return "(" + n.left.String() + string(n.op) + n.right.String() + ")"
}

func (n *Math) Dfs(cb func(depth int, n old_tree.Node), extra ...int) {
	depth := extractDepth(extra)
	cb(depth, n)
	n.left.Dfs(cb, depth+1)
	n.right.Dfs(cb, depth+1)
}

func (n *Math) Mutate(g old_tree.Generator) {
	switch rand.Intn(5) {
	case 0:
		n.op = getRandom(AllMathOps)
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

func (n *Math) Grow(g old_tree.Generator) {
	n.left = g.FFunc()
	n.right = g.FFunc()
}

func (n *Math) Resolve(args old_tree.ResolveArguments) float64 {
	right := n.right.Resolve(args)
	left := n.left.Resolve(args)

	switch n.op {
	case Plus:
		return left + right
	case Minus:
		return left - right
	case Multiplication:
		return left * right
	case Division:
		if right == 0 {
			return 0
		}
		return left / right
	}

	panic("Unknown op: " + string(n.op))
}

func NewMath(op MathOp, left old_tree.FloatNode, right old_tree.FloatNode) *Math {
	return &Math{
		op:    op,
		left:  left,
		right: right,
	}
}
