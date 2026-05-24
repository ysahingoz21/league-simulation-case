package app

import (
	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/config"
	"league-simulation-case/internal/handler"
)

func NewRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()
	healthHandler := handler.NewHealthHandler(cfg)

	router.GET("/health", healthHandler.GetHealth)

	return router
}
