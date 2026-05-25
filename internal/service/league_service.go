package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

var (
	ErrInvalidFixtureWeek   = errors.New("invalid fixture week")
	ErrLeagueNotInitialized = errors.New("league is not initialized")
)

type LeagueBootstrap struct {
	League    domain.LeagueState
	Teams     []domain.Team
	Fixtures  []domain.Match
	Standings []domain.Standing
}

type LeagueService interface {
	Initialize(ctx context.Context) (LeagueBootstrap, error)
	Reset(ctx context.Context) (LeagueBootstrap, error)
	GetState(ctx context.Context) (domain.LeagueState, error)
	ListTeams(ctx context.Context) ([]domain.Team, error)
	ListFixtures(ctx context.Context) ([]domain.Match, error)
	ListFixturesByWeek(ctx context.Context, week int) ([]domain.Match, error)
}

type leagueService struct {
	teams       repository.TeamRepository
	matches     repository.MatchRepository
	league      repository.LeagueRepository
	standings   repository.StandingRepository
	predictions repository.PredictionRepository
}

func NewLeagueService(
	teams repository.TeamRepository,
	matches repository.MatchRepository,
	league repository.LeagueRepository,
	standings repository.StandingRepository,
	predictions repository.PredictionRepository,
) LeagueService {
	return &leagueService{
		teams:       teams,
		matches:     matches,
		league:      league,
		standings:   standings,
		predictions: predictions,
	}
}

func (s *leagueService) Initialize(ctx context.Context) (LeagueBootstrap, error) {
	if err := s.clearLeagueData(ctx); err != nil {
		return LeagueBootstrap{}, err
	}

	defaultTeams := defaultTeams()
	if err := s.teams.CreateMany(ctx, defaultTeams); err != nil {
		return LeagueBootstrap{}, fmt.Errorf("create default teams: %w", err)
	}

	teams, err := s.teams.List(ctx)
	if err != nil {
		return LeagueBootstrap{}, fmt.Errorf("list teams after create: %w", err)
	}

	standings := buildInitialStandings(teams)
	if err := s.standings.ReplaceAll(ctx, standings); err != nil {
		return LeagueBootstrap{}, fmt.Errorf("create initial standings: %w", err)
	}

	fixtures, err := GenerateDoubleRoundRobinFixtures(teams)
	if err != nil {
		return LeagueBootstrap{}, fmt.Errorf("generate fixtures: %w", err)
	}

	if err := s.matches.CreateMany(ctx, fixtures); err != nil {
		return LeagueBootstrap{}, fmt.Errorf("create fixtures: %w", err)
	}

	state := domain.LeagueState{
		CurrentWeek: 0,
		TotalWeeks:  domain.TotalWeeks,
		IsCompleted: false,
	}
	if err := s.league.UpsertState(ctx, state); err != nil {
		return LeagueBootstrap{}, fmt.Errorf("upsert league state: %w", err)
	}

	createdState, err := s.league.GetState(ctx)
	if err != nil {
		return LeagueBootstrap{}, fmt.Errorf("get created league state: %w", err)
	}

	createdFixtures, err := s.matches.List(ctx)
	if err != nil {
		return LeagueBootstrap{}, fmt.Errorf("list created fixtures: %w", err)
	}

	createdStandings, err := s.standings.List(ctx)
	if err != nil {
		return LeagueBootstrap{}, fmt.Errorf("list created standings: %w", err)
	}

	return LeagueBootstrap{
		League:    createdState,
		Teams:     teams,
		Fixtures:  createdFixtures,
		Standings: createdStandings,
	}, nil
}

func (s *leagueService) Reset(ctx context.Context) (LeagueBootstrap, error) {
	return s.Initialize(ctx)
}

func (s *leagueService) GetState(ctx context.Context) (domain.LeagueState, error) {
	state, err := s.league.GetState(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LeagueState{}, ErrLeagueNotInitialized
		}

		return domain.LeagueState{}, fmt.Errorf("get league state: %w", err)
	}

	return state, nil
}

func (s *leagueService) ListTeams(ctx context.Context) ([]domain.Team, error) {
	teams, err := s.teams.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	return teams, nil
}

func (s *leagueService) ListFixtures(ctx context.Context) ([]domain.Match, error) {
	fixtures, err := s.matches.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list fixtures: %w", err)
	}

	return fixtures, nil
}

func (s *leagueService) ListFixturesByWeek(ctx context.Context, week int) ([]domain.Match, error) {
	if week < 1 || week > domain.TotalWeeks {
		return nil, ErrInvalidFixtureWeek
	}

	fixtures, err := s.matches.ListByWeek(ctx, week)
	if err != nil {
		return nil, fmt.Errorf("list fixtures by week: %w", err)
	}

	return fixtures, nil
}

func (s *leagueService) clearLeagueData(ctx context.Context) error {
	steps := []func(context.Context) error{
		s.predictions.DeleteAll,
		s.standings.DeleteAll,
		s.matches.DeleteAll,
		s.league.DeleteAll,
		s.teams.DeleteAll,
	}

	for _, step := range steps {
		if err := step(ctx); err != nil {
			return fmt.Errorf("clear league data: %w", err)
		}
	}

	return nil
}

func defaultTeams() []domain.Team {
	return []domain.Team{
		{Name: "Turkey", Strength: 95},
		{Name: "USA", Strength: 85},
		{Name: "Australia", Strength: 80},
		{Name: "Paraguay", Strength: 75},
	}
}

func buildInitialStandings(teams []domain.Team) []domain.Standing {
	standings := make([]domain.Standing, 0, len(teams))
	for i, team := range teams {
		standings = append(standings, domain.Standing{
			TeamID:   team.ID,
			TeamName: team.Name,
			Rank:     i + 1,
		})
	}

	return standings
}
