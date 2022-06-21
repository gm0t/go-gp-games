package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/tree"
)

func Mutate(root *tree.Node, generator tree.Generator) {
	allNodes := make([]*tree.Node, 0)
	tree.Dfs(root, func(node *tree.Node, depth int) {
		allNodes = append(allNodes, node)
	})

	if len(allNodes) < 1 {
		fmt.Println("No suitable node for mutation")
		return
	}

	mutant := getRandom(allNodes)

	switch mutant.Type {
	case tree.Float:
		fallthrough
	case tree.Boolean:
		fallthrough
	case tree.Action:
		mutateTerm(mutant, generator)
	case tree.IF:
		mutateIf(mutant, generator)
	case tree.Plus:
		fallthrough
	case tree.Minus:
		fallthrough
	case tree.Divide:
		fallthrough
	case tree.Multiply:
		mutateMath(mutant, generator)
	case tree.Eq:
		fallthrough
	case tree.Gt:
		fallthrough
	case tree.Lt:
		mutateBFunc(mutant, generator)
	}
}

func mutateMath(node *tree.Node, generator tree.Generator) {
	switch rand.Intn(4) {
	case 0:
		newOp := getRandom(tree.MathOperators())
		node.Type = newOp
		node.Key = newOp
	case 1:
		Mutate(node.Children[0], generator)
	case 2:
		Mutate(node.Children[1], generator)
	case 3:
		overwrite(node, generator.Term(tree.Float))
	}
}
func mutateIf(node *tree.Node, generator tree.Generator) {
	switch rand.Intn(3) {
	case 0:
		mutateBFunc(node.Children[0], generator)
	case 1:
		Mutate(node.Children[1], generator)
	case 2:
		Mutate(node.Children[2], generator)
	}
}

func mutateBFunc(node *tree.Node, generator tree.Generator) {
	switch rand.Intn(3) {
	case 0:
		newOp := getRandom(tree.ComparisonOperators())
		node.Type = newOp
		node.Key = newOp
	case 1:
		Mutate(node.Children[0], generator)
	case 2:
		Mutate(node.Children[1], generator)
	}
}

func mutateTerm(node *tree.Node, generator tree.Generator) {
	if rand.Float32() < 0.8 {
		overwrite(node, generator.Term(node.Type))
	} else {
		overwrite(node, generator.Tree(node.Type, 1+rand.Intn(2)))
	}
}

func overwrite(dst *tree.Node, source *tree.Node) {
	dst.Type = source.Type
	dst.Children = source.Children
	dst.Key = source.Key
}
