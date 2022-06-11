package catcher

type Game struct {
	state  *State
	player Player
}

func NewGame(state *State, player Player) *Game {
	return &Game{state: state, player: player}
}

type Result struct {
	Iterations     int
	DistanceToGoal float64
}

func (g *Game) Run(maxIterations int) Result {
	state := g.state
	player := g.player
	distance := state.GetCurrentDistance()
	i := 0
	for i < maxIterations && distance > 0 {
		action := player.GetAction(state.Clone())
		state.Apply(action)
		distance = state.GetCurrentDistance()
		i += 1
	}

	return Result{
		Iterations:     i,
		DistanceToGoal: distance,
	}
}
