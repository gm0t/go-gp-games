package catcher

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

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
	lock := sync.Mutex{}
	//var lastWinner *tree.Node
	return func(agent *tree.Node, generation int) float64 {
		lock.Lock()
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
		value := accumulatedResult / float64(len(states))
		if math.IsInf(value, +1) {
			fmt.Printf("Fitness is infinite: %v / %v ", accumulatedResult, float64(len(states)))
		}
		lock.Unlock()
		return value
	}
}

func simulateGame(agent *tree.Node, state *State) Result {
	game := NewGame(state, NewPlayer(agent))
	return game.Run(200)
}
