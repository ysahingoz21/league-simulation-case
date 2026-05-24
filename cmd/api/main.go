package main

import (
	"log"
	"net/http"

	"league-simulation-case/internal/app"
	"league-simulation-case/internal/config"
	"league-simulation-case/internal/database"
)

func main() {
	cfg := config.Load()
	db, err := database.OpenSQLite(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.ApplySchema(db, database.DefaultSchemaPath); err != nil {
		log.Fatal(err)
	}

	router := app.NewRouter(cfg, db)
	addr := ":" + cfg.Port

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
