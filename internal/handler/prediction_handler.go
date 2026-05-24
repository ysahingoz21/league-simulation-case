package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/domain"
	"league-simulation-case/internal/service"
)

type PredictionHandler struct {
	predictionService service.PredictionService
}

func NewPredictionHandler(predictionService service.PredictionService) PredictionHandler {
	return PredictionHandler{predictionService: predictionService}
}

func (h PredictionHandler) GetPredictions(c *gin.Context) {
	result, err := h.predictionService.GetPredictions(c.Request.Context())
	if err != nil {
		if errors.Is(err, service.ErrLeagueNotInitialized) {
			c.JSON(http.StatusNotFound, gin.H{"error": "league is not initialized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"week":        result.Week,
		"predictions": predictionResponses(result.Predictions),
	}
	if result.Message != "" {
		response["message"] = result.Message
	}

	c.JSON(http.StatusOK, response)
}

func predictionResponses(predictions []domain.Prediction) []gin.H {
	response := make([]gin.H, 0, len(predictions))
	for _, prediction := range predictions {
		response = append(response, gin.H{
			"id":                       prediction.ID,
			"week":                     prediction.Week,
			"team_id":                  prediction.TeamID,
			"team_name":                prediction.TeamName,
			"championship_probability": prediction.ChampionshipProbability,
			"expected_points":          prediction.ExpectedPoints,
			"projected_rank":           prediction.ProjectedRank,
			"created_at":               prediction.CreatedAt,
		})
	}

	return response
}
