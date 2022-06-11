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
	strategy tree.FloatNode
}

func (ai *AiFPlayer) GetAction(state *State) actions.Action {
	allActions := []actions.Action{actions.Down, actions.Up, actions.Left, actions.Right}

	strategy := ai.strategy

	bestScore := math.Inf(1)
	bestAction := actions.Up

	for _, action := range allActions {
		args := state.Clone().Apply(action).BuildArgs()
		score := strategy.Resolve(args)
		if score < bestScore {
			bestScore = score
			bestAction = action
		}
	}

	return bestAction
}

func NewAiFPlayer(strategy tree.FloatNode) *AiFPlayer {
	return &AiFPlayer{strategy: strategy}
}

type AiAPlayer struct {
	strategy tree.ActionNode
}

func (ai *AiAPlayer) GetAction(state *State) actions.Action {
	action := ai.strategy.Resolve(state.BuildArgs())

	return action.(actions.Action)
}

func NewAiAPlayer(strategy tree.ActionNode) *AiAPlayer {
	return &AiAPlayer{strategy: strategy}
}