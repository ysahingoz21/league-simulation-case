package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/service"
)

type LeagueHandler struct {
	leagueService    service.LeagueService
	standingsService service.StandingsService
}

func NewLeagueHandler(leagueService service.LeagueService, standingsService service.StandingsService) LeagueHandler {
	return LeagueHandler{
		leagueService:    leagueService,
		standingsService: standingsService,
	}
}

func (h LeagueHandler) InitializeLeague(c *gin.Context) {
	result, err := h.leagueService.Initialize(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bootstrapResponse(result))
}

func (h LeagueHandler) ResetLeague(c *gin.Context) {
	result, err := h.leagueService.Reset(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bootstrapResponse(result))
}

func (h LeagueHandler) GetLeague(c *gin.Context) {
	state, err := h.leagueService.GetState(c.Request.Context())
	if err != nil {
		if errors.Is(err, service.ErrLeagueNotInitialized) {
			c.JSON(http.StatusNotFound, gin.H{"error": "league is not initialized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"league": leagueStateResponse(state)})
}

func (h LeagueHandler) ListTeams(c *gin.Context) {
	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teamResponses(teams)})
}

func (h LeagueHandler) ListStandings(c *gin.Context) {
	standings, err := h.standingsService.List(c.Request.Context())
	if err != nil {
		if errors.Is(err, service.ErrLeagueNotInitialized) {
			c.JSON(http.StatusNotFound, gin.H{"error": "league is not initialized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"standings": standingResponses(standings)})
}

func (h LeagueHandler) ListFixtures(c *gin.Context) {
	teams, fixtures, ok := h.loadTeamsAndFixtures(c, 0)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"fixtures": fixtureResponses(fixtures, teams)})
}

func (h LeagueHandler) ListFixturesByWeek(c *gin.Context) {
	week, err := strconv.Atoi(c.Param("week"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "week must be an integer"})
		return
	}

	teams, fixtures, ok := h.loadTeamsAndFixtures(c, week)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"week":     week,
		"fixtures": fixtureResponses(fixtures, teams),
	})
}

func (h LeagueHandler) loadTeamsAndFixtures(c *gin.Context, week int) ([]domain.Team, []domain.Match, bool) {
	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, nil, false
	}

	var fixtures []domain.Match
	if week == 0 {
		fixtures, err = h.leagueService.ListFixtures(c.Request.Context())
	} else {
		fixtures, err = h.leagueService.ListFixturesByWeek(c.Request.Context(), week)
	}
	if err != nil {
		if errors.Is(err, service.ErrInvalidFixtureWeek) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "week must be between 1 and 6"})
			return nil, nil, false
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, nil, false
	}

	return teams, fixtures, true
}

func bootstrapResponse(result service.LeagueBootstrap) gin.H {
	return gin.H{
		"league":    leagueStateResponse(result.League),
		"teams":     teamResponses(result.Teams),
		"fixtures":  fixtureResponses(result.Fixtures, result.Teams),
		"standings": standingResponses(result.Standings),
	}
}

func leagueStateResponse(state domain.LeagueState) gin.H {
	return gin.H{
		"current_week": state.CurrentWeek,
		"total_weeks":  state.TotalWeeks,
		"is_completed": state.IsCompleted,
		"created_at":   state.CreatedAt,
		"updated_at":   state.UpdatedAt,
	}
}

func teamResponses(teams []domain.Team) []gin.H {
	response := make([]gin.H, 0, len(teams))
	for _, team := range teams {
		response = append(response, gin.H{
			"id":         team.ID,
			"name":       team.Name,
			"strength":   team.Strength,
			"created_at": team.CreatedAt,
		})
	}

	return response
}

func fixtureResponses(fixtures []domain.Match, teams []domain.Team) []gin.H {
	teamNames := make(map[int64]string, len(teams))
	for _, team := range teams {
		teamNames[team.ID] = team.Name
	}

	response := make([]gin.H, 0, len(fixtures))
	for _, fixture := range fixtures {
		response = append(response, gin.H{
			"id":             fixture.ID,
			"week":           fixture.Week,
			"home_team_id":   fixture.HomeTeamID,
			"home_team_name": teamNames[fixture.HomeTeamID],
			"away_team_id":   fixture.AwayTeamID,
			"away_team_name": teamNames[fixture.AwayTeamID],
			"home_goals":     fixture.HomeGoals,
			"away_goals":     fixture.AwayGoals,
			"status":         fixture.Status,
			"played_at":      fixture.PlayedAt,
			"created_at":     fixture.CreatedAt,
			"updated_at":     fixture.UpdatedAt,
		})
	}

	return response
}

func standingResponses(standings []domain.Standing) []gin.H {
	response := make([]gin.H, 0, len(standings))
	for _, standing := range standings {
		response = append(response, gin.H{
			"team_id":         standing.TeamID,
			"team_name":       standing.TeamName,
			"played":          standing.Played,
			"wins":            standing.Wins,
			"draws":           standing.Draws,
			"losses":          standing.Losses,
			"goals_for":       standing.GoalsFor,
			"goals_against":   standing.GoalsAgainst,
			"goal_difference": standing.GoalDifference,
			"points":          standing.Points,
			"rank":            standing.Rank,
		})
	}

	return response
}
