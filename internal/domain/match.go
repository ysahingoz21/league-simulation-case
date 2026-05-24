package domain

import "time"

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusPlayed    MatchStatus = "played"
)

type Match struct {
	ID         int64
	Week       int
	HomeTeamID int64
	AwayTeamID int64
	HomeGoals  *int
	AwayGoals  *int
	Status     MatchStatus
	PlayedAt   *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m Match) IsPlayed() bool {
	return m.Status == MatchStatusPlayed && m.HomeGoals != nil && m.AwayGoals != nil
}

func (m Match) HasTeam(teamID int64) bool {
	return m.HomeTeamID == teamID || m.AwayTeamID == teamID
}

func (m Match) WinnerTeamID() *int64 {
	if !m.IsPlayed() {
		return nil
	}

	if *m.HomeGoals > *m.AwayGoals {
		winnerID := m.HomeTeamID
		return &winnerID
	}

	if *m.AwayGoals > *m.HomeGoals {
		winnerID := m.AwayTeamID
		return &winnerID
	}

	return nil
}
