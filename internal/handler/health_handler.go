package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/config"
)

type HealthHandler struct {
	config config.Config
}

func NewHealthHandler(cfg config.Config) HealthHandler {
	return HealthHandler{
		config: cfg,
	}
}

func (h HealthHandler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"message":     "League Simulation API is running",
		"environment": h.config.AppEnv,
	})
}
