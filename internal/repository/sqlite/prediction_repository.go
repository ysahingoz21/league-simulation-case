package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

type PredictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) repository.PredictionRepository {
	return &PredictionRepository{db: db}
}

func (r *PredictionRepository) ReplaceForWeek(ctx context.Context, week int, predictions []domain.Prediction) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin prediction transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM predictions WHERE week = ?`, week); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete predictions for week %d: %w", week, err)
	}

	if len(predictions) == 0 {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit empty prediction transaction: %w", err)
		}
		return nil
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO predictions (
			week, team_id, championship_probability, expected_points, projected_rank, created_at
		)
		VALUES (?, ?, ?, ?, ?, COALESCE(?, CURRENT_TIMESTAMP))
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare prediction insert: %w", err)
	}
	defer stmt.Close()

	for _, prediction := range predictions {
		insertWeek := week
		if prediction.Week != 0 {
			insertWeek = prediction.Week
		}

		if _, err := stmt.ExecContext(
			ctx,
			insertWeek,
			prediction.TeamID,
			prediction.ChampionshipProbability,
			prediction.ExpectedPoints,
			prediction.ProjectedRank,
			nullableTime(prediction.CreatedAt),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert prediction for team %d: %w", prediction.TeamID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit prediction transaction: %w", err)
	}

	return nil
}

func (r *PredictionRepository) ListLatest(ctx context.Context) ([]domain.Prediction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			p.id,
			p.week,
			p.team_id,
			t.name,
			p.championship_probability,
			p.expected_points,
			p.projected_rank,
			p.created_at
		FROM predictions AS p
		JOIN teams AS t ON t.id = p.team_id
		WHERE p.week = (SELECT MAX(week) FROM predictions)
		ORDER BY p.projected_rank, t.name
	`)
	if err != nil {
		return nil, fmt.Errorf("list latest predictions: %w", err)
	}
	defer rows.Close()

	var predictions []domain.Prediction
	for rows.Next() {
		var prediction domain.Prediction
		if err := rows.Scan(
			&prediction.ID,
			&prediction.Week,
			&prediction.TeamID,
			&prediction.TeamName,
			&prediction.ChampionshipProbability,
			&prediction.ExpectedPoints,
			&prediction.ProjectedRank,
			&prediction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan prediction: %w", err)
		}
		predictions = append(predictions, prediction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate predictions: %w", err)
	}

	return predictions, nil
}

func (r *PredictionRepository) DeleteAll(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM predictions`); err != nil {
		return fmt.Errorf("delete predictions: %w", err)
	}

	return nil
}
