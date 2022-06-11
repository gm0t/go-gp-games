package functions

import (
	"math/rand"
)

func extractDepth(extra []int) int {
	if len(extra) > 0 {
		return extra[0]
	}

	return 0
}

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}
