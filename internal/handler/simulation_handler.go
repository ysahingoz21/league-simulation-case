package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/service"
)

type SimulationHandler struct {
	simulationService service.SimulationService
	leagueService     service.LeagueService
}

func NewSimulationHandler(simulationService service.SimulationService, leagueService service.LeagueService) SimulationHandler {
	return SimulationHandler{
		simulationService: simulationService,
		leagueService:     leagueService,
	}
}

func (h SimulationHandler) PlayNextWeek(c *gin.Context) {
	result, err := h.simulationService.PlayNextWeek(c.Request.Context())
	if err != nil {
		h.writeSimulationError(c, err)
		return
	}

	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"week":      result.Week,
		"league":    leagueStateResponse(result.League),
		"matches":   fixtureResponses(result.Matches, teams),
		"standings": standingResponses(result.Standings),
	})
}

func (h SimulationHandler) PlayWeek(c *gin.Context) {
	week, err := strconv.Atoi(c.Param("week"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "week must be an integer"})
		return
	}

	result, err := h.simulationService.PlayWeek(c.Request.Context(), week)
	if err != nil {
		h.writeSimulationError(c, err)
		return
	}

	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"week":      result.Week,
		"league":    leagueStateResponse(result.League),
		"matches":   fixtureResponses(result.Matches, teams),
		"standings": standingResponses(result.Standings),
	})
}

func (h SimulationHandler) PlayAll(c *gin.Context) {
	result, err := h.simulationService.PlayAll(c.Request.Context())
	if err != nil {
		h.writeSimulationError(c, err)
		return
	}

	teams, err := h.leagueService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	weeks := make([]gin.H, 0, len(result.Weeks))
	for _, weekResult := range result.Weeks {
		weeks = append(weeks, gin.H{
			"week":    weekResult.Week,
			"matches": fixtureResponses(weekResult.Matches, teams),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"league":    leagueStateResponse(result.League),
		"weeks":     weeks,
		"standings": standingResponses(result.Standings),
	})
}

func (h SimulationHandler) writeSimulationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrLeagueNotInitialized):
		c.JSON(http.StatusNotFound, gin.H{"error": "league is not initialized"})
	case errors.Is(err, service.ErrLeagueCompleted):
		c.JSON(http.StatusBadRequest, gin.H{"error": "league is already completed"})
	case errors.Is(err, service.ErrInvalidFixtureWeek):
		c.JSON(http.StatusBadRequest, gin.H{"error": "week must be between 1 and 6"})
	case errors.Is(err, service.ErrWeekAlreadyPlayed):
		c.JSON(http.StatusBadRequest, gin.H{"error": "week has already been played"})
	case errors.Is(err, service.ErrWeekOutOfOrder):
		c.JSON(http.StatusBadRequest, gin.H{"error": "week must be played sequentially"})
	case errors.Is(err, service.ErrWeekHasPlayedMatch):
		c.JSON(http.StatusBadRequest, gin.H{"error": "week already contains played matches"})
	case errors.Is(err, service.ErrWeekHasNoMatches):
		c.JSON(http.StatusBadRequest, gin.H{"error": "week has no matches"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
