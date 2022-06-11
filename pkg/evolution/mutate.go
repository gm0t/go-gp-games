package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/tree"
)

func Mutate(node tree.Node, generator tree.Generator) {
	allNodes := make([]tree.FunctionNode, 0)
	node.Dfs(func(depth int, n tree.Node) {
		if f, isFunc := n.(tree.FunctionNode); isFunc {
			allNodes = append(allNodes, f)
		}
	})

	if len(allNodes) < 1 {
		fmt.Println("No suitable node for mutation")
		return
	}

	candidate := getRandom(allNodes)
	candidate.Mutate(generator)

	node.Dfs(func(depth int, n tree.Node) {
		if f, isFunc := n.(tree.FunctionNode); isFunc && rand.Float64() < 0.1 {
			f.Mutate(generator)
		}
	})
}
