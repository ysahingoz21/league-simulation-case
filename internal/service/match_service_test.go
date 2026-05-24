package service

import (
	"context"
	"testing"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository/sqlite"
)

func TestMatchServiceUpdatePlayedMatchRecalculatesStandings(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	simSvc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
		&fakeSimulator{results: [][2]int{{2, 0}, {1, 1}}},
	)
	if _, err := simSvc.PlayWeek(context.Background(), 1); err != nil {
		t.Fatalf("play week 1: %v", err)
	}

	matchSvc := NewMatchService(
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
	)

	result, err := matchSvc.UpdateMatchResult(context.Background(), 1, 0, 3)
	if err != nil {
		t.Fatalf("update match result: %v", err)
	}

	if result.Match.HomeGoals == nil || result.Match.AwayGoals == nil || *result.Match.HomeGoals != 0 || *result.Match.AwayGoals != 3 {
		t.Fatalf("unexpected updated match: %+v", result.Match)
	}

	if result.Standings[0].TeamName != "B Team" || result.Standings[0].Points != 3 {
		t.Fatalf("expected B Team to lead after edit, got %+v", result.Standings[0])
	}
}

func TestMatchServiceUpdateAfterWeekFourRefreshesPredictions(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	simSvc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
		&fakeSimulator{results: [][2]int{
			{1, 0}, {0, 0},
			{2, 1}, {1, 1},
			{0, 1}, {3, 0},
			{1, 2}, {2, 2},
		}},
	)
	for week := 1; week <= 4; week++ {
		if _, err := simSvc.PlayWeek(context.Background(), week); err != nil {
			t.Fatalf("play week %d: %v", week, err)
		}
	}

	matchSvc := NewMatchService(
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
	)

	result, err := matchSvc.UpdateMatchResult(context.Background(), 1, 4, 0)
	if err != nil {
		t.Fatalf("update match result after week 4: %v", err)
	}

	if len(result.Predictions) != domain.TotalTeams {
		t.Fatalf("expected %d refreshed predictions, got %d", domain.TotalTeams, len(result.Predictions))
	}
}

func TestMatchServiceUpdateCompletedLeagueReturnsDeterministicPredictions(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	simSvc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
		&fakeSimulator{results: [][2]int{
			{3, 0}, {0, 0},
			{2, 0}, {1, 0},
			{2, 1}, {2, 0},
			{1, 0}, {1, 0},
			{2, 0}, {0, 0},
			{1, 0}, {1, 0},
		}},
	)
	if _, err := simSvc.PlayAll(context.Background()); err != nil {
		t.Fatalf("play all: %v", err)
	}

	matchSvc := NewMatchService(
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
	)

	result, err := matchSvc.UpdateMatchResult(context.Background(), 12, 0, 0)
	if err != nil {
		t.Fatalf("update match in completed league: %v", err)
	}

	if len(result.Predictions) != domain.TotalTeams {
		t.Fatalf("expected %d predictions, got %d", domain.TotalTeams, len(result.Predictions))
	}

	if result.Predictions[0].ChampionshipProbability != 100 {
		t.Fatalf("expected top prediction probability 100, got %+v", result.Predictions[0])
	}
}

func TestMatchServiceRejectsInvalidGoalsAndMissingMatch(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	matchSvc := NewMatchService(
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		NewPredictionService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewPredictionRepository(db),
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
			10,
		),
	)

	if _, err := matchSvc.UpdateMatchResult(context.Background(), 1, -1, 2); err != ErrInvalidGoals {
		t.Fatalf("expected ErrInvalidGoals, got %v", err)
	}

	if _, err := matchSvc.UpdateMatchResult(context.Background(), 999, 1, 2); err != ErrMatchNotFound {
		t.Fatalf("expected ErrMatchNotFound, got %v", err)
	}
}
