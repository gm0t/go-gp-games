package tree

import (
	"math/rand"

	"lr1Go/pkg/tree/generator"
)

type Generator interface {
	FFunc() *Node
	BFunc() *Node
	AFunc() *Node
	FTerm() *Node
	BTerm() *Node
	ATerm() *Node
	FTree(depth int) *Node
	ATree(depth int) *Node
	BTree(depth int) *Node
}

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}

type BasicGenerator struct {
	floats   []string
	booleans []string
	actions  []string
}

func (b BasicGenerator) FFunc() *Node {
	return generator.NewIf(
		b.BTerm(),
		b.FTerm(),
		b.FTerm(),
	)
}

func (b BasicGenerator) BFunc() *Node {
	nType := getRandom(comparisonNodeTypes)
	return &Node{
		Key:      nType,
		Type:     nType,
		Children: []*Node{},
	}
}

func (b BasicGenerator) AFunc() *Node {
	//TODO implement me
	panic("implement me")
}

func (b BasicGenerator) FTerm() *Node {
	//TODO implement me
	panic("implement me")
}

func (b BasicGenerator) BTerm() *Node {
	//TODO implement me
	panic("implement me")
}

func (b BasicGenerator) ATerm() *Node {
	//TODO implement me
	panic("implement me")
}

func (b BasicGenerator) FTree(depth int) *Node {
	//TODO implement me
	panic("implement me")
}

func (b BasicGenerator) ATree(depth int) *Node {
	//TODO implement me
	panic("implement me")
}

func (b BasicGenerator) BTree(depth int) *Node {
	//TODO implement me
	panic("implement me")
}

func NewGenerator(booleans []string, floats []string, actions []string) Generator {
	return &BasicGenerator{
		floats:   floats,
		booleans: booleans,
		actions:  actions,
	}
}
