package tree

import (
	"fmt"
)

type NodeType = string
type NodeKey = string

const (
	Float   NodeType = "float"
	Boolean NodeType = "bool"
	Action  NodeType = "action"
)

const (
	IF       NodeKey = "if"
	Plus     NodeKey = "+"
	Minus    NodeKey = "-"
	Multiply NodeKey = "*"
	Divide   NodeKey = "/"
	Eq       NodeKey = "=="
	Gt       NodeKey = ">"
	Lt       NodeKey = "<"
)

var comparisonNodes = []NodeType{Eq, Gt, Lt}

//var mathNodes = []NodeType{Plus, Minus, Divide, Multiply}
var mathNodes = []NodeType{}

func ComparisonOperators() []NodeType {
	return comparisonNodes
}

func MathOperators() []NodeType {
	return mathNodes
}

type Node struct {
	Key      string   `json:"key"`
	Type     NodeType `json:"type"`
	Children []*Node  `json:"children"`
}

func (n *Node) String() string {
	switch n.Key {
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

	return n.Key
}
