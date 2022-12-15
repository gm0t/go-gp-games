package catcher

import (
	"fmt"

	"lr1Go/pkg/evolution"
	"lr1Go/pkg/tree"
)

func NewTerminationCondition() func(population *evolution.Population) bool {
	return func(population *evolution.Population) bool {
		top := population.Top(10)

		if canSolveEverything(population.Best().Agent, population.Best().Fitness > 1) {
			fmt.Println("The best one solved it!")
		}

		for idx, g := range top {
			if canSolveEverything(g.Agent, false) {
				fmt.Println("Solved by: ", idx, g == population.Best())
				fmt.Println("Can the best one solve it?", canSolveEverything(population.Best().Agent, false))
				return true
			}
		}

		return false
	}
}

func allTestStates() []*State {
	states := make([]*State, 0)

	possibleValues := []float64{
		0, 1, 5, 10, 50,
	}

	for _, pX := range possibleValues {
		for _, pY := range possibleValues {
			for _, gX := range possibleValues {
				for _, gY := range possibleValues {
					states = append(states, NewState(
						pX, pY,
						gX, gY,
					))

				}
			}
		}
	}

	return states
}

func canSolveEverything(agent *tree.Node, debug bool) bool {
	states := allTestStates()

	for _, state := range states {
		result := simulateGame(agent, state.Clone())
		if result.DistanceToGoal >= 1 {
			if debug {
				fmt.Printf("============== \nFailed check: %f %v \n============== \n", result.DistanceToGoal, state)
			}
			return false
		}
	}

	return true
}
