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
	Actions        []string
	State          *State
}

func (g *Game) Run(maxIterations int) Result {
	state := g.state
	player := g.player
	distance := state.GetCurrentDistance()
	actions := make([]string, 0)
	i := 1
	for i < maxIterations && distance > 0 {
		action := player.GetAction(state.Clone())
		state.Apply(action)
		actions = append(actions, string(action))
		distance = state.GetCurrentDistance()
		i += 1
	}

	return Result{
		Iterations:     i,
		Actions:        actions,
		State:          state.Clone(),
		DistanceToGoal: distance,
	}
}
