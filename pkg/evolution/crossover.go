package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/old-tree"
)

func getAnyChild(root old_tree.FunctionNode) old_tree.Node {
	children := make([]old_tree.FunctionNode, 0)
	root.Dfs(func(depth int, n old_tree.Node) {
		if fn, isFunc := n.(old_tree.FunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getFloatChild(root old_tree.FunctionNode) old_tree.FloatNode {
	children := make([]old_tree.FloatNode, 0)
	root.Dfs(func(depth int, n old_tree.Node) {
		if fn, isFunc := n.(old_tree.FloatFunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getActionChild(root old_tree.FunctionNode) old_tree.ActionNode {
	children := make([]old_tree.ActionNode, 0)
	root.Dfs(func(depth int, n old_tree.Node) {
		if fn, isFunc := n.(old_tree.ActionFunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getBoolChild(root old_tree.FunctionNode) old_tree.BooleanNode {
	children := make([]old_tree.BooleanNode, 0)
	root.Dfs(func(depth int, n old_tree.Node) {
		if fn, isFunc := n.(old_tree.BooleanFunctionNode); isFunc && depth > 1 {
			children = append(children, fn)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func Crossover(parent1, parent2 old_tree.FunctionNode) (old_tree.FunctionNode, old_tree.FunctionNode) {
	if parent1 == parent2 {
		// there is no point in crossover with itself
		return nil, nil
	}

	node1 := parent1.Clone().(old_tree.FunctionNode)
	node2 := parent2.Clone().(old_tree.FunctionNode)

	point1 := getAnyChild(node1)
	if point1 == nil {
		fmt.Println("Failed to find a good point for crossover in parent 1")
		return nil, nil
	}

	var point2 old_tree.Node
	if _, isFloat := point1.(old_tree.FloatNode); isFloat {
		point2 = getFloatChild(node2)
	} else if _, isBool := point1.(old_tree.BooleanNode); isBool {
		point2 = getBoolChild(node2)
	} else if _, isAction := point1.(old_tree.ActionNode); isAction {
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
//func findEqNodes(parent1 old-tree.FunctionNode, parent2 old-tree.FunctionNode) {
//	visited := make(map[old-tree.Node]bool)
//	parent1.Dfs(func(depth int, n old-tree.Node) {
//		visited[n] = true
//	})
//
//	parent2.Dfs(func(depth int, n old-tree.Node) {
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

func replace(root old_tree.FunctionNode, cNode old_tree.Node, nNode old_tree.Node) {
	replaced := false
	root.Dfs(func(depth int, n old_tree.Node) {
		if f, isFunc := n.(old_tree.FunctionNode); !replaced && isFunc {
			if fNode, isFloat := cNode.(old_tree.FloatNode); isFloat {
				replaced = f.ReplaceF(fNode, nNode.(old_tree.FloatNode))
			} else if aNode, isAction := cNode.(old_tree.ActionNode); isAction {
				replaced = f.ReplaceA(aNode, nNode.(old_tree.ActionNode))
			} else if bNode, isBoolean := cNode.(old_tree.BooleanNode); isBoolean {
				replaced = f.ReplaceB(bNode, nNode.(old_tree.BooleanNode))
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
