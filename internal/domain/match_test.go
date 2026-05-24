package domain

import "testing"

func TestMatchWinnerTeamID(t *testing.T) {
	homeWinGoals := 2
	homeWinAgainst := 1
	awayWinGoals := 1
	awayWinAgainst := 3
	drawGoalsHome := 1
	drawGoalsAway := 1

	tests := []struct {
		name     string
		match    Match
		wantNil  bool
		wantTeam int64
	}{
		{
			name: "home win",
			match: Match{
				HomeTeamID: 1,
				AwayTeamID: 2,
				HomeGoals:  &homeWinGoals,
				AwayGoals:  &homeWinAgainst,
				Status:     MatchStatusPlayed,
			},
			wantTeam: 1,
		},
		{
			name: "away win",
			match: Match{
				HomeTeamID: 1,
				AwayTeamID: 2,
				HomeGoals:  &awayWinGoals,
				AwayGoals:  &awayWinAgainst,
				Status:     MatchStatusPlayed,
			},
			wantTeam: 2,
		},
		{
			name: "draw",
			match: Match{
				HomeTeamID: 1,
				AwayTeamID: 2,
				HomeGoals:  &drawGoalsHome,
				AwayGoals:  &drawGoalsAway,
				Status:     MatchStatusPlayed,
			},
			wantNil: true,
		},
		{
			name: "unplayed",
			match: Match{
				HomeTeamID: 1,
				AwayTeamID: 2,
				Status:     MatchStatusScheduled,
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.match.WinnerTeamID()
			if tt.wantNil {
				if got != nil {
					t.Fatalf("expected nil winner, got %d", *got)
				}

				return
			}

			if got == nil {
				t.Fatal("expected winner team ID, got nil")
			}

			if *got != tt.wantTeam {
				t.Fatalf("expected winner %d, got %d", tt.wantTeam, *got)
			}
		})
	}
}
