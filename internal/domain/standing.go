package domain

import "sort"

type Standing struct {
	TeamID         int64
	TeamName       string
	Played         int
	Wins           int
	Draws          int
	Losses         int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	Points         int
	Rank           int
}

func (s Standing) ApplyMatch(match Match, team Team) Standing {
	if !match.IsPlayed() || !match.HasTeam(team.ID) {
		return s
	}

	s.TeamID = team.ID
	s.TeamName = team.Name
	s.Played++

	var goalsFor int
	var goalsAgainst int

	if match.HomeTeamID == team.ID {
		goalsFor = *match.HomeGoals
		goalsAgainst = *match.AwayGoals
	} else {
		goalsFor = *match.AwayGoals
		goalsAgainst = *match.HomeGoals
	}

	s.GoalsFor += goalsFor
	s.GoalsAgainst += goalsAgainst
	s.GoalDifference = s.GoalsFor - s.GoalsAgainst

	switch {
	case goalsFor > goalsAgainst:
		s.Wins++
		s.Points += WinPoints
	case goalsFor < goalsAgainst:
		s.Losses++
	default:
		s.Draws++
		s.Points += DrawPoints
	}

	return s
}

func SortStandings(standings []Standing) {
	sort.Slice(standings, func(i, j int) bool {
		left := standings[i]
		right := standings[j]

		switch {
		case left.Points != right.Points:
			return left.Points > right.Points
		case left.GoalDifference != right.GoalDifference:
			return left.GoalDifference > right.GoalDifference
		case left.GoalsFor != right.GoalsFor:
			return left.GoalsFor > right.GoalsFor
		default:
			return left.TeamName < right.TeamName
		}
	})
}

func AssignRanks(standings []Standing) {
	for i := range standings {
		standings[i].Rank = i + 1
	}
}
