package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) repository.MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) CreateMany(ctx context.Context, matches []domain.Match) error {
	if len(matches) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin match transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO matches (
			week, home_team_id, away_team_id, home_goals, away_goals, status, played_at, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, COALESCE(?, CURRENT_TIMESTAMP), COALESCE(?, CURRENT_TIMESTAMP))
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare match insert: %w", err)
	}
	defer stmt.Close()

	for _, match := range matches {
		if _, err := stmt.ExecContext(
			ctx,
			match.Week,
			match.HomeTeamID,
			match.AwayTeamID,
			nullableInt(match.HomeGoals),
			nullableInt(match.AwayGoals),
			matchStatusOrDefault(match.Status),
			nullableTimePtr(match.PlayedAt),
			nullableTime(match.CreatedAt),
			nullableTime(match.UpdatedAt),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert match for week %d: %w", match.Week, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit match transaction: %w", err)
	}

	return nil
}

func (r *MatchRepository) List(ctx context.Context) ([]domain.Match, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, status, played_at, created_at, updated_at
		FROM matches
		ORDER BY week, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()

	return scanMatches(rows)
}

func (r *MatchRepository) ListByWeek(ctx context.Context, week int) ([]domain.Match, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, status, played_at, created_at, updated_at
		FROM matches
		WHERE week = ?
		ORDER BY id
	`, week)
	if err != nil {
		return nil, fmt.Errorf("list matches for week %d: %w", week, err)
	}
	defer rows.Close()

	return scanMatches(rows)
}

func (r *MatchRepository) GetByID(ctx context.Context, id int64) (domain.Match, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, status, played_at, created_at, updated_at
		FROM matches
		WHERE id = ?
	`, id)

	match, err := scanMatch(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Match{}, fmt.Errorf("match %d not found: %w", id, err)
		}

		return domain.Match{}, fmt.Errorf("get match %d: %w", id, err)
	}

	return match, nil
}

func (r *MatchRepository) UpdateResult(ctx context.Context, id int64, homeGoals int, awayGoals int, playedAt time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE matches
		SET
			home_goals = ?,
			away_goals = ?,
			status = ?,
			played_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, homeGoals, awayGoals, domain.MatchStatusPlayed, playedAt, id)
	if err != nil {
		return fmt.Errorf("update match %d result: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read match %d update result: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("match %d not found: %w", id, sql.ErrNoRows)
	}

	return nil
}

func (r *MatchRepository) ListUnplayed(ctx context.Context) ([]domain.Match, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, status, played_at, created_at, updated_at
		FROM matches
		WHERE status = ?
		ORDER BY week, id
	`, domain.MatchStatusScheduled)
	if err != nil {
		return nil, fmt.Errorf("list unplayed matches: %w", err)
	}
	defer rows.Close()

	return scanMatches(rows)
}

func (r *MatchRepository) DeleteAll(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM matches`); err != nil {
		return fmt.Errorf("delete matches: %w", err)
	}

	return nil
}

func scanMatches(rows *sql.Rows) ([]domain.Match, error) {
	var matches []domain.Match
	for rows.Next() {
		match, err := scanMatch(rows)
		if err != nil {
			return nil, fmt.Errorf("scan match: %w", err)
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate matches: %w", err)
	}

	return matches, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanMatch(row scanner) (domain.Match, error) {
	var match domain.Match
	var homeGoals sql.NullInt64
	var awayGoals sql.NullInt64
	var playedAt sql.NullTime
	var status string

	err := row.Scan(
		&match.ID,
		&match.Week,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&homeGoals,
		&awayGoals,
		&status,
		&playedAt,
		&match.CreatedAt,
		&match.UpdatedAt,
	)
	if err != nil {
		return domain.Match{}, err
	}

	match.Status = domain.MatchStatus(status)
	if homeGoals.Valid {
		value := int(homeGoals.Int64)
		match.HomeGoals = &value
	}
	if awayGoals.Valid {
		value := int(awayGoals.Int64)
		match.AwayGoals = &value
	}
	if playedAt.Valid {
		value := playedAt.Time
		match.PlayedAt = &value
	}

	return match, nil
}
