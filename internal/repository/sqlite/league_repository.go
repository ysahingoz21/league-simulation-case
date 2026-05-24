package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

const leagueStateID = 1

type LeagueRepository struct {
	db *sql.DB
}

func NewLeagueRepository(db *sql.DB) repository.LeagueRepository {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) GetState(ctx context.Context) (domain.LeagueState, error) {
	var state domain.LeagueState
	var isCompleted int

	err := r.db.QueryRowContext(ctx, `
		SELECT current_week, total_weeks, is_completed, created_at, updated_at
		FROM league_state
		WHERE id = ?
	`, leagueStateID).Scan(
		&state.CurrentWeek,
		&state.TotalWeeks,
		&isCompleted,
		&state.CreatedAt,
		&state.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LeagueState{}, fmt.Errorf("league state not found: %w", err)
		}

		return domain.LeagueState{}, fmt.Errorf("get league state: %w", err)
	}

	state.IsCompleted = isCompleted == 1
	return state, nil
}

func (r *LeagueRepository) UpsertState(ctx context.Context, state domain.LeagueState) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO league_state (id, current_week, total_weeks, is_completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, COALESCE(?, CURRENT_TIMESTAMP), COALESCE(?, CURRENT_TIMESTAMP))
		ON CONFLICT(id) DO UPDATE SET
			current_week = excluded.current_week,
			total_weeks = excluded.total_weeks,
			is_completed = excluded.is_completed,
			updated_at = COALESCE(excluded.updated_at, CURRENT_TIMESTAMP)
	`, leagueStateID, state.CurrentWeek, state.TotalWeeks, boolToInt(state.IsCompleted), nullableTime(state.CreatedAt), nullableTime(state.UpdatedAt))
	if err != nil {
		return fmt.Errorf("upsert league state: %w", err)
	}

	return nil
}

func (r *LeagueRepository) UpdateCurrentWeek(ctx context.Context, currentWeek int, isCompleted bool) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO league_state (id, current_week, total_weeks, is_completed)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			current_week = excluded.current_week,
			is_completed = excluded.is_completed,
			updated_at = CURRENT_TIMESTAMP
	`, leagueStateID, currentWeek, domain.TotalWeeks, boolToInt(isCompleted))
	if err != nil {
		return fmt.Errorf("update current week: %w", err)
	}

	return nil
}

func (r *LeagueRepository) Reset(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO league_state (id, current_week, total_weeks, is_completed)
		VALUES (?, 0, ?, 0)
		ON CONFLICT(id) DO UPDATE SET
			current_week = 0,
			total_weeks = ?,
			is_completed = 0,
			updated_at = CURRENT_TIMESTAMP
	`, leagueStateID, domain.TotalWeeks, domain.TotalWeeks)
	if err != nil {
		return fmt.Errorf("reset league state: %w", err)
	}

	return nil
}

func (r *LeagueRepository) DeleteAll(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM league_state`); err != nil {
		return fmt.Errorf("delete league state: %w", err)
	}

	return nil
}
