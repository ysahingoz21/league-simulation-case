package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

var (
	ErrLeagueCompleted    = errors.New("league is already completed")
	ErrWeekAlreadyPlayed  = errors.New("week has already been played")
	ErrWeekOutOfOrder     = errors.New("week must be played sequentially")
	ErrWeekHasPlayedMatch = errors.New("week already contains played matches")
	ErrWeekHasNoMatches   = errors.New("week has no matches")
)

type WeekSimulationResult struct {
	Week      int
	League    domain.LeagueState
	Matches   []domain.Match
	Standings []domain.Standing
}

type SimulationService interface {
	PlayNextWeek(ctx context.Context) (WeekSimulationResult, error)
	PlayWeek(ctx context.Context, week int) (WeekSimulationResult, error)
}

type simulationService struct {
	teams     repository.TeamRepository
	matches   repository.MatchRepository
	league    repository.LeagueRepository
	standings StandingsService
	simulator Simulator
	now       func() time.Time
}

func NewSimulationService(
	teams repository.TeamRepository,
	matches repository.MatchRepository,
	league repository.LeagueRepository,
	standings StandingsService,
	simulator Simulator,
) SimulationService {
	return &simulationService{
		teams:     teams,
		matches:   matches,
		league:    league,
		standings: standings,
		simulator: simulator,
		now:       time.Now,
	}
}

func (s *simulationService) PlayNextWeek(ctx context.Context) (WeekSimulationResult, error) {
	state, err := s.getLeagueState(ctx)
	if err != nil {
		return WeekSimulationResult{}, err
	}

	if state.IsCompleted {
		return WeekSimulationResult{}, ErrLeagueCompleted
	}

	nextWeek := state.CurrentWeek + 1
	if nextWeek < 1 || nextWeek > state.TotalWeeks {
		return WeekSimulationResult{}, ErrInvalidFixtureWeek
	}

	return s.playWeek(ctx, nextWeek, state)
}

func (s *simulationService) PlayWeek(ctx context.Context, week int) (WeekSimulationResult, error) {
	if week < 1 || week > domain.TotalWeeks {
		return WeekSimulationResult{}, ErrInvalidFixtureWeek
	}

	state, err := s.getLeagueState(ctx)
	if err != nil {
		return WeekSimulationResult{}, err
	}

	if state.IsCompleted {
		return WeekSimulationResult{}, ErrLeagueCompleted
	}

	if week <= state.CurrentWeek {
		return WeekSimulationResult{}, ErrWeekAlreadyPlayed
	}

	if week > state.CurrentWeek+1 {
		return WeekSimulationResult{}, ErrWeekOutOfOrder
	}

	return s.playWeek(ctx, week, state)
}

func (s *simulationService) playWeek(ctx context.Context, week int, state domain.LeagueState) (WeekSimulationResult, error) {
	matches, err := s.matches.ListByWeek(ctx, week)
	if err != nil {
		return WeekSimulationResult{}, fmt.Errorf("list matches for week %d: %w", week, err)
	}

	if len(matches) == 0 {
		return WeekSimulationResult{}, ErrWeekHasNoMatches
	}

	for _, match := range matches {
		if match.IsPlayed() {
			return WeekSimulationResult{}, ErrWeekHasPlayedMatch
		}
	}

	teams, err := s.teams.List(ctx)
	if err != nil {
		return WeekSimulationResult{}, fmt.Errorf("list teams for simulation: %w", err)
	}

	teamMap := make(map[int64]domain.Team, len(teams))
	for _, team := range teams {
		teamMap[team.ID] = team
	}

	playedAt := s.now().UTC()
	for _, match := range matches {
		homeTeam, ok := teamMap[match.HomeTeamID]
		if !ok {
			return WeekSimulationResult{}, fmt.Errorf("home team %d not found", match.HomeTeamID)
		}

		awayTeam, ok := teamMap[match.AwayTeamID]
		if !ok {
			return WeekSimulationResult{}, fmt.Errorf("away team %d not found", match.AwayTeamID)
		}

		homeGoals, awayGoals := s.simulator.SimulateMatch(homeTeam, awayTeam)
		if err := s.matches.UpdateResult(ctx, match.ID, homeGoals, awayGoals, playedAt); err != nil {
			return WeekSimulationResult{}, fmt.Errorf("update result for match %d: %w", match.ID, err)
		}
	}

	isCompleted := week == state.TotalWeeks
	if err := s.league.UpdateCurrentWeek(ctx, week, isCompleted); err != nil {
		return WeekSimulationResult{}, fmt.Errorf("update league current week: %w", err)
	}

	updatedStandings, err := s.standings.Recalculate(ctx)
	if err != nil {
		return WeekSimulationResult{}, fmt.Errorf("recalculate standings: %w", err)
	}

	updatedState, err := s.getLeagueState(ctx)
	if err != nil {
		return WeekSimulationResult{}, err
	}

	playedMatches, err := s.matches.ListByWeek(ctx, week)
	if err != nil {
		return WeekSimulationResult{}, fmt.Errorf("list updated matches for week %d: %w", week, err)
	}

	return WeekSimulationResult{
		Week:      week,
		League:    updatedState,
		Matches:   playedMatches,
		Standings: updatedStandings,
	}, nil
}

func (s *simulationService) getLeagueState(ctx context.Context) (domain.LeagueState, error) {
	return getLeagueState(ctx, s.league)
}

func getLeagueState(ctx context.Context, league repository.LeagueRepository) (domain.LeagueState, error) {
	state, err := league.GetState(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LeagueState{}, ErrLeagueNotInitialized
		}

		return domain.LeagueState{}, fmt.Errorf("get league state: %w", err)
	}

	return state, nil
}
