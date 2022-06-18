package utils

import (
	"fmt"
	"math"
	"math/rand"

	"lr1Go/pkg/old-tree"
	"lr1Go/pkg/old-tree/functions"
	"lr1Go/pkg/old-tree/terminals"
)

type fTermGenerator = func() old_tree.FloatNode
type bTermGenerator = func() old_tree.BooleanNode

type fFuncGenerator = func() old_tree.FloatFunctionNode
type bFuncGenerator = func() old_tree.BooleanFunctionNode

func getRandom[G any](terms []G) G {
	idx := rand.Intn(len(terms))
	return terms[idx]
}

type Generator struct {
	fFunc fFuncGenerator
	bFunc bFuncGenerator
	fTerm fTermGenerator
	bTerm bTermGenerator
	aFunc func() old_tree.ActionFunctionNode
	aTerm func() old_tree.ActionNode
}

func (g *Generator) BTree(maxDepth int) old_tree.BooleanFunctionNode {
	root := g.bFunc()

	root.Dfs(func(depth int, n old_tree.Node) {
		if f, isFunction := n.(old_tree.FunctionNode); isFunction && (depth < maxDepth) {
			f.Grow(g)
		}
	})

	return root
}

func (g *Generator) FFunc() old_tree.FloatFunctionNode {
	return g.fFunc()
}

func (g *Generator) AFunc() old_tree.ActionFunctionNode {
	return g.aFunc()
}

func (g *Generator) BFunc() old_tree.BooleanFunctionNode {
	return g.bFunc()
}

func (g *Generator) FTerm() old_tree.FloatNode {
	return g.fTerm()
}

func (g *Generator) ATerm() old_tree.ActionNode {
	return g.aTerm()
}

func (g *Generator) BTerm() old_tree.BooleanNode {
	return g.bTerm()
}

func (g *Generator) FTree(maxDepth int) old_tree.FloatFunctionNode {
	root := g.fFunc()

	root.Dfs(func(depth int, n old_tree.Node) {
		if f, isFunction := n.(old_tree.FunctionNode); isFunction && (depth < maxDepth) {
			f.Grow(g)
		}
	})

	return root
}

func (g *Generator) ATree(maxDepth int) old_tree.ActionFunctionNode {
	root := g.aFunc()

	root.Dfs(func(depth int, n old_tree.Node) {
		if f, isFunction := n.(old_tree.FunctionNode); isFunction && (depth < maxDepth) {
			f.Grow(g)
		}
	})

	return root
}

func NodeGenerator(floatKeys []string, booleanKeys []string, actions []old_tree.Action) *Generator {
	fTerm := func() old_tree.FloatNode {
		rndValue := math.Round(rand.Float64()*1000) / 100
		rndKey := fmt.Sprintf("rnd[%.2f]", rndValue)

		allFloatKeys := make([]string, len(floatKeys)+1)
		copy(allFloatKeys, floatKeys)
		allFloatKeys[len(allFloatKeys)-1] = rndKey

		chosenKey := getRandom(floatKeys)

		return terminals.NewFloatTerm(chosenKey, func(args old_tree.ResolveArguments) float64 {
			if chosenKey == rndKey {
				return rndValue
			}

			return args.Float(chosenKey)
		})
	}

	bFunc := func() old_tree.BooleanFunctionNode {
		return functions.NewComparison(getRandom(functions.AllCompOps), fTerm(), fTerm())
	}

	bTerm := func() old_tree.BooleanNode {
		if len(booleanKeys) == 0 {
			return bFunc()
		}

		chosenKey := getRandom(booleanKeys)
		return terminals.NewBooleanTerm(chosenKey, func(args old_tree.ResolveArguments) bool {
			return args.Boolean(chosenKey)
		})
	}

	aTerm := func() old_tree.ActionNode {
		action := getRandom(actions)
		return terminals.NewActionTerm(fmt.Sprintf("%v", action), func(args old_tree.ResolveArguments) old_tree.Action {
			return action
		})
	}

	fFunc := func() old_tree.FloatFunctionNode {
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

	aFunc := func() old_tree.ActionFunctionNode {
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
