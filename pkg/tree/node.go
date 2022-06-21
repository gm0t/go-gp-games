package tree

import (
	"fmt"
)

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

func ComparisonOperators() []NodeType {
	return comparisonNodeTypes
}

func MathOperators() []NodeType {
	return mathNodeTypes
}

type Node struct {
	Key      string   `json:"key"`
	Type     NodeType `json:"type"`
	Children []*Node  `json:"children"`
}

func (n *Node) String() string {
	switch n.Type {
	case Float:
		fallthrough
	case Boolean:
		fallthrough
	case Action:
		return n.Key
	case IF:
		return fmt.Sprintf("(%v ? %v : %v)", n.Children[0], n.Children[1], n.Children[2])
	case Plus:
		fallthrough
	case Minus:
		fallthrough
	case Multiply:
		fallthrough
	case Divide:
		fallthrough
	case Gt:
		fallthrough
	case Lt:
		fallthrough
	case Eq:
		return fmt.Sprintf("(%v %v %v)", n.Children[0], n.Key, n.Children[1])
	}

	return "(Unknown: " + n.Type + ")"
}
