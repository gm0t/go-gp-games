package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/tree"
)

func getAnyChild(root *tree.Node) *tree.Node {
	children := make([]*tree.Node, 0)
	tree.Dfs(root, func(node *tree.Node, depth int) {
		if depth > 1 {
			children = append(children, node)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func getMatchingChild(root *tree.Node, childType tree.NodeType) *tree.Node {
	children := make([]*tree.Node, 0)
	tree.Dfs(root, func(node *tree.Node, depth int) {
		if depth > 1 && node.Type == childType {
			children = append(children, node)
		}
	})

	if len(children) < 1 {
		return nil
	}

	return getRandom(children)
}

func Crossover(parent1, parent2 *tree.Node) (*tree.Node, *tree.Node) {
	if parent1 == parent2 {
		// there is no point in crossover with itself
		return nil, nil
	}

	node1 := tree.Clone(parent1)
	node2 := tree.Clone(parent2)

	point1 := getAnyChild(node1)
	if point1 == nil {
		fmt.Println("Failed to find a good point for crossover in parent 1")
		return nil, nil
	}

	point2 := getMatchingChild(node2, point1.Type)
	if point2 == nil {
		fmt.Println("Failed to find a good point for crossover in parent 2")
		return nil, nil
	}

	// swapping parameters
	p1Key := point1.Key
	p1Children := point1.Children
	point1.Key = point2.Key
	point1.Children = point2.Children
	point2.Key = p1Key
	point2.Children = p1Children

	return node1, node2
}

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}
