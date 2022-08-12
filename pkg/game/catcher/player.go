package catcher

import (
	"math"

	"lr1Go/pkg/game/catcher/actions"
	"lr1Go/pkg/tree"
)

type Player interface {
	GetAction(state *State) actions.Action
}

type AiFPlayer struct {
	strategy *tree.Node
}

func (ai *AiFPlayer) GetAction(state *State) actions.Action {
	allActions := []actions.Action{actions.Down, actions.Up, actions.Left, actions.Right}

	strategy := ai.strategy

	bestScore := math.Inf(1)
	bestAction := actions.Up

	for _, action := range allActions {
		args := state.Clone().Apply(action).BuildArgs()
		score := tree.Resolve(strategy, args).(float64)
		if score < bestScore {
			bestScore = score
			bestAction = action
		}
	}

	return bestAction
}

func Normalize(node *tree.Node) *tree.Node {
	if node.Key == tree.IF {
		return normalizeIf(node)
	}

	normalized := *node
	normalized.Children = normalizeChildren(node.Children)
	return &normalized
}

func normalizeChildren(children []*tree.Node) []*tree.Node {
	newChildren := make([]*tree.Node, len(children))
	for i, child := range children {
		newChildren[i] = Normalize(child)
	}

	return newChildren
}

func normalizeIf(ifNode *tree.Node) *tree.Node {
	normalized := *ifNode
	cond := normalized.Children[0]
	if isAlwaysTrue(cond) {
		return Normalize(normalized.Children[1])
	} else if isAlwaysFalse(cond) {
		return Normalize(normalized.Children[2])
	}

	normalized.Children = normalizeChildren(ifNode.Children)
	return &normalized
}

func isAlwaysTrue(cond *tree.Node) bool {
	if cond.Key == tree.Eq {
		return cond.Children[0].String() == cond.Children[1].String()
	}

	return false
}

func isAlwaysFalse(cond *tree.Node) bool {
	if cond.Key == tree.Lt || cond.Key == tree.Gt {
		return cond.Children[0].String() == cond.Children[1].String()
	}

	return false
}

func NewAiFPlayer(strategy *tree.Node) *AiFPlayer {
	return &AiFPlayer{strategy: strategy}
}

type AiAPlayer struct {
	strategy *tree.Node
}

func (ai *AiAPlayer) GetAction(state *State) actions.Action {
	return tree.Resolve(ai.strategy, state.BuildArgs()).(actions.Action)
}

func NewAiAPlayer(strategy *tree.Node) *AiAPlayer {
	return &AiAPlayer{strategy: strategy}
}
