package main

import (
	"context"
	"net/http"
	"os"

	"github.com/spikycham/feedme/internal/db"
	"github.com/spikycham/feedme/internal/middleware"
	"github.com/spikycham/feedme/internal/router"
	"github.com/spikycham/feedme/pkg/dotenv"
	"github.com/spikycham/feedme/pkg/logger"
)

func main() {
	ctx := context.Background()

	log := logger.New()

	// Load related environment variables.
	dotenv.Load(".env")
	port := os.Getenv("PORT")
	psqlUri := os.Getenv("POSTGRES_URI")

	// Connect to postgres database.
	pool, err := db.Connect(ctx, psqlUri)
	if err != nil {
		log.Error("failed to connect to postgres", err)
	}

	// Create the multiplexer and bind the routes.
	routr := router.New()
	routr.Register(pool)
	// Add all middlewares.
	handler := middleware.Register(routr)

	log.Info("start serving at port:", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Error("failed to start service")
	}
}
