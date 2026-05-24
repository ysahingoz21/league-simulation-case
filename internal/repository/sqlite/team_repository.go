package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) repository.TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateMany(ctx context.Context, teams []domain.Team) error {
	if len(teams) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin team transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO teams (name, strength, created_at)
		VALUES (?, ?, COALESCE(?, CURRENT_TIMESTAMP))
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare team insert: %w", err)
	}
	defer stmt.Close()

	for _, team := range teams {
		if _, err := stmt.ExecContext(ctx, team.Name, team.Strength, nullableTime(team.CreatedAt)); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert team %s: %w", team.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit team transaction: %w", err)
	}

	return nil
}

func (r *TeamRepository) List(ctx context.Context) ([]domain.Team, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, strength, created_at
		FROM teams
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.Strength, &team.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams: %w", err)
	}

	return teams, nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id int64) (domain.Team, error) {
	var team domain.Team

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, strength, created_at
		FROM teams
		WHERE id = ?
	`, id).Scan(&team.ID, &team.Name, &team.Strength, &team.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Team{}, fmt.Errorf("team %d not found: %w", id, err)
		}

		return domain.Team{}, fmt.Errorf("get team %d: %w", id, err)
	}

	return team, nil
}

func (r *TeamRepository) DeleteAll(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM teams`); err != nil {
		return fmt.Errorf("delete teams: %w", err)
	}

	return nil
}
