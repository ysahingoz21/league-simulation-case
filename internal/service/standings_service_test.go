package service

import (
	"context"
	"testing"
	"time"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository/sqlite"
)

func TestStandingsServiceRecalculate(t *testing.T) {
	db := newServiceTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	matchRepo := sqlite.NewMatchRepository(db)
	matchOne, err := matchRepo.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("get match 1: %v", err)
	}
	matchTwo, err := matchRepo.GetByID(context.Background(), 2)
	if err != nil {
		t.Fatalf("get match 2: %v", err)
	}
	playedAt := time.Date(2026, time.May, 24, 14, 0, 0, 0, time.UTC)
	if err := matchRepo.UpdateResult(context.Background(), 1, 2, 1, playedAt); err != nil {
		t.Fatalf("update match 1: %v", err)
	}
	if err := matchRepo.UpdateResult(context.Background(), 2, 1, 1, playedAt); err != nil {
		t.Fatalf("update match 2: %v", err)
	}
	if err := matchRepo.UpdateResult(context.Background(), 3, 3, 0, playedAt); err != nil {
		t.Fatalf("update match 3: %v", err)
	}

	svc := NewStandingsService(
		sqlite.NewTeamRepository(db),
		matchRepo,
		sqlite.NewLeagueRepository(db),
		sqlite.NewStandingRepository(db),
	)

	standings, err := svc.Recalculate(context.Background())
	if err != nil {
		t.Fatalf("recalculate standings: %v", err)
	}

	if len(standings) != domain.TotalTeams {
		t.Fatalf("expected %d standings, got %d", domain.TotalTeams, len(standings))
	}

	if standings[0].TeamID != matchOne.HomeTeamID || standings[0].Points != 6 || standings[0].GoalDifference != 4 {
		t.Fatalf("unexpected first standing: %+v", standings[0])
	}

	if standings[1].TeamID != matchTwo.AwayTeamID || standings[1].Points != 1 || standings[1].Draws != 1 || standings[1].GoalDifference != 0 {
		t.Fatalf("unexpected second standing: %+v", standings[1])
	}

	if standings[2].TeamID != matchTwo.HomeTeamID || standings[2].Points != 1 || standings[2].GoalsFor != 1 || standings[2].GoalsAgainst != 4 || standings[2].GoalDifference != -3 {
		t.Fatalf("unexpected third standing: %+v", standings[2])
	}

	if standings[3].TeamID != matchOne.AwayTeamID || standings[3].Points != 0 || standings[3].Losses != 1 {
		t.Fatalf("unexpected fourth standing: %+v", standings[3])
	}

	for i, standing := range standings {
		if standing.Rank != i+1 {
			t.Fatalf("expected rank %d, got %d", i+1, standing.Rank)
		}
	}
}

func TestStandingsServiceListBlankStandingsAfterInit(t *testing.T) {
	db := newServiceTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	svc := NewStandingsService(
		sqlite.NewTeamRepository(db),
		sqlite.NewMatchRepository(db),
		sqlite.NewLeagueRepository(db),
		sqlite.NewStandingRepository(db),
	)

	standings, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("list standings: %v", err)
	}

	if len(standings) != domain.TotalTeams {
		t.Fatalf("expected %d standings, got %d", domain.TotalTeams, len(standings))
	}

	for i, standing := range standings {
		if standing.Points != 0 || standing.Played != 0 || standing.Rank != i+1 {
			t.Fatalf("unexpected blank standing: %+v", standing)
		}
	}
}

func TestStandingsServiceTotalsRemainConsistent(t *testing.T) {
	db := newServiceTestDB(t)
	leagueSvc := newTestLeagueService(db)
	if _, err := leagueSvc.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize league: %v", err)
	}

	matchRepo := sqlite.NewMatchRepository(db)
	playedAt := time.Date(2026, time.May, 24, 15, 0, 0, 0, time.UTC)
	updates := []struct {
		matchID   int64
		homeGoals int
		awayGoals int
	}{
		{1, 2, 1},
		{2, 0, 0},
		{3, 3, 2},
		{4, 1, 4},
	}
	for _, update := range updates {
		if err := matchRepo.UpdateResult(context.Background(), update.matchID, update.homeGoals, update.awayGoals, playedAt); err != nil {
			t.Fatalf("update match %d: %v", update.matchID, err)
		}
	}

	svc := NewStandingsService(
		sqlite.NewTeamRepository(db),
		matchRepo,
		sqlite.NewLeagueRepository(db),
		sqlite.NewStandingRepository(db),
	)

	standings, err := svc.Recalculate(context.Background())
	if err != nil {
		t.Fatalf("recalculate standings: %v", err)
	}

	var totalWins, totalDraws, totalLosses int
	var totalGoalsFor, totalGoalsAgainst int
	for _, standing := range standings {
		totalWins += standing.Wins
		totalDraws += standing.Draws
		totalLosses += standing.Losses
		totalGoalsFor += standing.GoalsFor
		totalGoalsAgainst += standing.GoalsAgainst
	}

	if totalWins != totalLosses {
		t.Fatalf("expected total wins (%d) to equal total losses (%d)", totalWins, totalLosses)
	}

	if totalGoalsFor != totalGoalsAgainst {
		t.Fatalf("expected total goals for (%d) to equal total goals against (%d)", totalGoalsFor, totalGoalsAgainst)
	}

	if totalDraws%2 != 0 {
		t.Fatalf("expected total draws to be even, got %d", totalDraws)
	}
}
