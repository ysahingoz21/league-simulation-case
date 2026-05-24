package service

import (
	"fmt"

	"league-simulation-case/internal/domain"
)

type fixturePair struct {
	homeIndex int
	awayIndex int
}

var fixturePattern = [][]fixturePair{
	{{homeIndex: 0, awayIndex: 1}, {homeIndex: 2, awayIndex: 3}},
	{{homeIndex: 0, awayIndex: 2}, {homeIndex: 1, awayIndex: 3}},
	{{homeIndex: 0, awayIndex: 3}, {homeIndex: 1, awayIndex: 2}},
	{{homeIndex: 1, awayIndex: 0}, {homeIndex: 3, awayIndex: 2}},
	{{homeIndex: 2, awayIndex: 0}, {homeIndex: 3, awayIndex: 1}},
	{{homeIndex: 3, awayIndex: 0}, {homeIndex: 2, awayIndex: 1}},
}

func GenerateDoubleRoundRobinFixtures(teams []domain.Team) ([]domain.Match, error) {
	if len(teams) != domain.TotalTeams {
		return nil, fmt.Errorf("expected %d teams, got %d", domain.TotalTeams, len(teams))
	}

	fixtures := make([]domain.Match, 0, domain.TotalWeeks*domain.MatchesPerWeek)
	for weekIndex, weekPairs := range fixturePattern {
		week := weekIndex + 1
		for _, pair := range weekPairs {
			fixtures = append(fixtures, domain.Match{
				Week:       week,
				HomeTeamID: teams[pair.homeIndex].ID,
				AwayTeamID: teams[pair.awayIndex].ID,
				Status:     domain.MatchStatusScheduled,
			})
		}
	}

	return fixtures, nil
}
