package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

type StandingRepository struct {
	db *sql.DB
}

func NewStandingRepository(db *sql.DB) repository.StandingRepository {
	return &StandingRepository{db: db}
}

func (r *StandingRepository) ReplaceAll(ctx context.Context, standings []domain.Standing) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin standing transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM standings`); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete standings: %w", err)
	}

	if len(standings) == 0 {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit empty standings transaction: %w", err)
		}
		return nil
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO standings (
			team_id, played, wins, draws, losses, goals_for, goals_against, goal_difference, points, rank, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare standing insert: %w", err)
	}
	defer stmt.Close()

	for _, standing := range standings {
		if _, err := stmt.ExecContext(
			ctx,
			standing.TeamID,
			standing.Played,
			standing.Wins,
			standing.Draws,
			standing.Losses,
			standing.GoalsFor,
			standing.GoalsAgainst,
			standing.GoalDifference,
			standing.Points,
			standing.Rank,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert standing for team %d: %w", standing.TeamID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit standing transaction: %w", err)
	}

	return nil
}

func (r *StandingRepository) List(ctx context.Context) ([]domain.Standing, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			s.team_id,
			t.name,
			s.played,
			s.wins,
			s.draws,
			s.losses,
			s.goals_for,
			s.goals_against,
			s.goal_difference,
			s.points,
			s.rank
		FROM standings AS s
		JOIN teams AS t ON t.id = s.team_id
		ORDER BY s.rank, t.name
	`)
	if err != nil {
		return nil, fmt.Errorf("list standings: %w", err)
	}
	defer rows.Close()

	var standings []domain.Standing
	for rows.Next() {
		var standing domain.Standing
		if err := rows.Scan(
			&standing.TeamID,
			&standing.TeamName,
			&standing.Played,
			&standing.Wins,
			&standing.Draws,
			&standing.Losses,
			&standing.GoalsFor,
			&standing.GoalsAgainst,
			&standing.GoalDifference,
			&standing.Points,
			&standing.Rank,
		); err != nil {
			return nil, fmt.Errorf("scan standing: %w", err)
		}
		standings = append(standings, standing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate standings: %w", err)
	}

	return standings, nil
}

func (r *StandingRepository) DeleteAll(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM standings`); err != nil {
		return fmt.Errorf("delete standings: %w", err)
	}

	return nil
}
