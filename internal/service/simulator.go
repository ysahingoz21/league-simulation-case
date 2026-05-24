package service

import (
	"math/rand"
	"time"

	"league-simulation-case/internal/domain"
)

type Simulator interface {
	SimulateMatch(homeTeam domain.Team, awayTeam domain.Team) (homeGoals int, awayGoals int)
}

type StrengthBasedSimulator struct {
	rng *rand.Rand
}

func NewStrengthBasedSimulator(source rand.Source) *StrengthBasedSimulator {
	if source == nil {
		source = rand.NewSource(time.Now().UnixNano())
	}

	return &StrengthBasedSimulator{
		rng: rand.New(source),
	}
}

func (s *StrengthBasedSimulator) SimulateMatch(homeTeam domain.Team, awayTeam domain.Team) (int, int) {
	homeGoals := s.generateGoals(homeTeam.Strength+5, awayTeam.Strength)
	awayGoals := s.generateGoals(awayTeam.Strength, homeTeam.Strength)

	return homeGoals, awayGoals
}

func (s *StrengthBasedSimulator) generateGoals(attackStrength, defenseStrength int) int {
	goals := 0
	for chance := 0; chance < 5; chance++ {
		threshold := (attackStrength / 2) - (defenseStrength / 4) + 10 - (chance * 8)
		if s.rng.Intn(100) < clamp(threshold, 0, 70) {
			goals++
		}
	}

	return goals
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
