package catcher

import (
	"math/rand"

	"lr1Go/pkg/evolution"
	"lr1Go/pkg/tree"
)

func buildTestStates(count int) []*State {
	states := make([]*State, count)
	for i := 0; i < count; i += 1 {
		states[i] = NewState(
			float64(rand.Intn(40)),
			float64(rand.Intn(40)),
			float64(rand.Intn(40)),
			float64(rand.Intn(40)),
		)
	}
	return states
}

func NewFitness() evolution.Fitness {
	states := buildTestStates(5)
	lastGen := 0
	//var lastWinner *tree.Node
	return func(agent *tree.Node, generation int) float64 {
		if lastGen != generation {
			lastGen = generation
			states = buildTestStates(5)
		}

		accumulatedResult := float64(0)
		for _, state := range states {
			result := simulateGame(agent, state.Clone())
			if result.DistanceToGoal < 1 {
				accumulatedResult += 1 + 1/float64(result.Iterations)
			} else {
				accumulatedResult += 1 / result.DistanceToGoal
			}
		}

		// average fitness across multiple runs
		return accumulatedResult / float64(len(states))
	}
}

func simulateGame(agent *tree.Node, state *State) Result {
	game := NewGame(state, buildPlayer(agent))
	return game.Run(200)
}

func buildPlayer(agent *tree.Node) Player {
	return NewAiFPlayer(agent)
	//
	//switch agent.Type {
	//case tree.Float:
	//	return NewAiFPlayer(agent)
	//case tree.Action:
	//	return NewAiAPlayer(agent)
	//}

	//panic("Unknown type of agent: " + agent.Type)
}
