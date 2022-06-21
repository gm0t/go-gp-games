package catcher

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/evolution"
	"lr1Go/pkg/tree"
)

func NewFitness() evolution.Fitness {
	solved := 0
	lastGen := 0
	playerX := 0.0
	playerY := 0.0
	goalX := 40.0
	goalY := 40.0

	state := func() *State {
		return NewState(playerX, playerY, goalX, goalY)
	}

	var lastWinner *tree.Node
	return func(agent *tree.Node, generation int) float64 {
		if lastGen != generation {
			lastGen = generation
			playerX = float64(rand.Intn(30))
			playerY = float64(rand.Intn(30))
			goalX = float64(rand.Intn(30)) + 40
			goalY = float64(rand.Intn(30)) + 40
			if solved > 0 {
				fmt.Printf("\n-----------------\nPuzzle was solved at gen %v %v times\n+++++++++++++++++\n", lastGen, solved)
				winnerFormula := lastWinner.String()
				fmt.Printf("Last winner is: %v \n %v -----------------\v\v", len(winnerFormula), winnerFormula)
			}
			//fmt.Printf("\n++++++\nPuzzle was solved at gen %v %v times\n++++++\n", lastGen, solved)
			solved = 0
		}

		s := state()
		result := simulateGame(agent, s)
		if result.DistanceToGoal < 1 {
			solved += 1
			lastWinner = agent
			return 1 + 1/float64(result.Iterations)
		}

		return 1 / result.DistanceToGoal
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
