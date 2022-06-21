package tree

import (
	"fmt"
	"strings"
)

func Print(root *Node) {
	Dfs(root, func(n *Node, depth int) {
		padding := strings.Repeat("| ", depth)
		fmt.Printf("%v%v\n", padding, n.Key)
	})
}
