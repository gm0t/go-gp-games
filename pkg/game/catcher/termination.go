package catcher

import (
	"lr1Go/pkg/evolution"
)

func NewTerminationCondition() func(population *evolution.Population) bool {
	return func(population *evolution.Population) bool {
		best := population.Best()
		if best.Fitness > 1 && best.Generations > 1000 && best.FitnessOver10Gens > 1 {
			// it seems like we have a winner
			return true
		}

		return false
	}
}
