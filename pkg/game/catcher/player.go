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
