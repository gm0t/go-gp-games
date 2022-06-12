package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/tree"
)

func getAnyChild(root tree.FunctionNode) tree.Node {
	children := make([]tree.FunctionNode, 0)
	root.Dfs(func(depth int, n tree.Node) {
		if fn, isFunc := n.(tree.FunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getFloatChild(root tree.FunctionNode) tree.FloatNode {
	children := make([]tree.FloatNode, 0)
	root.Dfs(func(depth int, n tree.Node) {
		if fn, isFunc := n.(tree.FloatFunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getActionChild(root tree.FunctionNode) tree.ActionNode {
	children := make([]tree.ActionNode, 0)
	root.Dfs(func(depth int, n tree.Node) {
		if fn, isFunc := n.(tree.ActionFunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getBoolChild(root tree.FunctionNode) tree.BooleanNode {
	children := make([]tree.BooleanNode, 0)
	root.Dfs(func(depth int, n tree.Node) {
		if fn, isFunc := n.(tree.BooleanFunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func Crossover(parent1, parent2 tree.FunctionNode) (tree.FunctionNode, tree.FunctionNode) {
	if parent1 == parent2 {
		// there is no point in crossover with itself
		return nil, nil
	}

	node1 := parent1.Clone().(tree.FunctionNode)
	node2 := parent2.Clone().(tree.FunctionNode)

	point1 := getAnyChild(node1)
	if point1 == nil {
		fmt.Println("Failed to find a good point for crossover in parent 1")
		return nil, nil
	}

	var point2 tree.Node
	if _, isFloat := point1.(tree.FloatNode); isFloat {
		point2 = getFloatChild(node2)
	} else if _, isBool := point1.(tree.BooleanNode); isBool {
		point2 = getBoolChild(node2)
	} else if _, isAction := point1.(tree.ActionNode); isAction {
		point2 = getActionChild(node2)
	}
	if point2 == nil {
		fmt.Println("Failed to find a good point for crossover in parent 2")
		return nil, nil
	}

	replace(node1, point1, point2.Clone())
	replace(node2, point2, point1.Clone())

	return node1, node2
}

//
//func findEqNodes(parent1 tree.FunctionNode, parent2 tree.FunctionNode) {
//	visited := make(map[tree.Node]bool)
//	parent1.Dfs(func(depth int, n tree.Node) {
//		visited[n] = true
//	})
//
//	parent2.Dfs(func(depth int, n tree.Node) {
//		if visited[n] {
//			fmt.Println("+++++++++++++++++++++++++++++")
//			fmt.Println(parent1)
//			fmt.Println(parent2)
//			fmt.Println("Matching node at level ", depth, "\n", n)
//			fmt.Println("+++++++++++++++++++++++++++++")
//			panic("wttf is this")
//		}
//	})
//}

func replace(root tree.FunctionNode, cNode tree.Node, nNode tree.Node) {
	replaced := false
	root.Dfs(func(depth int, n tree.Node) {
		if f, isFunc := n.(tree.FunctionNode); !replaced && isFunc {
			if fNode, isFloat := cNode.(tree.FloatNode); isFloat {
				replaced = f.ReplaceF(fNode, nNode.(tree.FloatNode))
			} else if aNode, isAction := cNode.(tree.ActionNode); isAction {
				replaced = f.ReplaceA(aNode, nNode.(tree.ActionNode))
			} else if bNode, isBoolean := cNode.(tree.BooleanNode); isBoolean {
				replaced = f.ReplaceB(bNode, nNode.(tree.BooleanNode))
			}
			//if replaced {
			//	fmt.Println("Replacing ", cNode, "for", nNode, "at", depth)
			//}
		}
	})
}

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}
