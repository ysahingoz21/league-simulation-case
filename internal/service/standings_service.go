package service

import (
	"context"
	"fmt"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

type StandingsService interface {
	Recalculate(ctx context.Context) ([]domain.Standing, error)
	List(ctx context.Context) ([]domain.Standing, error)
}

type standingsService struct {
	teams     repository.TeamRepository
	matches   repository.MatchRepository
	league    repository.LeagueRepository
	standings repository.StandingRepository
}

func NewStandingsService(
	teams repository.TeamRepository,
	matches repository.MatchRepository,
	league repository.LeagueRepository,
	standings repository.StandingRepository,
) StandingsService {
	return &standingsService{
		teams:     teams,
		matches:   matches,
		league:    league,
		standings: standings,
	}
}

func (s *standingsService) Recalculate(ctx context.Context) ([]domain.Standing, error) {
	if _, err := getLeagueState(ctx, s.league); err != nil {
		return nil, err
	}

	teams, err := s.teams.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams for standings: %w", err)
	}

	matches, err := s.matches.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list matches for standings: %w", err)
	}

	standingsByTeam := make(map[int64]domain.Standing, len(teams))
	for _, team := range teams {
		standingsByTeam[team.ID] = domain.Standing{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	for _, match := range matches {
		if !match.IsPlayed() {
			continue
		}

		for _, team := range teams {
			if !match.HasTeam(team.ID) {
				continue
			}

			standing := standingsByTeam[team.ID]
			standingsByTeam[team.ID] = standing.ApplyMatch(match, team)
		}
	}

	standings := make([]domain.Standing, 0, len(teams))
	for _, team := range teams {
		standings = append(standings, standingsByTeam[team.ID])
	}

	domain.SortStandings(standings)
	domain.AssignRanks(standings)

	if err := s.standings.ReplaceAll(ctx, standings); err != nil {
		return nil, fmt.Errorf("replace standings: %w", err)
	}

	return standings, nil
}

func (s *standingsService) List(ctx context.Context) ([]domain.Standing, error) {
	if _, err := getLeagueState(ctx, s.league); err != nil {
		return nil, err
	}

	standings, err := s.standings.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list standings: %w", err)
	}

	return standings, nil
}
