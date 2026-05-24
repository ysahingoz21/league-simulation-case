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
	ErrMatchNotFound       = errors.New("match not found")
	ErrInvalidGoals        = errors.New("invalid goals")
	ErrMatchEditNotAllowed = errors.New("match edit not allowed")
)

type MatchUpdateResult struct {
	Match       domain.Match
	League      domain.LeagueState
	Standings   []domain.Standing
	Predictions []domain.Prediction
}

type MatchService interface {
	GetMatch(ctx context.Context, matchID int64) (domain.Match, error)
	UpdateMatchResult(ctx context.Context, matchID int64, homeGoals int, awayGoals int) (MatchUpdateResult, error)
}

type matchService struct {
	matches     repository.MatchRepository
	league      repository.LeagueRepository
	standings   StandingsService
	predictions PredictionService
	now         func() time.Time
}

func NewMatchService(
	matches repository.MatchRepository,
	league repository.LeagueRepository,
	standings StandingsService,
	predictions PredictionService,
) MatchService {
	return &matchService{
		matches:     matches,
		league:      league,
		standings:   standings,
		predictions: predictions,
		now:         time.Now,
	}
}

func (s *matchService) GetMatch(ctx context.Context, matchID int64) (domain.Match, error) {
	if matchID <= 0 {
		return domain.Match{}, ErrMatchNotFound
	}

	if _, err := getLeagueState(ctx, s.league); err != nil {
		return domain.Match{}, err
	}

	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Match{}, ErrMatchNotFound
		}

		return domain.Match{}, fmt.Errorf("get match: %w", err)
	}

	return match, nil
}

func (s *matchService) UpdateMatchResult(ctx context.Context, matchID int64, homeGoals int, awayGoals int) (MatchUpdateResult, error) {
	if matchID <= 0 {
		return MatchUpdateResult{}, ErrMatchNotFound
	}

	if homeGoals < 0 || awayGoals < 0 || homeGoals > 20 || awayGoals > 20 {
		return MatchUpdateResult{}, ErrInvalidGoals
	}

	state, err := getLeagueState(ctx, s.league)
	if err != nil {
		return MatchUpdateResult{}, err
	}

	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MatchUpdateResult{}, ErrMatchNotFound
		}

		return MatchUpdateResult{}, fmt.Errorf("get match for update: %w", err)
	}

	if !match.IsPlayed() && match.Week > state.CurrentWeek {
		return MatchUpdateResult{}, ErrMatchEditNotAllowed
	}

	playedAt := s.now().UTC()
	if match.PlayedAt != nil {
		playedAt = *match.PlayedAt
	}

	if err := s.matches.UpdateResult(ctx, matchID, homeGoals, awayGoals, playedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MatchUpdateResult{}, ErrMatchNotFound
		}

		return MatchUpdateResult{}, fmt.Errorf("update match result: %w", err)
	}

	updatedMatch, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return MatchUpdateResult{}, fmt.Errorf("get updated match: %w", err)
	}

	standings, err := s.standings.Recalculate(ctx)
	if err != nil {
		return MatchUpdateResult{}, fmt.Errorf("recalculate standings after match update: %w", err)
	}

	predictions := []domain.Prediction{}
	if state.CurrentWeek >= predictionStartWeek {
		predictions, err = s.predictions.GenerateForCurrentWeek(ctx)
		if err != nil {
			return MatchUpdateResult{}, fmt.Errorf("refresh predictions after match update: %w", err)
		}
	}

	updatedLeague, err := getLeagueState(ctx, s.league)
	if err != nil {
		return MatchUpdateResult{}, err
	}

	return MatchUpdateResult{
		Match:       updatedMatch,
		League:      updatedLeague,
		Standings:   standings,
		Predictions: predictions,
	}, nil
}
