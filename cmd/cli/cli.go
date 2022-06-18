package main

import (
	"fmt"

	"lr1Go/pkg/evolution"
	"lr1Go/pkg/game/catcher/agent"
	"lr1Go/pkg/old-tree"
)

func main() {
	generator := agent.NewGenerator()

	population := evolution.NewPopulation(50, generator, createFitness(), 0.5)
	population.Evolve(1500000)
	agent, f := population.Best()
	fmt.Printf("Best agent: \n-fitness: %v\n-score: %v\n", f, simulateGame(agent.(old_tree.FloatFunctionNode), buildRandomState()))
	//old-tree.Print(agent)
}
