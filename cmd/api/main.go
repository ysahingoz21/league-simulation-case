package main

import (
	"log"
	"net/http"

	"league-simulation-case/internal/app"
	"league-simulation-case/internal/config"
)

func main() {
	cfg := config.Load()
	router := app.NewRouter(cfg)
	addr := ":" + cfg.Port

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
