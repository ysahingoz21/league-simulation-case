package service

import (
	"context"
	"testing"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository/sqlite"
)

func TestPredictionServiceBeforeWeekFourReturnsMessage(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	svc := NewPredictionService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		sqlite.NewPredictionRepository(db),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
		10,
	)

	result, err := svc.GetPredictions(context.Background())
	if err != nil {
		t.Fatalf("get predictions before week 4: %v", err)
	}

	if result.Message == "" {
		t.Fatal("expected informative message before week 4")
	}

	if len(result.Predictions) != 0 {
		t.Fatalf("expected no predictions before week 4, got %d", len(result.Predictions))
	}
}

func TestPredictionServiceGeneratesAtWeekFour(t *testing.T) {
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
			&fakeSimulator{results: [][2]int{
				{1, 0}, {0, 0},
				{2, 1}, {1, 1},
				{0, 1}, {3, 0},
				{1, 2}, {2, 2},
				{0, 0}, {1, 0},
				{2, 2}, {1, 3},
				{1, 0}, {0, 1},
				{2, 0}, {1, 1},
			}},
			20,
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

	predictionSvc := NewPredictionService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		sqlite.NewPredictionRepository(db),
		&fakeSimulator{results: [][2]int{
			{1, 0}, {0, 1},
			{1, 1}, {2, 0},
		}},
		20,
	)

	result, err := predictionSvc.GetPredictions(context.Background())
	if err != nil {
		t.Fatalf("get predictions at week 4: %v", err)
	}

	if len(result.Predictions) != domain.TotalTeams {
		t.Fatalf("expected %d predictions, got %d", domain.TotalTeams, len(result.Predictions))
	}

	totalProbability := 0.0
	for _, prediction := range result.Predictions {
		if prediction.ChampionshipProbability < 0 || prediction.ChampionshipProbability > 100 {
			t.Fatalf("probability out of range: %+v", prediction)
		}
		totalProbability += prediction.ChampionshipProbability
	}

	if totalProbability < 99 || totalProbability > 101 {
		t.Fatalf("expected probabilities to sum approximately 100, got %f", totalProbability)
	}
}

func TestPredictionServiceCompletedLeagueReturnsDeterministicChampion(t *testing.T) {
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
			&fakeSimulator{results: [][2]int{{1, 0}, {0, 0}}},
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

	predictionSvc := NewPredictionService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		sqlite.NewPredictionRepository(db),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
		10,
	)

	result, err := predictionSvc.GetPredictions(context.Background())
	if err != nil {
		t.Fatalf("get predictions after completion: %v", err)
	}

	if len(result.Predictions) != domain.TotalTeams {
		t.Fatalf("expected %d predictions, got %d", domain.TotalTeams, len(result.Predictions))
	}

	if result.Predictions[0].ChampionshipProbability != 100 {
		t.Fatalf("expected champion probability 100, got %+v", result.Predictions[0])
	}

	for i := 1; i < len(result.Predictions); i++ {
		if result.Predictions[i].ChampionshipProbability != 0 {
			t.Fatalf("expected non-champions to have 0 probability, got %+v", result.Predictions[i])
		}
	}
}

func TestPredictionServiceReturnsPersistedLatestPredictions(t *testing.T) {
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

	predictionSvc := NewPredictionService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		sqlite.NewPredictionRepository(db),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
		20,
	)

	first, err := predictionSvc.GetPredictions(context.Background())
	if err != nil {
		t.Fatalf("get first predictions: %v", err)
	}

	second, err := predictionSvc.GetPredictions(context.Background())
	if err != nil {
		t.Fatalf("get second predictions: %v", err)
	}

	if len(first.Predictions) != len(second.Predictions) {
		t.Fatalf("expected same number of predictions, got %d and %d", len(first.Predictions), len(second.Predictions))
	}

	for i := range first.Predictions {
		if first.Predictions[i].Week != second.Predictions[i].Week ||
			first.Predictions[i].TeamID != second.Predictions[i].TeamID ||
			first.Predictions[i].ChampionshipProbability != second.Predictions[i].ChampionshipProbability {
			t.Fatalf("expected persisted latest predictions to be stable, first=%+v second=%+v", first.Predictions[i], second.Predictions[i])
		}
	}
}
