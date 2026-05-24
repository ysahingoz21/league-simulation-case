package service

import (
	"fmt"
	"testing"

	"league-simulation-case/internal/domain"
)

func TestGenerateDoubleRoundRobinFixtures(t *testing.T) {
	teams := []domain.Team{
		{ID: 1, Name: "A Team"},
		{ID: 2, Name: "B Team"},
		{ID: 3, Name: "C Team"},
		{ID: 4, Name: "D Team"},
	}

	fixtures, err := GenerateDoubleRoundRobinFixtures(teams)
	if err != nil {
		t.Fatalf("generate fixtures: %v", err)
	}

	if len(fixtures) != 12 {
		t.Fatalf("expected 12 fixtures, got %d", len(fixtures))
	}

	weeks := map[int]int{}
	weekTeamCounts := map[int]map[int64]int{}
	pairCounts := map[string]int{}
	reversedFound := map[string]bool{}

	for _, fixture := range fixtures {
		weeks[fixture.Week]++

		if weekTeamCounts[fixture.Week] == nil {
			weekTeamCounts[fixture.Week] = map[int64]int{}
		}
		weekTeamCounts[fixture.Week][fixture.HomeTeamID]++
		weekTeamCounts[fixture.Week][fixture.AwayTeamID]++

		pairKey := unorderedPairKey(fixture.HomeTeamID, fixture.AwayTeamID)
		pairCounts[pairKey]++

		reverseKey := fmt.Sprintf("%d-%d", fixture.AwayTeamID, fixture.HomeTeamID)
		if reversedFound[reverseKey] {
			reversedFound[fmt.Sprintf("%d-%d", fixture.HomeTeamID, fixture.AwayTeamID)] = true
		}
	}

	if len(weeks) != domain.TotalWeeks {
		t.Fatalf("expected %d weeks, got %d", domain.TotalWeeks, len(weeks))
	}

	for week := 1; week <= domain.TotalWeeks; week++ {
		if weeks[week] != domain.MatchesPerWeek {
			t.Fatalf("expected %d matches in week %d, got %d", domain.MatchesPerWeek, week, weeks[week])
		}

		if len(weekTeamCounts[week]) != domain.TotalTeams {
			t.Fatalf("expected %d teams in week %d, got %d", domain.TotalTeams, week, len(weekTeamCounts[week]))
		}

		for teamID, count := range weekTeamCounts[week] {
			if count != 1 {
				t.Fatalf("expected team %d to play once in week %d, got %d", teamID, week, count)
			}
		}
	}

	if len(pairCounts) != 6 {
		t.Fatalf("expected 6 unique pairs, got %d", len(pairCounts))
	}

	for pairKey, count := range pairCounts {
		if count != 2 {
			t.Fatalf("expected pair %s to play twice, got %d", pairKey, count)
		}
	}

	expectedReversePairs := [][2]int64{
		{1, 2},
		{3, 4},
		{1, 3},
		{2, 4},
		{1, 4},
		{2, 3},
	}

	for _, pair := range expectedReversePairs {
		if !containsFixture(fixtures, pair[0], pair[1]) || !containsFixture(fixtures, pair[1], pair[0]) {
			t.Fatalf("expected reversed home/away fixtures for pair %d-%d", pair[0], pair[1])
		}
	}
}

func unorderedPairKey(a, b int64) string {
	if a < b {
		return fmt.Sprintf("%d-%d", a, b)
	}

	return fmt.Sprintf("%d-%d", b, a)
}

func containsFixture(fixtures []domain.Match, homeID, awayID int64) bool {
	for _, fixture := range fixtures {
		if fixture.HomeTeamID == homeID && fixture.AwayTeamID == awayID {
			return true
		}
	}

	return false
}
