package tree

func dfsLoop(node *Node, callback func(node *Node, depth int), depth int) {
	callback(node, depth)
	for _, child := range node.Children {
		dfsLoop(child, callback, depth+1)
	}
}

func Dfs(root *Node, callback func(node *Node, depth int)) {
	dfsLoop(root, callback, 0)
}
