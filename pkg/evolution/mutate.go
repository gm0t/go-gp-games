package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/old-tree"
)

func Mutate(node old_tree.Node, generator old_tree.Generator) {
	allNodes := make([]old_tree.FunctionNode, 0)
	node.Dfs(func(depth int, n old_tree.Node) {
		if f, isFunc := n.(old_tree.FunctionNode); isFunc {
			allNodes = append(allNodes, f)
		}
	})

	if len(allNodes) < 1 {
		fmt.Println("No suitable node for mutation")
		return
	}

	candidate := getRandom(allNodes)
	candidate.Mutate(generator)

	node.Dfs(func(depth int, n old_tree.Node) {
		if f, isFunc := n.(old_tree.FunctionNode); isFunc && rand.Float64() < 0.1 {
			f.Mutate(generator)
		}
	})
}
