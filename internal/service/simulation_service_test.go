package service

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	dbpkg "league-simulation-case/internal/database"
	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository/sqlite"
)

type fakeSimulator struct {
	results [][2]int
	index   int
}

func (f *fakeSimulator) SimulateMatch(homeTeam domain.Team, awayTeam domain.Team) (int, int) {
	result := f.results[f.index]
	f.index++
	return result[0], result[1]
}

func TestSimulationServicePlayNextWeek(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{{2, 1}, {0, 0}}},
	).(*simulationService)
	svc.now = func() time.Time { return time.Date(2026, time.May, 24, 13, 0, 0, 0, time.UTC) }

	result, err := svc.PlayNextWeek(context.Background())
	if err != nil {
		t.Fatalf("play next week: %v", err)
	}

	if result.Week != 1 {
		t.Fatalf("expected week 1, got %d", result.Week)
	}

	if result.League.CurrentWeek != 1 || result.League.IsCompleted {
		t.Fatalf("unexpected league state: %+v", result.League)
	}

	if len(result.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(result.Matches))
	}

	if len(result.Standings) != domain.TotalTeams {
		t.Fatalf("expected %d standings, got %d", domain.TotalTeams, len(result.Standings))
	}

	if result.Standings[0].Points != 3 || result.Standings[0].TeamName != "A Team" {
		t.Fatalf("unexpected leading standing: %+v", result.Standings[0])
	}

	for _, match := range result.Matches {
		if !match.IsPlayed() {
			t.Fatalf("expected played match, got %+v", match)
		}
	}
}

func TestSimulationServiceWeekCannotBeReplayed(t *testing.T) {
	db := newSimulationTestDB(t)
	initializeAndPlayWeek(t, db, 1)

	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{{1, 0}, {1, 1}}},
	)

	if _, err := svc.PlayWeek(context.Background(), 1); err != ErrWeekAlreadyPlayed {
		t.Fatalf("expected ErrWeekAlreadyPlayed, got %v", err)
	}
}

func TestSimulationServiceCannotPlayWeekOutOfOrder(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
	)

	if _, err := svc.PlayWeek(context.Background(), 3); err != ErrWeekOutOfOrder {
		t.Fatalf("expected ErrWeekOutOfOrder, got %v", err)
	}
}

func TestSimulationServicePlayingFinalWeekMarksLeagueCompleted(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{
			{1, 0}, {0, 0},
			{2, 1}, {1, 1},
			{0, 1}, {3, 0},
			{1, 2}, {2, 2},
			{0, 0}, {1, 0},
			{2, 2}, {1, 3},
		}},
	).(*simulationService)
	svc.now = func() time.Time { return time.Date(2026, time.May, 24, 13, 0, 0, 0, time.UTC) }

	for week := 1; week <= domain.TotalWeeks; week++ {
		result, err := svc.PlayWeek(context.Background(), week)
		if err != nil {
			t.Fatalf("play week %d: %v", week, err)
		}

		if week < domain.TotalWeeks && result.League.IsCompleted {
			t.Fatalf("league should not be completed at week %d", week)
		}
	}

	finalState := sqlite.NewLeagueRepository(db)
	state, err := finalState.GetState(context.Background())
	if err != nil {
		t.Fatalf("get final league state: %v", err)
	}

	if !state.IsCompleted || state.CurrentWeek != domain.TotalWeeks {
		t.Fatalf("unexpected completed state: %+v", state)
	}
}

func TestSimulationServiceCompletedLeagueCannotPlayNextWeek(t *testing.T) {
	db := newSimulationTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	leagueRepo := sqlite.NewLeagueRepository(db)
	if err := leagueRepo.UpdateCurrentWeek(context.Background(), domain.TotalWeeks, true); err != nil {
		t.Fatalf("set completed state: %v", err)
	}

	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
	)

	if _, err := svc.PlayNextWeek(context.Background()); err != ErrLeagueCompleted {
		t.Fatalf("expected ErrLeagueCompleted, got %v", err)
	}
}

func TestSimulationServiceRequiresInitializedLeague(t *testing.T) {
	db := newSimulationTestDB(t)
	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 1}}},
	)

	if _, err := svc.PlayNextWeek(context.Background()); err != ErrLeagueNotInitialized {
		t.Fatalf("expected ErrLeagueNotInitialized, got %v", err)
	}
}

func newSimulationTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "simulation-test.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open simulation test db: %v", err)
	}

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		db.Close()
		t.Fatalf("apply schema: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func initializeAndPlayWeek(t *testing.T, db *sql.DB, week int) {
	t.Helper()

	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	svc := NewSimulationService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		NewStandingsService(
			sqlite.NewTeamRepository(db),
			sqlite.NewMatchRepository(db),
			sqlite.NewLeagueRepository(db),
			sqlite.NewStandingRepository(db),
		),
		&fakeSimulator{results: [][2]int{{1, 0}, {0, 0}}},
	)

	if _, err := svc.PlayWeek(context.Background(), week); err != nil {
		t.Fatalf("play week %d: %v", week, err)
	}
}
