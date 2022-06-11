package catcher

import (
	"math"

	"lr1Go/pkg/game/catcher/actions"
	"lr1Go/pkg/tree"
)

type Player interface {
	GetAction(state *State) actions.Action
}

type AiPlayer struct {
	strategy tree.FloatNode
}

func NewAiPlayer(strategy tree.FloatNode) *AiPlayer {
	return &AiPlayer{strategy: strategy}
}

func (ai *AiPlayer) GetAction(state *State) actions.Action {
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
