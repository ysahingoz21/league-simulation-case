package domain

import "time"

const (
	TotalTeams     = 4
	TotalWeeks     = 6
	MatchesPerWeek = 2
	WinPoints      = 3
	DrawPoints     = 1
)

type LeagueState struct {
	CurrentWeek int
	TotalWeeks  int
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
