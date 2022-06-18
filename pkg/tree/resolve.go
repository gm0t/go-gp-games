package tree

import (
	"math"
)

type ResolveArguments interface {
	Float(key string) float64
	Boolean(key string) bool
	Action(key string) interface{}
}

func Resolve(node *Node, args ResolveArguments) interface{} {
	children := node.Children
	switch node.Type {
	case Float:
		return args.Float(node.Key)
	case Boolean:
		return args.Boolean(node.Key)
	case Action:
		return args.Action(node.Key)
	case IF:
		if resolveBoolean(children[0], args) {
			return Resolve(children[1], args)
		}
		return Resolve(children[2], args)
	case Plus:
		return resolveFloat(children[0], args) + resolveFloat(children[1], args)
	case Minus:
		return resolveFloat(children[0], args) - resolveFloat(children[1], args)
	case Multiply:
		return resolveFloat(children[0], args) * resolveFloat(children[1], args)
	case Divide:
		right := resolveFloat(children[1], args)
		if right == 0 {
			return 0
		}
		return resolveFloat(children[0], args) / right
	case Gt:
		return resolveFloat(children[0], args) > resolveFloat(children[1], args)
	case Lt:
		return resolveFloat(children[0], args) < resolveFloat(children[1], args)
	case Eq:
		return math.Abs(resolveFloat(children[0], args)-resolveFloat(children[1], args)) < 0.01
	}

	panic("Unknown node type :" + node.Type)
}

func resolveFloat(node *Node, args ResolveArguments) float64 {
	return Resolve(node, args).(float64)
}

func resolveBoolean(node *Node, args ResolveArguments) bool {
	return Resolve(node, args).(bool)
}
