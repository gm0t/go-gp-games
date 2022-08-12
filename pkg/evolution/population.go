package evolution

import (
	"fmt"
	"math/rand"

	"lr1Go/pkg/tree"
)

type Gene struct {
	Agent       *tree.Node `json:"agent"`
	Generations int        `json:"generations"`
	Fitness     float64    `json:"fitness"`
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
	currentGen     int
}

func (p *Population) Genes() []*Gene {
	return p.genes
}

func (p *Population) CurrentGeneration() int {
	return p.currentGen
}

func (p *Population) Evolve(generations int) {
	for currentGeneration := 0; currentGeneration < generations; currentGeneration += 1 {
		for _, gen := range p.genes {
			gen.Fitness = p.fitness(gen.Agent, currentGeneration)
		}

		mutants := buildMutants(currentGeneration, p.genes, p.mutationChance, p.fitness, p.generator)
		children := buildChildren(currentGeneration, p.genes, p.fitness, p.size)
		truncate(mutants, p.generator, 5)
		truncate(children, p.generator, 5)

		pool := make([]*Gene, len(mutants)+len(children)+len(p.genes))
		copy(pool, p.genes)
		copy(pool[len(p.genes):], mutants)
		copy(pool[len(p.genes)+len(mutants):], children)

		p.genes = NewTournamentSelector(pool).Select(p.size)
		var longestSurvivor *Gene
		for _, g := range p.genes {
			g.Generations += 1
			if longestSurvivor == nil || longestSurvivor.Generations < g.Generations {
				longestSurvivor = g
			}
		}
		p.currentGen = currentGeneration
	}
}

func truncate(genes []*Gene, generator tree.Generator, maxDepth int) {
	for _, gene := range genes {
		agents := make([]*tree.Node, 0)
		tree.Dfs(gene.Agent, func(n *tree.Node, depth int) {
			if depth == maxDepth {
				agents = append(agents, n)
			}
		})

		for _, agent := range agents {
			generator.Truncate(agent)
		}
	}
}

func (p *Population) Best() *Gene {
	var best *Gene
	for _, g := range p.genes {
		if best == nil || best.Fitness < g.Fitness {
			best = g
		}
	}
	return best
}

func (p *Population) Worst() *Gene {
	var worst *Gene
	for _, g := range p.genes {
		if worst == nil || worst.Fitness > g.Fitness {
			worst = g
		}
	}
	return worst
}

func buildChildren(generation int, genes []*Gene, fitness Fitness, size int) []*Gene {
	children := make([]*Gene, 0)
	roulette := NewRouletteSelector(genes)
	for i := 0; i < size/2; i += 1 {
		child1, child2 := Crossover(roulette.Pick().Agent, roulette.Pick().Agent)
		if child1 != nil {
			children = append(children, &Gene{
				Agent:   child1,
				Fitness: fitness(child1, generation),
			})
		}
		if child2 != nil {
			children = append(children, &Gene{
				Agent:   child2,
				Fitness: fitness(child2, generation),
			})
		}
	}

	return children
}

func buildMutants(generation int, genes []*Gene, chance float64, fitness Fitness, generator tree.Generator) []*Gene {
	mutants := make([]*Gene, 0)
	for _, g := range genes {
		if rand.Float64() <= chance {
			newAgent := tree.Clone(g.Agent)
			Mutate(newAgent, generator)
			mutants = append(mutants, &Gene{
				Agent:   newAgent,
				Fitness: fitness(newAgent, generation),
			})
		}
	}
	return mutants
}

var perfectAgentX = tree.NewIf(
	tree.NewFunction(
		tree.Lt,
		tree.Boolean,
		[]*tree.Node{{
			Key:      "MyX",
			Type:     "float",
			Children: nil,
		}, {
			Key:      "GoalX",
			Type:     "float",
			Children: nil,
		}},
	),
	tree.NewFunction(
		tree.Minus,
		tree.Float,
		[]*tree.Node{{
			Key:      "GoalX",
			Type:     "float",
			Children: nil,
		}, {
			Key:      "MyX",
			Type:     "float",
			Children: nil,
		}},
	),
	tree.NewFunction(
		tree.Minus,
		tree.Float,
		[]*tree.Node{{
			Key:      "MyX",
			Type:     "float",
			Children: nil,
		}, &tree.Node{
			Key:      "GoalX",
			Type:     "float",
			Children: nil,
		}},
	),
)
var perfectAgentY = tree.NewIf(
	tree.NewFunction(
		tree.Lt,
		tree.Boolean,
		[]*tree.Node{{
			Key:      "MyY",
			Type:     "float",
			Children: nil,
		}, {
			Key:      "GoalY",
			Type:     "float",
			Children: nil,
		}},
	),
	tree.NewFunction(
		tree.Minus,
		tree.Float,
		[]*tree.Node{{
			Key:      "GoalY",
			Type:     "float",
			Children: nil,
		}, {
			Key:      "MyY",
			Type:     "float",
			Children: nil,
		}},
	),
	tree.NewFunction(
		tree.Minus,
		tree.Float,
		[]*tree.Node{{
			Key:      "MyY",
			Type:     "float",
			Children: nil,
		}, {
			Key:      "GoalY",
			Type:     "float",
			Children: nil,
		}},
	),
)

var PerfectAgent = tree.NewFunction(tree.Plus, tree.Float, []*tree.Node{perfectAgentX, perfectAgentY})

func NewPopulation(
	size int,
	generator tree.Generator,
	fitness Fitness,
	mutationChance float64,
) *Population {
	initialGenes := make([]*Gene, 0)
	//initialGenes = append(initialGenes, &Gene{
	//	Agent:   PerfectAgent,
	//	Fitness: 0,
	//})

	for i := 0; i < size; i += 1 {
		agent := generator.FTree(4)
		initialGenes = append(initialGenes, &Gene{
			Agent:   agent,
			Fitness: 0,
		})
	}

	fmt.Println(PerfectAgent, fitness(PerfectAgent, 0))

	return &Population{
		size:           size,
		generator:      generator,
		fitness:        fitness,
		genes:          initialGenes,
		mutationChance: mutationChance,
	}
}
