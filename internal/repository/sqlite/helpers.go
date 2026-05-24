package sqlite

import (
	"database/sql"
	"time"

	"league-simulation-case/internal/domain"
)

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}

	return value
}

func nullableTimePtr(value *time.Time) any {
	if value == nil || value.IsZero() {
		return nil
	}

	return *value
}

func nullableInt(value *int) any {
	if value == nil {
		return nil
	}

	return *value
}

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

func matchStatusOrDefault(status domain.MatchStatus) domain.MatchStatus {
	if status == "" {
		return domain.MatchStatusScheduled
	}

	return status
}

func isNoRows(err error) bool {
	return err != nil && err == sql.ErrNoRows
}
