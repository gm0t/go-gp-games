package game

type State interface {
	FloatValues() map[string]float64
	BooleanValues() map[string]bool
}
