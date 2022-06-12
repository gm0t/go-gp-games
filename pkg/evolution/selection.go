package evolution

import (
	"math/rand"

	"lr1Go/pkg/tree"
)

type RouletteSelector struct {
	genes   []*Gene
	fitness []float64
}

func NewRouletteSelector(genes []*Gene) *RouletteSelector {
	var totalFitness float64
	normalizedFitness := make([]float64, len(genes))
	for _, gen := range genes {
		totalFitness += gen.fitness
	}
	for i, gen := range genes {
		normalizedFitness[i] = gen.fitness / totalFitness
	}

	return &RouletteSelector{
		genes:   genes,
		fitness: normalizedFitness,
	}
}

func (s *RouletteSelector) Pick() *Gene {
	chance := rand.Float64()
	var acc float64
	for i, f := range s.fitness {
		acc += f
		if acc >= chance {
			return s.genes[i]
		}
	}

	return s.genes[len(s.genes)-1]

}

func (s *RouletteSelector) Select(size int) []*Gene {
	selection := make([]*Gene, size)
	for i := 0; i < size; i += 1 {
		chosen := s.Pick()
		selection[i] = &Gene{
			agent:   chosen.agent.Clone().(tree.FunctionNode),
			fitness: chosen.fitness,
		}
	}

	return selection
}
