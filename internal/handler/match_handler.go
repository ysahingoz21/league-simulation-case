package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/service"
)

type MatchHandler struct {
	matchService  service.MatchService
	leagueService service.LeagueService
}

type updateMatchRequest struct {
	HomeGoals *int `json:"home_goals"`
	AwayGoals *int `json:"away_goals"`
}

func NewMatchHandler(matchService service.MatchService, leagueService service.LeagueService) MatchHandler {
	return MatchHandler{
		matchService:  matchService,
		leagueService: leagueService,
	}
}

func (h MatchHandler) GetMatch(c *gin.Context) {
	matchID, ok := parseMatchID(c)
	if !ok {
		return
	}

	match, err := h.matchService.GetMatch(c.Request.Context(), matchID)
	if err != nil {
		h.writeMatchError(c, err)
		return
	}

	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match": singleMatchResponse(match, teams),
	})
}

func (h MatchHandler) UpdateMatch(c *gin.Context) {
	matchID, ok := parseMatchID(c)
	if !ok {
		return
	}

	var req updateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.HomeGoals == nil || req.AwayGoals == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "home_goals and away_goals are required"})
		return
	}

	result, err := h.matchService.UpdateMatchResult(c.Request.Context(), matchID, *req.HomeGoals, *req.AwayGoals)
	if err != nil {
		h.writeMatchError(c, err)
		return
	}

	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match":       singleMatchResponse(result.Match, teams),
		"league":      leagueStateResponse(result.League),
		"standings":   standingResponses(result.Standings),
		"predictions": predictionResponses(result.Predictions),
	})
}

func (h MatchHandler) writeMatchError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrLeagueNotInitialized):
		c.JSON(http.StatusNotFound, gin.H{"error": "league is not initialized"})
	case errors.Is(err, service.ErrMatchNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
	case errors.Is(err, service.ErrInvalidGoals):
		c.JSON(http.StatusBadRequest, gin.H{"error": "goals must be between 0 and 20"})
	case errors.Is(err, service.ErrMatchEditNotAllowed):
		c.JSON(http.StatusBadRequest, gin.H{"error": "editing future scheduled matches is not allowed"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func parseMatchID(c *gin.Context) (int64, bool) {
	matchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || matchID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match id must be a positive integer"})
		return 0, false
	}

	return matchID, true
}

func singleMatchResponse(match domain.Match, teams []domain.Team) gin.H {
	teamNames := make(map[int64]string, len(teams))
	for _, team := range teams {
		teamNames[team.ID] = team.Name
	}

	return gin.H{
		"id":             match.ID,
		"week":           match.Week,
		"home_team_id":   match.HomeTeamID,
		"home_team_name": teamNames[match.HomeTeamID],
		"away_team_id":   match.AwayTeamID,
		"away_team_name": teamNames[match.AwayTeamID],
		"home_goals":     match.HomeGoals,
		"away_goals":     match.AwayGoals,
		"status":         match.Status,
		"played_at":      match.PlayedAt,
		"created_at":     match.CreatedAt,
		"updated_at":     match.UpdatedAt,
	}
}
