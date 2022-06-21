package main

import (
	"lr1Go/pkg/evolution"
	"lr1Go/pkg/game/catcher"
	"lr1Go/pkg/game/catcher/actions"
	"lr1Go/pkg/tree"
)

func main() {
	generator := tree.NewGenerator(catcher.BoolKeys, catcher.FloatKeys, []string{
		string(actions.Up),
		string(actions.Down),
		string(actions.Left),
		string(actions.Right),
	})

	//generator := agent.NewGenerator()
	population := evolution.NewPopulation(50, generator, catcher.NewFitness(), 0.5)
	population.Evolve(1500000)
	//agent, f := population.Best()
	//fmt.Printf("Best agent: \n-fitness: %v\n-score: %v\n", f, simulateGame(agent.(old_tree.FloatFunctionNode), buildRandomState()))
	////old-tree.Print(agent)
}
