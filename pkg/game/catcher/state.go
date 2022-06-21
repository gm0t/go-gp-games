package catcher

import (
	"math"

	"lr1Go/pkg/game/catcher/actions"
	"lr1Go/pkg/tree"
)

const (
	MyX   = "MyX"
	MyY   = "MyY"
	GoalX = "GoalX"
	GoalY = "GoalY"
)

var FloatKeys = []string{
	MyX, MyY, GoalY, GoalX,
}

var BoolKeys = []string{}

type State struct {
	PlayerX float64
	PlayerY float64
	GoalX   float64
	GoalY   float64
}

func NewState(playerX float64, playerY float64, goalX float64, goalY float64) *State {
	return &State{PlayerX: playerX, PlayerY: playerY, GoalX: goalX, GoalY: goalY}
}

type Args struct {
	floatArgs map[string]float64
	boolArgs  map[string]bool
}

func (a Args) Action(key string) interface{} {
	return key
}

func (a Args) Float(key string) float64 {
	return a.floatArgs[key]
}

func (a Args) Boolean(key string) bool {
	return a.boolArgs[key]
}

func (s *State) Apply(action actions.Action) *State {
	switch action {
	case actions.Up:
		s.PlayerY += 1
	case actions.Down:
		s.PlayerY -= 1
	case actions.Right:
		s.PlayerX += 1
	case actions.Left:
		s.PlayerX -= 1
	}

	return s
}

func (s *State) Clone() *State {
	return &State{
		PlayerX: s.PlayerX,
		PlayerY: s.PlayerY,
		GoalX:   s.GoalY,
		GoalY:   s.GoalY,
	}
}

func (s *State) BuildArgs() tree.ResolveArguments {
	return Args{
		floatArgs: map[string]float64{
			MyX:   s.PlayerX,
			MyY:   s.PlayerY,
			GoalX: s.PlayerX,
			GoalY: s.PlayerY,
		},
		boolArgs: map[string]bool{},
	}
}

func (s *State) GetCurrentDistance() float64 {
	return math.Abs(s.GoalX-s.PlayerX) + math.Abs(s.GoalY-s.PlayerY)
}
