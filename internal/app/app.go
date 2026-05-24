package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"league-simulation-case/internal/config"
	"league-simulation-case/internal/handler"
	"league-simulation-case/internal/repository/sqlite"
	"league-simulation-case/internal/service"
)

func NewRouter(cfg config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()
	healthHandler := handler.NewHealthHandler(cfg)
	teamRepository := sqlite.NewTeamRepository(db)
	matchRepository := sqlite.NewMatchRepository(db)
	leagueRepository := sqlite.NewLeagueRepository(db)
	standingRepository := sqlite.NewStandingRepository(db)
	predictionRepository := sqlite.NewPredictionRepository(db)
	leagueService := service.NewLeagueService(
		teamRepository,
		matchRepository,
		leagueRepository,
		standingRepository,
		predictionRepository,
	)
	simulationService := service.NewSimulationService(
		teamRepository,
		matchRepository,
		leagueRepository,
		service.NewStrengthBasedSimulator(nil),
	)
	leagueHandler := handler.NewLeagueHandler(leagueService)
	simulationHandler := handler.NewSimulationHandler(simulationService, leagueService)

	router.GET("/health", healthHandler.GetHealth)
	router.POST("/api/v1/league/init", leagueHandler.InitializeLeague)
	router.POST("/api/v1/league/reset", leagueHandler.ResetLeague)
	router.GET("/api/v1/league", leagueHandler.GetLeague)
	router.GET("/api/v1/teams", leagueHandler.ListTeams)
	router.GET("/api/v1/fixtures", leagueHandler.ListFixtures)
	router.GET("/api/v1/fixtures/:week", leagueHandler.ListFixturesByWeek)
	router.POST("/api/v1/simulation/week/next", simulationHandler.PlayNextWeek)
	router.POST("/api/v1/simulation/week/:week", simulationHandler.PlayWeek)

	return router
}
