package agent

import (
	"lr1Go/pkg/game/catcher"
	"lr1Go/pkg/game/catcher/actions"
	"lr1Go/pkg/tree"
	"lr1Go/pkg/utils"
)

func NewGenerator() tree.Generator {
	return utils.NodeGenerator(
		catcher.FloatKeys,
		catcher.BoolKeys,
		[]tree.Action{
			actions.Up,
			actions.Down,
			actions.Left,
			actions.Right,
		},
	)
}
