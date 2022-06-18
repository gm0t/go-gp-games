package tree

func Clone(node *Node) *Node {
	// this will create a clone of the original struct
	newNode := *node

	// and we need to make a clone of each child separately
	newChildren := make([]*Node, len(node.Children))
	for i, child := range node.Children {
		newChildren[i] = Clone(child)
	}
	newNode.Children = newChildren
	return &newNode
}
