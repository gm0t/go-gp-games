package tree

type NodeType = string

const (
	Float    NodeType = "float"
	Boolean  NodeType = "bool"
	Action   NodeType = "action"
	IF       NodeType = "if"
	Plus     NodeType = "+"
	Minus    NodeType = "-"
	Multiply NodeType = "*"
	Divide   NodeType = "/"
	Eq       NodeType = "=="
	Gt       NodeType = ">"
	Lt       NodeType = "<"
)

var comparisonNodeTypes = []NodeType{Eq, Gt, Lt}
var mathNodeTypes = []NodeType{Plus, Minus, Minus, Divide, Multiply}

type Node struct {
	Key      string   `json:"key"`
	Type     NodeType `json:"type"`
	Children []*Node  `json:"children"`
}
