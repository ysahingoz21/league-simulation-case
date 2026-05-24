package domain

import "time"

type Prediction struct {
	ID                      int64
	Week                    int
	TeamID                  int64
	TeamName                string
	ChampionshipProbability float64
	ExpectedPoints          float64
	ProjectedRank           float64
	CreatedAt               time.Time
}
