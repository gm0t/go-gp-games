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
	node1 := parent1.Clone().(tree.FunctionNode)
	node2 := parent2.Clone().(tree.FunctionNode)

	child1 := getAnyChild(node1)
	if child1 == nil {
		fmt.Println("Failed to find a good point for crossover in parent 1")
		return nil, nil
	}

	if fNode, isFloat := child1.(tree.FloatNode); isFloat {
		fCrossover(node1, node2, fNode)
	}

	if bNode, isBool := child1.(tree.BooleanNode); isBool {
		bCrossover(node1, node2, bNode)
	}

	if aNode, isAction := child1.(tree.ActionNode); isAction {
		aCrossover(node1, node2, aNode)
	}

	return node1, node2
}

func fCrossover(node1 tree.FunctionNode, node2 tree.FunctionNode, point1 tree.FloatNode) {
	point2 := getFloatChild(node2)
	if point2 == nil {
		return
	}

	replaced := false
	visited := make(map[tree.Node]bool)

	node1.Dfs(func(depth int, n tree.Node) {
		if visited[n] {
			panic("LOOP DETECTED!")
		}
		//visited[n] = true
		//if depth > 1000 {
		//	tree.Print(node1)
		//	panic("Something went wrong and we are too deep")
		//}
		if f, isFunc := n.(tree.FunctionNode); !replaced && isFunc {
			replaced = f.ReplaceF(point1, point2)
		}
	})

	replaced = false
	node2.Dfs(func(depth int, n tree.Node) {
		if f, isFunc := n.(tree.FunctionNode); !replaced && isFunc {
			replaced = f.ReplaceF(point2, point1)
		}
	})
}

func bCrossover(node1 tree.FunctionNode, node2 tree.FunctionNode, point1 tree.BooleanNode) {
	point2 := getBoolChild(node2)
	if point2 == nil {
		return
	}

	node1.Dfs(func(depth int, n tree.Node) {
		if f, isFunc := n.(tree.FunctionNode); isFunc {
			f.ReplaceB(point1, point2)
		}
	})
}

func aCrossover(node1 tree.FunctionNode, node2 tree.FunctionNode, point1 tree.ActionNode) {
	point2 := getActionChild(node2)
	if point2 == nil {
		return
	}

	node1.Dfs(func(depth int, n tree.Node) {
		if f, isFunc := n.(tree.FunctionNode); isFunc {
			f.ReplaceA(point1, point2)
		}
	})
}

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}
