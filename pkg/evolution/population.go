package evolution

import (
	"math"
	"math/rand"
	"sort"

	"lr1Go/pkg/tree"
)

type Gene struct {
	Agent             *tree.Node `json:"agent"`
	Generations       int        `json:"generations"`
	Born              int        `json:"born"`
	Fitness           float64    `json:"fitness"`
	FitnessOver10Gens float64    `json:"fitnessOver10Gens"`
	fitnessLog        []float64
	fitnessLogIdx     int
}

func (g *Gene) AddFitness(value float64) {
	if g.fitnessLog == nil {
		g.fitnessLog = make([]float64, 10)
	}

	g.fitnessLog[g.getFitnessLogIdx()] = value
	g.FitnessOver10Gens = averageNonZero(g.fitnessLog)
	g.Fitness = value
}

func averageNonZero(arr []float64) float64 {
	total := float64(0)
	cnt := 0
	for _, value := range arr {
		if value > 0 {
			total += value
			cnt += 1
		}
	}

	return total / float64(cnt)
}

func (g *Gene) getFitnessLogIdx() int {
	idx := g.fitnessLogIdx
	if g.fitnessLogIdx == 9 {
		g.fitnessLogIdx = 0
	} else {
		g.fitnessLogIdx += 1
	}

	return idx
}

type Fitness func(node *tree.Node, generation int) float64

type Options struct {
	MaxGenerations int
	MutationChance float64
}

type StatRow struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Median  float64 `json:"median"`
	Average float64 `json:"average"`
}

type Population struct {
	size              int
	generator         tree.Generator
	fitness           Fitness
	genes             []*Gene
	elites            []*Gene
	statsLog          []*StatRow
	mutationChance    float64
	currentGeneration int
	isFinished        bool
	childrenSize      int
	eliteSize         int
	isStopping        bool
	waitingToStop     chan interface{}
	params            Params
}

func (p *Population) Genes() []*Gene {
	return p.genes
}
func (p *Population) Elites() []*Gene {
	return p.elites
}

func (p *Population) CurrentGeneration() int {
	return p.currentGeneration
}

func (p *Population) Evolve(terminate func(population *Population) bool) {
	p.elites = make([]*Gene, 0)
	p.isStopping = false

	for !terminate(p) {
		for _, gen := range p.genes {
			gen.AddFitness(p.fitness(gen.Agent, p.currentGeneration))
			gen.Generations += 1
		}
		for _, gen := range p.elites {
			gen.AddFitness(p.fitness(gen.Agent, p.currentGeneration))
			gen.Generations += 1
		}

		mutants := buildMutants(p.currentGeneration, p.genes, p.mutationChance, p.fitness, p.generator)
		children := buildChildren(p.currentGeneration, p.genes, p.fitness, p.childrenSize)
		truncate(mutants, p.generator, 5)
		truncate(children, p.generator, 5)

		poolSize := len(mutants) + len(children) + len(p.genes) + len(p.elites)
		pool := make([]*Gene, poolSize)
		copy(pool, p.genes)
		copy(pool[len(p.genes):], mutants)
		copy(pool[len(p.genes)+len(mutants):], children)
		copy(pool[len(p.genes)+len(mutants)+len(children):], p.elites)

		p.genes = NewTournamentSelector(pool).Select(p.size)
		sort.Sort(ByFitness(pool))
		p.elites = make([]*Gene, p.eliteSize)
		for i := 0; i < p.eliteSize; i += 1 {
			g := pool[len(pool)-(i+1)]
			// TODO: elites should be unique!
			p.elites[i] = &Gene{
				Agent:             tree.Clone(g.Agent),
				Generations:       g.Generations,
				Born:              g.Born,
				Fitness:           g.Fitness,
				FitnessOver10Gens: g.FitnessOver10Gens,
				fitnessLog:        g.fitnessLog,
				fitnessLogIdx:     g.fitnessLogIdx,
			}
		}

		p.currentGeneration += 1
		p.statsLog = append(p.statsLog, buildStatsFor(p.genes))
		if p.isStopping {
			p.isFinished = true
			return
		}
	}

	p.isFinished = true
}

