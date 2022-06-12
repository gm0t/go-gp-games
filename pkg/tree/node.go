package tree

type Node interface {
	Dfs(cb func(depth int, n Node), extra ...int)
	String() string
	Clone() Node
}

type BooleanNode interface {
	Node
	Resolve(args ResolveArguments) bool
}

type FloatNode interface {
	Node
	Resolve(args ResolveArguments) float64
}

type Action interface{}

type ActionNode interface {
	Node
	Resolve(args ResolveArguments) Action
}

type Generator interface {
	FFunc() FloatFunctionNode
	BFunc() BooleanFunctionNode
	AFunc() ActionFunctionNode
	FTerm() FloatNode
	BTerm() BooleanNode
	ATerm() ActionNode
	FTree(depth int) FloatFunctionNode
	ATree(depth int) ActionFunctionNode
	BTree(depth int) BooleanFunctionNode
}

type FunctionNode interface {
	Node
	Mutate(generator Generator)
	Grow(generator Generator)
	Truncate(generator Generator)
	ReplaceF(cNode FloatNode, nNode FloatNode) bool
	ReplaceB(cNode BooleanNode, nNode BooleanNode) bool
	ReplaceA(cNode ActionNode, nNode ActionNode) bool
}

type FloatFunctionNode interface {
	FunctionNode
	FloatNode
}

type BooleanFunctionNode interface {
	FunctionNode
	BooleanNode
}

type ActionFunctionNode interface {
	FunctionNode
	ActionNode
}

type ResolveArguments interface {
	Float(key string) float64
	Boolean(key string) bool
}
