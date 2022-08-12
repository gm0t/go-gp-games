package evolution

import (
	"math/rand"
)

type TournamentSelector struct {
	genes []*Gene
}

func NewTournamentSelector(genes []*Gene) *TournamentSelector {
	return &TournamentSelector{
		genes: genes,
	}
}

func (s *TournamentSelector) Select(size int) []*Gene {
	return tournamentSelection(s.genes, size)
}

func runTournament(genes []*Gene, size int) ([]*Gene, []*Gene) {
	if len(genes)%2 == 1 {
		genes = append(genes, genes[0])
	}

	losers := make([]*Gene, 0)
	winners := make([]*Gene, 0)
	for i := 0; i < len(genes) && len(winners) <= size; i += 2 {
		if genes[i].Fitness > genes[i+1].Fitness {
			losers = append(losers, genes[i+1])
			winners = append(winners, genes[i])
		} else {
			losers = append(losers, genes[i])
			winners = append(winners, genes[i+1])
		}
	}

	return winners, losers
}

func tournamentSelection(genes []*Gene, size int) []*Gene {
	shuffled := make([]*Gene, len(genes))
	copy(shuffled, genes)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	winners, losers := runTournament(shuffled, len(shuffled))

	if len(winners) > size {
		return tournamentSelection(winners, size)
	}

	var extra []*Gene
	for len(winners) < size && len(losers) > 2 {
		extra, losers = runTournament(losers, size-len(winners))
		winners = append(winners, extra...)
	}

	return winners
}
