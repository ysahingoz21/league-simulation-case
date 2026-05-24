package domain

import (
	"reflect"
	"testing"
)

func TestSortStandings(t *testing.T) {
	standings := []Standing{
		{TeamID: 1, TeamName: "D Team", Points: 7, GoalDifference: 2, GoalsFor: 6},
		{TeamID: 2, TeamName: "A Team", Points: 9, GoalDifference: 1, GoalsFor: 5},
		{TeamID: 3, TeamName: "C Team", Points: 9, GoalDifference: 3, GoalsFor: 4},
		{TeamID: 4, TeamName: "B Team", Points: 9, GoalDifference: 3, GoalsFor: 7},
		{TeamID: 5, TeamName: "E Team", Points: 9, GoalDifference: 3, GoalsFor: 7},
	}

	SortStandings(standings)

	got := []string{
		standings[0].TeamName,
		standings[1].TeamName,
		standings[2].TeamName,
		standings[3].TeamName,
		standings[4].TeamName,
	}

	want := []string{"B Team", "E Team", "C Team", "A Team", "D Team"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected order %v, got %v", want, got)
	}
}

func TestAssignRanks(t *testing.T) {
	standings := []Standing{
		{TeamName: "A Team"},
		{TeamName: "B Team"},
		{TeamName: "C Team"},
	}

	AssignRanks(standings)

	for i, standing := range standings {
		wantRank := i + 1
		if standing.Rank != wantRank {
			t.Fatalf("expected team %s to have rank %d, got %d", standing.TeamName, wantRank, standing.Rank)
		}
	}
}

func TestStandingApplyMatch(t *testing.T) {
	homeGoals := 2
	awayGoals := 1
	team := Team{ID: 10, Name: "A Team"}
	match := Match{
		HomeTeamID: 10,
		AwayTeamID: 20,
		HomeGoals:  &homeGoals,
		AwayGoals:  &awayGoals,
		Status:     MatchStatusPlayed,
	}

	got := (Standing{}).ApplyMatch(match, team)

	want := Standing{
		TeamID:         10,
		TeamName:       "A Team",
		Played:         1,
		Wins:           1,
		Draws:          0,
		Losses:         0,
		GoalsFor:       2,
		GoalsAgainst:   1,
		GoalDifference: 1,
		Points:         3,
		Rank:           0,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected standing %+v, got %+v", want, got)
	}
}
