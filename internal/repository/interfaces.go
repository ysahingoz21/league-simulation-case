package repository

import (
	"context"
	"time"

	"league-simulation-case/internal/domain"
)

type TeamRepository interface {
	CreateMany(ctx context.Context, teams []domain.Team) error
	List(ctx context.Context) ([]domain.Team, error)
	GetByID(ctx context.Context, id int64) (domain.Team, error)
	DeleteAll(ctx context.Context) error
}

type MatchRepository interface {
	CreateMany(ctx context.Context, matches []domain.Match) error
	List(ctx context.Context) ([]domain.Match, error)
	ListByWeek(ctx context.Context, week int) ([]domain.Match, error)
	GetByID(ctx context.Context, id int64) (domain.Match, error)
	UpdateResult(ctx context.Context, id int64, homeGoals int, awayGoals int, playedAt time.Time) error
	ListUnplayed(ctx context.Context) ([]domain.Match, error)
	DeleteAll(ctx context.Context) error
}

type LeagueRepository interface {
	GetState(ctx context.Context) (domain.LeagueState, error)
	UpsertState(ctx context.Context, state domain.LeagueState) error
	UpdateCurrentWeek(ctx context.Context, currentWeek int, isCompleted bool) error
	Reset(ctx context.Context) error
	DeleteAll(ctx context.Context) error
}

type StandingRepository interface {
	ReplaceAll(ctx context.Context, standings []domain.Standing) error
	List(ctx context.Context) ([]domain.Standing, error)
	DeleteAll(ctx context.Context) error
}

type PredictionRepository interface {
	ReplaceForWeek(ctx context.Context, week int, predictions []domain.Prediction) error
	ListLatest(ctx context.Context) ([]domain.Prediction, error)
	DeleteAll(ctx context.Context) error
}
