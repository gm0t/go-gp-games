package old_tree

import (
	"fmt"
	"strings"
)

func Print(root Node) {
	root.Dfs(func(depth int, n Node) {
		padding := strings.Repeat("| ", depth)
		fmt.Printf("%v%v\n", padding, n.String())
	})
}
