package domain

import "time"

type Team struct {
	ID        int64
	Name      string
	Strength  int
	CreatedAt time.Time
}
