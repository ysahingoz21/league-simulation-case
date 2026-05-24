package service

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	dbpkg "league-simulation-case/internal/database"
	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository/sqlite"
)

func TestLeagueServiceInitialize(t *testing.T) {
	ctx := context.Background()
	db := newServiceTestDB(t)

	svc := newTestLeagueService(db)
	result, err := svc.Initialize(ctx)
	if err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	if len(result.Teams) != domain.TotalTeams {
		t.Fatalf("expected %d teams, got %d", domain.TotalTeams, len(result.Teams))
	}

	if len(result.Fixtures) != domain.TotalWeeks*domain.MatchesPerWeek {
		t.Fatalf("expected %d fixtures, got %d", domain.TotalWeeks*domain.MatchesPerWeek, len(result.Fixtures))
	}

	if len(result.Standings) != domain.TotalTeams {
		t.Fatalf("expected %d standings, got %d", domain.TotalTeams, len(result.Standings))
	}

	if result.League.CurrentWeek != 0 || result.League.TotalWeeks != domain.TotalWeeks || result.League.IsCompleted {
		t.Fatalf("unexpected league state: %+v", result.League)
	}

	for _, fixture := range result.Fixtures {
		if fixture.Status != domain.MatchStatusScheduled {
			t.Fatalf("expected scheduled fixture, got %s", fixture.Status)
		}
		if fixture.HomeGoals != nil || fixture.AwayGoals != nil {
			t.Fatalf("expected nil scores for scheduled fixture: %+v", fixture)
		}
	}
}

func TestLeagueServiceListFixturesByWeekValidation(t *testing.T) {
	svc := newTestLeagueService(newServiceTestDB(t))

	if _, err := svc.ListFixturesByWeek(context.Background(), 0); err != ErrInvalidFixtureWeek {
		t.Fatalf("expected ErrInvalidFixtureWeek for week 0, got %v", err)
	}

	if _, err := svc.ListFixturesByWeek(context.Background(), 7); err != ErrInvalidFixtureWeek {
		t.Fatalf("expected ErrInvalidFixtureWeek for week 7, got %v", err)
	}
}

func newServiceTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "service-test.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open service test db: %v", err)
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

func newTestLeagueService(db *sql.DB) LeagueService {
	return NewLeagueService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		sqlite.NewStandingRepository(db),
		sqlite.NewPredictionRepository(db),
	)
}