func buildStatsFor(genes []*Gene) *StatRow {
	min := math.Inf(+1)
	max := math.Inf(-1)
	average := float64(0)

	for _, g := range genes {
		average += g.Fitness
		if g.Fitness > max {
			max = g.Fitness
		}
		if g.Fitness < min {
			min = g.Fitness
		}
	}
	sort.Sort(ByFitness(genes))
	median := genes[len(genes)/2].Fitness

	return &StatRow{
		Min:     min,
		Max:     max,
		Median:  median,
		Average: average / float64(len(genes)),
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
	best := func(current, next *Gene) bool {
		return next.Fitness > current.Fitness
	}

	if len(p.elites) == 0 {
		// we don't have any elites yet
		return findBestBy(p.genes, best)
	}

	// elite genes are always best
	return findBestBy(p.elites, best)
}

func (p *Population) Oldest() *Gene {
	oldest := func(current, next *Gene) bool {
		return next.Generations > current.Generations
	}

	if len(p.elites) == 0 {
		// we don't have any elites yet
		return findBestBy(p.genes, oldest)
	}

	oldestElite := findBestBy(p.elites, oldest)
	oldestBasic := findBestBy(p.genes, oldest)

	if oldestElite.Generations > oldestBasic.Generations {
		return oldestElite
	}

	return oldestBasic
}

func (p *Population) Worst() *Gene {
	return findBestBy(p.genes, func(current, next *Gene) bool {
		return next.Fitness < current.Fitness
	})
}

func (p *Population) IsFinished() bool {
	return p.isFinished
}

func (p *Population) Stop(wait chan interface{}) {
	p.isStopping = true
	go func() {
		for !p.isFinished {
		}
		wait <- true
	}()
}

func (p *Population) Params() Params {
	return p.params
}

func (p *Population) Stats() []*StatRow {
	return p.statsLog
}

func findBestBy(genes []*Gene, condition func(current, next *Gene) bool) *Gene {
	var current *Gene
	for _, next := range genes {
		if current == nil || condition(current, next) {
			current = next
		}
	}
	return current
}

func buildChildren(generation int, genes []*Gene, fitness Fitness, size int) []*Gene {
	children := make([]*Gene, 0)
	roulette := NewRouletteSelector(genes)
	for i := 0; i < size/2; i += 1 {
		child1, child2 := Crossover(roulette.Pick().Agent, roulette.Pick().Agent)
		if child1 != nil {
			children = append(children, NewGene(generation, child1, fitness))
		}
		if child2 != nil {
			children = append(children, NewGene(generation, child2, fitness))
		}
	}

	return children
}

func NewGene(generation int, agent *tree.Node, fitness Fitness) *Gene {
	gene := &Gene{
		Agent: agent,
		Born:  generation,
	}
	gene.AddFitness(fitness(agent, generation))
	return gene
}

func buildMutants(generation int, genes []*Gene, chance float64, fitness Fitness, generator tree.Generator) []*Gene {
	mutants := make([]*Gene, 0)
	for _, g := range genes {
		if rand.Float64() <= chance {
			newAgent := tree.Clone(g.Agent)
			Mutate(newAgent, generator)
			mutants = append(mutants, NewGene(generation, newAgent, fitness))
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

type Params struct {
	Size           int            `json:"size"`
	ElitesSize     int            `json:"elitesSize"`
	ChildrenSize   int            `json:"childrenSize"`
	Generator      tree.Generator `json:"-"`
	Fitness        Fitness        `json:"-"`
	MutationChance float64        `json:"mutationChance"`
}

func buildInitialGenes(size int, generator tree.Generator) []*Gene {
	genes := make([]*Gene, 0)
	for i := 0; i < size; i += 1 {
		agent := generator.ATree(4)
		genes = append(genes, &Gene{
			Agent:       agent,
			Born:        0,
			Generations: 0,
			Fitness:     0,
		})
	}

	return genes
}

func NewPopulation(params *Params) *Population {
	return &Population{
		size:              params.Size,
		childrenSize:      params.ChildrenSize,
		eliteSize:         params.ElitesSize,
		generator:         params.Generator,
		fitness:           params.Fitness,
		genes:             buildInitialGenes(params.Size, params.Generator),
		elites:            make([]*Gene, 0),
		mutationChance:    params.MutationChance,
		currentGeneration: 0,
		isFinished:        false,
		params:            *params,
	}
}

type ByFitness []*Gene

func (genes ByFitness) Len() int           { return len(genes) }
func (genes ByFitness) Less(i, j int) bool { return genes[i].Fitness < genes[j].Fitness }
func (genes ByFitness) Swap(i, j int)      { genes[i], genes[j] = genes[j], genes[i] }
