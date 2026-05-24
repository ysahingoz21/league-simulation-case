package service

import (
	"context"
	"fmt"
	"sort"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/repository"
)

const (
	defaultPredictionIterations = 1000
	predictionStartWeek         = 4
)

type PredictionQueryResult struct {
	Week        int
	Message     string
	Predictions []domain.Prediction
}

type PredictionService interface {
	GetPredictions(ctx context.Context) (PredictionQueryResult, error)
	GenerateForCurrentWeek(ctx context.Context) ([]domain.Prediction, error)
}

type predictionService struct {
	teams       repository.TeamRepository
	matches     repository.MatchRepository
	league      repository.LeagueRepository
	predictions repository.PredictionRepository
	simulator   Simulator
	iterations  int
}

func NewPredictionService(
	teams repository.TeamRepository,
	matches repository.MatchRepository,
	league repository.LeagueRepository,
	predictions repository.PredictionRepository,
	simulator Simulator,
	iterations int,
) PredictionService {
	if iterations <= 0 {
		iterations = defaultPredictionIterations
	}

	return &predictionService{
		teams:       teams,
		matches:     matches,
		league:      league,
		predictions: predictions,
		simulator:   simulator,
		iterations:  iterations,
	}
}

func (s *predictionService) GetPredictions(ctx context.Context) (PredictionQueryResult, error) {
	state, err := getLeagueState(ctx, s.league)
	if err != nil {
		return PredictionQueryResult{}, err
	}

	if state.CurrentWeek < predictionStartWeek {
		return PredictionQueryResult{
			Week:        state.CurrentWeek,
			Message:     "Predictions are available after week 4",
			Predictions: []domain.Prediction{},
		}, nil
	}

	latest, err := s.predictions.ListLatest(ctx)
	if err != nil {
		return PredictionQueryResult{}, fmt.Errorf("list latest predictions: %w", err)
	}

	if len(latest) == 0 || latest[0].Week != state.CurrentWeek {
		latest, err = s.GenerateForCurrentWeek(ctx)
		if err != nil {
			return PredictionQueryResult{}, err
		}
	}

	return PredictionQueryResult{
		Week:        state.CurrentWeek,
		Predictions: latest,
	}, nil
}

func (s *predictionService) GenerateForCurrentWeek(ctx context.Context) ([]domain.Prediction, error) {
	state, err := getLeagueState(ctx, s.league)
	if err != nil {
		return nil, err
	}

	if state.CurrentWeek < predictionStartWeek {
		return []domain.Prediction{}, nil
	}

	teams, err := s.teams.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams for predictions: %w", err)
	}

	matches, err := s.matches.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list matches for predictions: %w", err)
	}

	currentStandings := calculateStandings(teams, matches)
	remainingMatches := make([]domain.Match, 0)
	for _, match := range matches {
		if !match.IsPlayed() {
			remainingMatches = append(remainingMatches, match)
		}
	}

	var generated []domain.Prediction
	if len(remainingMatches) == 0 {
		generated = completedLeaguePredictions(state.CurrentWeek, currentStandings)
	} else {
		generated = s.runMonteCarlo(state.CurrentWeek, teams, currentStandings, remainingMatches)
	}

	if err := s.predictions.ReplaceForWeek(ctx, state.CurrentWeek, generated); err != nil {
		return nil, fmt.Errorf("persist predictions for week %d: %w", state.CurrentWeek, err)
	}

	latest, err := s.predictions.ListLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("list persisted predictions: %w", err)
	}

	return latest, nil
}

func (s *predictionService) runMonteCarlo(week int, teams []domain.Team, currentStandings []domain.Standing, remainingMatches []domain.Match) []domain.Prediction {
	type aggregate struct {
		titleCount  int
		totalPoints float64
		totalRanks  float64
	}

	teamMap := make(map[int64]domain.Team, len(teams))
	for _, team := range teams {
		teamMap[team.ID] = team
	}

	aggregates := make(map[int64]*aggregate, len(teams))
	for _, team := range teams {
		aggregates[team.ID] = &aggregate{}
	}

	for simulation := 0; simulation < s.iterations; simulation++ {
		projected := cloneStandings(currentStandings)
		projectedMap := make(map[int64]domain.Standing, len(projected))
		for _, standing := range projected {
			projectedMap[standing.TeamID] = standing
		}

		for _, match := range remainingMatches {
			homeTeam := teamMap[match.HomeTeamID]
			awayTeam := teamMap[match.AwayTeamID]
			homeGoals, awayGoals := s.simulator.SimulateMatch(homeTeam, awayTeam)
			simulatedMatch := match
			simulatedMatch.Status = domain.MatchStatusPlayed
			simulatedMatch.HomeGoals = intPtr(homeGoals)
			simulatedMatch.AwayGoals = intPtr(awayGoals)

			homeStanding := projectedMap[homeTeam.ID]
			projectedMap[homeTeam.ID] = homeStanding.ApplyMatch(simulatedMatch, homeTeam)

			awayStanding := projectedMap[awayTeam.ID]
			projectedMap[awayTeam.ID] = awayStanding.ApplyMatch(simulatedMatch, awayTeam)
		}

		finalStandings := make([]domain.Standing, 0, len(teams))
		for _, team := range teams {
			finalStandings = append(finalStandings, projectedMap[team.ID])
		}

		domain.SortStandings(finalStandings)
		domain.AssignRanks(finalStandings)

		for _, standing := range finalStandings {
			aggregate := aggregates[standing.TeamID]
			if standing.Rank == 1 {
				aggregate.titleCount++
			}
			aggregate.totalPoints += float64(standing.Points)
			aggregate.totalRanks += float64(standing.Rank)
		}
	}

	predictions := make([]domain.Prediction, 0, len(teams))
	for _, team := range teams {
		aggregate := aggregates[team.ID]
		predictions = append(predictions, domain.Prediction{
			Week:                    week,
			TeamID:                  team.ID,
			TeamName:                team.Name,
			ChampionshipProbability: (float64(aggregate.titleCount) / float64(s.iterations)) * 100,
			ExpectedPoints:          aggregate.totalPoints / float64(s.iterations),
			ProjectedRank:           aggregate.totalRanks / float64(s.iterations),
		})
	}

	sortPredictions(predictions)
	return predictions
}

func completedLeaguePredictions(week int, standings []domain.Standing) []domain.Prediction {
	predictions := make([]domain.Prediction, 0, len(standings))
	for _, standing := range standings {
		probability := 0.0
		if standing.Rank == 1 {
			probability = 100
		}

		predictions = append(predictions, domain.Prediction{
			Week:                    week,
			TeamID:                  standing.TeamID,
			TeamName:                standing.TeamName,
			ChampionshipProbability: probability,
			ExpectedPoints:          float64(standing.Points),
			ProjectedRank:           float64(standing.Rank),
		})
	}

	sortPredictions(predictions)
	return predictions
}

func sortPredictions(predictions []domain.Prediction) {
	sort.Slice(predictions, func(i, j int) bool {
		left := predictions[i]
		right := predictions[j]

		switch {
		case left.ChampionshipProbability != right.ChampionshipProbability:
			return left.ChampionshipProbability > right.ChampionshipProbability
		case left.ProjectedRank != right.ProjectedRank:
			return left.ProjectedRank < right.ProjectedRank
		default:
			return left.TeamName < right.TeamName
		}
	})
}

func cloneStandings(standings []domain.Standing) []domain.Standing {
	cloned := make([]domain.Standing, len(standings))
	copy(cloned, standings)
	return cloned
}

func intPtr(value int) *int {
	v := value
	return &v
}
