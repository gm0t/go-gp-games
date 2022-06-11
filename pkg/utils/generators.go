package utils

import (
	"fmt"
	"math"
	"math/rand"

	"lr1Go/pkg/tree"
	"lr1Go/pkg/tree/functions"
	"lr1Go/pkg/tree/terminals"
)

type fTermGenerator = func() tree.FloatNode
type bTermGenerator = func() tree.BooleanNode

type fFuncGenerator = func() tree.FloatFunctionNode
type bFuncGenerator = func() tree.BooleanFunctionNode

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}

type Generator struct {
	fFunc fFuncGenerator
	bFunc bFuncGenerator
	fTerm fTermGenerator
	bTerm bTermGenerator
	aFunc func() tree.ActionFunctionNode
	aTerm func() tree.ActionNode
}

func (g *Generator) BTree(maxDepth int) tree.BooleanFunctionNode {
	root := g.bFunc()

	root.Dfs(func(depth int, n tree.Node) {
		if f, isFunction := n.(tree.FunctionNode); isFunction && (depth < maxDepth) {
			f.Grow(g)
		}
	})

	return root
}

func (g *Generator) FFunc() tree.FloatFunctionNode {
	return g.fFunc()
}

func (g *Generator) AFunc() tree.ActionFunctionNode {
	return g.aFunc()
}

func (g *Generator) BFunc() tree.BooleanFunctionNode {
	return g.bFunc()
}

func (g *Generator) FTerm() tree.FloatNode {
	return g.fTerm()
}

func (g *Generator) ATerm() tree.ActionNode {
	return g.aTerm()
}

func (g *Generator) BTerm() tree.BooleanNode {
	return g.bTerm()
}

func (g *Generator) FTree(maxDepth int) tree.FloatFunctionNode {
	root := g.fFunc()

	root.Dfs(func(depth int, n tree.Node) {
		if f, isFunction := n.(tree.FunctionNode); isFunction && (depth < maxDepth) {
			f.Grow(g)
		}
	})

	return root
}

func (g *Generator) ATree(maxDepth int) tree.ActionFunctionNode {
	root := g.aFunc()

	root.Dfs(func(depth int, n tree.Node) {
		if f, isFunction := n.(tree.FunctionNode); isFunction && (depth < maxDepth) {
			f.Grow(g)
		}
	})

	return root
}

func NodeGenerator(floatKeys []string, booleanKeys []string, actions []tree.Action) *Generator {
	fTerm := func() tree.FloatNode {
		rndValue := math.Round(rand.Float64()*1000) / 100
		rndKey := fmt.Sprintf("rnd[%.2f]", rndValue)

		allFloatKeys := make([]string, len(floatKeys)+1)
		copy(allFloatKeys, floatKeys)
		allFloatKeys[len(allFloatKeys)-1] = rndKey

		chosenKey := getRandom(floatKeys)

		return terminals.NewFloatTerm(chosenKey, func(args tree.ResolveArguments) float64 {
			if chosenKey == rndKey {
				return rndValue
			}

			return args.Float(chosenKey)
		})
	}

	bFunc := func() tree.BooleanFunctionNode {
		return functions.NewComparison(getRandom(functions.AllCompOps), fTerm(), fTerm())
	}

	bTerm := func() tree.BooleanNode {
		if len(booleanKeys) == 0 {
			return bFunc()
		}

		chosenKey := getRandom(booleanKeys)
		return terminals.NewBooleanTerm(chosenKey, func(args tree.ResolveArguments) bool {
			return args.Boolean(chosenKey)
		})
	}

	aTerm := func() tree.ActionNode {
		action := getRandom(actions)
		return terminals.NewActionTerm(fmt.Sprintf("%v", action), func(args tree.ResolveArguments) tree.Action {
			return action
		})
	}

	fFunc := func() tree.FloatFunctionNode {
		chosenOption := rand.Intn(len(functions.AllMathOps) + 1)

		if chosenOption < len(functions.AllMathOps) {
			return functions.NewMath(functions.AllMathOps[chosenOption], fTerm(), fTerm())
		}

		return functions.NewIfF(
			bTerm(),
			fTerm(),
			fTerm(),
		)
	}

	aFunc := func() tree.ActionFunctionNode {
		return functions.NewIfA(bTerm(), aTerm(), aTerm())
	}

	return &Generator{
		fFunc: fFunc,
		bFunc: bFunc,
		aFunc: aFunc,
		fTerm: fTerm,
		bTerm: bTerm,
		aTerm: aTerm,
	}
}
