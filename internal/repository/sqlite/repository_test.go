package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	dbpkg "league-simulation-case/internal/database"
	"league-simulation-case/internal/domain"
)

func TestTeamRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	repo := NewTeamRepository(db)
	teams := []domain.Team{
		{Name: "A Team", Strength: 90},
		{Name: "B Team", Strength: 85},
	}

	if err := repo.CreateMany(ctx, teams); err != nil {
		t.Fatalf("create many teams: %v", err)
	}

	gotTeams, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list teams: %v", err)
	}

	if len(gotTeams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(gotTeams))
	}

	gotTeam, err := repo.GetByID(ctx, gotTeams[0].ID)
	if err != nil {
		t.Fatalf("get team by id: %v", err)
	}

	if gotTeam.Name != "A Team" {
		t.Fatalf("expected team name A Team, got %s", gotTeam.Name)
	}

	if err := repo.DeleteAll(ctx); err != nil {
		t.Fatalf("delete all teams: %v", err)
	}

	gotTeams, err = repo.List(ctx)
	if err != nil {
		t.Fatalf("list teams after delete: %v", err)
	}

	if len(gotTeams) != 0 {
		t.Fatalf("expected 0 teams after delete, got %d", len(gotTeams))
	}
}

func TestLeagueRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	repo := NewLeagueRepository(db)
	initialState := domain.LeagueState{
		CurrentWeek: 0,
		TotalWeeks:  domain.TotalWeeks,
		IsCompleted: false,
	}

	if err := repo.UpsertState(ctx, initialState); err != nil {
		t.Fatalf("upsert league state: %v", err)
	}

	state, err := repo.GetState(ctx)
	if err != nil {
		t.Fatalf("get league state: %v", err)
	}

	if state.TotalWeeks != domain.TotalWeeks {
		t.Fatalf("expected total weeks %d, got %d", domain.TotalWeeks, state.TotalWeeks)
	}

	if err := repo.UpdateCurrentWeek(ctx, 4, false); err != nil {
		t.Fatalf("update current week: %v", err)
	}

	state, err = repo.GetState(ctx)
	if err != nil {
		t.Fatalf("get updated league state: %v", err)
	}

	if state.CurrentWeek != 4 || state.IsCompleted {
		t.Fatalf("unexpected updated state: %+v", state)
	}

	if err := repo.Reset(ctx); err != nil {
		t.Fatalf("reset league state: %v", err)
	}

	state, err = repo.GetState(ctx)
	if err != nil {
		t.Fatalf("get reset league state: %v", err)
	}

	if state.CurrentWeek != 0 || state.IsCompleted || state.TotalWeeks != domain.TotalWeeks {
		t.Fatalf("unexpected reset state: %+v", state)
	}
}

func TestMatchRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	insertTeams(t, ctx, db)

	repo := NewMatchRepository(db)
	playedAt := time.Date(2026, time.May, 24, 12, 0, 0, 0, time.UTC)
	homeGoals := 2
	awayGoals := 1
	matches := []domain.Match{
		{
			Week:       1,
			HomeTeamID: 1,
			AwayTeamID: 2,
			Status:     domain.MatchStatusScheduled,
		},
		{
			Week:       1,
			HomeTeamID: 3,
			AwayTeamID: 4,
			HomeGoals:  &homeGoals,
			AwayGoals:  &awayGoals,
			Status:     domain.MatchStatusPlayed,
			PlayedAt:   &playedAt,
		},
		{
			Week:       2,
			HomeTeamID: 1,
			AwayTeamID: 3,
			Status:     domain.MatchStatusScheduled,
		},
	}

	if err := repo.CreateMany(ctx, matches); err != nil {
		t.Fatalf("create many matches: %v", err)
	}

	weekOneMatches, err := repo.ListByWeek(ctx, 1)
	if err != nil {
		t.Fatalf("list matches by week: %v", err)
	}

	if len(weekOneMatches) != 2 {
		t.Fatalf("expected 2 week-one matches, got %d", len(weekOneMatches))
	}

	match, err := repo.GetByID(ctx, weekOneMatches[0].ID)
	if err != nil {
		t.Fatalf("get match by id: %v", err)
	}

	if match.Week != 1 {
		t.Fatalf("expected week 1 match, got week %d", match.Week)
	}

	if err := repo.UpdateResult(ctx, weekOneMatches[0].ID, 3, 2, playedAt); err != nil {
		t.Fatalf("update match result: %v", err)
	}

	match, err = repo.GetByID(ctx, weekOneMatches[0].ID)
	if err != nil {
		t.Fatalf("get updated match: %v", err)
	}

	if !match.IsPlayed() || *match.HomeGoals != 3 || *match.AwayGoals != 2 {
		t.Fatalf("unexpected updated match: %+v", match)
	}

	unplayedMatches, err := repo.ListUnplayed(ctx)
	if err != nil {
		t.Fatalf("list unplayed matches: %v", err)
	}

	if len(unplayedMatches) != 1 {
		t.Fatalf("expected 1 unplayed match, got %d", len(unplayedMatches))
	}

	if err := repo.DeleteAll(ctx); err != nil {
		t.Fatalf("delete all matches: %v", err)
	}

	allMatches, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list matches after delete: %v", err)
	}

	if len(allMatches) != 0 {
		t.Fatalf("expected 0 matches after delete, got %d", len(allMatches))
	}
}

func TestStandingRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	insertTeams(t, ctx, db)

	repo := NewStandingRepository(db)
	standings := []domain.Standing{
		{TeamID: 1, Played: 1, Wins: 1, GoalsFor: 2, GoalsAgainst: 1, GoalDifference: 1, Points: 3, Rank: 1},
		{TeamID: 2, Played: 1, Losses: 1, GoalsFor: 1, GoalsAgainst: 2, GoalDifference: -1, Points: 0, Rank: 2},
	}

	if err := repo.ReplaceAll(ctx, standings); err != nil {
		t.Fatalf("replace standings: %v", err)
	}

	gotStandings, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list standings: %v", err)
	}

	if len(gotStandings) != 2 {
		t.Fatalf("expected 2 standings, got %d", len(gotStandings))
	}

	if gotStandings[0].TeamName != "A Team" {
		t.Fatalf("expected first standing team A Team, got %s", gotStandings[0].TeamName)
	}

	if err := repo.DeleteAll(ctx); err != nil {
		t.Fatalf("delete standings: %v", err)
	}

	gotStandings, err = repo.List(ctx)
	if err != nil {
		t.Fatalf("list standings after delete: %v", err)
	}

	if len(gotStandings) != 0 {
		t.Fatalf("expected 0 standings after delete, got %d", len(gotStandings))
	}
}

func TestPredictionRepository(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)
	insertTeams(t, ctx, db)

	repo := NewPredictionRepository(db)
	weekFour := []domain.Prediction{
		{Week: 4, TeamID: 1, ChampionshipProbability: 60, ExpectedPoints: 11.5, ProjectedRank: 1},
		{Week: 4, TeamID: 2, ChampionshipProbability: 40, ExpectedPoints: 10.0, ProjectedRank: 2},
	}
	weekFive := []domain.Prediction{
		{Week: 5, TeamID: 2, ChampionshipProbability: 55, ExpectedPoints: 12.0, ProjectedRank: 1},
		{Week: 5, TeamID: 1, ChampionshipProbability: 45, ExpectedPoints: 11.0, ProjectedRank: 2},
	}

	if err := repo.ReplaceForWeek(ctx, 4, weekFour); err != nil {
		t.Fatalf("replace week 4 predictions: %v", err)
	}

	if err := repo.ReplaceForWeek(ctx, 5, weekFive); err != nil {
		t.Fatalf("replace week 5 predictions: %v", err)
	}

	gotPredictions, err := repo.ListLatest(ctx)
	if err != nil {
		t.Fatalf("list latest predictions: %v", err)
	}

	if len(gotPredictions) != 2 {
		t.Fatalf("expected 2 latest predictions, got %d", len(gotPredictions))
	}

	if gotPredictions[0].Week != 5 || gotPredictions[0].TeamName != "B Team" {
		t.Fatalf("unexpected latest prediction: %+v", gotPredictions[0])
	}

	if err := repo.DeleteAll(ctx); err != nil {
		t.Fatalf("delete predictions: %v", err)
	}

	gotPredictions, err = repo.ListLatest(ctx)
	if err != nil {
		t.Fatalf("list predictions after delete: %v", err)
	}

	if len(gotPredictions) != 0 {
		t.Fatalf("expected 0 predictions after delete, got %d", len(gotPredictions))
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open test sqlite db: %v", err)
	}

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "..", "database", "schema.sql")); err != nil {
		db.Close()
		t.Fatalf("apply schema: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func insertTeams(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	repo := NewTeamRepository(db)
	err := repo.CreateMany(ctx, []domain.Team{
		{Name: "A Team", Strength: 90},
		{Name: "B Team", Strength: 85},
		{Name: "C Team", Strength: 80},
		{Name: "D Team", Strength: 75},
	})
	if err != nil {
		t.Fatalf("insert teams: %v", err)
	}
}

func TestRepositoriesReturnNoRowsClearly(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	teamRepo := NewTeamRepository(db)
	if _, err := teamRepo.GetByID(ctx, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for missing team, got %v", err)
	}

	matchRepo := NewMatchRepository(db)
	if _, err := matchRepo.GetByID(ctx, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for missing match, got %v", err)
	}
}
