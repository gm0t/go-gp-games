package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/tree"
)

type Gene struct {
	agent   *tree.Node
	fitness float64
}

type Fitness func(node *tree.Node, generation int) float64

type Options struct {
	MaxGenerations int
	MutationChance float64
}

type Population struct {
	size           int
	generator      tree.Generator
	fitness        Fitness
	genes          []*Gene
	mutationChance float64
}

func (p *Population) Generation() {
}

func (p *Population) Evolve(generations int) {
	for g := 0; g < generations; g += 1 {
		for _, gen := range p.genes {
			gen.fitness = p.fitness(gen.agent, g)
		}

		mutants := buildMutants(g, p.genes, p.mutationChance, p.fitness, p.generator)
		children := buildChildren(g, p.genes, p.fitness, p.size)
		truncate(mutants, p.generator, 3)
		truncate(mutants, p.generator, 3)

		pool := make([]*Gene, len(mutants)+len(children)+len(p.genes))
		copy(pool, p.genes)
		copy(pool[len(p.genes):], mutants)
		copy(pool[len(p.genes)+len(mutants):], children)

		p.genes = NewRouletteSelector(pool).Select(p.size)
		if g%50 == 0 {
			_, bf := p.Best()
			_, wf := p.Worst()
			fmt.Printf("%d generation completed: %v / %v = %.4f - %.4f\n-----\n", g, len(mutants), len(children), wf, bf)
		}
	}
}

func truncate(genes []*Gene, generator tree.Generator, maxDepth int) {
	for _, gene := range genes {
		agents := make([]*tree.Node, 0)
		tree.Dfs(gene.agent, func(n *tree.Node, depth int) {
			if depth == maxDepth {
				agents = append(agents, n)
			}
		})

		for _, agent := range agents {
			generator.Truncate(agent)
		}
	}
}

func (p *Population) Best() (*tree.Node, float64) {
	var best *Gene
	for _, g := range p.genes {
		if best == nil || best.fitness < g.fitness {
			best = g
		}
	}
	return best.agent, best.fitness
}

func (p *Population) Worst() (*tree.Node, float64) {
	var worst *Gene
	for _, g := range p.genes {
		if worst == nil || worst.fitness > g.fitness {
			worst = g
		}
	}
	return worst.agent, worst.fitness
}

func buildChildren(generation int, genes []*Gene, fitness Fitness, size int) []*Gene {
	children := make([]*Gene, 0)
	roulette := NewRouletteSelector(genes)
	for i := 0; i < size/2; i += 1 {
		child1, child2 := Crossover(roulette.Pick().agent, roulette.Pick().agent)
		if child1 != nil {
			children = append(children, &Gene{
				agent:   child1,
				fitness: fitness(child1, generation),
			})
		}
		if child2 != nil {
			children = append(children, &Gene{
				agent:   child2,
				fitness: fitness(child2, generation),
			})
		}
	}

	return children
}

func buildMutants(generation int, genes []*Gene, chance float64, fitness Fitness, generator tree.Generator) []*Gene {
	mutants := make([]*Gene, 0)
	for _, g := range genes {
		if rand.Float64() <= chance {
			newAgent := tree.Clone(g.agent)
			Mutate(newAgent, generator)
			mutants = append(mutants, &Gene{
				agent:   newAgent,
				fitness: fitness(newAgent, generation),
			})
		}
	}
	return mutants
}

func NewPopulation(
	size int,
	generator tree.Generator,
	fitness Fitness,
	mutationChance float64,
) *Population {
	initialGenes := make([]*Gene, 0)
	for i := 0; i < size; i += 1 {
		agent := generator.FTree(2)
		initialGenes = append(initialGenes, &Gene{
			agent:   agent,
			fitness: 0,
		})
	}

	return &Population{
		size:           size,
		generator:      generator,
		fitness:        fitness,
		genes:          initialGenes,
		mutationChance: mutationChance,
	}
}
