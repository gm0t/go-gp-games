package main

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/evolution"
	"lr1Go/pkg/game/catcher"
	"lr1Go/pkg/game/catcher/agent"
	"lr1Go/pkg/tree"
)

func main() {
	generator := agent.NewGenerator()

	population := evolution.NewPopulation(50, generator, createFitness(), 0.5)
	population.Evolve(1500000)
	agent, f := population.Best()
	fmt.Printf("Best agent: \n-fitness: %v\n-score: %v\n", f, simulateGame(agent.(tree.FloatFunctionNode), buildRandomState()))
	//tree.Print(agent)
}

func buildRandomState() *catcher.State {
	return catcher.NewState(
		float64(rand.Intn(30)),
		float64(rand.Intn(30)),
		float64(40+rand.Intn(30)),
		float64(40+rand.Intn(30)),
	)
}

func createFitness() evolution.Fitness {
	solved := 0
	lastGen := 0
	playerX := 0.0
	playerY := 0.0
	goalX := 40.0
	goalY := 40.0

	state := func() *catcher.State {
		return catcher.NewState(playerX, playerY, goalX, goalY)
	}

	return func(agent tree.FunctionNode, generation int) float64 {
		if lastGen != generation {
			lastGen = generation
			if solved > 5 {
				playerX = float64(rand.Intn(30))
				playerY = float64(rand.Intn(30))
				goalX = float64(rand.Intn(30)) + 40
				goalY = float64(rand.Intn(30)) + 40
				fmt.Printf("\n-----------------\nPuzzle was solved at gen %v %v times, building new state... %v\n-----------------\n", lastGen, solved, state())
			}
			//fmt.Printf("\n++++++\nPuzzle was solved at gen %v %v times\n++++++\n", lastGen, solved)
			solved = 0
		}

		s := state()
		result := simulateGame(agent, s)
		if result.DistanceToGoal < 1 {
			solved += 1
			return 1 + 1/float64(result.Iterations)
		}

		return 1 / result.DistanceToGoal
	}
}

func simulateGame(agent tree.FunctionNode, state *catcher.State) catcher.Result {

	game := catcher.NewGame(state, buildPlayer(agent))
	return game.Run(200)
}

func buildPlayer(agent tree.FunctionNode) catcher.Player {
	if fAgent, isFloat := agent.(tree.FloatFunctionNode); isFloat {
		return catcher.NewAiFPlayer(fAgent)
	}

	if aAgent, isAction := agent.(tree.ActionFunctionNode); isAction {
		return catcher.NewAiAPlayer(aAgent)
	}

	panic("Unknown type of agent!")
}
