package domain

type SimulationResult struct {
	Week        int
	Matches     []Match
	Standings   []Standing
	Predictions []Prediction
}
