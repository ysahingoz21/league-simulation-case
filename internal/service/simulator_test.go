package service

import (
	"math/rand"
	"testing"

	"league-simulation-case/internal/domain"
)

func TestStrengthBasedSimulatorProducesBoundedScores(t *testing.T) {
	simulator := NewStrengthBasedSimulator(rand.NewSource(42))
	home := domain.Team{ID: 1, Name: "A Team", Strength: 90}
	away := domain.Team{ID: 2, Name: "B Team", Strength: 80}

	for i := 0; i < 100; i++ {
		homeGoals, awayGoals := simulator.SimulateMatch(home, away)
		if homeGoals < 0 || awayGoals < 0 {
			t.Fatalf("expected non-negative goals, got %d-%d", homeGoals, awayGoals)
		}

		if homeGoals > 5 || awayGoals > 5 {
			t.Fatalf("expected realistic score bounds, got %d-%d", homeGoals, awayGoals)
		}
	}
}
